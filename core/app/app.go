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

// Application ...
type Application struct {
	cycle      *cycle.Cycle
	smu        *sync.RWMutex
	initOnce   sync.Once
	stopOnce   sync.Once
	servers    []server.Server
	workers    []worker.Worker
	jobs       map[string]job.Runner
	disableMap map[Disable]bool
	HideBanner bool
	stopped    chan struct{}
	components []component.Component
	sig        chan os.Signal
}

type Option func(a *Application)

type Disable int

const (
	DisableParserFlag      Disable = 1
	DisableLoadConfig      Disable = 2
	DisableDefaultGovernor Disable = 3
)

func (a *Application) WithOptions(options ...Option) {
	for _, option := range options {
		option(a)
	}
}

func WithDisable(d Disable) Option {
	return func(a *Application) {
		a.disableMap[d] = true
	}
}

// New create a new Application instance
func New(fns ...func() error) (*Application, error) {
	app := &Application{}
	if err := app.Startup(fns...); err != nil {
		return nil, err
	}
	return app, nil
}

// DefaultApp ...
func DefaultApp() *Application {
	app := &Application{}
	app.initialize()
	return app
}

// run hooks
func (app *Application) runHooks(stage hooks.Stage) {
	hooks.Do(stage)
}

// RegisterHooks register a stage Hook
func (app *Application) RegisterHooks(stage hooks.Stage, fns ...func()) {
	hooks.Register(stage, fns...)
}

// initialize application
func (app *Application) initialize() {
	app.initOnce.Do(func() {
		//assign
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

// Startup ...
func (app *Application) Startup(fns ...func() error) error {
	app.initialize()
	return goroutine.SerialUntilError(fns...)()
}

// Serve start server
func (app *Application) Serve(s ...server.Server) error {
	app.smu.Lock()
	defer app.smu.Unlock()
	app.servers = append(app.servers, s...)
	return nil
}

// Schedule ..
func (app *Application) Schedule(w worker.Worker) error {
	app.workers = append(app.workers, w)
	return nil
}

// Job ..
func (app *Application) Job(runner job.Runner) error {
	namedJob, ok := runner.(interface{ GetJobName() string })
	// job runner must implement GetJobName
	if !ok {
		return nil
	}
	jobName := namedJob.GetJobName()
	app.jobs[jobName] = runner
	return nil
}

// Executor ...
func (app *Application) Executor(e executor.Executor) {
	executor.Register(e.GetAddress(), e)
}

// SetRegistry set customize registry
func (app *Application) SetRegistry(reg registry.Registry) {
	registry.DefaultRegisterer = reg
}

// clean after app quit
func (app *Application) clean() {

}

// Run run application
func (app *Application) Run(servers ...server.Server) error {
	app.smu.Lock()
	app.servers = append(app.servers, servers...)
	app.smu.Unlock()

	hooks.Do(hooks.Stage_BeforeRun)

	app.waitSignals() //start signal listen task in goroutine
	defer app.clean()

	// todo jobs not graceful
	_ = app.startJobs()

	// start servers and govern server
	app.cycle.Run(app.startServers)
	// start workers
	app.cycle.Run(app.startWorkers)
	//blocking and wait quit
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

// Stop application immediately after necessary cleanup
func (app *Application) Stop() (err error) {
	app.stopOnce.Do(func() {
		app.stopped <- struct{}{}
		app.runHooks(hooks.Stage_BeforeStop)
		//stop servers
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
		//stop workers
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.WorkerStop)
			}(w)
		}

		<-app.cycle.Done()
		// run hook
		app.runHooks(hooks.Stage_AfterStop)
		app.cycle.Close()
	})
	return
}

// GracefulStop application after necessary cleanup
func (app *Application) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.stopped <- struct{}{}
		app.runHooks(hooks.Stage_BeforeStop)
		//stop servers
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
		//stop workers
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.WorkerStop)
			}(w)
		}
		// stop executor
		app.cycle.Run(executor.GracefulStop)
		<-app.cycle.Done()
		// run hooks
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
	// start multi servers
	app.smu.Lock()
	for _, s := range app.servers {
		s := s
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
	// start multi workers
	for _, w := range app.workers {
		w := w
		eg.Go(func() error {
			return w.WorkerStart()
		})
	}
	return eg.Wait()
}

// startJobs starts jobs
func (app *Application) startJobs() error {
	if len(app.jobs) == 0 {
		return nil
	}
	var jobs = make([]func(), 0)
	//warp jobs
	for name, runner := range app.jobs {
		jobs = append(jobs, func() {
			// runner.Run panic 错误在更上层抛出
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

// start executor
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

// printBanner init
func (app *Application) printBanner() error {
	if app.HideBanner {
		return nil
	}
	env.PrintVersion()
	return nil
}
