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

var (
	CmdNew = &cobra.Command{
		Use:   "api",
		Short: "接口对象",
		Long:  "创建接口对象: toolkit api",
		Run:   Run,
	}
)

func Run(cmd *cobra.Command, args []string) {
	// 数据初始化
	if err := base.InitBase(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	//接口文件夹
	dir := ""
	//表名
	tableName := ""
	//驱动类型
	apiType := ""
	if len(args) == 0 {
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
	} else {
		dir = args[0]
		tableName = args[1]
		apiType = args[2]
	}

	if tableName == "" || dir == "" || apiType == "" {
		fmt.Println("初始化:", color.RedString("接口路径 & 表名称 & 驱动类型 必须填写"))
		return
	}
	// 文件夹的路径
	fullApiDir := path.Join(base.Path, "api", dir)
	_ = os.MkdirAll(fullApiDir, os.ModePerm)

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
	var methodList []base.Method

	//获取 go.mod
	mod, err := base.ModulePath(path.Join(base.Path, "go.mod"))
	if err != nil {
		fmt.Println("go.mod文件不存在:", color.RedString(err.Error()))
		return
	}

	// 数字长度
	strLen := strings.LastIndex(dir, "/")

	//获取的索引信息没有
	if err != nil {
		method := base.Method{
			Table:          tableItem,
			TableColumn:    tableColumn,
			Type:           "List",
			Name:           "List",
			Default:        true,
			Condition:      nil,
			ConditionTotal: 0,
			Primary:        tablePrimary,
			Pkg:            dir[strLen+1:],
			PkgPath:        dir,
			ModName:        mod,
		}
		methodList = append(methodList, method)
	} else {
		//存储条件信息
		field := make([]string, 0)
		//有索引信息
		for _, v := range tableIndex {
			//查询条件
			condition := make([]base.Column, 0)
			//数据库索引
			indexField := v.Field
			indexFieldSlice := util.Explode(",", indexField)
			for _, fieldValue := range indexFieldSlice {
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
			method := base.Method{
				Table:          tableItem,
				TableColumn:    tableColumn,
				Type:           customIndexType,
				Name:           customIndexName,
				Default:        false,
				Condition:      condition,
				ConditionTotal: len(condition),
				Primary:        tablePrimary,
				Pkg:            dir[strLen+1:],
				PkgPath:        dir,
				ModName:        mod,
			}
			//添加到集合中
			methodList = append(methodList, method)
		}
		condition := make([]base.Column, 0)
		for _, fieldValue := range field {
			//构造查询条件
			positionIndex := cast.ToInt64(len(condition)) + 1
			currentColumn := tableColumnMap[fieldValue]
			currentColumn.PosiTion = positionIndex
			condition = append(condition, currentColumn)
			// condition = append(condition, tableColumnMap[fieldValue])
		}
		method := base.Method{
			Table:          tableItem,
			TableColumn:    tableColumn,
			Type:           "List",
			Name:           "List",
			Default:        true,
			Condition:      condition,
			ConditionTotal: len(condition),
			Primary:        tablePrimary,
			Pkg:            dir[strLen+1:],
			PkgPath:        dir,
			ModName:        mod,
		}
		methodList = append(methodList, method)
	}
	// 数据模型
	moduleParam := base.ModuleParam{
		Pkg:         dir[strLen+1:],
		PkgPath:     dir,
		Primary:     tablePrimary,
		Table:       tableItem,
		TableColumn: tableColumn,
		Method:      methodList,
		ModName:     mod,
	}
	Generate(moduleParam, fullApiDir, tableName, dir[strLen+1:], dir, apiType)
}

func Generate(moduleParam base.ModuleParam, fullApiDir, tableName, pkg, pkgPath, apiType string) {

	var tplString string

	//"gin", "hertz"

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
	// 文件夹路径
	outApiFile := path.Join(fullApiDir, tableName+".go")
	if util.FileExists(outApiFile) {
		util.Delete(outApiFile)
	}
	file, err := os.OpenFile(outApiFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}
	//渲染输出
	err = tpl.Execute(file, moduleParam)
	if err != nil {
		fmt.Println("模板解析错误:", color.RedString(err.Error()))
		return
	}
	//格式化代码
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
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "创建"))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("handle.POST(\"%s\",%s)", "/"+pkgPath+"/"+base.Helper(moduleParam.Table.TableName), pkg+"."+base.CamelStr(moduleParam.Table.TableName)+"ItemCreate"))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "更新"))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("handle.PUT(\"%s\",%s)", "/"+pkgPath+"/"+base.Helper(moduleParam.Table.TableName)+"/:"+base.Helper(moduleParam.Primary.AlisaColumnName), pkg+"."+base.CamelStr(moduleParam.Table.TableName)+"ItemUpdate"))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "单条数据信息查看"))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("handle.GET(\"%s\",%s)", "/"+pkgPath+"/"+base.Helper(moduleParam.Table.TableName)+"/:"+base.Helper(moduleParam.Primary.AlisaColumnName), pkg+"."+base.CamelStr(moduleParam.Table.TableName)+"Item"))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "单条数据信息删除"))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("handle.DELETE(\"%s\",%s)", "/"+pkgPath+"/"+base.Helper(moduleParam.Table.TableName)+"/:"+base.Helper(moduleParam.Primary.AlisaColumnName), pkg+"."+base.CamelStr(moduleParam.Table.TableName)+"ItemDelete"))
	builder.WriteString("\n")
	for _, item := range moduleParam.Method {
		if item.Type == "List" {
			if item.Default {
				builder.WriteString(fmt.Sprintf("// %s->%s->%s", moduleParam.Table.TableName, moduleParam.Table.TableComment, "列表"))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("handle.GET(\"%s\",%s)", "/"+pkgPath+"/"+base.Helper(moduleParam.Table.TableName), pkg+"."+base.CamelStr(item.Table.TableName)+base.CamelStr(item.Name)))
				builder.WriteString("\n")
			}
		}
	}
	fmt.Println(builder.String())
}
