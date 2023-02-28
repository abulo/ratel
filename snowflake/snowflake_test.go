package snowflake

import (
	"net"
	"reflect"
	"testing"
)

func TestConstructConfig(t *testing.T) {
	type args struct {
		timeBits    uint
		seqBits     uint
		machineBits uint
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConstructConfig(tt.args.timeBits, tt.args.seqBits, tt.args.machineBits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConstructConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstructConfigWithMachineID(t *testing.T) {
	type args struct {
		timeBits    uint
		seqBits     uint
		machineBits uint
		machineID   int64
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConstructConfigWithMachineID(tt.args.timeBits, tt.args.seqBits, tt.args.machineBits, tt.args.machineID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConstructConfigWithMachineID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_getCurrentTimestamp(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.getCurrentTimestamp(); got != tt.want {
				t.Errorf("Config.getCurrentTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GenInt64ID(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GenInt64ID(); got != tt.want {
				t.Errorf("Config.GenInt64ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetSeqFromID(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name string
		c    *Config
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetSeqFromID(tt.args.id); got != tt.want {
				t.Errorf("Config.GetSeqFromID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetTimeFromID(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name string
		c    *Config
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetTimeFromID(tt.args.id); got != tt.want {
				t.Errorf("Config.GetTimeFromID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMachineID(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMachineID(); got != tt.want {
				t.Errorf("getMachineID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStartEpochFromEnv(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStartEpochFromEnv(); got != tt.want {
				t.Errorf("getStartEpochFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIP(t *testing.T) {
	tests := []struct {
		name    string
		want    net.IP
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIP()
			if (err != nil) != tt.wantErr {
				t.Errorf("getIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnv(t *testing.T) {
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnv(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
