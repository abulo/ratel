package stack

var (
	globalDefers = NewStack()
)

// Register 将传入的清理函数追加到全局栈，由 Clean 统一逆序执行。
func Register(fns ...func() error) {
	globalDefers.Push(fns...)
}

// Clean 按注册逆序执行全局栈中的全部清理函数；单个函数返回的错误被忽略。
func Clean() {
	globalDefers.Clean()
}
