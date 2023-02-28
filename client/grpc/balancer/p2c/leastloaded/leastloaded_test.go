package leastloaded

import (
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/client/grpc/balancer/p2c/basep2c"
	"google.golang.org/grpc/balancer"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want basep2c.P2c
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_leastLoaded_Add(t *testing.T) {
	type args struct {
		item interface{}
	}
	tests := []struct {
		name string
		p    *leastLoaded
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Add(tt.args.item)
		})
	}
}

func Test_leastLoaded_Next(t *testing.T) {
	tests := []struct {
		name  string
		p     *leastLoaded
		want  interface{}
		want1 func(balancer.DoneInfo)
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.p.Next()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("leastLoaded.Next() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("leastLoaded.Next() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
