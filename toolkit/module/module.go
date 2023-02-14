package module

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
	"github.com/abulo/ratel/toolkit/base"
	"github.com/abulo/ratel/util"
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
	if len(args) == 0 {
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
	} else {
		dir = args[0]
		tableName = args[1]
	}
	if tableName == "" || dir == "" {
		fmt.Println("初始化:", color.RedString("模型层名称 & 表名称 必须填写"))
		return
	}
	// 文件夹的路径
	fullModuleDir := path.Join(base.Path, "module", dir)
	_ = os.MkdirAll(fullModuleDir, os.ModePerm)

	// 文件夹的路径
	fullProtoDir := path.Join(base.Path, "proto", dir)
	_ = os.MkdirAll(fullProtoDir, os.ModePerm)

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
	tableColumn, err := base.TableColumn(ctx, base.Config.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表结构信息:", color.RedString(err.Error()))
		return
	}
	tableColumnMap := make(map[string]base.Column)
	for _, item := range tableColumn {
		tableColumnMap[item.ColumnName] = item
	}
	// 表信息
	tableItem, err := base.TableItem(ctx, base.Config.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表信息:", color.RedString(err.Error()))
		return
	}
	// 表索引
	tableIndex, err := base.TableIndex(ctx, base.Config.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表索引:", color.RedString(err.Error()))
		return
	}
	// 表主键
	tablePrimary, err := base.TablePrimary(ctx, base.Config.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表主键:", color.RedString(err.Error()))
		return
	}
	var methodList []base.Method

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
		}
		methodList = append(methodList, method)
	}
	//获取 go.mod
	mod, err := base.ModulePath(path.Join(base.Path, "go.mod"))
	if err != nil {
		fmt.Println("go.mod文件不存在:", color.RedString(err.Error()))
		mod = "test"
	}
	// 数字长度
	strLen := strings.LastIndex(dir, "/")
	// 数据模型
	moduleParam := base.ModuleParam{
		Pkg:         dir[strLen+1:],
		Primary:     tablePrimary,
		Table:       tableItem,
		TableColumn: tableColumn,
		Method:      methodList,
		ModName:     mod,
	}
	GenerateModule(moduleParam, fullModuleDir, tableName)
	GenerateProto(moduleParam, fullProtoDir, tableName)
}

func GenerateProto(moduleParam base.ModuleParam, fullProtoDir, tableName string) {
	//protoc --go-grpc_out=../../api/v1 --go_out=../../api/v1 *proto
	// 模板变量
	tpl := template.Must(template.New("proto").Funcs(template.FuncMap{
		"Convert":    base.Convert,
		"SymbolChar": base.SymbolChar,
		"Char":       base.Char,
		"Helper":     base.Helper,
		"CamelStr":   base.CamelStr,
		"Add":        base.Add,
	}).Parse(ProtoTemplate()))

	// 文件夹路径
	outProtoFile := path.Join(fullProtoDir, tableName+".proto")
	if util.FileExists(outProtoFile) {
		util.Delete(outProtoFile)
	}
	file, err := os.OpenFile(outProtoFile, os.O_CREATE|os.O_WRONLY, 0755)
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
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outProtoFile))
}

func GenerateModule(moduleParam base.ModuleParam, fullModuleDir, tableName string) {
	// 模板变量
	tpl := template.Must(template.New("module").Funcs(template.FuncMap{
		"Convert":    base.Convert,
		"SymbolChar": base.SymbolChar,
		"Char":       base.Char,
		"Helper":     base.Helper,
		"CamelStr":   base.CamelStr,
	}).Parse(ModuleTemplate()))
	// 文件夹路径
	outModuleFile := path.Join(fullModuleDir, tableName+".go")
	if util.FileExists(outModuleFile) {
		util.Delete(outModuleFile)
	}
	file, err := os.OpenFile(outModuleFile, os.O_CREATE|os.O_WRONLY, 0755)
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
	_ = os.Chdir(fullModuleDir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("代码格式化错误:", color.RedString(err.Error()))
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(fullModuleDir, "*.go"))
	cmdImport.CombinedOutput()
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outModuleFile))
}

