// Package worker 定义应用框架托管的后台工作单元接口。
//
// 该包约定 Worker 具备启动与停止两个生命周期方法，应用在启动阶段调用 WorkerStart、
// 停止阶段调用 WorkerStop，由框架或外部调度器统一管理其运行周期。
package worker
