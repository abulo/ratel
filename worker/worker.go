package worker

// Worker 是可被应用框架统一启停的后台工作单元。
//
// 应用启动时依次调用 WorkerStart 拉起工作逻辑，停止时依次调用 WorkerStop 释放资源；
// 方法返回的错误由应用生命周期流程收集处理。除框架托管外，实现者也可由自定义调度器驱动。
type Worker interface {
	WorkerStart() error
	WorkerStop() error
}
