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
		Use:   "module",
		Short: "数据模型层",
		Long:  "创建数据库模型层: ratel module dir table_name",
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
		fmt.Println("数据库配置文件不存在")
		return
	}

	//加载配置文件
	AppConfig = config.New("dao")
	AppConfig.AddDriver(toml.Driver)
	AppConfig.LoadFiles(configFile)

	moduleDir := ""
	tableName := ""
	if len(args) == 0 {
		if err = survey.AskOne(&survey.Input{
			Message: "表名称",
			Help:    "数据库中某个表名称",
		}, &tableName); err != nil || tableName == "" {
			return
		}
		if err = survey.AskOne(&survey.Input{
			Message: "模型层名称",
			Help:    "模块的文件夹路径",
		}, &moduleDir); err != nil || moduleDir == "" {
			return
		}
	} else {
		tableName = args[0]
		moduleDir = args[1]
	}

	if tableName == "" || moduleDir == "" {
		fmt.Println("模型层名称 & 表名称 必须填写")
		return
	}

	newModuleDir := wd + "/" + moduleDir
	_ = os.MkdirAll(newModuleDir, os.ModePerm)

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
			for _, field := range fieldList {
				tmp := Argument{}
				arg := util.Explode(":", field)
				tmp.Field = arg[0]
				tmp.FieldInput = Helper(arg[0])
				tmp.FieldType = arg[1]
				argument = append(argument, tmp)
			}
		}
		function.Argument = argument
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
			for _, field := range fields {
				tmp := Argument{}
				arg := util.Explode(":", field)
				tmp.Field = arg[0]
				tmp.FieldInput = Helper(arg[0])
				tmp.FieldType = arg[1]
				argument = append(argument, tmp)
			}
			function.Argument = argument
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
			for _, field := range fieldList {
				tmp := Argument{}
				arg := util.Explode(":", field)
				tmp.Field = arg[0]
				tmp.FieldInput = Helper(arg[0])
				tmp.FieldType = arg[1]
				argument = append(argument, tmp)
			}
		}
		function.Argument = argument
		functionList = append(functionList, function)
	}

	n := strings.LastIndex(moduleDir, "/")
	newModule := ModuleArg{}
	newModule.PackageName = moduleDir[n+1:]
	newModule.TableName = tableName
	newModule.Mark = CamelStr(tableName)
	newModule.FunctionList = functionList
	priInfo, err := QueryColumnPrimaryKey(ctx, AppConfig.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("该数据没有主键类型:", err)
		return
	}
	newModule.PrimaryKey = Helper(tableName + "_" + priInfo.ColumnName)
	newModule.ColumnName = priInfo.ColumnName
	//go文件生成地址
	tpl := template.Must(template.New("name").Funcs(template.FuncMap{"Characters": Characters}).Parse(strings.TrimSpace(moduleTemplate)))
	//输出文件
	outFile := path.Join(newModuleDir, tableName+".go")
	if util.FileExists(outFile) {
		util.Delete(outFile)
	}
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", err)
		return
	}
	//渲染输出
	err = tpl.Execute(file, newModule)
	if err != nil {
		fmt.Println("模板解析错误:", err)
		return
	}
	_ = os.Chdir(newModuleDir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("代码格式化错误:", err)
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(newModuleDir, "*.go"))
	cmdImport.CombinedOutput()
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outFile))
}

func Characters(in string) string {
	return "`" + in + "`"
}

// CamelStr 下划线转驼峰
func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = util.UCWords(name)
	return strings.Replace(name, " ", "", -1)
}

func Helper(name string) string {
	name = CamelStr(name)
	return strings.ToLower(string(name[0])) + name[1:]
}

