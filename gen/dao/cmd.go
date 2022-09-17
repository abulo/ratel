package dao

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
)

// Run ...
//
//	 db		sql链接
//	 DbName  数据库名称
//		outputDir: "输出目录",
//		outputPackage: "struct文件的包名
func Run(db *query.Query, DbName, outputDir, outputPackage string) {
	_ = os.MkdirAll(outputDir, os.ModePerm)
	tables, err := queryTables(db, DbName)
	if err != nil {
		panic(err)
	}

	builder := strings.Builder{}
	for _, table := range tables {
		columns, err := queryColumns(db, DbName, table.TableName)
		if err != nil {
			continue
		}
		//转换表名
		builder.Reset()
		packageTime := false
		packageSQL := false
		builder.WriteString(fmt.Sprintf("//%s\ntype %s struct {\n", table.TableComment, CamelStr(table.TableName)))
		for _, column := range columns {
			//转换列名
			dataType := strings.ToUpper(column.DataType)
			value, ok := DataTypeMap[dataType]
			if ok {
				if column.IsNullable == "YES" {
					dataType = value[1]
					packageSQL = true
				} else {
					dataType = value[0]
				}
				//是否需要 sql 包
				if dataType == "time.Time" {
					packageTime = true
				}
			} else {
				dataType = "string"
			}
			//拼接字符串
			camelStr := CamelStr(column.ColumnName)
			builder.WriteString(fmt.Sprintf("	%s %s `db:\"%s\" json:\"%s\"  form:\"%s\"` //%s", camelStr, dataType, column.ColumnName, strings.ToLower(string(camelStr[0]))+camelStr[1:], strings.ToLower(string(camelStr[0]))+camelStr[1:], column.ColumnComment))
			if column.ColumnKey != "" {
				builder.WriteString("(" + column.ColumnKey + ")")
			}
			builder.WriteString("\n")
		}

		builder.WriteString("}\n")
		fileStr := "package " + outputPackage + "\nimport ("
		if packageSQL {
			fileStr += "\"github.com/abulo/ratel/v3/stores/query\"\n\n"
		}
		if packageTime {
			fileStr += "\"time\"\n\n"
		}
		fileStr += ")\n"
		fileStr += builder.String()
		outFile := path.Join(outputDir, table.TableName+".go")
		if util.FileExists(outFile) {
			util.Delete(outFile)
		}
		_ = os.WriteFile(outFile, []byte(fileStr), os.ModePerm)
	}

	_ = os.Chdir(outputDir)
	cmd := exec.Command("go", "fmt")
	out, e := cmd.CombinedOutput()
	if e != nil {
		panic(e)
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(outputDir, "*.go"))
	cmdImport.CombinedOutput()
	fmt.Printf("gofmt结果:\n%s\n", string(out))
}

func queryColumns(db *query.Query, DbName, tableName string) ([]Column, error) {
	var columns []Column
	sql := "SELECT COLUMN_NAME,IS_NULLABLE,DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '" + DbName + "' and TABLE_NAME = '" + tableName + "'  ORDER BY ORDINAL_POSITION ASC"
	err := db.NewBuilder(context.Background()).QueryRows(sql).ToStruct(&columns)
	return columns, err
}

func queryTables(db *query.Query, DbName string) ([]Table, error) {
	var tables []Table
	sql := "SELECT TABLE_NAME ,TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = '" + DbName + "'"
	err := db.NewBuilder(context.Background()).QueryRows(sql).ToStruct(&tables)
	return tables, err
}

// CamelStr 下划线转驼峰
func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = util.UCWords(name)
	return strings.Replace(name, " ", "", -1)
}
