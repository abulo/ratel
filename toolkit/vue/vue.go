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
)

// Vue represents the upgrade command.
var Vue = &cobra.Command{
	Use:   "vue",
	Short: "前端指令",
	Long:  "前端手架命令 : toolkit vue",
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {
	// 数据初始化
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
		Message: "查询条件",
		Help:    "条件",
		Options: selectIndex,
	}, &useIndex); err != nil {
		return
	}

	var methodList []base.Method
	needPageMethodList := make([]string, 0)

	// 是否需要分页
	pageBool := false

	var deleteBool bool
	if delete == "yes" {
		deleteBool = true
	}
	//查询条件
	AllCondition := make([]base.Column, 0)

	// 数字长度
	strLen := strings.LastIndex(dir, "/")

	//添加默认方法
	methodList = append(methodList, base.Method{
		Table:          tableItem,
		TableColumn:    tableColumn,
		Type:           "Create",
		Name:           base.CamelStr(tableItem.TableName) + "Create",
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
		Name:           base.CamelStr(tableItem.TableName) + "Update",
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
		Name:           base.CamelStr(tableItem.TableName) + "Delete",
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
		Name:           base.CamelStr(tableItem.TableName),
		Condition:      nil,
		ConditionTotal: 0,
		Pkg:            dir[strLen+1:],
		Primary:        tablePrimary,
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
				Pkg:            dir[strLen+1:],
				Primary:        tablePrimary,
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
				Pkg:            dir[strLen+1:],
				Primary:        tablePrimary,
				Page:           pageBool,
				SoftDelete:     deleteBool,
			},
		)
		multiSelect = append(multiSelect,
			base.CamelStr(tableItem.TableName)+"Recover",
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
			Pkg:            dir[strLen+1:],
			Primary:        tablePrimary,
			Page:           pageBool,
			SoftDelete:     deleteBool,
		}
		multiSelect = append(multiSelect, methodName)
		methodList = append(methodList, method)
		needPageMethodList = append(needPageMethodList, methodName)
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
				Pkg:            dir[strLen+1:],
				ConditionTotal: len(condition),
				Primary:        tablePrimary,
				Page:           pageBool,
				SoftDelete:     deleteBool,
			}
			AllCondition = append(AllCondition, condition...)
			multiSelect = append(multiSelect, methodName)
			//添加到集合中
			methodList = append(methodList, method)
			if base.CamelStr(customIndexType) == "List" {
				needPageMethodList = append(needPageMethodList, methodName)
			}
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

	var newMethodList []base.Method
	for _, val := range methodList {
		newMethod := val
		// 判断是否分页
		if util.InArray(newMethod.Name, multiPageSelected) {
			newMethod.Page = true
		}
		if util.InArray(newMethod.Name, multiSelected) {
			newMethodList = append(newMethodList, newMethod)
		}
	}
	if len(multiPageSelected) > 0 {
		pageBool = true
	}

	// 数据模型
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
	if util.InArray("interface", tplSelected) {
		if err := survey.AskOne(&survey.Input{
			Message: "接口地址",
			Help:    "地址",
		}, &apiUrl); err != nil || apiUrl == "" {
			return
		}
	}
	viewUrl := ""
	if util.InArray("page", tplSelected) {
		if err := survey.AskOne(&survey.Input{
			Message: "Vue页面地址",
			Help:    "地址",
		}, &viewUrl); err != nil || viewUrl == "" {
			return
		}
	}

	if util.InArray("interface", tplSelected) {
		interfaceDir := path.Join(base.Path, base.Config.String("vue.InterfaceDir"))
		GenerateInterface(moduleParam, interfaceDir, tableName)
	}
	if util.InArray("module", tplSelected) {
		methodDir := path.Join(base.Path, base.Config.String("vue.ModulesDir"))
		GenerateMethod(moduleParam, apiUrl, methodDir, tableName)
	}
	if util.InArray("page", tplSelected) {
		pageDir := path.Join(base.Path, base.Config.String("vue.PageDir"))
		GeneratePage(moduleParam, pageDir, viewUrl, tableName)
	}
}
