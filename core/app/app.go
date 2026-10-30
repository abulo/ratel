package app

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/abulo/ratel/v3/core/component"
	"github.com/abulo/ratel/v3/core/cycle"
	"github.com/abulo/ratel/v3/core/ecode"
	"github.com/abulo/ratel/v3/core/env"
	"github.com/abulo/ratel/v3/core/executor"
	"github.com/abulo/ratel/v3/core/goroutine"
	"github.com/abulo/ratel/v3/core/hooks"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/registry"
	"github.com/abulo/ratel/v3/server"
	"github.com/abulo/ratel/v3/worker"
	"github.com/abulo/ratel/v3/worker/job"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// Application 是应用程序实例，持有待启动的 Server、Worker、定时任务、执行器与生命周期钩子，
// 负责协调启动、信号监听、停止与清理的完整流程。
//
// 零值不可直接使用，须通过 New 或 DefaultApp 创建；初始化流程与停止流程都只执行一次。
type Application struct {
	cycle      *cycle.Cycle
	smu        *sync.RWMutex
	initOnce   sync.Once
	stopOnce   sync.Once
	servers    []server.Server
	workers    []worker.Worker
	jobs       map[string]job.Runner
	disableMap map[Disable]bool
	// HideBanner 为 true 时启动阶段不打印版本横幅信息。
	HideBanner bool
	stopped    chan struct{}
	components []component.Component
	sig        chan os.Signal
}

// Option 是用于配置 Application 的函数式选项。
type Option func(a *Application)

// Disable 标识应用启动流程中可被关闭的环节。
type Disable int

// 可关闭环节常量，配合 WithDisable 使用。
const (
	// DisableParserFlag 关闭命令行参数解析。
	DisableParserFlag Disable = 1
	// DisableLoadConfig 关闭配置加载。
	DisableLoadConfig Disable = 2
	// DisableDefaultGovernor 关闭默认治理能力。
	DisableDefaultGovernor Disable = 3
)

// WithOptions 将传入的选项依次应用到当前应用实例。
func (a *Application) WithOptions(options ...Option) {
	for _, option := range options {
		option(a)
	}
}

// WithDisable 返回一个选项，用于关闭指定的启动环节。
func WithDisable(d Disable) Option {
	return func(a *Application) {
		a.disableMap[d] = true
	}
}

// New 创建应用实例并按串行顺序依次执行传入的启动函数；任一启动函数返回错误时，
// 返回 nil 与该错误。
func New(fns ...func() error) (*Application, error) {
	app := &Application{}
	if err := app.Startup(fns...); err != nil {
		return nil, err
	}
	return app, nil
}

// DefaultApp 初始化并返回一个不带自定义启动函数的应用实例。
func DefaultApp() *Application {
	app := &Application{}
	app.initialize()
	return app
}

func (app *Application) runHooks(stage hooks.Stage) {
	hooks.Do(stage)
}

// RegisterHooks 为指定生命周期阶段注册钩子函数，钩子在对应阶段由 Do 触发执行。
func (app *Application) RegisterHooks(stage hooks.Stage, fns ...func()) {
	hooks.Register(stage, fns...)
}

func (app *Application) initialize() {
	app.initOnce.Do(func() {
		app.cycle = cycle.NewCycle()
		app.smu = &sync.RWMutex{}
		app.servers = make([]server.Server, 0)
		app.workers = make([]worker.Worker, 0)
		app.jobs = make(map[string]job.Runner)
		app.disableMap = make(map[Disable]bool)
		app.stopped = make(chan struct{})
		app.components = make([]component.Component, 0)
		_ = app.printBanner()
		// app.initLogger()
	})
}

// Startup 初始化应用实例，并串行执行传入的启动函数；任一函数返回错误时立即中断并返回该错误。
func (app *Application) Startup(fns ...func() error) error {
	app.initialize()
	return goroutine.SerialUntilError(fns...)()
}

// Serve 将传入的 Server 追加到应用托管列表，供 Run 统一启动与停止；当前实现始终返回 nil。
func (app *Application) Serve(s ...server.Server) error {
	app.smu.Lock()
	defer app.smu.Unlock()
	app.servers = append(app.servers, s...)
	return nil
}

// Schedule 将传入的 Worker 追加到应用托管列表，供 Run 统一启动与停止；当前实现始终返回 nil。
func (app *Application) Schedule(w worker.Worker) error {
	app.workers = append(app.workers, w)
	return nil
}

// Job 注册一个定时任务；任务须实现 GetJobName 以提供唯一名称，未实现该方法的任务会被静默忽略。
func (app *Application) Job(runner job.Runner) error {
	namedJob, ok := runner.(interface{ GetJobName() string })
	if !ok {
		return nil
	}
	jobName := namedJob.GetJobName()
	app.jobs[jobName] = runner
	return nil
}

// Executor 将执行器按其地址登记到全局执行器注册表，由应用统一管理启动与停止。
func (app *Application) Executor(e executor.Executor) {
	executor.Register(e.GetAddress(), e)
}

// SetRegistry 替换应用默认使用的服务注册中心实现。
func (app *Application) SetRegistry(reg registry.Registry) {
	registry.DefaultRegisterer = reg
}

func (app *Application) clean() {

}

// Run 是应用的阻塞式入口：先追加传入的 Server，触发 BeforeRun 钩子并启动信号监听，
// 随后并行启动定时任务、Server 与 Worker，阻塞等待整个生命周期结束。
//
// 生命周期以错误结束时记录错误日志并返回该错误；正常退出时返回 nil。
func (app *Application) Run(servers ...server.Server) error {
	app.smu.Lock()
	app.servers = append(app.servers, servers...)
	app.smu.Unlock()

	hooks.Do(hooks.Stage_BeforeRun)

	app.waitSignals()
	defer app.clean()

	// TODO(owner): 定时任务当前不支持优雅停止。
	_ = app.startJobs()

	app.cycle.Run(app.startServers)
	app.cycle.Run(app.startWorkers)
	if err := <-app.cycle.Wait(); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"app": ecode.ModApp,
			"err": err,
		}).Error("ratel shutdown with error")
		return err
	}
	logger.Logger.Info("shutdown ratel, bye!", ecode.ModApp)
	return nil
}

