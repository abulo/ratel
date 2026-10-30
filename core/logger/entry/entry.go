package entry

import (
	"time"
)

// Entry 是写入各日志后端的标准化日志条目，由后端钩子在触发时构造。
//
// Entry 为纯数据载体，零值即可使用；各后端仅读取其字段，不持有跨条目共享状态。
type Entry struct {
	// Host 为产生日志的主机名；主机名获取失败时由后端回填 "unknown"。
	Host      string         `json:"host"`
	// Timestamp 为日志产生时间。
	Timestamp time.Time      `json:"timestamp"`
	// File 为调用位置，格式为 文件路径:行号；调用信息缺失时为空字符串。
	File      string         `json:"file"`
	// Func 为产生日志的函数全限定名；调用信息缺失时为空字符串。
	Func      string         `json:"func"`
	// Message 为日志正文。
	Message   string         `json:"message"`
	// Data 为日志附加字段；钩子会先并入自定义扩展字段，再提交后端。
	Data      map[string]any `json:"data"`
	// Level 为大写形式的日志级别名称（如 INFO、ERROR）。
	Level     string         `json:"level"`
}
