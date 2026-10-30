// Package config 提供基于 Viper 的配置读取封装。
//
// 该包支持从单个配置文件或目录批量加载 JSON、TOML、YAML 配置，并提供
// 字符串、整数、无符号整数、int64、时长、浮点数、布尔等类型化读取方法，
// 读取缺失的键时可由调用方传入默认值。
package config