// Stop 在必要清理完成后立即停止应用：执行 BeforeStop 钩子，注销并停止全部 Server 与 Worker，
// 等待后台任务结束后执行 AfterStop 钩子并关闭生命周期；重复调用安全，停止流程只执行一次。
func (app *Application) Stop() (err error) {
	app.stopOnce.Do(func() {
		app.stopped <- struct{}{}
		app.runHooks(hooks.Stage_BeforeStop)
		for _, s := range app.servers {
			func(s server.Server) {
				app.smu.RLock()
				if err := registry.DefaultRegisterer.UnregisterService(context.Background(), s.Info()); err != nil {
					logger.Logger.WithFields(logrus.Fields{
						"ModApp": ecode.ModApp,
						"Name":   s.Info().Name,
						"Label":  s.Info().Label(),
						"err":    err,
					}).Info("exit server stop")
				}
				app.cycle.Run(s.Stop)
				app.smu.RUnlock()
			}(s)
		}
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.WorkerStop)
			}(w)
		}

		<-app.cycle.Done()
		app.runHooks(hooks.Stage_AfterStop)
		app.cycle.Close()
	})
	return
}

// GracefulStop 使用给定上下文优雅停止应用：注销服务时将上下文透传给注册中心，
// 依次停止 Server、Worker 与执行器，等待后台任务结束后执行 AfterStop 钩子并关闭生命周期；
// 重复调用安全，停止流程只执行一次。
func (app *Application) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.stopped <- struct{}{}
		app.runHooks(hooks.Stage_BeforeStop)
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(func() error {
					app.smu.RLock()
					defer app.smu.RUnlock()
					if err := registry.DefaultRegisterer.UnregisterService(ctx, s.Info()); err != nil {
						logger.Logger.WithFields(logrus.Fields{
							"ModApp": ecode.ModApp,
							"Name":   s.Info().Name,
							"Label":  s.Info().Label(),
							"err":    err,
						}).Info("exit server graceful stop")
					}
					return s.GracefulStop(ctx)
				})
			}(s)
		}
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.WorkerStop)
			}(w)
		}
		app.cycle.Run(executor.GracefulStop)
		<-app.cycle.Done()
		app.runHooks(hooks.Stage_AfterStop)
		app.cycle.Close()
	})
	return err
}

func (app *Application) startServers() error {
	var eg errgroup.Group
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	go func() {
		<-app.stopped
		cancel()
	}()
	app.smu.Lock()
	for _, s := range app.servers {
		eg.Go(func() (err error) {
			time.AfterFunc(time.Second, func() {
				_ = registry.DefaultRegisterer.RegisterService(ctx, s.Info())
				logger.Logger.WithFields(logrus.Fields{
					"app":    ecode.ModApp,
					"action": "init",
					"name":   s.Info().Name,
					"label":  s.Info().Label(),
				}).Info("start server")
			})
			err = s.Serve()
			return
		})
	}
	app.smu.Unlock()
	return eg.Wait()
}

func (app *Application) startWorkers() error {
	var eg errgroup.Group
	for _, w := range app.workers {
		eg.Go(func() error {
			return w.WorkerStart()
		})
	}
	return eg.Wait()
}

func (app *Application) startJobs() error {
	if len(app.jobs) == 0 {
		return nil
	}
	var jobs = make([]func(), 0)
	for name, runner := range app.jobs {
		jobs = append(jobs, func() {
			// runner.Run 内部 panic 不在此层捕获，由更上层统一处理。
			logger.Logger.WithFields(logrus.Fields{
				"jobName": name,
			}).Info("job run begin")
			defer func() {
				logger.Logger.WithFields(logrus.Fields{
					"jobName": name,
				}).Info("job run end")
			}()
			runner.Run()
		})
	}
	goroutine.Parallel(jobs...)()
	return nil
}

func (app *Application) startExecutors() error {
	return executor.Run()
}

func (app *Application) isDisable(d Disable) bool {
	b, ok := app.disableMap[d]
	if !ok {
		return false
	}
	return b
}

func (app *Application) waitSignals() {
	logger.Logger.WithFields(logrus.Fields{
		"pid": syscall.Getpid(),
	}).Info("init listen signal")
	signals := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP, syscall.SIGUSR1, syscall.SIGQUIT, os.Interrupt}
	app.sig = make(chan os.Signal)
	signal.Notify(app.sig, signals...)
	go app.exitHandler()
}

func (app *Application) exitHandler() {
	debug.FreeOSMemory()
	sig := <-app.sig
	logger.Logger.WithFields(logrus.Fields{
		"pid": syscall.Getpid(),
		"SIG": sig,
	}).Info("Received SIG")

	switch sig {
	case syscall.SIGHUP:
		if err := app.Run(); err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"pid": syscall.Getpid(),
				"SIG": sig,
			}).Info("Received SIG")
		}
		_ = app.GracefulStop(context.Background())
	case syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP, syscall.SIGUSR1, syscall.SIGQUIT, os.Interrupt:
		logger.Logger.Info("close server")
		os.Exit(128 + int(sig.(syscall.Signal)))
	}
}

func (app *Application) printBanner() error {
	if app.HideBanner {
		return nil
	}
	env.PrintVersion()
	return nil
}
