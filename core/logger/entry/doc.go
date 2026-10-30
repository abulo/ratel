// Package entry 定义日志钩子写入后端时使用的标准化日志条目结构。
//
// 条目携带主机名、时间戳、调用位置、日志级别与附加字段，由各后端钩子
// 在 Fire 阶段从 logrus.Entry 转换得到，再提交给对应存储。
package entry