// ModuleTemplate 模板
func ModuleTemplate() string {
	outString := `
package {{.Pkg}}

import (
	"context"
	"{{.ModName}}/dao"
	"{{.ModName}}/initial"

	"github.com/abulo/ratel/stores/query"
	"github.com/abulo/ratel/util"
	"github.com/spf13/cast"
)
// {{.Table.TableName}} {{.Table.TableComment}}


// {{CamelStr .Table.TableName}}ItemCreate 创建数据
func {{CamelStr .Table.TableName}}ItemCreate(ctx context.Context,data dao.{{CamelStr .Table.TableName}})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Insert(data)
}

// {{CamelStr .Table.TableName}}ItemUpdate 更新数据
func {{CamelStr .Table.TableName}}ItemUpdate(ctx context.Context,{{.Primary.ColumnName}} {{.Primary.DataTypeMap.Default}},data dao.{{CamelStr .Table.TableName}})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{.Primary.ColumnName}}).Update(data)
}

// {{CamelStr .Table.TableName}}Item 获取数据
func {{CamelStr .Table.TableName}}Item(ctx context.Context,{{.Primary.ColumnName}} {{.Primary.DataTypeMap.Default}})(dao.{{CamelStr .Table.TableName}},error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{CamelStr .Table.TableName}}
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{.Primary.ColumnName}}).Row().ToStruct(&res)
}

// {{CamelStr .Table.TableName}}ItemDelete 删除数据
func {{CamelStr .Table.TableName}}ItemDelete(ctx context.Context,{{.Primary.ColumnName}} {{.Primary.DataTypeMap.Default}})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{.Primary.ColumnName}}).Delete()
}
{{- range .Method}}
{{- if eq .Type "List"}}
{{- if .Default}}

// {{CamelStr .Table.TableName}}{{CamelStr .Name}} 列表数据
func {{CamelStr .Table.TableName}}{{CamelStr .Name}}(ctx context.Context,condition map[string]interface{})([]dao.{{CamelStr .Table.TableName}},error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res []dao.{{CamelStr .Table.TableName}}
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	if !util.Empty(condition["pageOffset"]) {
		builder.Offset(cast.ToInt64(condition["pageOffset"]))
	}
	if !util.Empty(condition["pageSize"]) {
		builder.Limit(cast.ToInt64(condition["pageSize"]))
	}
	err := builder.OrderBy("{{Char .Primary.ColumnName}}", query.DESC).Rows().ToStruct(&res)
	return res, err
}

// {{CamelStr .Table.TableName}}{{CamelStr .Name}}Total 列表数据总量
func {{CamelStr .Table.TableName}}{{CamelStr .Name}}Total(ctx context.Context,condition map[string]interface{})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	return builder.Count()
}
{{- else}}

// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}} 列表数据
func {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}(ctx context.Context,condition map[string]interface{})([]dao.{{CamelStr .Table.TableName}},error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res []dao.{{CamelStr .Table.TableName}}
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	if !util.Empty(condition["pageOffset"]) {
		builder.Offset(cast.ToInt64(condition["pageOffset"]))
	}
	if !util.Empty(condition["pageSize"]) {
		builder.Limit(cast.ToInt64(condition["pageSize"]))
	}
	err := builder.OrderBy("{{Char .Primary.ColumnName}}", query.DESC).Rows().ToStruct(&res)
	return res, err
}

// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Total 列表数据总量
func {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Total(ctx context.Context,condition map[string]interface{})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	return builder.Count()
}
{{- end}}
{{- else}}

// {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}} 单列数据
func {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}(ctx context.Context,condition map[string]interface{})(dao.{{CamelStr .Table.TableName}},error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{CamelStr .Table.TableName}}
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	err := builder.Row().ToStruct(&res)
	return res, err
}
{{- end}}
{{- end}}
`
	return outString
}

