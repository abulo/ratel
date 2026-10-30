// Package dao 实现 toolkit 的数据访问对象（DAO）代码生成命令。
//
// 该包读取配置数据库中的表清单（命令行参数可指定单表或逗号分隔的多表），逐表读取字段
// 元数据，按内置或配置文件指定的模板生成 GORM 模型结构体文件，并在生成后执行 go fmt
// 与 goimports 整理格式。
package dao
