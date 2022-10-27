package executor

import (
	"context"
	"sync"

	"github.com/abulo/ratel/v3/core/goroutine"
	"github.com/abulo/ratel/v3/core/logger"
)

// 全局执行器存储
var _instances = sync.Map{}

// Executor 执行器接口
type Executor interface {
	GetAddress() string   // 执行性标识
	RegXJob(jobs ...XJob) // 注册执行器任务
	Run() error           // 启动执行器
	Stop()                // 退出执行器，向调度中心取消注册
	GracefulStop()        // 优雅地退出执行器，向调度中心取消注册
}

// XJob 定时任务接口
type XJob interface {
	Run(ctx context.Context, param *RunReq) (msg string, err error)
	GetJobName() string
}

// RunReq 触发任务请求参数
type RunReq struct {
	JobID                 int64  `json:"jobId"`                 // 任务ID
	ExecutorHandler       string `json:"executorHandler"`       // 任务标识
	ExecutorParams        string `json:"executorParams"`        // 任务参数
	ExecutorBlockStrategy string `json:"executorBlockStrategy"` // 任务阻塞策略
	ExecutorTimeout       int64  `json:"executorTimeout"`       // 任务超时时间，单位秒，大于零时生效
	LogID                 int64  `json:"logId"`                 // 本次调度日志ID
	LogDateTime           int64  `json:"logDateTime"`           // 本次调度日志时间
	GlueType              string `json:"glueType"`              // 任务模式，可选值参考 com.xxl.job.core.glue.GlueTypeEnum
	GlueSource            string `json:"glueSource"`            // GLUE脚本代码
	GlueUpdateTime        int64  `json:"glueUpdatetime"`        // GLUE脚本更新时间，用于判定脚本是否变更以及是否需要刷新
	BroadcastIndex        int64  `json:"broadcastIndex"`        // 分片参数：当前分片
	BroadcastTotal        int64  `json:"broadcastTotal"`        // 分片参数：总分片
}

// Register 注册执行器到ratel
func Register(address string, e Executor) {
	if _, ok := _instances.Load(address); ok {
		panic("duplicate executor address " + address)
	}
	logger.Logger.Info("[ratel]", "add executor for", address)
	_instances.Store(address, e)
}

// Run ...
func Run() error {
	var executors = make([]func(), 0)
	_instances.Range(func(key, val interface{}) bool {
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

// Stop ...
func Stop() error {
	var executors = make([]func(), 0)
	_instances.Range(func(key, val interface{}) bool {
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

// GracefulStop ...
func GracefulStop() error {
	var executors = make([]func(), 0)
	_instances.Range(func(key, val interface{}) bool {
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
