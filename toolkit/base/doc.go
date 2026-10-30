// Package base 为 toolkit 代码生成工具提供共享的基础能力。
//
// 该包负责初始化运行环境（配置文件、数据库连接），读取 MySQL 与 PostgreSQL
// 的表、字段、索引、主键等元数据，并提供命名转换、类型映射以及 Git 仓库、
// Go 模块路径等辅助函数，供 module 等模板生成包复用。
package base
