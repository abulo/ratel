package queue

var (
	internalQueue Queuer
)

// Queuer 是一个任务队列，用于在高并发情况下缓解服务器压力并改进任务处理
type Queuer interface {
	Run()
	Push(job Jober)
	Terminate()
}

// Run 运行start running queue，指定缓冲区数和工作线程数
func Run(maxCapacity, maxThread int) {
	if internalQueue == nil {
		internalQueue = NewQueue(maxCapacity, maxThread)
	}
	internalQueue.Run()
}

// RunListQueue 启动运行列表队列，指定工作线程数
func RunListQueue(maxThread int) {
	if internalQueue == nil {
		internalQueue = NewListQueue(maxThread)
	}
	internalQueue.Run()
}

// Push 将可执行任务推入队列
func Push(job Jober) {
	if internalQueue == nil {
		return
	}
	internalQueue.Push(job)
}

// Terminate 终止队列以接收任务并释放资源
func Terminate() {
	if internalQueue == nil {
		return
	}
	internalQueue.Terminate()
}
