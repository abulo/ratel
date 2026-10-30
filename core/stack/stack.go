package stack

import (
	"slices"
	"sync"
)

// NewStack 创建一个空的延迟函数栈；返回的栈可并发使用。
func NewStack() *DeferStack {
	return &DeferStack{
		fns: make([]func() error, 0),
		mu:  sync.RWMutex{},
	}
}

// DeferStack 保存待延迟执行的函数，清理时按后进先出的逆序运行；
// 须由 NewStack 创建，入栈与清理操作均为并发安全。
type DeferStack struct {
	fns []func() error
	mu  sync.RWMutex
}

// Push 将传入的函数追加到栈顶。
func (ds *DeferStack) Push(fns ...func() error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.fns = append(ds.fns, fns...)
}

// Clean 按逆序依次执行栈中全部函数；单个函数返回的错误被忽略。
//
// 执行后函数切片不会被清空，重复调用 Clean 会再次执行已有函数。
func (ds *DeferStack) Clean() {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for _, v := range slices.Backward(ds.fns) {
		_ = v()
	}
}
