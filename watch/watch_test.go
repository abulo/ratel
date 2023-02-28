package watch

import (
	"context"
	"os/exec"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestWatch(t *testing.T) {
	type args struct {
		ctx context.Context
		opt *Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Watch(tt.args.ctx, tt.args.opt); (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_builder_watch(t *testing.T) {
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
		ctx   context.Context
		paths []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
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
			if err := b.watch(tt.args.ctx, tt.args.paths); (err != nil) != tt.wantErr {
				t.Errorf("builder.watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_builder_initWatcher(t *testing.T) {
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
		name    string
		fields  fields
		args    args
		want    *fsnotify.Watcher
		wantErr bool
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
			got, err := b.initWatcher(tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("builder.initWatcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.initWatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoVersion(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("GoVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GoVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
