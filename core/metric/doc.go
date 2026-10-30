// Package metric 提供基于 Prometheus 的指标采集封装。
//
// 包内预置服务端、客户端、定时任务、缓存、第三方库等场景的计数器、
// 直方图与 Summary 向量，并提供选项式构造方法；所有向量在构造时
// 自动注册到 Prometheus 默认注册器，命名空间缺省为 ratel。
package metric
