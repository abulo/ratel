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
	fmt.Printf("\nGenerate %s Proto Command\n", color.GreenString(tableName))
	fmt.Printf("\ncd %s\n", color.GreenString(protoParentDir))
	fmt.Printf("\nprotoc --go-grpc_out=%s --go_out=%s %s/%s.proto\n", color.GreenString(serviceParentDir), color.GreenString(serviceParentDir), color.GreenString(currentDir), color.GreenString(tableName))
	fmt.Printf("\nprotoc-go-inject-tag -input=%s/%s.pb.go\n", color.GreenString(fullServiceDir), color.GreenString(tableName))
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
	{{.DataTypeMap.Proto}} {{.ColumnName}} = {{.PosiTion}}; //{{.ColumnComment}}
	{{- end}}
}

{{- if .Page}}
// {{CamelStr .Table.TableName}}ListObject 列表数据对象
message {{CamelStr .Table.TableName}}ListObject {
	int64 total = 1;
	repeated {{CamelStr .Table.TableName}}Object list = 2;
}
{{- end}}

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
	{{.Primary.DataTypeMap.Proto}} {{ .Primary.AlisaColumnName}} = 1; //{{.Primary.ColumnComment}}
	{{CamelStr .Table.TableName}}Object data = 2;
}

// {{CamelStr .Table.TableName}}ItemUpdateResponse 更新数据响应
message {{CamelStr .Table.TableName}}ItemUpdateResponse {
	int64 code = 1;
	string msg = 2;
}

// {{CamelStr .Table.TableName}}ItemDeleteRequest 删除数据
message {{CamelStr .Table.TableName}}ItemDeleteRequest {
	{{.Primary.DataTypeMap.Proto}} {{ .Primary.AlisaColumnName}} = 1; //{{.Primary.ColumnComment}}
}

// {{CamelStr .Table.TableName}}ItemDeleteResponse 删除数据响应
message {{CamelStr .Table.TableName}}ItemDeleteResponse {
	int64 code = 1;
	string msg = 2;
}


// {{CamelStr .Table.TableName}}ItemRequest 数据
message {{CamelStr .Table.TableName}}ItemRequest {
	{{.Primary.DataTypeMap.Proto}} {{ .Primary.AlisaColumnName}} = 1; //{{.Primary.ColumnComment}}
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
	{{ProtoRequest .Condition}}
	{{- if .Page}}
	int64 page_num = {{Add .ConditionTotal 1}};
  	int64 page_size = {{Add .ConditionTotal 2}};
	{{- end}}
}

// {{CamelStr .Table.TableName}}{{CamelStr .Name}}Response 数据响应
message {{CamelStr .Table.TableName}}{{CamelStr .Name}}Response {
	int64 code = 1;
  	string msg = 2;
	{{- if .Page}}
	{{CamelStr .Table.TableName}}ListObject data = 3;
	{{- else}}
	repeated {{CamelStr .Table.TableName}}Object data = 3;
	{{- end }}
}
{{- else}}

// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Request 列表数据
message {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Request {
	{{ProtoRequest .Condition}}
	{{- if .Page}}
	int64 page_num = {{Add .ConditionTotal 1}};
  	int64 page_size = {{Add .ConditionTotal 2}};
	{{- end}}
}

// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Response 数据响应
message {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Response {
	int64 code = 1;
  	string msg = 2;
	{{- if .Page}}
	{{CamelStr .Table.TableName}}ListObject data = 3;
	{{- else}}
	repeated {{CamelStr .Table.TableName}}Object data = 3;
	{{- end }}
}
{{- end}}
{{- else}}

// {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Request 单列数据
message {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Request {
	{{ProtoRequest .Condition}}
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
