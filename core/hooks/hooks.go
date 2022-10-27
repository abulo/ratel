package hooks

import (
	"fmt"
	"sync"
)

var (
	globalHooks = make(map[Stage][]func())
	mu          = sync.RWMutex{}
)

// Stage ...
type Stage int

// String ...
func (s Stage) String() string {
	switch s {
	case StageBeforeLoadConfig:
		return "BeforeLoadConfig"
	case StageAfterLoadConfig:
		return "AfterLoadStart"
	case StageBeforeStop:
		return "BeforeStop"
	case StageAfterStop:
		return "AfterStop"
	}

	return "Unknown"
}

// Stage_BeforeLoadConfig ...
const (
	StageBeforeLoadConfig Stage = iota + 1
	StageAfterLoadConfig
	StageBeforeStop
	StageAfterStop
)

// Register 注册一个defer函数
func Register(stage Stage, fns ...func()) {
	mu.Lock()
	defer mu.Unlock()

	globalHooks[stage] = append(globalHooks[stage], fns...)
}

// Do 执行
func Do(stage Stage) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("[ratel] %+v\n", fmt.Sprintf("hook stage (%s)...", stage))
	for i := len(globalHooks[stage]) - 1; i >= 0; i-- {
		fn := globalHooks[stage][i]
		if fn != nil {
			fn()
		}
	}
	fmt.Printf("[ratel] %+v\n", fmt.Sprintf("hook stage (%s)... done", stage))
}
