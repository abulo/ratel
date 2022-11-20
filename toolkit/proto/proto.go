package proto

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abulo/ratel/v3/config"
	"github.com/abulo/ratel/v3/config/toml"
	"github.com/abulo/ratel/v3/stores/mysql"
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	// CmdNew represents the new command.
	CmdNew = &cobra.Command{
		Use:   "proto",
		Short: "gRPC",
		Long:  "创建proto文件. Example: ratel proto dir table_name",
		Run:   run,
	}
	AppConfig *config.Config
	Link      *query.Query
)

func run(cmd *cobra.Command, args []string) {
	timeout := "60s"
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	mysqlConfig := "mysql.toml"
	configFile := wd + "/" + mysqlConfig
	if !util.FileExists(configFile) {
		fmt.Println("The mysql configuration file does not exist.")
		return
	}

	//加载配置文件
	AppConfig = config.New("dao")
	AppConfig.AddDriver(toml.Driver)
	AppConfig.LoadFiles(configFile)

	//创建数据链接
	opt := &mysql.Config{}

	if Username := cast.ToString(AppConfig.String("mysql.Username")); Username != "" {
		opt.Username = Username
	}
	if Password := cast.ToString(AppConfig.String("mysql.Password")); Password != "" {
		opt.Password = Password
	}
	if Host := cast.ToString(AppConfig.String("mysql.Host")); Host != "" {
		opt.Host = Host
	}
	if Port := cast.ToString(AppConfig.String("mysql.Port")); Port != "" {
		opt.Port = Port
	}
	if Charset := cast.ToString(AppConfig.String("mysql.Charset")); Charset != "" {
		opt.Charset = Charset
	}
	if Database := cast.ToString(AppConfig.String("mysql.Database")); Database != "" {
		opt.Database = Database
	}

	// # MaxOpenConns 连接池最多同时打开的连接数
	// MaxOpenConns = 128
	// # MaxIdleConns 连接池里最大空闲连接数。必须要比maxOpenConns小
	// MaxIdleConns = 32
	// # MaxLifetime 连接池里面的连接最大存活时长(分钟)
	// MaxLifetime = 10
	// # MaxIdleTime 连接池里面的连接最大空闲时长(分钟)
	// MaxIdleTime = 5

	if MaxLifetime := cast.ToInt(AppConfig.Int("mysql.MaxLifetime")); MaxLifetime > 0 {
		opt.MaxLifetime = time.Duration(MaxLifetime) * time.Minute
	}
	if MaxIdleTime := cast.ToInt(AppConfig.Int("mysql.MaxIdleTime")); MaxIdleTime > 0 {
		opt.MaxIdleTime = time.Duration(MaxIdleTime) * time.Minute
	}
	if MaxIdleConns := cast.ToInt(AppConfig.Int("mysql.MaxIdleConns")); MaxIdleConns > 0 {
		opt.MaxIdleConns = cast.ToInt(MaxIdleConns)
	}
	if MaxOpenConns := cast.ToInt(AppConfig.Int("mysql.MaxOpenConns")); MaxOpenConns > 0 {
		opt.MaxOpenConns = cast.ToInt(MaxOpenConns)
	}
	opt.DriverName = "mysql"
	opt.DisableMetric = cast.ToBool(AppConfig.Bool("mysql.DisableMetric"))
	opt.DisableTrace = cast.ToBool(AppConfig.Bool("mysql.DisableTrace"))
	Link = mysql.NewClient(opt)

	tableName := ""
	moduleDir := ""
	if len(args) == 0 {
		if err = survey.AskOne(&survey.Input{
			Message: "表名称",
			Help:    "数据库中某个表名称",
		}, &tableName); err != nil || tableName == "" {
			return
		}
		if err = survey.AskOne(&survey.Input{
			Message: "模块的文件夹",
			Help:    "模块的文件夹路径",
		}, &moduleDir); err != nil || moduleDir == "" {
			return
		}
	} else {
		if len(args) < 2 {
			fmt.Println("模块的文件夹 & 表名称 必须填写")
			return
		}
		tableName = args[0]
		moduleDir = args[1]
	}
	if tableName == "" || moduleDir == "" {
		fmt.Println("模块的文件夹 & 表名称 必须填写")
		return
	}
	tableInfo, err := QueryTable(ctx, AppConfig.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表信息不准确")
		return
	}
	ColumnInfo, err := QueryColumn(ctx, AppConfig.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表字段信息不准确")
		return
	}
	res := Proto{
		Table:  tableInfo,
		Column: ColumnInfo,
	}

	ColumnInfoMap := make(map[string]string)
	for _, v := range ColumnInfo {
		ColumnInfoMap[v.ColumnName] = v.ColumnComment
	}

	//获取表信息
	functionList := make([]Function, 0)
	fieldList := make([]string, 0)
	indexList, err := QueryIndex(ctx, AppConfig.String("mysql.Database"), tableName)
	if err != nil {
		function := Function{}
		function.Type = "list"
		function.Name = CamelStr("list")
		function.TableName = tableName
		function.Mark = CamelStr(tableName)
		function.Default = true
		argument := make([]Argument, 0)
		if len(fieldList) > 0 {
			for fieldIndex, field := range fieldList {
				tmp := Argument{}
				arg := util.Explode(":", field)
				tmp.Field = arg[0]
				tmp.FieldInput = arg[0]
				tmp.FieldType = arg[1]
				tmp.ProtoType = ProtoType(arg[1])
				tmp.ColumnComment = ColumnInfoMap[arg[0]]
				newKey := fieldIndex + 1
				tmp.Seq = cast.ToInt64(newKey)
				argument = append(argument, tmp)
			}
		}
		function.Argument = argument
		function.ArgumentNumber = len(function.Argument)
		functionList = append(functionList, function)
	} else {

		for _, index := range indexList {
			// index.IndexName
			indexName := util.Explode(":", index.IndexName)
			if len(indexName) < 2 {
				continue
			}
			function := Function{}
			function.Type = indexName[0]
			function.Name = CamelStr(indexName[1])
			function.TableName = tableName
			function.Mark = CamelStr(tableName)
			function.Default = false
			//获取参数
			fields := util.Explode(",", index.Field)
			for _, field := range fields {
				if !util.InArray(field, fieldList) {
					fieldList = append(fieldList, field)
				}
			}

			argument := make([]Argument, 0)
			for fieldIndex, field := range fields {
				tmp := Argument{}
				arg := util.Explode(":", field)
				tmp.Field = arg[0]
				tmp.FieldInput = arg[0]
				tmp.FieldType = arg[1]
				tmp.ProtoType = ProtoType(arg[1])
				tmp.ColumnComment = ColumnInfoMap[arg[0]]
				newKey := fieldIndex + 1
				tmp.Seq = cast.ToInt64(newKey)
				argument = append(argument, tmp)
			}
			function.Argument = argument
			function.ArgumentNumber = len(function.Argument)
			functionList = append(functionList, function)
		}

		function := Function{}
		function.Type = "list"
		function.Name = CamelStr("list")
		function.TableName = tableName
		function.Mark = CamelStr(tableName)
		function.Default = true
		argument := make([]Argument, 0)
		if len(fieldList) > 0 {
			for fieldIndex, field := range fieldList {
				tmp := Argument{}
				arg := util.Explode(":", field)
				tmp.Field = arg[0]
				tmp.FieldInput = arg[0]
				tmp.FieldType = arg[1]
				tmp.ProtoType = ProtoType(arg[1])
				tmp.ColumnComment = ColumnInfoMap[arg[0]]
				newKey := fieldIndex + 1
				tmp.Seq = cast.ToInt64(newKey)
				argument = append(argument, tmp)
			}
		}
		function.Argument = argument
		function.ArgumentNumber = len(function.Argument)
		functionList = append(functionList, function)
	}

	res.FunctionList = functionList

	newModuleDir := wd + "/" + moduleDir
	_ = os.MkdirAll(newModuleDir, os.ModePerm)

	n := strings.LastIndex(moduleDir, "/")
	res.PackageName = moduleDir[n+1:]
	res.Mark = CamelStr(tableName)

	//go文件生成地址
	// tpl := template.Must(template.New("name").Parse(ProtoTpl))
	tpl := template.Must(template.New("name").Funcs(template.FuncMap{"Helper": Helper, "CamelStr": CamelStr, "Add": Add}).Parse(strings.TrimSpace(ProtoTpl)))
	//输出文件
	outFile := path.Join(newModuleDir, tableName+".proto")
	if util.FileExists(outFile) {
		util.Delete(outFile)
	}
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", err)
		return
	}
	//渲染输出
	err = tpl.Execute(file, res)
	if err != nil {
		fmt.Println("模板解析错误:", err)
		return
	}
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outFile))
}

