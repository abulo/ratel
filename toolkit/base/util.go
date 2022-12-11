package base

import (
	"fmt"
	"strings"

	"github.com/abulo/ratel/v3/util"
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

// 函数转换
func Convert(Condition []Column) string {
	builder := strings.Builder{}
	for _, item := range Condition {
		builder.WriteString(fmt.Sprintf("	if !util.Empty(condition[\"%s\"]){", Helper(item.ColumnName)))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("		builder.Where(\"%s\",%s)", Char(item.ColumnName),
			fmt.Sprintf(item.DataTypeMap.Convert+"(condition[\"%s\"])", Helper(item.ColumnName))))
		builder.WriteString("	}")
		builder.WriteString("\n")
	}
	return builder.String()
}
