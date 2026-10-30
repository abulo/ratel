// Package convert 提供数据库驱动值到目标类型的赋值转换，供 null 包各可空类型的 Scan 方法复用。
//
// 转换逻辑移植自标准库 database/sql 内部实现，覆盖字符串、字节切片、时间、数值与布尔等常见类型；
// 目标必须为指针，无法安全转换时返回错误。
package convert
