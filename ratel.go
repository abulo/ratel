package ratel

import (
	"context"
	"sync"

	"github.com/abulo/ratel/cycle"
	"github.com/abulo/ratel/goroutine"
	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/registry"
	"github.com/abulo/ratel/server"
	"github.com/abulo/ratel/signals"
	"github.com/abulo/ratel/trace"
	"github.com/abulo/ratel/worker"
	"golang.org/x/sync/errgroup"
)

// Ratel Create an instance of Application, by using &Ratel{}
type Ratel struct {
	cycle       *cycle.Cycle
	smu         *sync.RWMutex
	initOnce    sync.Once
	startupOnce sync.Once
	stopOnce    sync.Once
	servers     []server.Server
	workers     []worker.Worker
	registerer  registry.Registry
}

//New new a Application
func New(fns ...func() error) (*Ratel, error) {
	app := &Ratel{}
	if err := app.Startup(fns...); err != nil {
		return nil, err
	}
	return app, nil
}

//Startup ..
func (app *Ratel) Startup(fns ...func() error) error {
	app.initialize()
	if err := app.startup(); err != nil {
		return err
	}
	return goroutine.SerialUntilError(fns...)()
}

// initialize application
func (app *Ratel) initialize() {
	app.initOnce.Do(func() {
		//assign
		app.cycle = cycle.NewCycle()
		app.smu = &sync.RWMutex{}
		app.servers = make([]server.Server, 0)
		app.workers = make([]worker.Worker, 0)
		app.Registry(registry.Nop{})
	})
}

// start up application
// By default the startup composition is:
// - parse config, watch, version flags
// - load config
// - init default biz logger, jupiter frame logger
// - init procs
func (app *Ratel) startup() (err error) {
	app.startupOnce.Do(func() {
		err = goroutine.SerialUntilError()()
	})
	return
}

// Serve start server
func (app *Ratel) Serve(s ...server.Server) error {
	app.smu.Lock()
	defer app.smu.Unlock()
	app.servers = append(app.servers, s...)
	return nil
}

// Registry set customize registry
func (app *Ratel) Registry(reg registry.Registry) *Ratel {
	app.registerer = reg
	return app
}

// Worker ..
func (app *Ratel) Worker(w worker.Worker) *Ratel {
	app.workers = append(app.workers, w)
	return app
}

//Tracer ...
func (app *Ratel) Tracer(name, host string) *Ratel {
	cfg := trace.InitConfig(host)
	cfg.ServiceName = name
	if err := trace.New(cfg).Setup(); err != nil {
		logger.Logger.Panic(err)
	}
	return app
}

// Run run application
func (app *Ratel) Run(servers ...server.Server) error {
	app.smu.Lock()
	app.servers = append(app.servers, servers...)
	app.smu.Unlock()
	app.waitSignals() //start signal listen task in goroutine

	// start servers and govern server
	app.cycle.Run(app.startServers)

	// start workers
	app.cycle.Run(app.startWorkers)

	//blocking and wait quit
	if err := <-app.cycle.Wait(); err != nil {
		logger.Logger.Error("shutdown with error", err)
		return err
	}
	logger.Logger.Info("shutdown, bye!")
	return nil
}

// waitSignals wait signal
func (app *Ratel) waitSignals() {
	logger.Logger.Info("init listen signal")
	signals.Shutdown(func(grace bool) { //when get shutdown signal
		//todo: support timeout
		if grace {
			app.GracefulStop(context.TODO())
		} else {
			app.Stop()
		}
	})
}

// GracefulStop application after necessary cleanup
func (app *Ratel) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		if app.registerer != nil {
			err = app.registerer.Close()
			if err != nil {
				logger.Logger.Error("stop register close err", err)
			}
		}
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

		<-app.cycle.Done()
		app.cycle.Close()
	})
	return nil
}

// Stop application immediately after necessary cleanup
func (app *Ratel) Stop() (err error) {
	app.stopOnce.Do(func() {

		if app.registerer != nil {
			err = app.registerer.Close()
			if err != nil {
				logger.Logger.Error("stop register close err", err)
			}
		}

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
		app.cycle.Close()
	})
	return nil
}

func (app *Ratel) startServers() error {
	var eg errgroup.Group
	// start multi servers
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			_ = app.registerer.RegisterService(context.TODO(), s.Info())
			defer app.registerer.UnregisterService(context.TODO(), s.Info())
			logger.Logger.Info("start server:", s.Info().Name, ":", s.Info().Label(), ":", s.Info().Scheme)
			err = s.Serve()
			return
		})
	}
	return eg.Wait()
}

func (app *Ratel) startWorkers() error {
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
