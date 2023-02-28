package query

import (
	"reflect"
	"testing"
)

func TestNewSingleFlight(t *testing.T) {
	tests := []struct {
		name string
		want SingleFlight
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSingleFlight(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSingleFlight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSharedCalls(t *testing.T) {
	tests := []struct {
		name string
		want SingleFlight
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSharedCalls(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSharedCalls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flightGroup_Do(t *testing.T) {
	type args struct {
		key string
		fn  func() (interface{}, error)
	}
	tests := []struct {
		name    string
		g       *flightGroup
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Do(tt.args.key, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("flightGroup.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flightGroup.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flightGroup_DoEx(t *testing.T) {
	type args struct {
		key string
		fn  func() (interface{}, error)
	}
	tests := []struct {
		name      string
		g         *flightGroup
		args      args
		wantVal   interface{}
		wantFresh bool
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotFresh, err := tt.g.DoEx(tt.args.key, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("flightGroup.DoEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("flightGroup.DoEx() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotFresh != tt.wantFresh {
				t.Errorf("flightGroup.DoEx() gotFresh = %v, want %v", gotFresh, tt.wantFresh)
			}
		})
	}
}

func Test_flightGroup_createCall(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name     string
		g        *flightGroup
		args     args
		wantC    *call
		wantDone bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotDone := tt.g.createCall(tt.args.key)
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("flightGroup.createCall() gotC = %v, want %v", gotC, tt.wantC)
			}
			if gotDone != tt.wantDone {
				t.Errorf("flightGroup.createCall() gotDone = %v, want %v", gotDone, tt.wantDone)
			}
		})
	}
}

func Test_flightGroup_makeCall(t *testing.T) {
	type args struct {
		c   *call
		key string
		fn  func() (interface{}, error)
	}
	tests := []struct {
		name string
		g    *flightGroup
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.makeCall(tt.args.c, tt.args.key, tt.args.fn)
		})
	}
}