var moduleTemplate = `
package {{.PackageName}}

import (
	"context"

	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// {{.Mark}}Create 创建数据
func {{.Mark}}Create(ctx context.Context, data dao.{{.Mark}}) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Characters .TableName}}").Insert(data)
}

// {{.Mark}}Update 更新数据
func {{.Mark}}Update(ctx context.Context, {{.PrimaryKey}} int64, data dao.{{.Mark}}) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Characters .TableName}}").Where("{{Characters .ColumnName}}", {{.PrimaryKey}}).Update(data)
}

// {{.Mark}}Delete 删除数据
func {{.Mark}}Delete(ctx context.Context, {{.PrimaryKey}} int64) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Characters .TableName}}").Where("{{Characters .ColumnName}}", {{.PrimaryKey}}).Delete()
}

// {{.Mark}}Item 获取数据
func {{.Mark}}Item(ctx context.Context, {{.PrimaryKey}} int64) (dao.{{.Mark}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{.Mark}}
	err := db.NewBuilder(ctx).Table("{{Characters .TableName}}").Where("{{Characters .ColumnName}}", {{.PrimaryKey}}).Row().ToStruct(&res)
	return res, err
}

{{- range .FunctionList}}
{{- if eq .Type "one"}}
// {{.Mark}}ItemBy{{.Name}} 获取数据
func {{.Mark}}ItemBy{{.Name}}(ctx context.Context, condition map[string]interface{}) (dao.{{.Mark}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{.Mark}}
	builder := db.NewBuilder(ctx).Table("{{Characters .TableName}}")
	{{- range .Argument}}
	if !util.Empty(condition["{{.FieldInput}}"]) {
		builder.Where("{{Characters .Field}}", condition["{{.FieldInput}}"])
	}
	{{- end}}
	err := builder.Row().ToStruct(&res)
	return res, err
}
{{- end}}
{{- if eq .Type "list"}}
{{- if .Default}}
// {{.Mark}}List 获取数据
func {{.Mark}}List(ctx context.Context, condition map[string]interface{}) ([]dao.{{.Mark}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res []dao.{{.Mark}}
	builder := db.NewBuilder(ctx).Table("{{Characters .TableName}}")
	{{- range .Argument}}
	if !util.Empty(condition["{{.FieldInput}}"]) {
		builder.Where("{{Characters .Field}}", condition["{{.FieldInput}}"])
	}
	{{- end}}
	if !util.Empty(condition["pageOffset"]) {
		builder.Offset(cast.ToInt64(condition["pageOffset"]))
	}
	if !util.Empty(condition["pageSize"]) {
		builder.Limit(cast.ToInt64(condition["pageSize"]))
	}
	err := builder.OrderBy("{{Characters $.ColumnName}}", query.DESC).Rows().ToStruct(&res)
	return res, err
}
// {{.Mark}}Total 获取数据数量
func {{.Mark}}Total(ctx context.Context, condition map[string]interface{}) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := db.NewBuilder(ctx).Table("{{Characters .TableName}}")
	{{- range .Argument}}
	if !util.Empty(condition["{{.FieldInput}}"]) {
		builder.Where("{{Characters .Field}}", condition["{{.FieldInput}}"])
	}
	{{- end}}
	return builder.Count()
}
{{- else}}



// {{.Mark}}ListBy{{.Name}} 获取数据
func {{.Mark}}ListBy{{.Name}}(ctx context.Context, condition map[string]interface{}) ([]dao.{{.Mark}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res []dao.{{.Mark}}
	builder := db.NewBuilder(ctx).Table("{{Characters .TableName}}")
	{{- range .Argument}}
	if !util.Empty(condition["{{.FieldInput}}"]) {
		builder.Where("{{Characters .Field}}", condition["{{.FieldInput}}"])
	}
	{{- end}}
	if !util.Empty(condition["pageOffset"]) {
		builder.Offset(cast.ToInt64(condition["pageOffset"]))
	}
	if !util.Empty(condition["pageSize"]) {
		builder.Limit(cast.ToInt64(condition["pageSize"]))
	}
	err := builder.OrderBy("{{Characters $.ColumnName}}", query.DESC).Rows().ToStruct(&res)
	return res, err
}
// {{.Mark}}TotalBy{{.Name}} 获取数据数量
func {{.Mark}}TotalBy{{.Name}}(ctx context.Context, condition map[string]interface{}) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := db.NewBuilder(ctx).Table("{{Characters .TableName}}")
	{{- range .Argument}}
	if !util.Empty(condition["{{.FieldInput}}"]) {
		builder.Where("{{Characters .Field}}", condition["{{.FieldInput}}"])
	}
	{{- end}}
	return builder.Count()
}

{{- end}}
{{- end}}
{{- end}}
`
