package hooks

import (
	"fmt"
	"slices"

	"github.com/fatih/color"
)

var (
	globalHooks = make([][]func(), StageMax)
)

// Stage 标识应用生命周期中的钩子触发阶段。
type Stage int

// String 返回阶段的可读名称；未识别编号返回 Unknown。
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

// 生命周期阶段常量，作为钩子注册与执行的分组编号。
const (
	// Stage_BeforeLoadConfig 为配置加载之前。
	Stage_BeforeLoadConfig Stage = iota
	// Stage_AfterLoadConfig 为配置加载之后。
	Stage_AfterLoadConfig
	// Stage_BeforeRun 为应用正式运行之前。
	Stage_BeforeRun
	// Stage_BeforeStop 为应用停止之前。
	Stage_BeforeStop
	// Stage_AfterStop 为应用停止之后。
	Stage_AfterStop
	// StageMax 为阶段编号上界，用作钩子表容量与越界判断。
	StageMax
)

// Register 将传入的钩子函数追加到指定阶段；同一阶段的函数由 Do 按注册逆序执行。
func Register(stage Stage, fns ...func()) {
	globalHooks[stage] = append(globalHooks[stage], fns...)
}

// Do 触发指定阶段已注册的全部钩子，按注册逆序依次执行；
// stage 超出 [0, StageMax) 范围时直接返回；钩子函数内部 panic 不会被本函数捕获。
func Do(stage Stage) {
	fmt.Printf("[ratel] %+v\n", color.GreenString(fmt.Sprintf("hook stage (%s)...", stage)))

	if stage >= StageMax {
		return
	}

	for _, fn := range slices.Backward(globalHooks[stage]) {

		if fn != nil {
			fn()
		}
	}

	fmt.Printf("[ratel] %+v\n", color.GreenString(fmt.Sprintf("hook stage (%s)... done", stage)))
}
