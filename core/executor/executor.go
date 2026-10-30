package executor

import (
	"context"
	"sync"

	"github.com/abulo/ratel/v3/core/goroutine"
	"github.com/abulo/ratel/v3/core/logger"
)

var _instances = sync.Map{}

// Executor 是分布式任务执行器契约，由应用统一注册并驱动其生命周期。
//
// 实现需提供自身网络地址、任务注册、启动、立即停止与优雅停止能力；
// 停止时实现应向调度中心取消注册。
type Executor interface {
	GetAddress() string
	RegXJob(jobs ...XJob)
	Run() error
	Stop()
	GracefulStop()
}

// XJob 是单个定时任务的执行契约，由执行器在调度触发时调用。
type XJob interface {
	Run(ctx context.Context, param *RunReq) (msg string, err error)
	GetJobName() string
}

// RunReq 为调度中心触发任务执行时透传的请求参数，字段语义与 xxl-job 调度参数一一对应。
type RunReq struct {
	// JobID 为调度任务在调度中心的编号。
	JobID int64 `json:"jobId"`
	// ExecutorHandler 为任务标识，执行器据此定位对应任务实现。
	ExecutorHandler string `json:"executorHandler"`
	// ExecutorParams 为调度时传入的任务参数。
	ExecutorParams string `json:"executorParams"`
	// ExecutorBlockStrategy 为任务阻塞策略，决定上次任务未结束时本次调度的处理方式。
	ExecutorBlockStrategy string `json:"executorBlockStrategy"`
	// ExecutorTimeout 为任务超时时间，单位秒，大于零时生效。
	ExecutorTimeout int64 `json:"executorTimeout"`
	// LogID 为本次调度日志编号。
	LogID int64 `json:"logId"`
	// LogDateTime 为本次调度日志时间。
	LogDateTime int64 `json:"logDateTime"`
	// GlueType 为任务模式，可选值参考 com.xxl.job.core.glue.GlueTypeEnum。
	GlueType string `json:"glueType"`
	// GlueSource 为 GLUE 模式下的脚本代码。
	GlueSource string `json:"glueSource"`
	// GlueUpdateTime 为 GLUE 脚本更新时间，用于判定脚本是否变更以及是否需要刷新。
	GlueUpdateTime int64 `json:"glueUpdatetime"`
	// BroadcastIndex 为分片参数：当前分片序号。
	BroadcastIndex int64 `json:"broadcastIndex"`
	// BroadcastTotal 为分片参数：总分片数量。
	BroadcastTotal int64 `json:"broadcastTotal"`
}

// Register 将执行器按网络地址登记到全局注册表；地址已被占用时直接 panic。
func Register(address string, e Executor) {
	if _, ok := _instances.Load(address); ok {
		panic("duplicate executor address " + address)
	}
	logger.Logger.Info("[ratel]", "add executor for", address)
	_instances.Store(address, e)
}

// Run 并行启动全部已注册执行器；单个执行器启动失败的错误被忽略，本函数始终返回 nil。
func Run() error {
	var executors = make([]func(), 0)
	_instances.Range(func(key, val any) bool {
		address := key.(string)
		if executor, ok := val.(Executor); ok {
			executors = append(executors, func() {
				logger.Logger.Info("xxl-job executor run begin for ", address)
				defer logger.Logger.Info("xxl-job executor run exit for ", address)
				_ = executor.Run()
			})
		}
		return true
	})

	goroutine.Parallel(executors...)()
	return nil
}

// Stop 并行调用全部已注册且实现 Stop 的执行器停止方法；本函数始终返回 nil。
func Stop() error {
	var executors = make([]func(), 0)
	_instances.Range(func(key, val any) bool {
		address := key.(string)
		if stopper, ok := val.(interface{ Stop() }); ok {
			executors = append(executors, func() {
				logger.Logger.Info("xxl-job executor stop for ", address)
				defer logger.Logger.Info("xxl-job executor exit for ", address)
				stopper.Stop()
			})
		}
		return true
	})

	goroutine.Parallel(executors...)()
	return nil
}

// GracefulStop 并行调用全部已注册且实现 GracefulStop 的执行器优雅停止方法；本函数始终返回 nil。
func GracefulStop() error {
	var executors = make([]func(), 0)
	_instances.Range(func(key, val any) bool {
		address := key.(string)
		if stopper, ok := val.(interface{ GracefulStop() }); ok {
			executors = append(executors, func() {
				logger.Logger.Info("xxl-job executor stop begin for ", address)
				defer logger.Logger.Info("xxl-job executor stop exit for", address)
				stopper.GracefulStop()
			})
		}
		return true
	})

	goroutine.Parallel(executors...)()
	return nil
}
