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
)

var (
	// CmdNew represents the new command.
	CmdNew = &cobra.Command{
		Use:   "mp",
		Short: "数据模型层",
		Long:  "创建数据库模型层: toolkit mp dir table_name",
		Run:   Run,
	}
)

func Run(cmd *cobra.Command, args []string) {
	// 数据初始化
	if err := base.InitBase(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	// 创建文件夹
	dirModule := path.Join(base.Path, "module")
	_ = os.MkdirAll(dirModule, os.ModePerm)
	// 创建文件夹
	dirProto := path.Join(base.Path, "proto")
	_ = os.MkdirAll(dirProto, os.ModePerm)

	//创建数据
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
	// 文件夹的路径
	fullModuleDir := path.Join(base.Path, "module", dir)
	_ = os.MkdirAll(fullModuleDir, os.ModePerm)

	// 文件夹的路径
	fullProtoDir := path.Join(base.Path, "proto")
	_ = os.MkdirAll(fullProtoDir, os.ModePerm)

	// 文件夹的路径
	fullServiceDir := path.Join(base.Path, "service", dir)
	_ = os.MkdirAll(fullServiceDir, os.ModePerm)

	fullConvertDir := path.Join(base.Path, "service", dir)
	_ = os.MkdirAll(fullConvertDir, os.ModePerm)

	// 初始化上下文
	timeout := "60s"
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	// 表结构信息
	tableColumn, err := base.TableColumn(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表结构信息:", color.RedString(err.Error()))
		return
	}
	tableColumnMap := make(map[string]base.Column)
	for _, item := range tableColumn {
		tableColumnMap[item.ColumnName] = item
	}
	// 表信息
	tableItem, err := base.TableItem(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表信息:", color.RedString(err.Error()))
		return
	}
	// 表索引
	tableIndex, err := base.TableIndex(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表索引:", color.RedString(err.Error()))
		return
	}
	// 表主键
	tablePrimary, err := base.TablePrimary(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表主键:", color.RedString(err.Error()))
		return
	}
	selectIndex := make([]string, 0)
	for _, item := range tableIndex {
		tmpSlice := util.Explode(",", item.Field)
		for _, v := range tmpSlice {
			if util.InArray(v, selectIndex) {
				continue
			}
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
	//获取 go.mod
	mod, err := base.ModulePath(path.Join(base.Path, "go.mod"))
	if err != nil {
		fmt.Println("go.mod文件不存在:", color.RedString(err.Error()))
		return
	}
	// 数字长度
	strLen := strings.LastIndex(dir, "/")

	var methodList []base.Method
	// needPageMethodList := make([]string, 0)

	// 是否需要分页
	pageBool := false

	var deleteBool bool
	if delete == "yes" {
		deleteBool = true
	}

	//添加默认方法
	methodList = append(methodList, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Create",
		Name:           base.CamelStr(tableItem.TableName) + "Create",
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
		Name:           base.CamelStr(tableItem.TableName) + "Update",
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
		Name:           base.CamelStr(tableItem.TableName) + "Delete",
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
		Name:           base.CamelStr(tableItem.TableName),
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
		base.CamelStr(tableItem.TableName)+"Create",
		base.CamelStr(tableItem.TableName)+"Update",
		base.CamelStr(tableItem.TableName)+"Delete",
		base.CamelStr(tableItem.TableName),
	)

	if deleteBool {
		methodList = append(methodList,
			base.Method{
				Table:          tableItem,
				TableColumn:    tableColumn,
				Type:           "Recover",
				Name:           base.CamelStr(tableItem.TableName) + "Recover",
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
			base.CamelStr(tableItem.TableName)+"Recover",
		)
		methodList = append(methodList,
			base.Method{
				Table:          tableItem,
				TableColumn:    tableColumn,
				Type:           "Drop",
				Name:           base.CamelStr(tableItem.TableName) + "Drop",
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
			base.CamelStr(tableItem.TableName)+"Drop",
		)
	}
	//获取的索引信息没有
	if err != nil {
		methodName := base.CamelStr(tableItem.TableName) + "List"
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
		// needPageMethodList = append(needPageMethodList, methodName)
	} else {
		// js, _ := json.Marshal(tableIndex)
		// fmt.Println(string(js))
		//存储条件信息
		field := make([]string, 0)
		idx := []string{"Item", "List"}
		//有索引信息
		for _, v := range tableIndex {
			//查询条件
			condition := make([]base.Column, 0)
			//数据库索引
			indexField := v.Field
			indexFieldSlice := util.Explode(",", indexField)
			for _, fieldValue := range indexFieldSlice {
				if !util.InArray(fieldValue, useIndex) {
					continue
				}
				//构造查询条件
				positionIndex := cast.ToInt64(len(condition)) + 1
				currentColumn := tableColumnMap[fieldValue]
				currentColumn.PosiTion = positionIndex
				condition = append(condition, currentColumn)
				if !util.InArray(fieldValue, field) {
					field = append(field, fieldValue)
				}
			}
			// 数据库中的索引名称
			indexName := v.IndexName

			// 拆分字符串,得到索引类型和索引名称
			indexNameSlice := util.Explode(":", indexName)
			if len(indexNameSlice) < 2 {
				continue
			}
			// 自定义函数名称和索引信息
			customIndexType := util.UCWords(indexNameSlice[0])
			customIndexName := util.UCWords(indexNameSlice[1])
			if !util.InArray(customIndexType, idx) {
				continue
			}
			methodName := base.CamelStr(tableItem.TableName) + base.CamelStr(customIndexName)
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
			//添加到集合中
			methodList = append(methodList, method)
			// if base.CamelStr(customIndexType) == "List" {
			// 	needPageMethodList = append(needPageMethodList, methodName)
			// }
		}
		condition := make([]base.Column, 0)
		for _, fieldValue := range field {
			//构造查询条件
			positionIndex := cast.ToInt64(len(condition)) + 1
			currentColumn := tableColumnMap[fieldValue]
			currentColumn.PosiTion = positionIndex
			condition = append(condition, currentColumn)
		}
		methodName := base.CamelStr(tableItem.TableName) + "List"
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
		// needPageMethodList = append(needPageMethodList, methodName)
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

	var newMethodList []base.Method
	for _, val := range methodList {
		if util.InArray(val.Name, multiSelected) {
			newMethod := val
			// 判断是否分页
			newMethod.Page = true
			newMethodList = append(newMethodList, newMethod)
		}
	}
	pageBool = true
	// 数据模型
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
	if util.InArray("module", tplSelected) {
		GenerateModule(moduleParam, fullModuleDir, tableName)
	}

	if util.InArray("proto", tplSelected) {
		GenerateProto(moduleParam, fullProtoDir, fullServiceDir, tableName)
	}
	if util.InArray("convert", tplSelected) {
		GenerateConvert(moduleParam, fullConvertDir, tableName)
	}
	if util.InArray("service", tplSelected) {
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
		builder.WriteString(fmt.Sprintf("	%s.Register%sServiceServer(server.Server, &%s.Srv%sServiceServer{", moduleParam.Pkg, base.CamelStr(moduleParam.Table.TableName), moduleParam.Pkg, base.CamelStr(moduleParam.Table.TableName)))
		builder.WriteString("\n")
		builder.WriteString("		Server: server,")
		builder.WriteString("\n")
		builder.WriteString("	})")
		builder.WriteString("\n")
		builder.WriteString("}")
		fmt.Println(builder.String())
	}
}
