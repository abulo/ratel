package etcdv3

import (
	"reflect"
	"testing"
	"time"
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

func TestConfig_Build(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    *Client
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
		want    *Client
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

func TestConfig_MustBuild(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   *Client
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

func TestConfig_MustSingleton(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   *Client
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

func TestConfig_SetEndpoints(t *testing.T) {
	type args struct {
		endpoint []string
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
			if got := tt.config.SetEndpoints(tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetEndpoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetCertFile(t *testing.T) {
	type args struct {
		cert string
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
			if got := tt.config.SetCertFile(tt.args.cert); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetCertFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetKeyFile(t *testing.T) {
	type args struct {
		key string
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
			if got := tt.config.SetKeyFile(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetKeyFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetCaCert(t *testing.T) {
	type args struct {
		ca string
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
			if got := tt.config.SetCaCert(tt.args.ca); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetCaCert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetBasicAuth(t *testing.T) {
	type args struct {
		auth bool
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
			if got := tt.config.SetBasicAuth(tt.args.auth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetBasicAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetUserName(t *testing.T) {
	type args struct {
		userName string
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
			if got := tt.config.SetUserName(tt.args.userName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetPassword(t *testing.T) {
	type args struct {
		pwd string
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
			if got := tt.config.SetPassword(tt.args.pwd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetConnectTimeout(t *testing.T) {
	type args struct {
		timeout time.Duration
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
			if got := tt.config.SetConnectTimeout(tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetConnectTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetSecure(t *testing.T) {
	type args struct {
		secure bool
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
			if got := tt.config.SetSecure(tt.args.secure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetSecure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetEnableTrace(t *testing.T) {
	type args struct {
		enableTrace bool
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
			if got := tt.config.SetEnableTrace(tt.args.enableTrace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetEnableTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetAutoSyncInterval(t *testing.T) {
	type args struct {
		autoSyncInterval time.Duration
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
			if got := tt.config.SetAutoSyncInterval(tt.args.autoSyncInterval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetAutoSyncInterval() = %v, want %v", got, tt.want)
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

func TestConfig_GetName(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.GetName(); got != tt.want {
				t.Errorf("Config.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}
