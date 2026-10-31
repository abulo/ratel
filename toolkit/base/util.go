package base

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
)

// CamelStr 将下划线分隔的名称转换为大驼峰（PascalCase）形式，供代码生成模板调用。
func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = util.UCWords(name)
	return strings.Replace(name, " ", "", -1)
}

// Helper 将大驼峰名称的首字母改为小写，返回小驼峰（camelCase）形式。
func Helper(name string) string {
	name = CamelStr(name)
	return strings.ToLower(string(name[0])) + name[1:]
}

// Char 返回按当前数据库驱动包裹标识符的引号：PostgreSQL 使用双引号，其余驱动使用反引号。
func Char(in string) string {
	if cast.ToString(Driver) == "pgx" {
		return "\\\"" + in + "\\\""
	}
	return "`" + in + "`"
}

// Pointer 接收类型名，名称中包含 null（可空类型）时返回空字符串，否则返回指针标记 "*"。
func Pointer(in string) string {
	if strings.Contains(in, "null") {
		return ""
	}
	return "*"
}

// SymbolChar 返回 SQL 标识符使用的反引号，供代码生成模板调用。
func SymbolChar() string {
	return "`"
}

// PathChar 将给定片段用反引号包裹后拼接返回，供代码生成模板调用。
func PathChar(str ...string) string {
	var pathList []string
	pathList = append(pathList, "`")
	pathList = append(pathList, str...)
	pathList = append(pathList, "`")
	return strings.Join(pathList, "")
}

// Add 将两个参数转换为整数后求和返回，供代码生成模板调用。
func Add(numberOne, numberTwo any) int {
	return cast.ToInt(numberOne) + cast.ToInt(numberTwo)
}

// InMethod 判断方法列表中是否存在类型匹配给定任一类型名的方法。
func InMethod(arr []Method, target ...string) bool {
	// PERF-8-12 [P2] 方法类型集合一次构建后 O(1) 查找，避免外层方法循环 × 内层线性/反射扫描；判定结果与原实现一致。
	typeSet := make(map[string]struct{}, len(arr))
	for _, item := range arr {
		typeSet[item.Type] = struct{}{}
	}
	for _, t := range target {
		if _, ok := typeSet[t]; ok {
			return true
		}
	}
	return false
}

// Convert 依据条件字段生成 DAO 层查询条件的 Go 代码片段，供代码生成模板调用。
func Convert(Condition []Column) string {
	builder := strings.Builder{}
	builder.Grow(len(Condition) * 128) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	for _, item := range Condition {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		lower := item.LowerName
		if lower == "" {
			lower = Helper(item.ColumnName)
		}
		builder.WriteString(fmt.Sprintf("	if val,ok := condition[\"%s\"] ;ok {", lower))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("		builder.Where(\"%s = ?\",val)", item.ColumnName))
		builder.WriteString("\n")
		builder.WriteString("	}")
		builder.WriteString("\n")
	}
	return builder.String()
}

// ModuleDaoConvertProto 生成 DAO 对象到 Proto 对象的赋值代码片段，可空字段按 null 包类型取值。
//
// res 为目标 Proto 对象变量名，resItem 为源 DAO 对象变量名，供代码生成模板调用。
func ModuleDaoConvertProto(Condition []Column, res, resItem string) string {
	builder := strings.Builder{}
	builder.Grow(len(Condition) * 256) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	for _, item := range Condition {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		camel := item.CamelName
		if camel == "" {
			camel = CamelStr(item.ColumnName)
		}
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s.Ptr()", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Time":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = *%s.%s.Ptr()", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = %s.%s", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("	%s.%s = %s.%s", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = *%s.%s.Ptr()", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Time":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("	if %s.%s.IsValid() {", resItem, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(*%s.%s.Ptr())", res, camel, resItem, camel))
				builder.WriteString("\n")
				builder.WriteString("		}")
				builder.WriteString("\n")

			}
		}
	}
	return builder.String()
}

