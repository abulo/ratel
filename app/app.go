package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/abulo/ratel/v3/cycle"
	"github.com/abulo/ratel/v3/ecode"
	"github.com/abulo/ratel/v3/goroutine"
	"github.com/abulo/ratel/v3/hooks"
	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/registry"
	"github.com/abulo/ratel/v3/server"
	"github.com/abulo/ratel/v3/worker"
	"golang.org/x/sync/errgroup"
)

type Application struct {
	cycle    *cycle.Cycle
	smu      *sync.RWMutex
	initOnce sync.Once
	stopOnce sync.Once
	servers  []server.Server
	workers  []worker.Worker
	stopped  chan struct{}
	sig      chan os.Signal
}

// New create a new Application instance
func New(fns ...func() error) (*Application, error) {
	app := &Application{}
	if err := app.Startup(fns...); err != nil {
		return nil, err
	}
	return app, nil
}

func DefaultApp() *Application {
	app := &Application{}
	app.initialize()
	return app
}

func (app *Application) Startup(fns ...func() error) error {
	app.initialize()
	return goroutine.SerialUntilError(fns...)()
}

//run hooks
func (app *Application) runHooks(stage hooks.Stage) {
	hooks.Do(stage)
}

//RegisterHooks register a stage Hook
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
		app.stopped = make(chan struct{})
	})
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

// SetRegistry set customize registry
// Deprecated, please use registry.DefaultRegisterer instead.
func (app *Application) SetRegistry(reg registry.Registry) {
	registry.DefaultRegisterer = reg
}

// Run run application
func (app *Application) Run(servers ...server.Server) error {
	app.smu.Lock()
	app.servers = append(app.servers, servers...)
	app.smu.Unlock()

	app.waitSignals() //start signal listen task in goroutine
	// defer app.clean()

	// start servers and govern server
	app.cycle.Run(app.startServers)
	// start workers
	app.cycle.Run(app.startWorkers)
	//blocking and wait quit
	if err := <-app.cycle.Wait(); err != nil {
		logger.Logger.Error("ratel shutdown with error", ecode.ModApp, err)
		return err
	}
	logger.Logger.Info("shutdown ratel, bye!", ecode.ModApp)
	return nil
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

func (app *Application) startServers() error {
	var eg errgroup.Group
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	go func() {
		<-app.stopped
		cancel()
	}()
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				_ = registry.DefaultRegisterer.UnregisterService(ctx, s.Info())
				logger.Logger.Info("exit server", ecode.ModApp, "exit", s.Info().Name, err, s.Info().Label())
			}()

			time.AfterFunc(time.Second, func() {
				_ = registry.DefaultRegisterer.RegisterService(ctx, s.Info())
				logger.Logger.Info("start server", ecode.ModApp, "init", s.Info().Name, s.Info().Label(), "scheme", s.Info().Scheme)
			})
			err = s.Serve()
			return
		})
	}
	return eg.Wait()
}

func (app *Application) waitSignals() {
	logger.Logger.Info("init listen signal, pid:", syscall.Getpid())
	signals := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP, syscall.SIGUSR1, syscall.SIGQUIT, os.Interrupt}
	app.sig = make(chan os.Signal)
	signal.Notify(app.sig, signals...)
	go app.exitHandler()
}

func (app *Application) exitHandler() {
	debug.FreeOSMemory()
	sig := <-app.sig
	logger.Logger.Info(fmt.Sprintf("Received SIG. [PID:%d, SIG:%v]", syscall.Getpid(), sig))
	switch sig {
	case syscall.SIGHUP:
		if err := app.Run(); err != nil {
			logger.Logger.Error(fmt.Sprintf("Received SIG. [PID:%d, SIG:%v]", syscall.Getpid(), sig))
		}
		app.GracefulStop(context.Background())
	case syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP, syscall.SIGUSR1, syscall.SIGQUIT, os.Interrupt:
		logger.Logger.Info("close server")
		os.Exit(128 + int(sig.(syscall.Signal)))
	}
}

// GracefulStop application after necessary cleanup
func (app *Application) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.stopped <- struct{}{}
		app.runHooks(hooks.Stage_BeforeStop)
		//stop servers
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(func() error {
					return s.GracefulStop(ctx)
				})
			}(s)
		}
		app.smu.RUnlock()
		//stop workers
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.WorkerStop)
			}(w)
		}
		// stop executor
		<-app.cycle.Done()
		// run hooks
		app.runHooks(hooks.Stage_AfterStop)
		app.cycle.Close()
	})
	return err
}

// Stop application immediately after necessary cleanup
func (app *Application) Stop() (err error) {
	app.stopOnce.Do(func() {
		app.stopped <- struct{}{}
		app.runHooks(hooks.Stage_BeforeStop)
		//stop servers
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(s.Stop)
			}(s)
		}
		app.smu.RUnlock()
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
