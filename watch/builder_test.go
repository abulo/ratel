package watch

import (
	"os/exec"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestOptions_newBuilder(t *testing.T) {
	type fields struct {
		Logger           Logger
		AutoTidy         bool
		MainFiles        string
		OutputName       string
		appName          string
		Flags            Flags
		Exts             []string
		anyExts          bool
		Excludes         []string
		AppArgs          string
		appArgs          []string
		Recursive        bool
		Dirs             []string
		paths            []string
		WatcherFrequency time.Duration
		goCmdArgs        []string
	}
	tests := []struct {
		name   string
		fields fields
		want   *builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := &Options{
				Logger:           tt.fields.Logger,
				AutoTidy:         tt.fields.AutoTidy,
				MainFiles:        tt.fields.MainFiles,
				OutputName:       tt.fields.OutputName,
				appName:          tt.fields.appName,
				Flags:            tt.fields.Flags,
				Exts:             tt.fields.Exts,
				anyExts:          tt.fields.anyExts,
				Excludes:         tt.fields.Excludes,
				AppArgs:          tt.fields.AppArgs,
				appArgs:          tt.fields.appArgs,
				Recursive:        tt.fields.Recursive,
				Dirs:             tt.fields.Dirs,
				paths:            tt.fields.paths,
				WatcherFrequency: tt.fields.WatcherFrequency,
				goCmdArgs:        tt.fields.goCmdArgs,
			}
			if got := opt.newBuilder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.newBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_logf(t *testing.T) {
	type fields struct {
		exts        []string
		anyExt      bool
		excludes    []string
		appName     string
		logs        Logger
		watcherFreq time.Duration
		appWD       string
		appCmd      *exec.Cmd
		appArgs     []string
		appKillMux  sync.Mutex
		goTidy      bool
		goCmd       *exec.Cmd
		goArgs      []string
		goKillMux   sync.Mutex
	}
	type args struct {
		typ int8
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				exts:        tt.fields.exts,
				anyExt:      tt.fields.anyExt,
				excludes:    tt.fields.excludes,
				appName:     tt.fields.appName,
				logs:        tt.fields.logs,
				watcherFreq: tt.fields.watcherFreq,
				appWD:       tt.fields.appWD,
				appCmd:      tt.fields.appCmd,
				appArgs:     tt.fields.appArgs,
				appKillMux:  tt.fields.appKillMux,
				goTidy:      tt.fields.goTidy,
				goCmd:       tt.fields.goCmd,
				goArgs:      tt.fields.goArgs,
				goKillMux:   tt.fields.goKillMux,
			}
			b.logf(tt.args.typ, tt.args.msg)
		})
	}
}

func Test_builder_isIgnore(t *testing.T) {
	type fields struct {
		exts        []string
		anyExt      bool
		excludes    []string
		appName     string
		logs        Logger
		watcherFreq time.Duration
		appWD       string
		appCmd      *exec.Cmd
		appArgs     []string
		appKillMux  sync.Mutex
		goTidy      bool
		goCmd       *exec.Cmd
		goArgs      []string
		goKillMux   sync.Mutex
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				exts:        tt.fields.exts,
				anyExt:      tt.fields.anyExt,
				excludes:    tt.fields.excludes,
				appName:     tt.fields.appName,
				logs:        tt.fields.logs,
				watcherFreq: tt.fields.watcherFreq,
				appWD:       tt.fields.appWD,
				appCmd:      tt.fields.appCmd,
				appArgs:     tt.fields.appArgs,
				appKillMux:  tt.fields.appKillMux,
				goTidy:      tt.fields.goTidy,
				goCmd:       tt.fields.goCmd,
				goArgs:      tt.fields.goArgs,
				goKillMux:   tt.fields.goKillMux,
			}
			if got := b.isIgnore(tt.args.path); got != tt.want {
				t.Errorf("builder.isIgnore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_tidy(t *testing.T) {
	type fields struct {
		exts        []string
		anyExt      bool
		excludes    []string
		appName     string
		logs        Logger
		watcherFreq time.Duration
		appWD       string
		appCmd      *exec.Cmd
		appArgs     []string
		appKillMux  sync.Mutex
		goTidy      bool
		goCmd       *exec.Cmd
		goArgs      []string
		goKillMux   sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				exts:        tt.fields.exts,
				anyExt:      tt.fields.anyExt,
				excludes:    tt.fields.excludes,
				appName:     tt.fields.appName,
				logs:        tt.fields.logs,
				watcherFreq: tt.fields.watcherFreq,
				appWD:       tt.fields.appWD,
				appCmd:      tt.fields.appCmd,
				appArgs:     tt.fields.appArgs,
				appKillMux:  tt.fields.appKillMux,
				goTidy:      tt.fields.goTidy,
				goCmd:       tt.fields.goCmd,
				goArgs:      tt.fields.goArgs,
				goKillMux:   tt.fields.goKillMux,
			}
			b.tidy()
		})
	}
}

