package redis

import (
	"reflect"
	"testing"
	"time"

	"github.com/abulo/ratel/v3/stores/redis"
)

func TestNewDriver(t *testing.T) {
	type args struct {
		client *redis.Client
	}
	tests := []struct {
		name    string
		args    args
		want    *RedisDriver
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDriver(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDriver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDriver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisDriver_Ping(t *testing.T) {
	tests := []struct {
		name    string
		rd      *RedisDriver
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rd.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("RedisDriver.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisDriver_SetTimeout(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		rd   *RedisDriver
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rd.SetTimeout(tt.args.timeout)
		})
	}
}

func TestRedisDriver_SetHeartBeat(t *testing.T) {
	type args struct {
		nodeID string
	}
	tests := []struct {
		name string
		rd   *RedisDriver
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rd.SetHeartBeat(tt.args.nodeID)
		})
	}
}

func TestRedisDriver_heartBeat(t *testing.T) {
	type args struct {
		nodeID string
	}
	tests := []struct {
		name string
		rd   *RedisDriver
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rd.heartBeat(tt.args.nodeID)
		})
	}
}

func TestRedisDriver_GetServiceNodeList(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name    string
		rd      *RedisDriver
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rd.GetServiceNodeList(tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisDriver.GetServiceNodeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisDriver.GetServiceNodeList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisDriver_RegisterServiceNode(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name       string
		rd         *RedisDriver
		args       args
		wantNodeID string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNodeID, err := tt.rd.RegisterServiceNode(tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisDriver.RegisterServiceNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNodeID != tt.wantNodeID {
				t.Errorf("RedisDriver.RegisterServiceNode() = %v, want %v", gotNodeID, tt.wantNodeID)
			}
		})
	}
}

func TestRedisDriver_registerServiceNode(t *testing.T) {
	type args struct {
		nodeID string
	}
	tests := []struct {
		name    string
		rd      *RedisDriver
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rd.registerServiceNode(tt.args.nodeID); (err != nil) != tt.wantErr {
				t.Errorf("RedisDriver.registerServiceNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisDriver_scan(t *testing.T) {
	type args struct {
		matchStr string
	}
	tests := []struct {
		name    string
		rd      *RedisDriver
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rd.scan(tt.args.matchStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisDriver.scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisDriver.scan() = %v, want %v", got, tt.want)
			}
		})
	}
}
