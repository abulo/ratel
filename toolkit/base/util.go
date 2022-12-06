package base

import (
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
