package mongodb

import (
	"reflect"
	"testing"
	"time"
)

func TestMongoDB_Close(t *testing.T) {
	tests := []struct {
		name    string
		cs      *MongoDB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cs.Close(); (err != nil) != tt.wantErr {
				t.Errorf("MongoDB.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetSlowThreshold(t *testing.T) {
	type args struct {
		threshold time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetSlowThreshold(tt.args.threshold)
		})
	}
}

func Test_defaultTimeoutOption(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultTimeoutOption(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultTimeoutOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		timeout time.Duration
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
			if got := WithTimeout(tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxConnIdleTime(t *testing.T) {
	type args struct {
		d time.Duration
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
			if got := WithMaxConnIdleTime(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxConnIdleTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxPoolSize(t *testing.T) {
	type args struct {
		size uint64
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
			if got := WithMaxPoolSize(tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxPoolSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMinPoolSize(t *testing.T) {
	type args struct {
		size uint64
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
			if got := WithMinPoolSize(tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMinPoolSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMongoDBClient(t *testing.T) {
	type args struct {
		uri  string
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    *MongoDB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMongoDBClient(tt.args.uri, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMongoDBClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMongoDBClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_SetDisableMetric(t *testing.T) {
	type args struct {
		disableMetric bool
	}
	tests := []struct {
		name string
		m    *MongoDB
		args args
		want *MongoDB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.SetDisableMetric(tt.args.disableMetric); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MongoDB.SetDisableMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_SetDisableTrace(t *testing.T) {
	type args struct {
		disableTrace bool
	}
	tests := []struct {
		name string
		m    *MongoDB
		args args
		want *MongoDB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.SetDisableTrace(tt.args.disableTrace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MongoDB.SetDisableTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_Ping(t *testing.T) {
	tests := []struct {
		name    string
		m       *MongoDB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("MongoDB.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDB_NewModel(t *testing.T) {
	type args struct {
		collection string
	}
	tests := []struct {
		name    string
		m       *MongoDB
		args    args
		want    *Model
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.NewModel(tt.args.collection)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoDB.NewModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MongoDB.NewModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_MustNewModel(t *testing.T) {
	type args struct {
		collection string
	}
	tests := []struct {
		name string
		m    *MongoDB
		args args
		want *Model
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.MustNewModel(tt.args.collection); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MongoDB.MustNewModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_NewCollection(t *testing.T) {
	type args struct {
		collection string
	}
	tests := []struct {
		name    string
		m       *MongoDB
		args    args
		want    *DecoratedCollection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.NewCollection(tt.args.collection)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoDB.NewCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MongoDB.NewCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDB_MustNewCollection(t *testing.T) {
	type args struct {
		collection string
	}
	tests := []struct {
		name string
		m    *MongoDB
		args args
		want *DecoratedCollection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.MustNewCollection(tt.args.collection); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MongoDB.MustNewCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
