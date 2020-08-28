package ratel

import (
	"fmt"
	"sync"

	"github.com/abulo/ratel/cycle"
	"github.com/abulo/ratel/goroutine"
	"github.com/abulo/ratel/terminal"
)

// Ratel Create an instance of Application, by using &Ratel{}
type Ratel struct {
	cycle       *cycle.Cycle
	smu         *sync.RWMutex
	initOnce    sync.Once
	startupOnce sync.Once
	stopOnce    sync.Once
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
		err = goroutine.SerialUntilError(
			app.printBanner,
		)()
	})
	return
}

//printBanner init
func (app *Ratel) printBanner() error {
	const banner = `
   (_)_   _ _ __ (_) |_ ___ _ __
   | | | | | '_ \| | __/ _ \ '__|
   | | |_| | |_) | | ||  __/ |
  _/ |\__,_| .__/|_|\__\___|_|
 |__/      |_|
 
 Welcome to jupiter, starting application ...
`
	fmt.Println(terminal.Green(banner))
	return nil
}
