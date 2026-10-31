package module

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// CmdNew 定义 toolkit mp 子命令，运行时由 Run 完成数据模型层代码的交互式生成。
var (
	CmdNew = &cobra.Command{
		Use:   "mp",
		Short: "数据模型层",
		Long:  "创建数据库模型层: toolkit mp dir table_name",
		Run:   Run,
	}
)

// Run 执行 toolkit mp 命令：初始化运行环境，交互式采集模型路径、表名、软删除策略与索引，
// 读取数据库元数据并组装方法列表后，按用户选择生成 module、proto、service、convert 代码。
//
// 任一前置步骤失败或交互输入为空时直接返回，不生成文件；数据库查询的超时时间固定为 60 秒。
func Run(cmd *cobra.Command, args []string) {
	if err := base.InitBase(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	dirModule := path.Join(cast.ToString(base.Path), "module")
	_ = os.MkdirAll(dirModule, os.ModePerm)
	dirProto := path.Join(cast.ToString(base.Path), "proto")
	_ = os.MkdirAll(dirProto, os.ModePerm)

	dir := ""
	tableName := ""
	multiSelect := make([]string, 0)
	delete := ""
	if err := survey.AskOne(&survey.Input{
		Message: "模型路径",
		Help:    "文件夹路径",
	}, &dir); err != nil || dir == "" {
		return
	}
	if err := survey.AskOne(&survey.Input{
		Message: "表名称",
		Help:    "数据库中某个表名称",
	}, &tableName); err != nil || tableName == "" {
		return
	}
	if err := survey.AskOne(&survey.Select{
		Message: "是否软删除",
		Help:    "是否软删除",
		Options: []string{"yes", "no"},
	}, &delete); err != nil || delete == "" {
		return
	}
	fullModuleDir := path.Join(cast.ToString(base.Path), "module", dir)
	_ = os.MkdirAll(fullModuleDir, os.ModePerm)

	fullProtoDir := path.Join(cast.ToString(base.Path), "proto")
	_ = os.MkdirAll(fullProtoDir, os.ModePerm)

	fullServiceDir := path.Join(cast.ToString(base.Path), "service", dir)
	_ = os.MkdirAll(fullServiceDir, os.ModePerm)

	fullConvertDir := path.Join(cast.ToString(base.Path), "service", dir)
	_ = os.MkdirAll(fullConvertDir, os.ModePerm)

	timeout := "60s"
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	// PERF-8-19 [P2] 四次元数据查询互不依赖，使用 errgroup 并发发起；错误按原查询顺序依次上报，
	// 错误提示文案与“首个失败即返回”的语义保持不变。（预期收益需实测验证）
	var (
		tableColumn                     []base.Column
		tableItem                       base.Table
		tableIndex                      []base.Index
		tablePrimary                    base.Column
		colErr, itemErr, idxErr, priErr error
	)
	var eg errgroup.Group
	eg.Go(func() error {
		tableColumn, colErr = base.TableColumn(ctx, base.Config.String("db.Database"), tableName)
		return nil
	})
	eg.Go(func() error {
		tableItem, itemErr = base.TableItem(ctx, base.Config.String("db.Database"), tableName)
		return nil
	})
	eg.Go(func() error {
		tableIndex, idxErr = base.TableIndex(ctx, base.Config.String("db.Database"), tableName)
		return nil
	})
	eg.Go(func() error {
		tablePrimary, priErr = base.TablePrimary(ctx, base.Config.String("db.Database"), tableName)
		return nil
	})
	_ = eg.Wait()
	if colErr != nil {
		fmt.Println("表结构信息:", color.RedString(colErr.Error()))
		return
	}
	if itemErr != nil {
		fmt.Println("表信息:", color.RedString(itemErr.Error()))
		return
	}
	if idxErr != nil {
		fmt.Println("表索引:", color.RedString(idxErr.Error()))
		return
	}
	if priErr != nil {
		fmt.Println("表主键:", color.RedString(priErr.Error()))
		return
	}

	tableColumnMap := make(map[string]base.Column)
	for _, item := range tableColumn {
		tableColumnMap[item.ColumnName] = item
	}

	// PERF-8-6 [P2] 表名大驼峰在命令级计算一次，后续方法名/注册片段统一复用，避免重复转换。
	tableCamel := base.CamelStr(tableItem.TableName)
	// PERF-8-7 [P2] 索引字段去重与成员判断预建集合，循环内 O(1) 判定，替代反射型 InArray 线性查找。
	selectIndex := make([]string, 0)
	selectIndexSet := make(map[string]struct{})
	for _, item := range tableIndex {
		tmpSlice := util.Explode(",", item.Field)
		for _, v := range tmpSlice {
			if _, ok := selectIndexSet[v]; ok {
				continue
			}
			selectIndexSet[v] = struct{}{}
			selectIndex = append(selectIndex, v)
		}
	}
	useIndex := make([]string, 0)
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "选择索引",
		Help:    "索引列表",
		Options: selectIndex,
	}, &useIndex); err != nil {
		return
	}
	useIndexSet := make(map[string]struct{}, len(useIndex))
	for _, name := range useIndex {
		useIndexSet[name] = struct{}{}
	}
	mod, err := base.ModulePath(path.Join(cast.ToString(base.Path), "go.mod"))
	if err != nil {
		fmt.Println("go.mod文件不存在:", color.RedString(err.Error()))
		return
	}
	strLen := strings.LastIndex(dir, "/")

	var methodList []base.Method

	pageBool := false

	var deleteBool bool
	if delete == "yes" {
		deleteBool = true
	}

	// 追加不依赖查询条件的默认方法：新增、更新、删除、详情。
	methodList = append(methodList, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Create",
		Name:           tableCamel + "Create",
		Condition:      nil,
		ConditionTotal: 0,
		Primary:        tablePrimary,
		Pkg:            dir[strLen+1:],
		PkgPath:        dir,
		ModName:        mod,
		Page:           pageBool,
		SoftDelete:     deleteBool,
	}, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Update",
		Name:           tableCamel + "Update",
		Condition:      nil,
		ConditionTotal: 0,
		Primary:        tablePrimary,
		Pkg:            dir[strLen+1:],
		PkgPath:        dir,
		ModName:        mod,
		Page:           pageBool,
		SoftDelete:     deleteBool,
	}, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Delete",
		Name:           tableCamel + "Delete",
		Condition:      nil,
		ConditionTotal: 0,
		Primary:        tablePrimary,
		Pkg:            dir[strLen+1:],
		PkgPath:        dir,
		ModName:        mod,
		Page:           pageBool,
		SoftDelete:     deleteBool,
	}, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Show",
		Name:           tableCamel,
		Condition:      nil,
		ConditionTotal: 0,
		Primary:        tablePrimary,
		Pkg:            dir[strLen+1:],
		PkgPath:        dir,
		ModName:        mod,
		Page:           pageBool,
		SoftDelete:     deleteBool,
	})

	multiSelect = append(multiSelect,
		tableCamel+"Create",
		tableCamel+"Update",
		tableCamel+"Delete",
		tableCamel,
	)

	if deleteBool {
		methodList = append(methodList,
			base.Method{
				Table:          tableItem,
				TableColumn:    tableColumn,
				Type:           "Recover",
				Name:           tableCamel + "Recover",
				Condition:      nil,
				ConditionTotal: 0,
				Primary:        tablePrimary,
				Pkg:            dir[strLen+1:],
				PkgPath:        dir,
				ModName:        mod,
				Page:           pageBool,
				SoftDelete:     deleteBool,
			},
		)
		multiSelect = append(multiSelect,
			tableCamel+"Recover",
		)
		methodList = append(methodList,
			base.Method{
				Table:          tableItem,
				TableColumn:    tableColumn,
				Type:           "Drop",
				Name:           tableCamel + "Drop",
				Condition:      nil,
				ConditionTotal: 0,
				Primary:        tablePrimary,
				Pkg:            dir[strLen+1:],
				PkgPath:        dir,
				ModName:        mod,
				Page:           pageBool,
				SoftDelete:     deleteBool,
			},
		)
		multiSelect = append(multiSelect,
			tableCamel+"Drop",
		)
	}
	if err != nil {
		methodName := tableCamel + "List"
		method := base.Method{
			Table:          tableItem,
			TableColumn:    tableColumn,
			Type:           "List",
			Name:           methodName,
			Condition:      nil,
			ConditionTotal: 0,
			Primary:        tablePrimary,
			Pkg:            dir[strLen+1:],
			PkgPath:        dir,
			ModName:        mod,
			Page:           pageBool,
			SoftDelete:     deleteBool,
		}
		multiSelect = append(multiSelect, methodName)
		methodList = append(methodList, method)
	} else {
		field := make([]string, 0)
		// PERF-8-7 [P2] 字段去重集合与索引类型白名单集合，替代循环内 InArray 线性扫描。
		fieldSet := make(map[string]struct{})
		idxSet := map[string]struct{}{"Item": {}, "List": {}}
		for _, v := range tableIndex {
			condition := make([]base.Column, 0)
			indexField := v.Field
			indexFieldSlice := util.Explode(",", indexField)
			for _, fieldValue := range indexFieldSlice {
				if _, ok := useIndexSet[fieldValue]; !ok {
					continue
				}
				positionIndex := cast.ToInt64(len(condition)) + 1
				currentColumn := tableColumnMap[fieldValue]
				currentColumn.PosiTion = positionIndex
				condition = append(condition, currentColumn)
				if _, ok := fieldSet[fieldValue]; !ok {
					fieldSet[fieldValue] = struct{}{}
					field = append(field, fieldValue)
				}
			}
			indexName := v.IndexName

			// 索引名按下划线拆分为方法类型与方法名两段，段数不足时跳过该索引。
			indexNameSlice := util.Explode("_", indexName)
			if len(indexNameSlice) < 2 {
				continue
			}
			customIndexType := util.UCWords(indexNameSlice[0])
			customIndexName := util.UCWords(indexNameSlice[1])
			if _, ok := idxSet[customIndexType]; !ok {
				continue
			}
			methodName := tableCamel + base.CamelStr(customIndexName)
			method := base.Method{
				Table:          tableItem,
				TableColumn:    tableColumn,
				Type:           customIndexType,
				Name:           methodName,
				Condition:      condition,
				ConditionTotal: len(condition),
				Primary:        tablePrimary,
				Pkg:            dir[strLen+1:],
				PkgPath:        dir,
				ModName:        mod,
				Page:           pageBool,
				SoftDelete:     deleteBool,
			}
			multiSelect = append(multiSelect, methodName)
			methodList = append(methodList, method)
		}
		condition := make([]base.Column, 0)
		for _, fieldValue := range field {
			positionIndex := cast.ToInt64(len(condition)) + 1
			currentColumn := tableColumnMap[fieldValue]
			currentColumn.PosiTion = positionIndex
			condition = append(condition, currentColumn)
		}
		methodName := tableCamel + "List"
		method := base.Method{
			Table:          tableItem,
			TableColumn:    tableColumn,
			Type:           "List",
			Name:           methodName,
			Condition:      condition,
			ConditionTotal: len(condition),
			Primary:        tablePrimary,
			Pkg:            dir[strLen+1:],
			PkgPath:        dir,
			ModName:        mod,
			Page:           pageBool,
			SoftDelete:     deleteBool,
		}
		multiSelect = append(multiSelect, methodName)
		methodList = append(methodList, method)
	}

	multiSelected := make([]string, 0)
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "方法",
		Help:    "方法列表",
		Options: multiSelect,
	}, &multiSelected); err != nil || len(multiSelected) == 0 {
		return
	}

	tpl := make([]string, 0)
	tpl = append(tpl, "module", "proto", "service", "convert")

	tplSelected := make([]string, 0)
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "模板",
		Help:    "模板列表",
		Options: tpl,
	}, &tplSelected); err != nil || len(tplSelected) == 0 {
		return
	}

	// PERF-8-7 [P2] 用户勾选结果预建集合，方法过滤与模板勾选循环内 O(1) 命中判断。
	multiSelectedSet := make(map[string]struct{}, len(multiSelected))
	for _, name := range multiSelected {
		multiSelectedSet[name] = struct{}{}
	}
	tplSelectedSet := make(map[string]struct{}, len(tplSelected))
	for _, name := range tplSelected {
		tplSelectedSet[name] = struct{}{}
	}

	var newMethodList []base.Method
	for _, val := range methodList {
		if _, ok := multiSelectedSet[val.Name]; ok {
			newMethod := val
			newMethod.Page = true
			newMethodList = append(newMethodList, newMethod)
		}
	}
	pageBool = true
	moduleParam := base.ModuleParam{
		Pkg:         dir[strLen+1:],
		PkgPath:     dir,
		ModName:     mod,
		Primary:     tablePrimary,
		Table:       tableItem,
		TableColumn: tableColumn,
		Method:      newMethodList,
		Page:        pageBool,
		SoftDelete:  deleteBool,
	}
	if _, ok := tplSelectedSet["module"]; ok {
		GenerateModule(moduleParam, fullModuleDir, tableName)
	}

	if _, ok := tplSelectedSet["proto"]; ok {
		GenerateProto(moduleParam, fullProtoDir, fullServiceDir, tableName)
	}
	if _, ok := tplSelectedSet["convert"]; ok {
		GenerateConvert(moduleParam, fullConvertDir, tableName)
	}
	if _, ok := tplSelectedSet["service"]; ok {
		GenerateService(moduleParam, fullServiceDir, tableName)
		strLenSlot := util.Explode("/", dir)
		pkgName := strLenSlot[len(strLenSlot)-1]
		if len(strLenSlot) > 1 {
			pkgName = strLenSlot[len(strLenSlot)-2]
		}

		builder := strings.Builder{}
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("package %s", pkgName))
		builder.WriteString("\n")
		builder.WriteString("import (")
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("	\"%s/service/%s\"", moduleParam.ModName, moduleParam.PkgPath))
		builder.WriteString("\n")
		builder.WriteString("\n")
		builder.WriteString("	\"github.com/abulo/ratel/v3/server/xgrpc\"")
		builder.WriteString("\n")
		builder.WriteString(")")
		builder.WriteString("\n")
		builder.WriteString("func Registry(server *xgrpc.Server) {")
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("	// %s->%s", moduleParam.Table.TableComment, moduleParam.Table.TableName))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("	%s.Register%sServiceServer(server.Server, &%s.Srv%sServiceServer{", moduleParam.Pkg, tableCamel, moduleParam.Pkg, tableCamel))
		builder.WriteString("\n")
		builder.WriteString("		Server: server,")
		builder.WriteString("\n")
		builder.WriteString("	})")
		builder.WriteString("\n")
		builder.WriteString("}")
		fmt.Println(builder.String())
	}
}
