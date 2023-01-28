package base

import (
	"fmt"
	"strings"

	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
)

// CamelStr 下划线转驼峰
func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = util.UCWords(name)
	return strings.Replace(name, " ", "", -1)
}

// Helper 将驼峰的首字母小写
func Helper(name string) string {
	name = CamelStr(name)
	return strings.ToLower(string(name[0])) + name[1:]
}

// Char 对数据库参数进行编码
func Char(in string) string {
	return "`" + in + "`"
}

// SymbolChar 模板变量函数
func SymbolChar() string {
	return "`"
}

func Add(numberOne, numberTwo interface{}) int {
	return cast.ToInt(numberOne) + cast.ToInt(numberTwo)
}

// 函数转换
func Convert(Condition []Column) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		builder.WriteString(fmt.Sprintf("	if !util.Empty(condition[\"%s\"]){", Helper(item.ColumnName)))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("		builder.Where(\"%s\",condition[\"%s\"])", Char(item.ColumnName), Helper(item.ColumnName)))
		builder.WriteString("	}")
		builder.WriteString("\n")
	}
	return builder.String()
}

// ConvertDao 条件转换
func ConvertDao(Condition []Column, res, request string) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.Int32From(%s.Data.Get%s()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.Int64From(%s.Data.Get%s()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.Float32From(%s.Data.Get%s()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.Float64From(%s.Data.Get%s()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.StringFrom(%s.Data.Get%s()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.BytesFrom(%s.Data.Get%s()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.BoolFrom(%s.Data.Get%s()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	if err := %s.Data.%s.CheckValid(); err == nil {", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.CTimeFrom(%s.Data.%s.AsTime()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	if err := %s.Data.%s.CheckValid(); err == nil {", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.DateFrom(%s.Data.%s.AsTime()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	if err := %s.Data.%s.CheckValid(); err == nil {", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.DateTimeFrom(%s.Data.%s.AsTime()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	if err := %s.Data.%s.CheckValid(); err == nil {", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.TimeStampFrom(%s.Data.%s.AsTime()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Data.Get%s() // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Data.Get%s() // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Data.Get%s() // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Data.Get%s() // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Data.Get%s() // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.BytesFrom(%s.Data.Get%s()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	if err := %s.Data.%s.CheckValid(); err == nil {", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.CTimeFrom(%s.Data.%s.AsTime()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	if err := %s.Data.%s.CheckValid(); err == nil {", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.DateFrom(%s.Data.%s.AsTime()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	if err := %s.Data.%s.CheckValid(); err == nil {", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.DateTimeFrom(%s.Data.%s.AsTime()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	if err := %s.Data.%s.CheckValid(); err == nil {", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.TimeStampFrom(%s.Data.%s.AsTime()) // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.Data.%s){", request, Helper(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Data.Get%s() // %s", res, Helper(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("	}")
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}
