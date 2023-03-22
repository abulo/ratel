package task

import (
	"reflect"
	"testing"
	"time"

	"github.com/abulo/ratel/v3/core/task/driver"
)

func Test_newNodePool(t *testing.T) {
	type args struct {
		serverName     string
		driver         driver.Driver
		task           *Task
		updateDuration time.Duration
		hashReplicas   int
	}
	tests := []struct {
		name    string
		args    args
		want    *NodePool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newNodePool(tt.args.serverName, tt.args.driver, tt.args.task, tt.args.updateDuration, tt.args.hashReplicas)
			if (err != nil) != tt.wantErr {
				t.Errorf("newNodePool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newNodePool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodePool_StartPool(t *testing.T) {
	tests := []struct {
		name    string
		np      *NodePool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.np.StartPool(); (err != nil) != tt.wantErr {
				t.Errorf("NodePool.StartPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodePool_updatePool(t *testing.T) {
	tests := []struct {
		name    string
		np      *NodePool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.np.updatePool(); (err != nil) != tt.wantErr {
				t.Errorf("NodePool.updatePool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodePool_tickerUpdatePool(t *testing.T) {
	tests := []struct {
		name string
		np   *NodePool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.np.tickerUpdatePool()
		})
	}
}

func TestNodePool_PickNodeByJobName(t *testing.T) {
	type args struct {
		jobName string
	}
	tests := []struct {
		name string
		np   *NodePool
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.np.PickNodeByJobName(tt.args.jobName); got != tt.want {
				t.Errorf("NodePool.PickNodeByJobName() = %v, want %v", got, tt.want)
			}
		})
	}
}
