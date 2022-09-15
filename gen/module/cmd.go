package module

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
)

func Run(db *query.Query, tableName, outputDir, outputPackage, dao, tplFile string) {

	columns, err := queryColumns(db, db.DBName, tableName)
	if err != nil {
		logger.Logger.Panic(err)
	}

	var res []string
	var funcList []Func
	for _, item := range columns {
		if item.IndexName != "PRIMARY" {
			res = append(res, util.Explode(",", item.ColumnName)...)
			tmpFunc := Func{}
			tmpFunc.FuncName = CamelStr(strings.Replace(item.ColumnName, ",", "_", -1))
			tmpFunc.CondiTion = util.Explode(",", item.ColumnName)
			tmpFunc.NonUnique = item.NonUnique
			funcList = append(funcList, tmpFunc)
		}
	}

	var parse Parse
	parse.TableName = tableName
	parse.Dao = dao
	parse.Package = outputPackage
	parse.CondiTion = res
	parse.Func = funcList
	_ = os.MkdirAll(outputDir, os.ModePerm)
	content, _ := util.FileGetContents(tplFile)
	tpl := template.Must(template.New("name").Funcs(template.FuncMap{"Helper": Helper}).Parse(content))
	//输出文件
	outFile := path.Join(outputDir, tableName+".go")
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	//渲染输出
	err = tpl.Execute(file, parse)
	if err != nil {
		panic(err)
	}
	_ = os.Chdir(outputDir)
	cmd := exec.Command("go", "fmt")
	out, e := cmd.CombinedOutput()
	if e != nil {
		panic(e)
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(outputDir, "*.go"))
	cmdImport.CombinedOutput()
	fmt.Printf("格式化结果:\n%s\n", string(out))
}
func Run1(tableName, outputDir, outputPackage, dao, tpl string) {

	content, _ := util.FileGetContents(tpl)

	content = util.StrReplace("{{.Package}}", outputPackage, content, -1)
	content = util.StrReplace("{{.Dao}}", dao, content, -1)
	content = util.StrReplace("{{.TableName}}", tableName, content, -1)
	builder := strings.Builder{}
	//转换表名
	builder.Reset()

	builder.WriteString(content)
	fileStr := builder.String()
	_ = os.WriteFile(path.Join(outputDir, tableName+".go"), []byte(fileStr), os.ModePerm)

	_ = os.Chdir(outputDir)
	cmd := exec.Command("go", "fmt")
	out, e := cmd.CombinedOutput()
	if e != nil {
		panic(e)
	}

	cmdImport := exec.Command("goimports", "-w", path.Join(outputDir, "*.go"))
	cmdImport.CombinedOutput()

	fmt.Printf("格式化结果:\n%s\n", string(out))
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

func queryColumns(db *query.Query, DbName, tableName string) ([]Column, error) {
	var columns []Column
	sql := "SELECT NON_UNIQUE,SEQ_IN_INDEX,INDEX_NAME,INDEX_TYPE,GROUP_CONCAT(COLUMN_NAME) AS COLUMN_NAME FROM `information_schema`.`STATISTICS` WHERE TABLE_SCHEMA = '" + DbName + "' and TABLE_NAME = '" + tableName + "'  GROUP BY TABLE_NAME, INDEX_NAME  ORDER BY NON_UNIQUE ASC,SEQ_IN_INDEX ASC"
	err := db.NewBuilder(context.Background()).QueryRows(sql).ToStruct(&columns)
	return columns, err
}