func ProtoTemplate() string {
	outString := `
syntax = "proto3";
// {{.Table.TableName}} {{.Table.TableComment}}
package {{.Pkg}};
option go_package = "./{{.Pkg}}";
import "google/protobuf/timestamp.proto";

// {{CamelStr .Table.TableName}}Object 数据对象
message {{CamelStr .Table.TableName}}Object {
	{{- range .TableColumn}}
	{{.DataTypeMap.Proto}} {{.ColumnName}} = {{.PosiTion}}; //{{.ColumnComment}}
	{{- end}}
}

// {{CamelStr .Table.TableName}}ListObject 列表数据对象
message {{CamelStr .Table.TableName}}ListObject {
	int64 total = 1;
	repeated {{CamelStr .Table.TableName}}Object list = 2;
}

// {{CamelStr .Table.TableName}}ItemCreateRequest 创建数据
message {{CamelStr .Table.TableName}}ItemCreateRequest {
	{{CamelStr .Table.TableName}}Object data = 1;
}

// {{CamelStr .Table.TableName}}ItemCreateResponse 创建数据响应
message {{CamelStr .Table.TableName}}ItemCreateResponse {
	int64 code = 1;
	string msg = 2;
}

// {{CamelStr .Table.TableName}}ItemUpdateRequest 更新数据
message {{CamelStr .Table.TableName}}ItemUpdateRequest {
	{{.Primary.DataTypeMap.Proto}} {{.Primary.ColumnName}} = 1; //{{.Primary.ColumnComment}}
	{{CamelStr .Table.TableName}}Object data = 2;
}

// {{CamelStr .Table.TableName}}ItemUpdateResponse 更新数据响应
message {{CamelStr .Table.TableName}}ItemUpdateResponse {
	int64 code = 1;
	string msg = 2;
}

// {{CamelStr .Table.TableName}}ItemDeleteRequest 删除数据
message {{CamelStr .Table.TableName}}ItemDeleteRequest {
	{{.Primary.DataTypeMap.Proto}} {{.Primary.ColumnName}} = 1; //{{.Primary.ColumnComment}}
}

// {{CamelStr .Table.TableName}}ItemDeleteResponse 删除数据响应
message {{CamelStr .Table.TableName}}ItemDeleteResponse {
	int64 code = 1;
	string msg = 2;
}


// {{CamelStr .Table.TableName}}ItemRequest 数据
message {{CamelStr .Table.TableName}}ItemRequest {
	{{.Primary.DataTypeMap.Proto}} {{.Primary.ColumnName}} = 1; //{{.Primary.ColumnComment}}
}

// {{CamelStr .Table.TableName}}ItemResponse 数据响应
message {{CamelStr .Table.TableName}}ItemResponse {
	int64 code = 1;
	string msg = 2;
	{{CamelStr .Table.TableName}}Object data = 3;
}

{{- range .Method}}
{{- if eq .Type "List"}}
{{- if .Default}}
// {{CamelStr .Table.TableName}}{{CamelStr .Name}}Request 列表数据
message {{CamelStr .Table.TableName}}{{CamelStr .Name}}Request {
	{{- range .Condition}}
	{{.DataTypeMap.Proto}} {{.ColumnName}} = {{.PosiTion}}; //{{.ColumnComment}}
	{{- end}}
	int64 page_number = {{Add .ConditionTotal 1}};
  	int64 result_per_page = {{Add .ConditionTotal 2}};
}

// {{CamelStr .Table.TableName}}{{CamelStr .Name}}Response 数据响应
message {{CamelStr .Table.TableName}}{{CamelStr .Name}}Response {
	int64 code = 1;
  	string msg = 2;
	{{CamelStr .Table.TableName}}ListObject data = 3;
}
{{- else}}

// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Request 列表数据
message {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Request {
	{{- range .Condition}}
	{{.DataTypeMap.Proto}} {{.ColumnName}} = {{.PosiTion}}; //{{.ColumnComment}}
	{{- end}}
	int64 page_number = {{Add .ConditionTotal 1}};
  	int64 result_per_page = {{Add .ConditionTotal 2}};
}

// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Response 数据响应
message {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Response {
	int64 code = 1;
  	string msg = 2;
	{{CamelStr .Table.TableName}}ListObject data = 3;
}
{{- end}}
{{- else}}

// {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Request 单列数据
message {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Request {
	{{- range .Condition}}
	{{.DataTypeMap.Proto}} {{.ColumnName}} = {{.PosiTion}}; //{{.ColumnComment}}
	{{- end}}
}

// {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Response 单列数据
message {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Response {
	int64 code = 1;
	string msg = 2;
	{{CamelStr .Table.TableName}}Object data = 3;
}
{{- end}}
{{- end}}

// {{CamelStr .Table.TableName}}Service 服务
service {{CamelStr .Table.TableName}}Service{
	rpc {{CamelStr .Table.TableName}}ItemCreate({{CamelStr .Table.TableName}}ItemCreateRequest) returns ({{CamelStr .Table.TableName}}ItemCreateResponse);
	rpc {{CamelStr .Table.TableName}}ItemUpdate({{CamelStr .Table.TableName}}ItemUpdateRequest) returns ({{CamelStr .Table.TableName}}ItemUpdateResponse);
	rpc {{CamelStr .Table.TableName}}ItemDelete({{CamelStr .Table.TableName}}ItemDeleteRequest) returns ({{CamelStr .Table.TableName}}ItemDeleteResponse);
	rpc {{CamelStr .Table.TableName}}Item({{CamelStr .Table.TableName}}ItemRequest) returns ({{CamelStr .Table.TableName}}ItemResponse);
	{{- range .Method}}
	{{- if eq .Type "List"}}
	{{- if .Default}}
	rpc {{CamelStr .Table.TableName}}{{CamelStr .Name}}({{CamelStr .Table.TableName}}{{CamelStr .Name}}Request) returns ({{CamelStr .Table.TableName}}{{CamelStr .Name}}Response);
	{{- else}}
	rpc {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}({{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Request) returns ({{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Response);
	{{- end}}
	{{- else}}
	rpc {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}({{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Request) returns ({{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Response);
	{{- end}}
	{{- end}}
}
`
	return outString
}
