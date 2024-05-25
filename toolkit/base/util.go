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

func Pointer(in string) string {
	if strings.Contains(in, "null") {
		return ""
	}
	return "*"
}

// SymbolChar 模板变量函数
func SymbolChar() string {
	return "`"
}

func Add(numberOne, numberTwo any) int {
	return cast.ToInt(numberOne) + cast.ToInt(numberTwo)
}

// 查询数组中是否包含某个元素
func InMethod(arr []Method, target ...string) bool {
	ret := false
	for _, item := range arr {
		if util.InArray(item.Type, target) {
			ret = true
			break
		}
	}
	return ret
}

// 函数转换
func Convert(Condition []Column) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		builder.WriteString(fmt.Sprintf("	if val,ok := condition[\"%s\"] ;ok {", Helper(item.ColumnName)))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("		builder.Where(\"%s\",val)", Char(item.ColumnName)))
		builder.WriteString("\n")
		builder.WriteString("	}")
		builder.WriteString("\n")
	}
	return builder.String()
}

// ModuleDaoConvertProto 条件转换
func ModuleDaoConvertProto(Condition []Column, res, resItem string) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = *%s.%s.Ptr()", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = *%s.%s.Ptr()", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, CamelStr(item.ColumnName), resItem, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")

			}
		}
	}
	return builder.String()
}

func ProtoRequest(condition []Column) string {
	builder := strings.Builder{}
	builder.WriteString("\n")
	for _, item := range condition {
		builder.WriteString(fmt.Sprintf("	// @inject_tag: db:\"%s\" json:\"%s\"", item.ColumnName, Helper(item.ColumnName)))
		builder.WriteString("\n")
		// if item.DataTypeMap.Proto == "int32" || item.DataTypeMap.Proto == "int64" {
		// item.DataTypeMap.Proto = "string"
		// }

		if item.DataTypeMap.OptionProto {
			builder.WriteString(fmt.Sprintf("	%s %s %s = %d; // %s",
				"optional",
				item.DataTypeMap.Proto,
				item.ColumnName,
				item.PosiTion,
				item.ColumnComment,
			))
		} else {
			builder.WriteString(fmt.Sprintf("	%s %s = %d; // %s",
				item.DataTypeMap.Proto,
				item.ColumnName,
				item.PosiTion,
				item.ColumnComment,
			))
		}

		builder.WriteString("\n")
	}
	return builder.String()
}

func TypeScriptCondition(Condition []Column) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		switch item.DataTypeMap.Default {
		case "null.Int32":
		case "int32":
			builder.WriteString(fmt.Sprintf("		%s?: number; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.Int64":
		case "int64":
			builder.WriteString(fmt.Sprintf("		%s?: number; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.Float32":
		case "float32":
			builder.WriteString(fmt.Sprintf("		%s?: number; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.Float64":
		case "float64":
			builder.WriteString(fmt.Sprintf("		%s?: number; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.String":
		case "string":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.Bytes":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.JSON":
			builder.WriteString(fmt.Sprintf("		%s?: any; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.Bool":
		case "bool":
			builder.WriteString(fmt.Sprintf("		%s?: boolean; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.CTime":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.Date":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.DateTime":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		case "null.TimeStamp":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func ApiToProto(Condition []Column, res, request string, page bool) string {
	builder := strings.Builder{}
	builder.WriteString("\n")
	for _, item := range Condition {
		switch item.DataTypeMap.Default {
		case "null.Int32":
		case "int32":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.Int32(cast.ToInt32(val)) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  proto.Int32(cast.ToInt32(val)) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Int64":
		case "int64":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.Int64(cast.ToInt64(val)) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  proto.Int64(cast.ToInt64(val)) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Float32":
		case "float32":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.Float32(cast.ToFloat32(val)) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Float64":
		case "float64":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.Float64(cast.ToFloat64(val)) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  proto.Float64(cast.ToFloat64(val)) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.String":
		case "string":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.String(val) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  proto.String(val) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Bytes":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  util.StringToBytes(cast.ToString(val) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  util.StringToBytes(cast.ToString(val) // %s", res, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.JSON":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = cast.ToString(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s = cast.ToString(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Bool":
		case "bool":
			builder.WriteString(fmt.Sprintf("	if !util.Empty(%s(\"%s\")){", request, Helper(item.ColumnName)))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = cast.ToBool(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s = cast.ToBool(%s(\"%s\")) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
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
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s = timestamppb.New(cast.ToTime(%s(\"%s\"))) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
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
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s = timestamppb.New(cast.ToTime(%s(\"%s\"))) // %s", res, CamelStr(item.ColumnName), request, Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
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
		builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", request, CamelStr(item.ColumnName)))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("		condition[\"%s\"] = %s.Get%s()", Helper(item.ColumnName), request, CamelStr(item.ColumnName)))
		builder.WriteString("	}")
		builder.WriteString("\n")
	}
	return builder.String()
}

func Rule(table []Column) string {
	builder := strings.Builder{}
	for _, item := range table {
		if item.IsNullable == "NO" {
			builder.WriteString(fmt.Sprintf("	%s: [{ required: true, message: \"%s不能为空\", trigger: \"blur\" }],", Helper(item.ColumnName), item.ColumnComment))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func Props(table, condition []Column) string {
	var conditionList []string
	for _, item := range condition {
		conditionList = append(conditionList, item.ColumnName)
	}
	builder := strings.Builder{}
	for _, item := range table {
		if util.InArray(item.ColumnName, conditionList) {
			builder.WriteString(fmt.Sprintf("	{ prop: \"%s\", label: \"%s\", search: { el: \"input\", span: 2, props: { placeholder: \"请输入%s\" } } },", Helper(item.ColumnName), item.ColumnComment, item.ColumnComment))
		} else {
			builder.WriteString(fmt.Sprintf("	{ prop: \"%s\", label: \"%s\" },", Helper(item.ColumnName), item.ColumnComment))
		}
		builder.WriteString("\n")

	}
	return builder.String()
}

func Json(Condition []Column) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("		%s: 0, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("		%s: 0, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("		%s: 0, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("		%s: 0, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		%s: {}, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("		%s: false, // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}

// TypeScript 条件转换
func TypeScript(Condition []Column) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("		%s: number | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("		%s: number | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("		%s: number | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("		%s: number | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		%s: any | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("		%s: boolean | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("		%s: number; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("		%s: number; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("		%s: number; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("		%s: number; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		%s: any; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("		%s: boolean; // %s", Helper(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}

// ModuleProtoConvertDao 条件转换
func ModuleProtoConvertDao(Condition []Column, res, request string) string {
	builder := strings.Builder{}
	builder.WriteString("\n")
	for _, item := range Condition {
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.Int32From(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.Int64From(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.Float32From(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.Float64From(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.StringFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.BytesFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.JSONFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.BoolFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.CTimeFrom(util.GrpcTime(%s.%s)) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.DateFrom(util.GrpcTime(%s.%s)) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.DateTimeFrom(util.GrpcTime(%s.%s)) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.TimeStampFrom(util.GrpcTime(%s.%s)) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.BytesFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.JSONFrom(%s.Get%s()) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.CTimeFrom(util.GrpcTime(%s.%s)) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.DateFrom(util.GrpcTime(%s.%s)) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.DateTimeFrom(util.GrpcTime(%s.%s)) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.TimeStampFrom(util.GrpcTime(%s.%s)) // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, CamelStr(item.ColumnName)))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.Get%s() // %s", res, CamelStr(item.ColumnName), request, CamelStr(item.ColumnName), item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}
