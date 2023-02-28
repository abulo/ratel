package grpc

import (
	"reflect"
	"testing"
	"time"

	"github.com/abulo/ratel/v3/registry/etcdv3"
	"google.golang.org/grpc"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Config
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

func TestConfig_SetEtcd(t *testing.T) {
	type args struct {
		option *etcdv3.Config
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetEtcd(tt.args.option); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetEtcd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetName(t *testing.T) {
	type args struct {
		Name string
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetName(tt.args.Name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetBalancerName(t *testing.T) {
	type args struct {
		BalancerName string
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetBalancerName(tt.args.BalancerName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetBalancerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetAddress(t *testing.T) {
	type args struct {
		Address string
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetAddress(tt.args.Address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetBlock(t *testing.T) {
	type args struct {
		Block bool
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetBlock(tt.args.Block); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDialTimeout(t *testing.T) {
	type args struct {
		DialTimeout time.Duration
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetDialTimeout(tt.args.DialTimeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetDialTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetReadTimeout(t *testing.T) {
	type args struct {
		ReadTimeout time.Duration
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetReadTimeout(tt.args.ReadTimeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetReadTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDirect(t *testing.T) {
	type args struct {
		Direct bool
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetDirect(tt.args.Direct); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetDirect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetSlowThreshold(t *testing.T) {
	type args struct {
		SlowThreshold time.Duration
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetSlowThreshold(tt.args.SlowThreshold); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetSlowThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDebug(t *testing.T) {
	type args struct {
		Debug bool
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetDebug(tt.args.Debug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDisableTraceInterceptor(t *testing.T) {
	type args struct {
		DisableTraceInterceptor bool
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetDisableTraceInterceptor(tt.args.DisableTraceInterceptor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetDisableTraceInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDisableAidInterceptor(t *testing.T) {
	type args struct {
		DisableAidInterceptor bool
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetDisableAidInterceptor(tt.args.DisableAidInterceptor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetDisableAidInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDisableTimeoutInterceptor(t *testing.T) {
	type args struct {
		DisableTimeoutInterceptor bool
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetDisableTimeoutInterceptor(tt.args.DisableTimeoutInterceptor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetDisableTimeoutInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDisableMetricInterceptor(t *testing.T) {
	type args struct {
		DisableMetricInterceptor bool
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetDisableMetricInterceptor(tt.args.DisableMetricInterceptor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetDisableMetricInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDisableAccessInterceptor(t *testing.T) {
	type args struct {
		DisableAccessInterceptor bool
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetDisableAccessInterceptor(tt.args.DisableAccessInterceptor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetDisableAccessInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetAccessInterceptorLevel(t *testing.T) {
	type args struct {
		AccessInterceptorLevel string
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetAccessInterceptorLevel(tt.args.AccessInterceptorLevel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetAccessInterceptorLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_WithDialOption(t *testing.T) {
	type args struct {
		opts []grpc.DialOption
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.WithDialOption(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.WithDialOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Build(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    *grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Singleton(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    *grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.Singleton()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Singleton() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.Singleton() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_MustSingleton(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   *grpc.ClientConn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.MustSingleton(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.MustSingleton() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_MustBuild(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   *grpc.ClientConn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.MustBuild(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.MustBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