func Test_builder_build(t *testing.T) {
	type fields struct {
		exts        []string
		anyExt      bool
		excludes    []string
		appName     string
		logs        Logger
		watcherFreq time.Duration
		appWD       string
		appCmd      *exec.Cmd
		appArgs     []string
		appKillMux  sync.Mutex
		goTidy      bool
		goCmd       *exec.Cmd
		goArgs      []string
		goKillMux   sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				exts:        tt.fields.exts,
				anyExt:      tt.fields.anyExt,
				excludes:    tt.fields.excludes,
				appName:     tt.fields.appName,
				logs:        tt.fields.logs,
				watcherFreq: tt.fields.watcherFreq,
				appWD:       tt.fields.appWD,
				appCmd:      tt.fields.appCmd,
				appArgs:     tt.fields.appArgs,
				appKillMux:  tt.fields.appKillMux,
				goTidy:      tt.fields.goTidy,
				goCmd:       tt.fields.goCmd,
				goArgs:      tt.fields.goArgs,
				goKillMux:   tt.fields.goKillMux,
			}
			b.build()
		})
	}
}

func Test_builder_killGo(t *testing.T) {
	type fields struct {
		exts        []string
		anyExt      bool
		excludes    []string
		appName     string
		logs        Logger
		watcherFreq time.Duration
		appWD       string
		appCmd      *exec.Cmd
		appArgs     []string
		appKillMux  sync.Mutex
		goTidy      bool
		goCmd       *exec.Cmd
		goArgs      []string
		goKillMux   sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				exts:        tt.fields.exts,
				anyExt:      tt.fields.anyExt,
				excludes:    tt.fields.excludes,
				appName:     tt.fields.appName,
				logs:        tt.fields.logs,
				watcherFreq: tt.fields.watcherFreq,
				appWD:       tt.fields.appWD,
				appCmd:      tt.fields.appCmd,
				appArgs:     tt.fields.appArgs,
				appKillMux:  tt.fields.appKillMux,
				goTidy:      tt.fields.goTidy,
				goCmd:       tt.fields.goCmd,
				goArgs:      tt.fields.goArgs,
				goKillMux:   tt.fields.goKillMux,
			}
			b.killGo()
		})
	}
}

func Test_builder_restartApp(t *testing.T) {
	type fields struct {
		exts        []string
		anyExt      bool
		excludes    []string
		appName     string
		logs        Logger
		watcherFreq time.Duration
		appWD       string
		appCmd      *exec.Cmd
		appArgs     []string
		appKillMux  sync.Mutex
		goTidy      bool
		goCmd       *exec.Cmd
		goArgs      []string
		goKillMux   sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				exts:        tt.fields.exts,
				anyExt:      tt.fields.anyExt,
				excludes:    tt.fields.excludes,
				appName:     tt.fields.appName,
				logs:        tt.fields.logs,
				watcherFreq: tt.fields.watcherFreq,
				appWD:       tt.fields.appWD,
				appCmd:      tt.fields.appCmd,
				appArgs:     tt.fields.appArgs,
				appKillMux:  tt.fields.appKillMux,
				goTidy:      tt.fields.goTidy,
				goCmd:       tt.fields.goCmd,
				goArgs:      tt.fields.goArgs,
				goKillMux:   tt.fields.goKillMux,
			}
			b.restartApp()
		})
	}
}

func Test_builder_killApp(t *testing.T) {
	type fields struct {
		exts        []string
		anyExt      bool
		excludes    []string
		appName     string
		logs        Logger
		watcherFreq time.Duration
		appWD       string
		appCmd      *exec.Cmd
		appArgs     []string
		appKillMux  sync.Mutex
		goTidy      bool
		goCmd       *exec.Cmd
		goArgs      []string
		goKillMux   sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				exts:        tt.fields.exts,
				anyExt:      tt.fields.anyExt,
				excludes:    tt.fields.excludes,
				appName:     tt.fields.appName,
				logs:        tt.fields.logs,
				watcherFreq: tt.fields.watcherFreq,
				appWD:       tt.fields.appWD,
				appCmd:      tt.fields.appCmd,
				appArgs:     tt.fields.appArgs,
				appKillMux:  tt.fields.appKillMux,
				goTidy:      tt.fields.goTidy,
				goCmd:       tt.fields.goCmd,
				goArgs:      tt.fields.goArgs,
				goKillMux:   tt.fields.goKillMux,
			}
			b.killApp()
		})
	}
}

func Test_builder_filterPaths(t *testing.T) {
	type fields struct {
		exts        []string
		anyExt      bool
		excludes    []string
		appName     string
		logs        Logger
		watcherFreq time.Duration
		appWD       string
		appCmd      *exec.Cmd
		appArgs     []string
		appKillMux  sync.Mutex
		goTidy      bool
		goCmd       *exec.Cmd
		goArgs      []string
		goKillMux   sync.Mutex
	}
	type args struct {
		paths []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				exts:        tt.fields.exts,
				anyExt:      tt.fields.anyExt,
				excludes:    tt.fields.excludes,
				appName:     tt.fields.appName,
				logs:        tt.fields.logs,
				watcherFreq: tt.fields.watcherFreq,
				appWD:       tt.fields.appWD,
				appCmd:      tt.fields.appCmd,
				appArgs:     tt.fields.appArgs,
				appKillMux:  tt.fields.appKillMux,
				goTidy:      tt.fields.goTidy,
				goCmd:       tt.fields.goCmd,
				goArgs:      tt.fields.goArgs,
				goKillMux:   tt.fields.goKillMux,
			}
			if got := b.filterPaths(tt.args.paths); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.filterPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
