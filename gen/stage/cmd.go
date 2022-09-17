package stage

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
)

func Run(db *query.Query, tableName, moduleName, packageName, goDir, goTpl, htmlDir, htmlTplDir, Title string) {
	columns, err := queryColumns(db, db.DBName, tableName)
	if err != nil {
		logger.Logger.Panic(err)
	}

	var parse Parse
	parse.Package = packageName
	parse.View = packageName
	parse.Table = tableName
	parse.Dao = CamelStr(tableName)
	parse.Module = moduleName
	camelStr := CamelStr(tableName + "_id")
	parse.Pri = strings.ToLower(string(camelStr[0])) + camelStr[1:]
	parse.Column = columns
	parse.ListTotal = cast.ToInt64(len(columns)) + 2
	parse.LayerTotal = cast.ToInt64(len(columns))
	parse.Title = Title
	_ = os.MkdirAll(goDir, os.ModePerm)
	content, _ := util.FileGetContents(goTpl)
	tpl := template.Must(template.New("name").Funcs(template.FuncMap{"Helper": Helper, "CamelStr": CamelStr, "Convert": Convert}).Parse(content))
	//输出文件
	outFile := path.Join(goDir, tableName+".go")
	if util.FileExists(outFile) {
		util.Delete(outFile)
	}
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	//渲染输出
	err = tpl.Execute(file, parse)
	if err != nil {
		panic(err)
	}
	_ = os.Chdir(goDir)
	cmd := exec.Command("go", "fmt")
	out, e := cmd.CombinedOutput()
	if e != nil {
		panic(e)
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(goDir, "*.go"))
	cmdImport.CombinedOutput()
	fmt.Printf("格式化结果:\n%s\n", string(out))

	//创建模板

	GenView(htmlTplDir, tableName, htmlDir, "add", parse)
	GenView(htmlTplDir, tableName, htmlDir, "list", parse)
	GenView(htmlTplDir, tableName, htmlDir, "edit", parse)
	GenView(htmlTplDir, tableName, htmlDir, "show", parse)
	GenView(htmlTplDir, tableName, htmlDir, "layer", parse)

	//路由
	RoutePermission(tableName, Title, packageName, parse)
}

func RoutePermission(tableName, title, packageName string, parse Parse) {
	out := `
{Parent: "system", Title: "` + title + `列表", Handle: "admin_` + tableName + `", Children: nil},
{Parent: "admin_` + tableName + `", Title: "` + title + `添加", Handle: "admin_` + tableName + `_add", Children: nil},
{Parent: "admin_` + tableName + `", Title: "` + title + `添加处理", Handle: "admin_` + tableName + `_create", Children: nil},
{Parent: "admin_` + tableName + `", Title: "` + title + `修改", Handle: "admin_` + tableName + `_edit", Children: nil},
{Parent: "admin_` + tableName + `", Title: "` + title + `修改处理", Handle: "admin_` + tableName + `_update", Children: nil},
{Parent: "admin_` + tableName + `", Title: "` + title + `删除", Handle: "admin_` + tableName + `_delete", Children: nil},
{Parent: "admin_` + tableName + `", Title: "` + title + `详情", Handle: "admin_` + tableName + `_show", Children: nil},
{Parent: "admin_` + tableName + `", Title: "` + title + `弹出层", Handle: "admin_` + tableName + `_layer", Children: nil},
`
	fmt.Println("权限添加")
	fmt.Println(out)

	fmt.Println("路由添加")
	route := `
// ` + title + `列表
route.GET("/` + tableName + `", "admin_` + tableName + `", ` + packageName + `.List` + parse.Dao + `)
// ` + title + `添加
route.GET("/` + tableName + `/add", "admin_` + tableName + `_add", ` + packageName + `.Add` + parse.Dao + `)
// ` + title + `添加处理
route.POST("/` + tableName + `", "admin_` + tableName + `_create", ` + packageName + `.Create` + parse.Dao + `)
// ` + title + `修改
route.GET("/` + tableName + `/:` + parse.Pri + `/edit", "admin_` + tableName + `_edit", ` + packageName + `.Edit` + parse.Dao + `)
// ` + title + `修改处理
route.POST("/` + tableName + `/:` + parse.Pri + `/update", "admin_` + tableName + `_update", ` + packageName + `.Update` + parse.Dao + `)
// ` + title + `删除
route.POST("/` + tableName + `/:` + parse.Pri + `/delete", "admin_` + tableName + `_delete", ` + packageName + `.Delete` + parse.Dao + `)
// ` + title + `详情
route.GET("/` + tableName + `/:` + parse.Pri + `/show", "admin_` + tableName + `_show", ` + packageName + `.Show` + parse.Dao + `)
// ` + title + `弹出层
route.GET("/` + tableName + `/layer", "admin_` + tableName + `_layer", ` + packageName + `.Layer` + parse.Dao + `)
`

	fmt.Println(route)

}

func GenView(viewTplDir, tableName, outputDir, action string, parse Parse) {

	_ = os.MkdirAll(outputDir, os.ModePerm)
	contentView, _ := util.FileGetContents(viewTplDir + "/" + action + ".tpl")
	viewTpl := template.Must(template.New("name").Funcs(template.FuncMap{"Helper": Helper, "CamelStr": CamelStr, "Convert": Convert}).Parse(contentView))
	outViewFile := path.Join(outputDir, action+"_"+tableName+".html")
	if util.FileExists(outViewFile) {
		util.Delete(outViewFile)
	}
	fileView, err := os.OpenFile(outViewFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	//渲染输出
	viewTpl.Execute(fileView, parse)
	if err != nil {
		panic(err)
	}
	content, _ := util.FileGetContents(outViewFile)
	content = util.StrReplace("@$", "{{", content, -1)
	content = util.StrReplace("$@", "}}", content, -1)
	util.FilePutContents(outViewFile, content, 0777)
}

func queryColumns(db *query.Query, DbName, tableName string) ([]Column, error) {
	var columns []Column
	sql := "SELECT COLUMN_NAME,IS_NULLABLE,DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '" + DbName + "' and TABLE_NAME = '" + tableName + "' AND IS_NULLABLE='NO' AND COLUMN_KEY != 'PRI' AND  COLUMN_NAME != 'create_at'  AND  COLUMN_NAME != 'update_at' ORDER BY ORDINAL_POSITION ASC"
	err := db.NewBuilder(context.Background()).QueryRows(sql).ToStruct(&columns)
	return columns, err
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

func Convert(column Column, name string) string {
	dataType := strings.ToUpper(column.DataType)
	value, ok := DataTypeMap[dataType]
	if ok {
		if column.IsNullable == "YES" {
			dataType = value[1]
		} else {
			dataType = value[0]
		}
	} else {
		dataType = "string"
	}
	var res string
	switch dataType {
	case "string":
		res = "cast.ToString(ctx.PostForm(\"" + Helper(name) + "\"))"
	case "int64":
		res = "cast.ToInt64(ctx.PostForm(\"" + Helper(name) + "\"))"
	}
	return res
}