var ProtoTpl = `syntax = "proto3";
//表名:{{.Table.TableName}},{{.Table.TableComment}}
package {{.PackageName}};
option go_package = "/{{.PackageName}}";

//{{.Table.TableComment}}数据对象
message {{.Mark}}Object {
	{{- range .Column}}
	{{.ProtoType}} {{.ColumnName}} = {{.Seq}}; //{{.ColumnComment}}
	{{- end}}
}

message {{.Mark}}ListObject {
	int64 total = 1;
	repeated {{.Mark}}Object list = 2;
}

// {{.Mark}}CreateRequest 创建数据请求
message {{.Mark}}CreateRequest {
	{{.Mark}}Object data = 1;
}
// {{.Mark}}CreateResponse 创建数据响应
message {{.Mark}}CreateResponse {
	int64 code = 1;
	string msg = 2;
	{{.Mark}}Object data = 3;
}

// {{.Mark}}UpdateRequest 更新数据请求
message {{.Mark}}UpdateRequest {
	{{.Mark}}Object data = 1;
}
// {{.Mark}}UpdateResponse 更新数据响应
message {{.Mark}}UpdateResponse {
	int64 code = 1;
	string msg = 2;
	{{.Mark}}Object data = 3;
}

// {{.Mark}}DeleteRequest 删除数据请求
message {{.Mark}}DeleteRequest {
	int64 id = 1;
}
// {{.Mark}}DeleteResponse 删除数据响应
message {{.Mark}}DeleteResponse {
	int64 code = 1;
	string msg = 2;
}

// {{.Mark}}ItemRequest 获取数据请求
message {{.Mark}}ItemRequest {
	int64 id = 1;
}
// {{.Mark}}ItemResponse 获取数据响应
message {{.Mark}}ItemResponse {
	int64 code = 1;
	string msg = 2;
	{{.Mark}}Object data = 3;
}

{{- range .FunctionList}}
{{if eq .Type "one"}}
// {{.Mark}}ItemBy{{.Name}}Request 数据请求
message {{.Mark}}ItemBy{{.Name}}Request {
	{{- range .Argument}}
	{{.ProtoType}} {{.Field}} = {{.Seq}};//{{.ColumnComment}}
	{{- end}}
}
// {{.Mark}}ItemBy{{.Name}}Response 数据响应
message {{.Mark}}ItemBy{{.Name}}Response {
	int64 code = 1;
	string msg = 2;
	{{.Mark}}Object data = 3;
}
{{- end}}
{{- if eq .Type "list"}}
{{- if .Default}}
// {{.Mark}}ListRequest 数据请求
message {{.Mark}}ListRequest {
	{{- range .Argument}}
	{{.ProtoType}} {{.Field}} = {{.Seq}};//{{.ColumnComment}}
	{{- end}}
	int64 page_number = {{Add .ArgumentNumber 1}};
  	int64 result_per_page = {{Add .ArgumentNumber 2}};
}
// {{.Mark}}ListResponse 数据响应
message {{.Mark}}ListResponse {
	int64 code = 1;
  	string msg = 2;
	{{.Mark}}ListObject data = 3;
}
{{- else}}
// {{.Mark}}ListBy{{.Name}}Request 获取数据
message {{.Mark}}ListBy{{.Name}}Request{
	{{- range .Argument}}
	{{.ProtoType}} {{.Field}} = {{.Seq}};//{{.ColumnComment}}
	{{- end}}
	int64 page_number = {{Add .ArgumentNumber 1}};
  	int64 result_per_page = {{Add .ArgumentNumber 2}};
}
// {{.Mark}}ListBy{{.Name}}Response 数据响应
message {{.Mark}}ListBy{{.Name}}Response {
	int64 code = 1;
  	string msg = 2;
	{{.Mark}}ListObject data = 3;
}
{{- end}}
{{- end}}
{{- end}}

service {{.Mark}}Service{
	rpc {{.Mark}}Create({{.Mark}}CreateRequest) returns ({{.Mark}}CreateResponse);
	rpc {{.Mark}}Update({{.Mark}}UpdateRequest) returns ({{.Mark}}UpdateResponse);
	rpc {{.Mark}}Delete({{.Mark}}DeleteRequest) returns ({{.Mark}}DeleteResponse);
	rpc {{.Mark}}Item({{.Mark}}ItemRequest) returns ({{.Mark}}ItemResponse);
{{- range .FunctionList}}
{{- if eq .Type "one"}}
	rpc {{.Mark}}ItemBy{{.Name}}({{.Mark}}ItemBy{{.Name}}Request) returns ({{.Mark}}ItemBy{{.Name}}Response);
{{- end}}
{{- if eq .Type "list"}}
{{- if .Default}}
	rpc {{.Mark}}List({{.Mark}}ListRequest) returns ({{.Mark}}ListResponse);
{{- else}}
	rpc {{.Mark}}ListBy{{.Name}}({{.Mark}}ListBy{{.Name}}Request) returns ({{.Mark}}ListBy{{.Name}}Response);
{{- end}}
{{- end}}
{{- end}}
}
`
