// Package filter 提供基于 Aho-Corasick 自动机的敏感词扫描器。
//
// 该包支持从词表切片或字典文件构建词库，按自定义跳过分隔符匹配文本，
// 提供命中查询、命中结果收集、字节替换与 rune 替换等能力；Search 为对外
// 使用入口，TrieWriter 负责词库构建。
package filter
