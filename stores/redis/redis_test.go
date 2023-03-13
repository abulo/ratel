package redis

import (
	"reflect"
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRedisClient(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedisClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRedisClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithClientType(t *testing.T) {
	type args struct {
		ClientType string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithClientType(tt.args.ClientType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithClientType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithHosts(t *testing.T) {
	type args struct {
		Hosts []string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithHosts(tt.args.Hosts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPassword(t *testing.T) {
	type args struct {
		Password string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithPassword(tt.args.Password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDatabase(t *testing.T) {
	type args struct {
		Database int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDatabase(tt.args.Database); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPoolSize(t *testing.T) {
	type args struct {
		PoolSize int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithPoolSize(tt.args.PoolSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPoolSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithKeyPrefix(t *testing.T) {
	type args struct {
		KeyPrefix string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithKeyPrefix(tt.args.KeyPrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithKeyPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableMetric(t *testing.T) {
	type args struct {
		DisableMetric bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableMetric(tt.args.DisableMetric); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableTrace(t *testing.T) {
	type args struct {
		DisableTrace bool
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableTrace(tt.args.DisableTrace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithAddr(t *testing.T) {
	type args struct {
		Addr string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAddr(tt.args.Addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithAddrs(t *testing.T) {
	type args struct {
		Addrs map[string]string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAddrs(tt.args.Addrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAddrs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMasterName(t *testing.T) {
	type args struct {
		MasterName string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMasterName(tt.args.MasterName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMasterName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRedis(t *testing.T) {
	type args struct {
		r *Client
	}
	tests := []struct {
		name    string
		args    args
		want    RedisNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRedis(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRedis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRedis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClient(t *testing.T) {
	type args struct {
		r *Client
	}
	tests := []struct {
		name    string
		args    args
		want    RedisNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getClient(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCluster(t *testing.T) {
	type args struct {
		r *Client
	}
	tests := []struct {
		name    string
		args    args
		want    RedisNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCluster(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFailover(t *testing.T) {
	type args struct {
		r *Client
	}
	tests := []struct {
		name    string
		args    args
		want    RedisNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFailover(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFailover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFailover() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRing(t *testing.T) {
	type args struct {
		r *Client
	}
	tests := []struct {
		name    string
		args    args
		want    RedisNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRing(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRing() = %v, want %v", got, tt.want)
			}
		})
	}
}
