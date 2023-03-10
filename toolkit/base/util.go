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

func Add(numberOne, numberTwo any) int {
	return cast.ToInt(numberOne) + cast.ToInt(numberTwo)
}

// 函数转换
func Convert(Condition []Column) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		builder.WriteString(fmt.Sprintf("	if !util.Empty(condition[\"%s\"]){", Helper(item.ColumnName)))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("		builder.Where(\"%s\",condition[\"%s\"])", Char(item.ColumnName), Helper(item.ColumnName)))
		builder.WriteString("\n")
		builder.WriteString("	}")
		builder.WriteString("\n")
	}
	return builder.String()
}

// ModuleDaoConvertProto 条件转换
func ModuleDaoConvertProto(Condition []Column, res string) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("	%s : *%s.%s.Ptr(),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("	%s : *%s.%s.Ptr(),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("	%s : *%s.%s.Ptr(),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("	%s : *%s.%s.Ptr(),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("	%s : *%s.%s.Ptr(),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	%s : *%s.%s.Ptr(),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("	%s : *%s.%s.Ptr(),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	%s :  timestamppb.New(*%s.%s.Ptr()),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	%s :  timestamppb.New(*%s.%s.Ptr()),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	%s :  timestamppb.New(*%s.%s.Ptr()),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	%s :  timestamppb.New(*%s.%s.Ptr()),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("	%s : %s.%s,", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("	%s : %s.%s,", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("	%s : %s.%s,", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("	%s : %s.%s,", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("	%s : %s.%s,", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	%s : *%s.%s.Ptr(),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	%s :  timestamppb.New(*%s.%s.Ptr()),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	%s :  timestamppb.New(*%s.%s.Ptr()),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	%s :  timestamppb.New(*%s.%s.Ptr()),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	%s :  timestamppb.New(*%s.%s.Ptr()),", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("	%s : %s.%s,", CamelStr(item.ColumnName), res, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}

func ApiToProto(Condition []Column, res, request string) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		switch item.DataTypeMap.Default {
		case "null.Int32":
		case "int32":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = cast.ToInt32(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Int64":
		case "int64":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = cast.ToInt64(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Float32":
		case "float32":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = cast.ToFloat32(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Float64":
		case "float64":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = cast.ToFloat64(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.String":
		case "string":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = cast.ToString(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Bytes":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = util.StringToBytes(cast.ToString(%s(\"%s\"))) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Bool":
		case "bool":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = cast.ToBool(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.CTime":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(cast.ToTime(%s(\"%s\"))) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Date":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(cast.ToTime(%s(\"%s\"))) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.DateTime":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(cast.ToTime(%s(\"%s\"))) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.TimeStamp":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(cast.ToTime(%s(\"%s\"))) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

// ModuleProtoConvertMap 条件转换
func ModuleProtoConvertMap(Condition []Column, request string) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("		condition[\"%s\"] = %s.Get%s()", Helper(item.ColumnName), request, CamelStr(item.ColumnName)))
		builder.WriteString("	}")
		builder.WriteString("\n")
	}
	return builder.String()
}

// ModuleProtoConvertDao 条件转换
func ModuleProtoConvertDao(Condition []Column, res, request string) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.Int32From(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.Int64From(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.Float32From(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.Float64From(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.StringFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.BytesFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.BoolFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	if err := %s.%s.CheckValid(); err == nil {", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.CTimeFrom(%s.%s.AsTime()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	if err := %s.%s.CheckValid(); err == nil {", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.DateFrom(%s.%s.AsTime()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	if err := %s.%s.CheckValid(); err == nil {", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.DateTimeFrom(%s.%s.AsTime()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	if err := %s.%s.CheckValid(); err == nil {", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.TimeStampFrom(%s.%s.AsTime()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Get%s() // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Get%s() // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Get%s() // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Get%s() // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Get%s() // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.BytesFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	if err := %s.%s.CheckValid(); err == nil {", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.CTimeFrom(%s.%s.AsTime()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	if err := %s.%s.CheckValid(); err == nil {", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.DateFrom(%s.%s.AsTime()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	if err := %s.%s.CheckValid(); err == nil {", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.DateTimeFrom(%s.%s.AsTime()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	if err := %s.%s.CheckValid(); err == nil {", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = null.TimeStampFrom(%s.%s.AsTime()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("	if !util.Empty(%s.%s){", request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.Get%s() // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}
