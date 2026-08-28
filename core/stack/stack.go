package stack

import (
	"slices"
	"sync"
)

// NewStack ...
func NewStack() *DeferStack {
	return &DeferStack{
		fns: make([]func() error, 0),
		mu:  sync.RWMutex{},
	}
}

// DeferStack ...
type DeferStack struct {
	fns []func() error
	mu  sync.RWMutex
}

// Push ...
func (ds *DeferStack) Push(fns ...func() error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.fns = append(ds.fns, fns...)
}

// Clean ...
func (ds *DeferStack) Clean() {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for _, v := range slices.Backward(ds.fns) {
		_ = v()
	}
}
