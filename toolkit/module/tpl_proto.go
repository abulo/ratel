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
	"github.com/spf13/cast"
)

// GenerateProto 使用 ProtoTemplate 渲染并写出 .proto 文件到 fullProtoDir，随后在 proto
// 父目录执行 protoc 生成 gRPC 代码，再运行 protoc-go-inject-tag 为生成的 pb.go 注入自定义字段标签。
//
// 模板渲染失败时打印错误并返回；protoc 与标签注入命令的执行错误不会中断后续流程。
func GenerateProto(moduleParam base.ModuleParam, fullProtoDir, fullServiceDir, tableName string) {
	// PERF-8-20 [P2] 公共模板函数复用 base 包级定义，仅追加本生成器专属的 ProtoRequest。
	funcs := base.CommonTemplateFuncs()
	funcs["ProtoRequest"] = base.ProtoRequest
	tpl := template.Must(template.New("proto").Funcs(funcs).Parse(ProtoTemplate()))
	outProtoFile := path.Join(fullProtoDir, tableName+".proto")
	if util.FileExists(outProtoFile) {
		util.Delete(outProtoFile)
	}
	file, err := os.OpenFile(outProtoFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}
	err = tpl.Execute(file, moduleParam)
	if err != nil {
		fmt.Println("模板解析错误:", color.RedString(err.Error()))
		return
	}
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outProtoFile))
	// protoc 在 proto 目录的父目录执行，生成代码输出到上级项目目录；命令错误仅回显，不中断生成流程。
	protoParentDir := util.GetParentDirectory(fullProtoDir)
	_ = os.Chdir(protoParentDir)
	cmdImportGrpc := exec.Command("protoc", "-I=proto", "--go-grpc_out=../", "--go_out=../", tableName+".proto")
	cmdImportGrpc.CombinedOutput()
	// 将 fullServiceDir 归一化为相对工作目录的路径，供 inject-tag 命令定位生成的 pb.go。
	fullServiceDir = strings.Replace(fullServiceDir, cast.ToString(base.Path)+"/", "", -1)
	cmdImportTag := exec.Command("protoc-go-inject-tag", "-input="+fullServiceDir+"/"+tableName+".pb.go")
	cmdImportTag.CombinedOutput()

	builder := strings.Builder{}
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("protoc  -I=proto  --go-grpc_out=../  --go_out=../  %s.proto ", tableName))
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("protoc-go-inject-tag -input=%s/%s.pb.go", fullServiceDir, tableName))
	builder.WriteString("\n")
	builder.WriteString("\n")
	fmt.Printf("\n🍺 Command   %s\n", color.GreenString(builder.String()))
}

// ProtoTemplate 返回 protobuf 服务定义模板；配置中 template.Proto 指定了自定义模板文件时优先读取该文件内容。
func ProtoTemplate() string {
	if exists := base.Config.Exists("template.Proto"); exists {
		filePath := path.Join(cast.ToString(base.Path), base.Config.String("template.Proto"))
		if util.FileExists(filePath) {
			if tplString, err := util.FileGetContents(filePath); err == nil {
				return tplString
			}
		}
	}
	outString := `
syntax = "proto3";
// {{.Table.TableName}} {{.Table.TableComment}}
package {{.Pkg}};
option go_package = "{{.ModName}}/service/{{.PkgPath}};{{.Pkg}}";
import "google/protobuf/timestamp.proto";
{{- if .Page}}
import "pagination.proto";
{{- end}}

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
	int64 data = 3;
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
{{- else if eq .Type "Drop"}}
// {{.Name}}Request 清理数据请求
message {{.Name}}Request {
	// @inject_tag: db:"{{.Primary.AlisaColumnName}}" json:"{{Helper .Primary.AlisaColumnName}}"
	{{.Primary.DataTypeMap.Proto}} {{ .Primary.AlisaColumnName}} = 1; //{{.Primary.ColumnComment}}
}
// {{.Name}}Response 清理数据响应
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
{{- else if eq .Type "Show"}}
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
	// @inject_tag: db:"pagination" json:"pagination"
	optional pagination.PaginationRequest pagination = {{Add .ConditionTotal 1}}; // 分页
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
}`
	return outString
}
