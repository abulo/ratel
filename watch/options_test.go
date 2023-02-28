package watch

import (
	"reflect"
	"testing"
	"time"
)

func TestOptions_sanitize(t *testing.T) {
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
		name    string
		fields  fields
		wantErr bool
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
			if err := opt.sanitize(); (err != nil) != tt.wantErr {
				t.Errorf("Options.sanitize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptions_sanitizeExts(t *testing.T) {
	tests := []struct {
		name string
		opt  *Options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opt.sanitizeExts()
		})
	}
}

func Test_getAppName(t *testing.T) {
	type args struct {
		outputName string
		wd         string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAppName(tt.args.outputName, tt.args.wd)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAppName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAppName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_recursivePaths(t *testing.T) {
	type args struct {
		recursive bool
		paths     []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := recursivePaths(tt.args.recursive, tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("recursivePaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("recursivePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitArgs(t *testing.T) {
	type args struct {
		args string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitArgs(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendArg(t *testing.T) {
	type args struct {
		args []string
		arg  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendArg(tt.args.args, tt.args.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
