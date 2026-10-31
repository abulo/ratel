package vue

import (
	"context"
	"fmt"
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

// Vue 是 toolkit vue 子命令，交互式采集表与模板选项后生成前端接口定义、API 模块与页面代码。
var Vue = &cobra.Command{
	Use:   "vue",
	Short: "前端指令",
	Long:  "前端手架命令 : toolkit vue",
	Run:   Run,
}

// Run 执行 vue 子命令：初始化运行环境，依次询问表名、软删除策略、模型路径、所用索引、
// 待生成方法、分页方法与模板类型，读取数据库元数据后，按勾选的模板调用 GenerateInterface、
// GenerateMethod 或 GeneratePage 生成前端文件。
//
// 任一交互式输入被取消，或环境初始化、元数据读取失败时直接返回，不生成文件。
func Run(cmd *cobra.Command, args []string) {
	if err := base.InitBase(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	tableName := ""
	delete := ""
	dir := ""

	multiSelect := make([]string, 0)
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
	if err := survey.AskOne(&survey.Input{
		Message: "模型路径",
		Help:    "文件夹路径",
	}, &dir); err != nil || dir == "" {
		return
	}

	timeout := "60s"
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	// PERF-8-23 [P2] 四次元数据查询互不依赖，使用 errgroup 并发发起；错误按原查询顺序依次上报，
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

	// PERF-8-6 [P2] 表名大驼峰在命令级计算一次，后续方法名统一复用，避免重复转换。
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
		Message: "查询条件",
		Help:    "条件",
		Options: selectIndex,
	}, &useIndex); err != nil {
		return
	}
	useIndexSet := make(map[string]struct{}, len(useIndex))
	for _, name := range useIndex {
		useIndexSet[name] = struct{}{}
	}

	var methodList []base.Method
	needPageMethodList := make([]string, 0)

	pageBool := false

	var deleteBool bool
	if delete == "yes" {
		deleteBool = true
	}
	AllCondition := make([]base.Column, 0)

	strLen := strings.LastIndex(dir, "/")

	methodList = append(methodList, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Create",
		Name:           tableCamel + "Create",
		Condition:      nil,
		ConditionTotal: 0,
		Pkg:            dir[strLen+1:],
		Primary:        tablePrimary,
		Page:           pageBool,
		SoftDelete:     deleteBool,
	}, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Update",
		Name:           tableCamel + "Update",
		Condition:      nil,
		ConditionTotal: 0,
		Pkg:            dir[strLen+1:],
		Primary:        tablePrimary,
		Page:           pageBool,
		SoftDelete:     deleteBool,
	}, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Delete",
		Name:           tableCamel + "Delete",
		Condition:      nil,
		ConditionTotal: 0,
		Pkg:            dir[strLen+1:],
		Primary:        tablePrimary,
		Page:           pageBool,
		SoftDelete:     deleteBool,
	}, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Show",
		Name:           tableCamel,
		Condition:      nil,
		ConditionTotal: 0,
		Pkg:            dir[strLen+1:],
		Primary:        tablePrimary,
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
				Pkg:            dir[strLen+1:],
				Primary:        tablePrimary,
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
				Pkg:            dir[strLen+1:],
				Primary:        tablePrimary,
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
			Pkg:            dir[strLen+1:],
			Primary:        tablePrimary,
			Page:           pageBool,
			SoftDelete:     deleteBool,
		}
		multiSelect = append(multiSelect, methodName)
		methodList = append(methodList, method)
		needPageMethodList = append(needPageMethodList, methodName)
	} else {
		field := make([]string, 0)
		// PERF-8-7 [P2] 字段去重集合与索引类型白名单集合，替代循环内 InArray 线性扫描。
		fieldSet := make(map[string]struct{})
		// 索引名约定为「类型_名称」，仅识别 Item、List 两类前缀，分别生成单条查询与列表方法。
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
				Pkg:            dir[strLen+1:],
				ConditionTotal: len(condition),
				Primary:        tablePrimary,
				Page:           pageBool,
				SoftDelete:     deleteBool,
			}
			AllCondition = append(AllCondition, condition...)
			multiSelect = append(multiSelect, methodName)
			methodList = append(methodList, method)
			if base.CamelStr(customIndexType) == "List" {
				needPageMethodList = append(needPageMethodList, methodName)
			}
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
			Pkg:            dir[strLen+1:],
			Condition:      condition,
			ConditionTotal: len(condition),
			Primary:        tablePrimary,
			Page:           pageBool,
			SoftDelete:     deleteBool,
		}
		AllCondition = append(AllCondition, condition...)
		multiSelect = append(multiSelect, methodName)
		methodList = append(methodList, method)
		needPageMethodList = append(needPageMethodList, methodName)
	}

	multiSelected := make([]string, 0)
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "方法",
		Help:    "方法列表",
		Options: multiSelect,
	}, &multiSelected); err != nil || len(multiSelected) == 0 {
		return
	}

	multiPageSelected := make([]string, 0)
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "分页方法",
		Help:    "方法列表",
		Options: needPageMethodList,
	}, &multiPageSelected); err != nil {
		return
	}

	tpl := make([]string, 0)
	tpl = append(tpl, "module", "interface", "page")

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
	multiPageSelectedSet := make(map[string]struct{}, len(multiPageSelected))
	for _, name := range multiPageSelected {
		multiPageSelectedSet[name] = struct{}{}
	}
	tplSelectedSet := make(map[string]struct{}, len(tplSelected))
	for _, name := range tplSelected {
		tplSelectedSet[name] = struct{}{}
	}

	var newMethodList []base.Method
	for _, val := range methodList {
		newMethod := val
		if _, ok := multiPageSelectedSet[newMethod.Name]; ok {
			newMethod.Page = true
		}
		if _, ok := multiSelectedSet[newMethod.Name]; ok {
			newMethodList = append(newMethodList, newMethod)
		}
	}
	if len(multiPageSelected) > 0 {
		pageBool = true
	}

	moduleParam := base.ModuleParam{
		Pkg:         dir[strLen+1:],
		Primary:     tablePrimary,
		Table:       tableItem,
		TableColumn: tableColumn,
		Method:      newMethodList,
		Page:        pageBool,
		SoftDelete:  deleteBool,
		Condition:   AllCondition,
	}
	apiUrl := ""
	if _, ok := tplSelectedSet["interface"]; ok {
		if err := survey.AskOne(&survey.Input{
			Message: "接口地址",
			Help:    "地址",
		}, &apiUrl); err != nil || apiUrl == "" {
			return
		}
	}
	viewUrl := ""
	if _, ok := tplSelectedSet["page"]; ok {
		if err := survey.AskOne(&survey.Input{
			Message: "Vue页面地址",
			Help:    "地址",
		}, &viewUrl); err != nil || viewUrl == "" {
			return
		}
	}

	if _, ok := tplSelectedSet["interface"]; ok {
		interfaceDir := path.Join(cast.ToString(base.Path), base.Config.String("vue.InterfaceDir"))
		GenerateInterface(moduleParam, interfaceDir, tableName)
	}
	if _, ok := tplSelectedSet["module"]; ok {
		methodDir := path.Join(cast.ToString(base.Path), base.Config.String("vue.ModuleDir"))
		GenerateMethod(moduleParam, apiUrl, methodDir, tableName)
	}
	if _, ok := tplSelectedSet["page"]; ok {
		pageDir := path.Join(cast.ToString(base.Path), base.Config.String("vue.PageDir"))
		GeneratePage(moduleParam, pageDir, viewUrl, tableName)
	}
}
