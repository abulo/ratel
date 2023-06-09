package module

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
)

func GenerateProto(moduleParam base.ModuleParam, fullProtoDir, fullServiceDir, tableName string) {
	//protoc --go-grpc_out=../../api/v1 --go_out=../../api/v1 *proto
	// 模板变量
	tpl := template.Must(template.New("proto").Funcs(template.FuncMap{
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
		"ProtoRequest":          base.ProtoRequest,
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
	//生成 grpc 代码
	serviceParentDir := util.GetParentDirectory(fullServiceDir)
	protoParentDir := util.GetParentDirectory(fullProtoDir)
	_ = os.Chdir(protoParentDir)
	strLen := strings.LastIndex(fullProtoDir, "/")
	currentDir := fullProtoDir[strLen+1:]
	cmdImportGrpc := exec.Command("protoc", "--go-grpc_out="+serviceParentDir, "--go_out="+serviceParentDir, currentDir+"/"+tableName+".proto")
	cmdImportGrpc.CombinedOutput()
	//修改自定义 tag
	cmdImportTag := exec.Command("protoc-go-inject-tag", "-input="+fullServiceDir+"/"+tableName+".pb.go")
	cmdImportTag.CombinedOutput()

	builder := strings.Builder{}
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("cd %s;", protoParentDir))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("protoc  --go-grpc_out=%s  --go_out=%s  %s/%s.proto;", serviceParentDir, serviceParentDir, currentDir, tableName))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("protoc-go-inject-tag -input=%s/%s.pb.go;", fullServiceDir, tableName))
	// return builder.String()
	shell := path.Join(fullProtoDir, tableName+".sh")
	if util.FileExists(shell) {
		util.Delete(shell)
	}
	_ = os.WriteFile(shell, []byte(builder.String()), os.ModePerm)
}

// @inject_tag: db:"{{.ColumnName}}" json:"{{Helper .ColumnName}}" form:"{{Helper .ColumnName}}" uri:"{{Helper .ColumnName}}" xml:"{{Helper .ColumnName}}" proto:"{{Helper .ColumnName}}"
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
	// @inject_tag: db:"{{.ColumnName}}" json:"{{Helper .ColumnName}}"
	{{- if .DataTypeMap.OptionProto}}
	optional {{.DataTypeMap.Proto}} {{.ColumnName}} = {{.PosiTion}}; //{{.ColumnComment}}
	{{- else }}
	{{.DataTypeMap.Proto}} {{.ColumnName}} = {{.PosiTion}}; //{{.ColumnComment}}
	{{- end}}
	{{- end}}
}

{{- if .Page}}
// {{CamelStr .Table.TableName}}TotalResponse 列表数据总量
message {{CamelStr .Table.TableName}}TotalResponse {
	int64 code = 1;
	string msg = 2;
	int64 data = 3;
}
{{- end}}


{{- range .Method}}
{{- if eq .Type "Create"}}
// {{.Name}}Request 创建数据请求
message {{.Name}}Request {
	{{CamelStr .Table.TableName}}Object data = 1;
}
// {{.Name}}Response 创建数据响应
message {{.Name}}Response {
	int64 code = 1;
	string msg = 2;
}
{{- else if eq .Type "Update"}}
// {{.Name}}Request 更新数据请求
message {{.Name}}Request {
	// @inject_tag: db:"{{.Primary.AlisaColumnName}}" json:"{{Helper .Primary.AlisaColumnName}}"
	{{.Primary.DataTypeMap.Proto}} {{ .Primary.AlisaColumnName}} = 1; //{{.Primary.ColumnComment}}
	{{CamelStr .Table.TableName}}Object data = 2;
}
// {{.Name}}Response 更新数据响应
message {{.Name}}Response {
	int64 code = 1;
	string msg = 2;
}
{{- else if eq .Type "Delete"}}
// {{.Name}}Request 删除数据请求
message {{.Name}}Request {
	// @inject_tag: db:"{{.Primary.AlisaColumnName}}" json:"{{Helper .Primary.AlisaColumnName}}"
	{{.Primary.DataTypeMap.Proto}} {{ .Primary.AlisaColumnName}} = 1; //{{.Primary.ColumnComment}}
}
// {{.Name}}Response 删除数据响应
message {{.Name}}Response {
	int64 code = 1;
	string msg = 2;
}
{{- else if eq .Type "Recover"}}
// {{.Name}}Request 恢复数据请求
message {{.Name}}Request {
	// @inject_tag: db:"{{.Primary.AlisaColumnName}}" json:"{{Helper .Primary.AlisaColumnName}}"
	{{.Primary.DataTypeMap.Proto}} {{ .Primary.AlisaColumnName}} = 1; //{{.Primary.ColumnComment}}
}
// {{.Name}}Response 删除数据响应
message {{.Name}}Response {
	int64 code = 1;
	string msg = 2;
}
{{- else if eq .Type "Only"}}
// {{.Name}}Request 查询单条数据请求
message {{.Name}}Request {
	// @inject_tag: db:"{{.Primary.AlisaColumnName}}" json:"{{Helper .Primary.AlisaColumnName}}"
	{{.Primary.DataTypeMap.Proto}} {{ .Primary.AlisaColumnName}} = 1; //{{.Primary.ColumnComment}}
}
// {{.Name}}Response 查询单条数据响应
message {{.Name}}Response {
	int64 code = 1;
	string msg = 2;
	{{CamelStr .Table.TableName}}Object data = 3;
}
{{- else if eq .Type "Item"}}
// {{.Name}}Request 查询单条数据请求
message {{.Name}}Request {
	{{ProtoRequest .Condition}}
}
// {{.Name}}Response 查询单条数据响应
message {{.Name}}Response {
	int64 code = 1;
	string msg = 2;
	{{CamelStr .Table.TableName}}Object data = 3;
}
{{- else if eq .Type "List"}}
// {{.Name}}Request 列表数据
message {{.Name}}Request {
	{{ProtoRequest .Condition}}
	{{- if .Page}}
	// @inject_tag: db:"page_num" json:"pageNum"
	optional int64 page_num = {{Add .ConditionTotal 1}};
	// @inject_tag: db:"page_size" json:"pageSize"
	optional int64 page_size = {{Add .ConditionTotal 2}};
	{{- end}}
}

// {{.Name}}Response 数据响应
message {{.Name}}Response {
	int64 code = 1;
  	string msg = 2;
	repeated {{CamelStr .Table.TableName}}Object data = 3;
}

{{- if .Page}}
// {{.Name}}TotalRequest 列表数据
message {{.Name}}TotalRequest {
	{{ProtoRequest .Condition}}
}
{{- end}}
{{- end}}
{{- end}}


// {{CamelStr .Table.TableName}}Service 服务
service {{CamelStr .Table.TableName}}Service{
	{{- range .Method}}
	rpc {{.Name}}({{.Name}}Request) returns ({{.Name}}Response);
	{{- if eq .Type "List"}}
	{{- if .Page}}
	rpc {{.Name}}Total({{.Name}}TotalRequest) returns ({{CamelStr .Table.TableName}}TotalResponse);
	{{- end}}
	{{- end}}
	{{- end}}
}
`
	return outString
}
