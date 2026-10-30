package queue

var (
	internalQueue Queuer
)

// Queuer 抽象任务队列的生命周期与任务提交契约：Run 启动消费，Push 提交任务，Terminate 停止并释放资源。
type Queuer interface {
	Run()
	Push(job Jober)
	Terminate()
}

// Run 启动包级共享队列；首次调用按给定缓冲容量与 worker 数量创建队列，队列已存在时忽略新参数直接复用。
func Run(maxCapacity, maxThread int) {
	if internalQueue == nil {
		internalQueue = NewQueue(maxCapacity, maxThread)
	}
	internalQueue.Run()
}

// RunListQueue 启动包级共享的链表队列；队列已创建时不再按新参数重建。
func RunListQueue(maxThread int) {
	if internalQueue == nil {
		internalQueue = NewListQueue(maxThread)
	}
	internalQueue.Run()
}

// Push 将任务提交到包级共享队列；队列尚未初始化时直接返回，任务被丢弃。
func Push(job Jober) {
	if internalQueue == nil {
		return
	}
	internalQueue.Push(job)
}

// Terminate 停止包级共享队列并等待已入队任务处理完成；队列未初始化时直接返回。
func Terminate() {
	if internalQueue == nil {
		return
	}
	internalQueue.Terminate()
}
