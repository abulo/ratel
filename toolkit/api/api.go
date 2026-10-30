package api

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// CmdNew 是 toolkit api 子命令，交互式采集接口与数据表参数后生成 HTTP 接口代码。
var (
	CmdNew = &cobra.Command{
		Use:   "api",
		Short: "接口对象",
		Long:  "创建接口对象: toolkit api",
		Run:   Run,
	}
)

// Run 执行 api 子命令：初始化运行环境，依次询问接口路径、表名、驱动类型、接口地址、
// 软删除策略、所用索引以及待生成的方法与分页选项，读取数据库元数据后调用 Generate
// 渲染接口文件。
//
// 任一交互式输入被取消，或环境初始化、元数据读取、模板渲染失败时直接返回，不生成文件。
func Run(cmd *cobra.Command, args []string) {
	if err := base.InitBase(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	dir := ""
	tableName := ""
	apiType := ""
	apiUrl := ""
	multiSelect := make([]string, 0)
	delete := ""
	if err := survey.AskOne(&survey.Input{
		Message: "接口路径",
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
		Message: "驱动类型",
		Help:    "选择驱动类型",
		Options: []string{"gin", "hertz"},
	}, &apiType); err != nil || apiType == "" {
		return
	}
	if err := survey.AskOne(&survey.Input{
		Message: "接口地址",
		Help:    "地址",
	}, &apiUrl); err != nil || apiUrl == "" {
		return
	}
	if err := survey.AskOne(&survey.Select{
		Message: "是否软删除",
		Help:    "是否软删除",
		Options: []string{"yes", "no"},
	}, &delete); err != nil || delete == "" {
		return
	}

	fullApiDir := path.Join(cast.ToString(cast.ToString(base.Path)), "api", dir)
	_ = os.MkdirAll(fullApiDir, os.ModePerm)

	timeout := "60s"
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	tableColumn, err := base.TableColumn(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表结构信息:", color.RedString(err.Error()))
		return
	}
	tableColumnMap := make(map[string]base.Column)
	for _, item := range tableColumn {
		tableColumnMap[item.ColumnName] = item
	}
	tableItem, err := base.TableItem(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表信息:", color.RedString(err.Error()))
		return
	}
	tableIndex, err := base.TableIndex(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表索引:", color.RedString(err.Error()))
		return
	}
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

	mod, err := base.ModulePath(path.Join(cast.ToString(base.Path), "go.mod"))
	if err != nil {
		fmt.Println("go.mod文件不存在:", color.RedString(err.Error()))
		return
	}
	strLen := strings.LastIndex(dir, "/")

	var methodList []base.Method
	needPageMethodList := make([]string, 0)

	pageBool := false

	var deleteBool bool
	if delete == "yes" {
		deleteBool = true
	}

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
		needPageMethodList = append(needPageMethodList, methodName)
	} else {
		field := make([]string, 0)
		// 索引名约定为「类型_名称」，仅识别 Item、List 两类前缀，分别生成单条查询与列表方法。
		idx := []string{"Item", "List"}
		for _, v := range tableIndex {
			condition := make([]base.Column, 0)
			indexField := v.Field
			indexFieldSlice := util.Explode(",", indexField)
			for _, fieldValue := range indexFieldSlice {
				if !util.InArray(fieldValue, useIndex) {
					continue
				}
				positionIndex := cast.ToInt64(len(condition)) + 1
				currentColumn := tableColumnMap[fieldValue]
				currentColumn.PosiTion = positionIndex
				condition = append(condition, currentColumn)
				if !util.InArray(fieldValue, field) {
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

	var newMethodList []base.Method
	for _, val := range methodList {
		newMethod := val
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
	moduleParam := base.ModuleParam{
		Pkg:         dir[strLen+1:],
		PkgPath:     dir,
		ModName:     mod,
		Primary:     tablePrimary,
		Table:       tableItem,
		TableColumn: tableColumn,
		Method:      newMethodList,
		Page:        pageBool,
	}
	Generate(moduleParam, fullApiDir, tableName, dir[strLen+1:], dir, apiType, apiUrl)
}

// Generate 按 apiType 选择 hertz 或 gin 模板，渲染 moduleParam 并写入 fullApiDir 下
// 以表名命名的 Go 文件；同名文件已存在时先删除再创建。
//
// 渲染完成后在 fullApiDir 目录执行 go fmt 与 goimports 整理格式，并在标准输出打印
// 各方法对应的路由注册代码片段。文件创建、模板解析或格式化失败时打印错误并返回，
// 不向调用方返回错误。apiType 非 hertz 时一律按 gin 模板处理。
func Generate(moduleParam base.ModuleParam, fullApiDir, tableName, pkg, pkgPath, apiType, apiUrl string) {
	var tplString string
	if apiType == "hertz" {
		tplString = HertzTemplate()
	} else {
		tplString = GinTemplate()
	}
	tpl := template.Must(template.New("api").Funcs(template.FuncMap{
		"Convert":               base.Convert,
		"SymbolChar":            base.SymbolChar,
		"Char":                  base.Char,
		"Helper":                base.Helper,
		"CamelStr":              base.CamelStr,
		"Add":                   base.Add,
		"ModuleProtoConvertDao": base.ModuleProtoConvertDao,
		"ModuleDaoConvertProto": base.ModuleDaoConvertProto,
		"ModuleProtoConvertMap": base.ModuleProtoConvertMap,
		"ApiToProto":            base.ApiToProto,
	}).Parse(tplString))
	outApiFile := path.Join(fullApiDir, tableName+".go")
	if util.FileExists(outApiFile) {
		util.Delete(outApiFile)
	}
	file, err := os.OpenFile(outApiFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}
	err = tpl.Execute(file, moduleParam)
	if err != nil {
		fmt.Println("模板解析错误:", color.RedString(err.Error()))
		return
	}
	_ = os.Chdir(fullApiDir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("代码格式化错误:", color.RedString(err.Error()))
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(fullApiDir, "*.go"))
	cmdImport.CombinedOutput()
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outApiFile))

	builder := strings.Builder{}

	for _, v := range moduleParam.Method {
		switch v.Type {
		case "Create":
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "创建"))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("handle.POST(\"%s\",%s)", apiUrl, pkg+"."+v.Name))
		case "Update":
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "更新"))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("handle.PUT(\"%s\",%s)", apiUrl+"/:"+base.Helper(moduleParam.Primary.AlisaColumnName)+"/update", pkg+"."+v.Name))
		case "Delete":
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "删除"))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("handle.DELETE(\"%s\",%s)", apiUrl+"/:"+base.Helper(moduleParam.Primary.AlisaColumnName)+"/delete", pkg+"."+v.Name))
		case "Recover":
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "恢复"))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("handle.PUT(\"%s\",%s)", apiUrl+"/:"+base.Helper(moduleParam.Primary.AlisaColumnName)+"/recover", pkg+"."+v.Name))
		case "Drop":
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "清理"))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("handle.DELETE(\"%s\",%s)", apiUrl+"/:"+base.Helper(moduleParam.Primary.AlisaColumnName)+"/drop", pkg+"."+v.Name))
		case "Show":
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "单条数据信息查看"))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("handle.GET(\"%s\",%s)", apiUrl+"/:"+base.Helper(moduleParam.Primary.AlisaColumnName)+"/item", pkg+"."+v.Name))
		case "Item":
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "单条数据信息查看"))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("handle.GET(\"%s\",%s)", apiUrl+"/"+v.Name+"/Item", pkg+"."+v.Name))
		case "List":
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "列表"))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("handle.GET(\"%s\",%s)", apiUrl, pkg+"."+v.Name))
		}
	}
	fmt.Println(builder.String())
}
