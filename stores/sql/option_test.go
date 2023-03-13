package sql

import (
	"reflect"
	"testing"
	"time"
)

func TestWithUsername(t *testing.T) {
	type args struct {
		username string
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
			if got := WithUsername(tt.args.username); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPassword(t *testing.T) {
	type args struct {
		password string
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
			if got := WithPassword(tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithHost(t *testing.T) {
	type args struct {
		host string
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
			if got := WithHost(tt.args.host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPort(t *testing.T) {
	type args struct {
		port string
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
			if got := WithPort(tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCharset(t *testing.T) {
	type args struct {
		charset string
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
			if got := WithCharset(tt.args.charset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCharset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDatabase(t *testing.T) {
	type args struct {
		database string
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
			if got := WithDatabase(tt.args.database); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTimeZone(t *testing.T) {
	type args struct {
		timeZone string
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
			if got := WithTimeZone(tt.args.timeZone); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTimeZone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxOpenConns(t *testing.T) {
	type args struct {
		maxOpenConns int
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
			if got := WithMaxOpenConns(tt.args.maxOpenConns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxOpenConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxIdleConns(t *testing.T) {
	type args struct {
		maxIdleConns int
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
			if got := WithMaxIdleConns(tt.args.maxIdleConns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxIdleConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxLifetime(t *testing.T) {
	type args struct {
		maxLifetime time.Duration
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
			if got := WithMaxLifetime(tt.args.maxLifetime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxLifetime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxIdleTime(t *testing.T) {
	type args struct {
		maxIdleTime time.Duration
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
			if got := WithMaxIdleTime(tt.args.maxIdleTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxIdleTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDriverName(t *testing.T) {
	type args struct {
		driverName string
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
			if got := WithDriverName(tt.args.driverName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDriverName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableMetric(t *testing.T) {
	type args struct {
		disableMetric bool
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
			if got := WithDisableMetric(tt.args.disableMetric); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableTrace(t *testing.T) {
	type args struct {
		disableTrace bool
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
			if got := WithDisableTrace(tt.args.disableTrace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithAddr(t *testing.T) {
	type args struct {
		addr []string
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
			if got := WithAddr(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDialTimeout(t *testing.T) {
	type args struct {
		dialTimeout string
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
			if got := WithDialTimeout(tt.args.dialTimeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDialTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOpenStrategy(t *testing.T) {
	type args struct {
		openStrategy string
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
			if got := WithOpenStrategy(tt.args.openStrategy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOpenStrategy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCompress(t *testing.T) {
	type args struct {
		compress bool
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
			if got := WithCompress(tt.args.compress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCompress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxExecutionTime(t *testing.T) {
	type args struct {
		maxExecutionTime string
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
			if got := WithMaxExecutionTime(tt.args.maxExecutionTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxExecutionTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableDebug(t *testing.T) {
	type args struct {
		disableDebug bool
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
			if got := WithDisableDebug(tt.args.disableDebug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisablePrepare(t *testing.T) {
	type args struct {
		disablePrepare bool
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
			if got := WithDisablePrepare(tt.args.disablePrepare); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisablePrepare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithParseTime(t *testing.T) {
	type args struct {
		parseTime bool
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
			if got := WithParseTime(tt.args.parseTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithParseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