// ProtoRequest 依据条件字段生成 Proto 请求消息的字段声明片段，供代码生成模板调用。
func ProtoRequest(condition []Column) string {
	builder := strings.Builder{}
	builder.Grow(len(condition) * 96) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	builder.WriteString("\n")
	for _, item := range condition {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		lower := item.LowerName
		if lower == "" {
			lower = Helper(item.ColumnName)
		}
		builder.WriteString(fmt.Sprintf("	// @inject_tag: db:\"%s\" json:\"%s\"", item.ColumnName, lower))
		builder.WriteString("\n")

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

// TypeScriptCondition 依据条件字段生成 TypeScript 查询条件对象的可选属性片段，供代码生成模板调用。
func TypeScriptCondition(Condition []Column) string {
	builder := strings.Builder{}
	builder.Grow(len(Condition) * 64) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	for _, item := range Condition {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		lower := item.LowerName
		if lower == "" {
			lower = Helper(item.ColumnName)
		}
		switch item.DataTypeMap.Default {
		case "null.Int32":
		case "int32":
			builder.WriteString(fmt.Sprintf("		%s?: number; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.Int64":
		case "int64":
			builder.WriteString(fmt.Sprintf("		%s?: number; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.Float32":
		case "float32":
			builder.WriteString(fmt.Sprintf("		%s?: number; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.Float64":
		case "float64":
			builder.WriteString(fmt.Sprintf("		%s?: number; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.String":
		case "string":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.Bytes":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.JSON":
			builder.WriteString(fmt.Sprintf("		%s?: any; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.Bool":
		case "bool":
			builder.WriteString(fmt.Sprintf("		%s?: boolean; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.CTime":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.Time":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.Date":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.DateTime":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		case "null.TimeStamp":
			builder.WriteString(fmt.Sprintf("		%s?: string; // %s", lower, item.ColumnComment))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

// ApiToProto 生成从请求参数到 Proto 对象的赋值代码片段，page 为 true 时同时写入分页对象。
//
// res 为目标 Proto 对象变量名，request 为请求参数取值函数名，供代码生成模板调用。
func ApiToProto(Condition []Column, res, request string, page bool) string {
	builder := strings.Builder{}
	builder.Grow(len(Condition) * 256) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	builder.WriteString("\n")
	for _, item := range Condition {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		camel := item.CamelName
		if camel == "" {
			camel = CamelStr(item.ColumnName)
		}
		lower := item.LowerName
		if lower == "" {
			lower = Helper(item.ColumnName)
		}
		switch item.DataTypeMap.Default {
		case "null.Int32":
		case "int32":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.Int32(cast.ToInt32(val)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  proto.Int32(cast.ToInt32(val)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Int64":
		case "int64":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.Int64(cast.ToInt64(val)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  proto.Int64(cast.ToInt64(val)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Float32":
		case "float32":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.Float32(cast.ToFloat32(val)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Float64":
		case "float64":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.Float64(cast.ToFloat64(val)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  proto.Float64(cast.ToFloat64(val)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.String":
		case "string":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  proto.String(val) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  proto.String(val) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Bytes":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  util.StringToBytes(cast.ToString(val)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  util.StringToBytes(cast.ToString(val)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.JSON":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  util.StringToBytes(cast.ToString(val)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  util.StringToBytes(cast.ToString(val)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Bool":
		case "bool":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s =  cast.ToBool(val) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s =  cast.ToBool(val) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.CTime":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Time":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.Date":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.DateTime":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		case "null.TimeStamp":
			builder.WriteString(fmt.Sprintf("	if val, ok := %s(\"%s\"); ok {", request, lower))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("		%s.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
			builder.WriteString("\n")
			if page {
				builder.WriteString(fmt.Sprintf("		%sTotal.%s = timestamppb.New(cast.ToTimeInDefaultLocation(val, time.Local)) // %s", res, camel, item.ColumnComment))
				builder.WriteString("\n")
			}
			builder.WriteString("	}")
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

// ModuleProtoConvertMap 生成 Proto 请求到 condition 映射表的赋值代码片段，供代码生成模板调用。
func ModuleProtoConvertMap(Condition []Column, request string) string {
	builder := strings.Builder{}
	builder.Grow(len(Condition) * 128) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	for _, item := range Condition {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		camel := item.CamelName
		if camel == "" {
			camel = CamelStr(item.ColumnName)
		}
		lower := item.LowerName
		if lower == "" {
			lower = Helper(item.ColumnName)
		}
		builder.WriteString(fmt.Sprintf("	if %s.%s != nil {", request, camel))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("		condition[\"%s\"] = %s.Get%s()", lower, request, camel))
		builder.WriteString("	}")
		builder.WriteString("\n")
	}
	return builder.String()
}

// Rule 依据表中不可空字段生成前端必填校验规则片段，供代码生成模板调用。
func Rule(table []Column) string {
	builder := strings.Builder{}
	builder.Grow(len(table) * 96) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	for _, item := range table {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		lower := item.LowerName
		if lower == "" {
			lower = Helper(item.ColumnName)
		}
		if item.IsNullable == "NO" {
			builder.WriteString(fmt.Sprintf("	%s: [{ required: true, message: \"%s不能为空\", trigger: \"blur\" }],", lower, item.ColumnComment))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

// Props 依据表字段与查询条件生成前端表格列定义片段，属于查询条件的字段附带搜索输入框配置。
func Props(table, condition []Column) string {
	// PERF-8-3 [P1] 条件字段预建集合，成员判断由 O(n·m) 反射线性扫描降为 O(1) map 查找，渲染结果与原来一致。
	conditionSet := make(map[string]struct{}, len(condition))
	for _, item := range condition {
		conditionSet[item.ColumnName] = struct{}{}
	}
	// PERF-8-10 [P2] 按字段数预估 builder 容量，避免宽表输出期间多次扩容搬迁。
	var builder strings.Builder
	builder.Grow(len(table) * 96)
	for _, item := range table {
		// PERF-8-4 [P1] 优先复用元数据加载阶段预计算的小驼峰名称，兜底保持与原 Helper 行为一致。
		lower := item.LowerName
		if lower == "" {
			lower = Helper(item.ColumnName)
		}
		if _, ok := conditionSet[item.ColumnName]; ok {
			builder.WriteString(fmt.Sprintf("	{ prop: \"%s\", label: \"%s\", search: { el: \"input\" , span: 2 } },", lower, item.ColumnComment))
		} else {
			builder.WriteString(fmt.Sprintf("	{ prop: \"%s\", label: \"%s\" },", lower, item.ColumnComment))
		}
		builder.WriteString("\n")

	}
	return builder.String()
}

// PropsSearch 依据表字段与查询条件生成前端搜索表单字段片段，供代码生成模板调用。
func PropsSearch(table, condition []Column) string {
	// PERF-8-3 [P1] 条件字段预建集合，成员判断由 O(n·m) 反射线性扫描降为 O(1) map 查找，渲染结果与原来一致。
	conditionSet := make(map[string]struct{}, len(condition))
	for _, item := range condition {
		conditionSet[item.ColumnName] = struct{}{}
	}
	// PERF-8-10 [P2] 按字段数预估 builder 容量。
	var builder strings.Builder
	builder.Grow(len(condition) * 64)
	for _, item := range table {
		if _, ok := conditionSet[item.ColumnName]; ok {
			// PERF-8-4 [P1] 复用预计算小驼峰名称，兜底保持原 Helper 行为。
			lower := item.LowerName
			if lower == "" {
				lower = Helper(item.ColumnName)
			}
			builder.WriteString(fmt.Sprintf("	{ prop: \"%s\", label: \"%s\", valueType: \"input\" },", lower, item.ColumnComment))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// Json 依据字段可空性与类型生成 JSON 对象的默认值片段，供代码生成模板调用。
func Json(Condition []Column) string {
	builder := strings.Builder{}
	builder.Grow(len(Condition) * 96) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	for _, item := range Condition {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		lower := item.LowerName
		if lower == "" {
			lower = Helper(item.ColumnName)
		}
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Time":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		%s: undefined, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("		%s: 0, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("		%s: 0, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("		%s: 0, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("		%s: 0, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		%s: {}, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Time":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		%s: \"\", // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("		%s: false, // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}

// TypeScript 依据字段可空性与类型生成 TypeScript 接口属性声明片段，供代码生成模板调用。
func TypeScript(Condition []Column) string {
	builder := strings.Builder{}
	builder.Grow(len(Condition) * 96) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	for _, item := range Condition {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		lower := item.LowerName
		if lower == "" {
			lower = Helper(item.ColumnName)
		}
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("		%s: number | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("		%s: number | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("		%s: number | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("		%s: number | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		%s: any | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("		%s: boolean | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Time":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		%s: string | undefined; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("		%s: number; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("		%s: number; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("		%s: number; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("		%s: number; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		%s: any; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Time":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		%s: string; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("		%s: boolean; // %s", lower, item.ColumnComment))
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}

// ModuleProtoConvertDao 生成 Proto 对象到 DAO 对象的赋值代码片段，可空字段按 null 包类型还原。
//
// res 为目标 DAO 对象变量名，request 为源 Proto 对象变量名，供代码生成模板调用。
func ModuleProtoConvertDao(Condition []Column, res, request string) string {
	builder := strings.Builder{}
	builder.Grow(len(Condition) * 256) // PERF-8-10 [P2] 预估输出容量，减少宽表渲染期间 builder 扩容搬迁
	builder.WriteString("\n")
	for _, item := range Condition {
		// PERF-8-4 [P1] 复用元数据加载阶段预计算的驼峰名称；缺失时兜底现场计算，转换结果与原实现一致。
		camel := item.CamelName
		if camel == "" {
			camel = CamelStr(item.ColumnName)
		}
		if item.IsNullable == "YES" {
			switch item.DataTypeMap.Empty {
			case "null.Int32":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.Int32From(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Int64":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.Int64From(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Float32":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.Float32From(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Float64":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.Float64From(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.String":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.StringFrom(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.BytesFrom(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.JSONFrom(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bool":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.BoolFrom(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.CTimeFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Time":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.TimeFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.DateFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.DateTimeFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.TimeStampFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			}
		} else {
			switch item.DataTypeMap.Default {
			case "int32":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "int64":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "float32":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "float64":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "string":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.%s // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Bytes":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.BytesFrom(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.JSON":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.JSONFrom(%s.Get%s()) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.CTime":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.CTimeFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Time":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.TimeFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.Date":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.DateFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.DateTime":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.DateTimeFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "null.TimeStamp":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = null.TimeStampFrom(util.GrpcTime(%s.%s)) // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			case "bool":
				builder.WriteString(fmt.Sprintf("		if %s != nil && %s.%s != nil {", request, request, camel))
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("			%s.%s = %s.Get%s() // %s", res, camel, request, camel, item.ColumnComment))
				builder.WriteString("\n")
				builder.WriteString("	}")
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}

// commonTemplateFuncs 为各代码生成模板共享的 10 个基础模板函数；
// PERF-8-20 [P2] 公共映射在包级定义一次，消除 module/vue 共 7 个生成函数每次调用重复构建同一 FuncMap 字面量。
var commonTemplateFuncs = template.FuncMap{
	"Convert":               Convert,
	"SymbolChar":            SymbolChar,
	"Char":                  Char,
	"Helper":                Helper,
	"CamelStr":              CamelStr,
	"Add":                   Add,
	"ModuleProtoConvertDao": ModuleProtoConvertDao,
	"ModuleDaoConvertProto": ModuleDaoConvertProto,
	"ModuleProtoConvertMap": ModuleProtoConvertMap,
	"ApiToProto":            ApiToProto,
}

// CommonTemplateFuncs 返回公共模板函数映射的浅拷贝；调用方可以在返回的映射上追加本生成器
// 专属函数而不会污染包级公共映射。
func CommonTemplateFuncs() template.FuncMap {
	funcs := make(template.FuncMap, len(commonTemplateFuncs))
	for name, fn := range commonTemplateFuncs {
		funcs[name] = fn
	}
	return funcs
}
