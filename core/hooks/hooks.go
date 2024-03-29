package hooks

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	globalHooks = make([][]func(), StageMax)
)

type Stage int

func (s Stage) String() string {
	switch s {
	case Stage_BeforeLoadConfig:
		return "BeforeLoadConfig"
	case Stage_AfterLoadConfig:
		return "AfterLoadConfig"
	case Stage_BeforeRun:
		return "BeforeRun"
	case Stage_BeforeStop:
		return "BeforeStop"
	case Stage_AfterStop:
		return "AfterStop"
	}

	return "Unknown"
}

const (
	Stage_BeforeLoadConfig Stage = iota
	Stage_AfterLoadConfig
	Stage_BeforeRun
	Stage_BeforeStop
	Stage_AfterStop
	StageMax
)

// Register 注册一个defer函数
func Register(stage Stage, fns ...func()) {
	globalHooks[stage] = append(globalHooks[stage], fns...)
}

// Do 执行
func Do(stage Stage) {
	fmt.Printf("[ratel] %+v\n", color.GreenString(fmt.Sprintf("hook stage (%s)...", stage)))

	if stage >= StageMax {
		return
	}

	for i := len(globalHooks[stage]) - 1; i >= 0; i-- {
		fn := globalHooks[stage][i]
		if fn != nil {
			fn()
		}
	}

	fmt.Printf("[ratel] %+v\n", color.GreenString(fmt.Sprintf("hook stage (%s)... done", stage)))
}
