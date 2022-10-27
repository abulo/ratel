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

	indexs, err := queryIndex(db, db.DBName, tableName)
	if err != nil {
		logger.Logger.Panic(err)
	}

	column, err := queryColumn(db, db.DBName, tableName)
	if err != nil {
		logger.Logger.Panic(err)
	}

	var res []string
	var funcList []Func
	for _, item := range indexs {
		if item.IndexName != "PRIMARY" {
			tmp := util.Explode(",", item.ColumnName)
			if len(tmp) > 0 {
				for _, v := range tmp {
					if !util.InArray(v, res) {
						res = append(res, v)
					}
				}
			}
			// res = append(res, util.Explode(",", item.ColumnName)...)
			tmpFunc := Func{}
			tmpFunc.FuncName = CamelStr(strings.Replace(item.ColumnName, ",", "_", -1))
			tmpFunc.CondiTion = ToCondiTion(util.Explode(",", item.ColumnName), column)
			tmpFunc.NonUnique = item.NonUnique
			funcList = append(funcList, tmpFunc)
		}
	}

	var parse Parse
	parse.TableName = tableName
	parse.Dao = dao
	parse.Package = outputPackage
	parse.CondiTion = ToCondiTion(res, column)
	parse.Func = funcList
	_ = os.MkdirAll(outputDir, os.ModePerm)
	content, _ := util.FileGetContents(tplFile)
	tpl := template.Must(template.New("name").Funcs(template.FuncMap{"Helper": Helper}).Parse(content))
	//输出文件
	outFile := path.Join(outputDir, tableName+".go")
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

func queryIndex(db *query.Query, DbName, tableName string) ([]Index, error) {
	var columns []Index
	sql := "SELECT NON_UNIQUE,SEQ_IN_INDEX,INDEX_NAME,INDEX_TYPE,GROUP_CONCAT(COLUMN_NAME) AS COLUMN_NAME FROM `information_schema`.`STATISTICS` WHERE TABLE_SCHEMA = '" + DbName + "' and TABLE_NAME = '" + tableName + "'  GROUP BY TABLE_NAME, INDEX_NAME  ORDER BY NON_UNIQUE ASC,SEQ_IN_INDEX ASC"
	err := db.NewBuilder(context.Background()).QueryRows(sql).ToStruct(&columns)
	return columns, err
}

func queryColumn(db *query.Query, DbName, tableName string) ([]Column, error) {
	var columns []Column
	sql := "SELECT COLUMN_NAME,IS_NULLABLE,DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '" + DbName + "' and TABLE_NAME = '" + tableName + "'  ORDER BY ORDINAL_POSITION ASC"
	err := db.NewBuilder(context.Background()).QueryRows(sql).ToStruct(&columns)
	return columns, err
}

func ToCondiTion(item []string, column []Column) []CondiTion {
	var condition []CondiTion
	for _, v := range column {
		if util.InArray(v.ColumnName, item) {
			dataType := strings.ToUpper(v.DataType)
			value, ok := DataTypeMap[dataType]
			if ok {
				dataType = value[0]
			} else {
				dataType = "string"
			}
			tmp := CondiTion{
				ItemType: dataType,
				ItemName: v.ColumnName,
			}
			condition = append(condition, tmp)
		}
	}
	return condition
}

// SELECT
// 	statistics.NON_UNIQUE,
// 	statistics.SEQ_IN_INDEX,
// 	statistics.INDEX_NAME,
// 	statistics.INDEX_TYPE,
// 	statistics.COLUMN_NAME,
// 	`columns`.DATA_TYPE,
// 	GROUP_CONCAT(CONCAT( statistics.COLUMN_NAME, ":", `columns`.DATA_TYPE )) AS field
// FROM
// 	`information_schema`.`STATISTICS` AS statistics
// 	LEFT JOIN information_schema.`COLUMNS` AS `columns` ON statistics.COLUMN_NAME = `columns`.COLUMN_NAME
// WHERE
// 	statistics.TABLE_SCHEMA = 'jxm_online'
// 	AND statistics.TABLE_NAME = 'video'
// 	AND `columns`.TABLE_SCHEMA = 'jxm_online'
// 	AND `columns`.TABLE_NAME = 'video'
// 	AND statistics.INDEX_NAME != 'PRIMARY'
// 	GROUP BY statistics.TABLE_NAME, statistics.INDEX_NAME
// ORDER BY
// 	statistics.NON_UNIQUE ASC,
// 	statistics.SEQ_IN_INDEX ASC
