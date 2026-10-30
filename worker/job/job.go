package job

// Runner 是应用启动后并发执行的任务执行器。
//
// Run 阻塞执行任务本体，通常为常驻循环；框架在启动阶段并发调用各已注册 Runner。
// 若需被框架按名称管理，实现者还应额外提供 GetJobName 方法。
type Runner interface {
	Run()
}
