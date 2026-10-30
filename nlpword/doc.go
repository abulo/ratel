// Package nlpword 提供基于 Aho-Corasick 自动机的敏感词处理工具。
//
// 该包支持从本地文件、网络地址或程序内批量加载敏感词，完成敏感词检测、
// 整词过滤、字符替换、全部命中收集与噪声干扰字符去除；Filter 为对外使用入口。
package nlpword
