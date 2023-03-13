package redis

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestClient_Prefix(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		r    *Client
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Prefix(tt.args.key); got != tt.want {
				t.Errorf("Client.Prefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_k(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		r    *Client
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.k(tt.args.key); got != tt.want {
				t.Errorf("Client.k() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ks(t *testing.T) {
	type args struct {
		key []string
	}
	tests := []struct {
		name string
		r    *Client
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.ks(tt.args.key...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_acceptable(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := acceptable(tt.args.err); got != tt.want {
				t.Errorf("acceptable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Pipeline(t *testing.T) {
	tests := []struct {
		name    string
		r       *Client
		wantVal redis.Pipeliner
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Pipeline()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Pipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.Pipeline() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_getCtx(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		r    *Client
		args args
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.getCtx(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.getCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Pipelined(t *testing.T) {
	type args struct {
		ctx context.Context
		fn  func(redis.Pipeliner) error
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Cmder
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Pipelined(tt.args.ctx, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Pipelined() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.Pipelined() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_TxPipelined(t *testing.T) {
	type args struct {
		ctx context.Context
		fn  func(redis.Pipeliner) error
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Cmder
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.TxPipelined(tt.args.ctx, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TxPipelined() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.TxPipelined() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_TxPipeline(t *testing.T) {
	tests := []struct {
		name    string
		r       *Client
		wantVal redis.Pipeliner
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.TxPipeline()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TxPipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.TxPipeline() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Command(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal map[string]*redis.CommandInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Command(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Command() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.Command() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClientGetName(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClientGetName(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClientGetName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClientGetName() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Echo(t *testing.T) {
	type args struct {
		ctx     context.Context
		message any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Echo(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Echo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Echo() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Ping(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Ping(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Ping() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Quit(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Quit(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Quit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Quit() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Del(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Del(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Del() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Unlink(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Unlink(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Unlink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Unlink() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Dump(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Dump(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Dump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Dump() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Exists(t *testing.T) {
	type args struct {
		ctx context.Context
		key []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Exists(tt.args.ctx, tt.args.key...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Exists() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Expire(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Expire(tt.args.ctx, tt.args.key, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Expire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Expire() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ExpireAt(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		tm  time.Time
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ExpireAt(tt.args.ctx, tt.args.key, tt.args.tm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ExpireAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ExpireAt() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ExpireNX(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		tm  time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ExpireNX(tt.args.ctx, tt.args.key, tt.args.tm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ExpireNX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ExpireNX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ExpireXX(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		tm  time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ExpireXX(tt.args.ctx, tt.args.key, tt.args.tm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ExpireXX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ExpireXX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ExpireGT(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		tm  time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ExpireGT(tt.args.ctx, tt.args.key, tt.args.tm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ExpireGT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ExpireGT() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ExpireLT(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		tm  time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ExpireLT(tt.args.ctx, tt.args.key, tt.args.tm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ExpireLT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ExpireLT() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Keys(t *testing.T) {
	type args struct {
		ctx     context.Context
		pattern string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Keys(tt.args.ctx, tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Keys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.Keys() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Migrate(t *testing.T) {
	type args struct {
		ctx     context.Context
		host    string
		port    string
		key     string
		db      int
		timeout time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Migrate(tt.args.ctx, tt.args.host, tt.args.port, tt.args.key, tt.args.db, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Migrate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Migrate() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Move(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		db  int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Move(tt.args.ctx, tt.args.key, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Move() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Move() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ObjectRefCount(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ObjectRefCount(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ObjectRefCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ObjectRefCount() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ObjectEncoding(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ObjectEncoding(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ObjectEncoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ObjectEncoding() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ObjectIdleTime(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal time.Duration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ObjectIdleTime(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ObjectIdleTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ObjectIdleTime() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Persist(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Persist(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Persist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Persist() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PExpire(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PExpire(tt.args.ctx, tt.args.key, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PExpire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.PExpire() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PExpireAt(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		tm  time.Time
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PExpireAt(tt.args.ctx, tt.args.key, tt.args.tm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PExpireAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.PExpireAt() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PTTL(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal time.Duration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PTTL(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PTTL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.PTTL() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_RandomKey(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.RandomKey(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RandomKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.RandomKey() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Rename(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		newkey string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Rename(tt.args.ctx, tt.args.key, tt.args.newkey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Rename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Rename() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_RenameNX(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		newkey string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.RenameNX(tt.args.ctx, tt.args.key, tt.args.newkey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RenameNX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.RenameNX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Restore(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		ttl   time.Duration
		value string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Restore(tt.args.ctx, tt.args.key, tt.args.ttl, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Restore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Restore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_RestoreReplace(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		ttl   time.Duration
		value string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.RestoreReplace(tt.args.ctx, tt.args.key, tt.args.ttl, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RestoreReplace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.RestoreReplace() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Sort(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		sort *redis.Sort
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Sort(tt.args.ctx, tt.args.key, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Sort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.Sort() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SortRO(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		sort *redis.Sort
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SortRO(tt.args.ctx, tt.args.key, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SortRO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SortRO() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SortStore(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		store string
		sort  *redis.Sort
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SortStore(tt.args.ctx, tt.args.key, tt.args.store, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SortStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SortStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SortInterfaces(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		sort *redis.Sort
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SortInterfaces(tt.args.ctx, tt.args.key, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SortInterfaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SortInterfaces() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Touch(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Touch(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Touch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Touch() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_TTL(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal time.Duration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.TTL(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TTL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.TTL() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Type(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Type(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Type() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Type() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Append(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Append(tt.args.ctx, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Append() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Append() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Decr(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Decr(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Decr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Decr() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_DecrBy(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.DecrBy(tt.args.ctx, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DecrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.DecrBy() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Get() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GetRange(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		end   int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GetRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GetRange() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GetSet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GetSet(tt.args.ctx, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GetSet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GetEx(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		ts  time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GetEx(tt.args.ctx, tt.args.key, tt.args.ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GetEx() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GetDel(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GetDel(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetDel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GetDel() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Incr(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Incr(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Incr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Incr() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_IncrBy(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.IncrBy(tt.args.ctx, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.IncrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.IncrBy() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_IncrByFloat(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value float64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.IncrByFloat(tt.args.ctx, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.IncrByFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.IncrByFloat() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_MGet(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.MGet(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.MGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.MGet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_MSet(t *testing.T) {
	type args struct {
		ctx    context.Context
		values []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.MSet(tt.args.ctx, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.MSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.MSet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_MSetNX(t *testing.T) {
	type args struct {
		ctx    context.Context
		values []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.MSetNX(tt.args.ctx, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.MSetNX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.MSetNX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Set(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		value      any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Set() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SetArgs(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value any
		a     redis.SetArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SetArgs(tt.args.ctx, tt.args.key, tt.args.value, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SetArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SetArgs() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SetEx(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		value      any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SetEx(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SetEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SetEx() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SetNX(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		value      any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SetNX(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SetNX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SetNX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SetXX(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		value      any
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SetXX(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SetXX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SetXX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SetRange(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		offset int64
		value  string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SetRange(tt.args.ctx, tt.args.key, tt.args.offset, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SetRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SetRange() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_StrLen(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.StrLen(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.StrLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.StrLen() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Copy(t *testing.T) {
	type args struct {
		ctx       context.Context
		sourceKey string
		destKey   string
		db        int
		replace   bool
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Copy(tt.args.ctx, tt.args.sourceKey, tt.args.destKey, tt.args.db, tt.args.replace)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Copy() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GetBit(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		offset int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GetBit(tt.args.ctx, tt.args.key, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetBit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GetBit() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SetBit(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		offset int64
		value  int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SetBit(tt.args.ctx, tt.args.key, tt.args.offset, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SetBit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SetBit() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BitCount(t *testing.T) {
	type args struct {
		ctx      context.Context
		key      string
		bitCount *redis.BitCount
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BitCount(tt.args.ctx, tt.args.key, tt.args.bitCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BitCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BitCount() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BitOpAnd(t *testing.T) {
	type args struct {
		ctx     context.Context
		destKey string
		keys    []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BitOpAnd(tt.args.ctx, tt.args.destKey, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BitOpAnd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BitOpAnd() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BitOpOr(t *testing.T) {
	type args struct {
		ctx     context.Context
		destKey string
		keys    []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BitOpOr(tt.args.ctx, tt.args.destKey, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BitOpOr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BitOpOr() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BitOpXor(t *testing.T) {
	type args struct {
		ctx     context.Context
		destKey string
		keys    []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BitOpXor(tt.args.ctx, tt.args.destKey, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BitOpXor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BitOpXor() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BitOpNot(t *testing.T) {
	type args struct {
		ctx     context.Context
		destKey string
		key     string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BitOpNot(tt.args.ctx, tt.args.destKey, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BitOpNot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BitOpNot() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BitPos(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		bit int64
		pos []int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BitPos(tt.args.ctx, tt.args.key, tt.args.bit, tt.args.pos...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BitPos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BitPos() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BitField(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		args []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BitField(tt.args.ctx, tt.args.key, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BitField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.BitField() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Scan(t *testing.T) {
	type args struct {
		ctx      context.Context
		cursorIn uint64
		match    string
		count    int64
	}
	tests := []struct {
		name       string
		r          *Client
		args       args
		wantVal    []string
		wantCursor uint64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotCursor, err := tt.r.Scan(tt.args.ctx, tt.args.cursorIn, tt.args.match, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.Scan() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotCursor != tt.wantCursor {
				t.Errorf("Client.Scan() gotCursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestClient_ScanType(t *testing.T) {
	type args struct {
		ctx      context.Context
		cursorIn uint64
		match    string
		count    int64
		keyType  string
	}
	tests := []struct {
		name       string
		r          *Client
		args       args
		wantVal    []string
		wantCursor uint64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotCursor, err := tt.r.ScanType(tt.args.ctx, tt.args.cursorIn, tt.args.match, tt.args.count, tt.args.keyType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ScanType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ScanType() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotCursor != tt.wantCursor {
				t.Errorf("Client.ScanType() gotCursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestClient_SScan(t *testing.T) {
	type args struct {
		ctx      context.Context
		key      string
		cursorIn uint64
		match    string
		count    int64
	}
	tests := []struct {
		name       string
		r          *Client
		args       args
		wantVal    []string
		wantCursor uint64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotCursor, err := tt.r.SScan(tt.args.ctx, tt.args.key, tt.args.cursorIn, tt.args.match, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SScan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SScan() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotCursor != tt.wantCursor {
				t.Errorf("Client.SScan() gotCursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestClient_HScan(t *testing.T) {
	type args struct {
		ctx      context.Context
		key      string
		cursorIn uint64
		match    string
		count    int64
	}
	tests := []struct {
		name       string
		r          *Client
		args       args
		wantVal    []string
		wantCursor uint64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotCursor, err := tt.r.HScan(tt.args.ctx, tt.args.key, tt.args.cursorIn, tt.args.match, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HScan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.HScan() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotCursor != tt.wantCursor {
				t.Errorf("Client.HScan() gotCursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestClient_ZScan(t *testing.T) {
	type args struct {
		ctx      context.Context
		key      string
		cursorIn uint64
		match    string
		count    int64
	}
	tests := []struct {
		name       string
		r          *Client
		args       args
		wantVal    []string
		wantCursor uint64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotCursor, err := tt.r.ZScan(tt.args.ctx, tt.args.key, tt.args.cursorIn, tt.args.match, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZScan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZScan() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotCursor != tt.wantCursor {
				t.Errorf("Client.ZScan() gotCursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestClient_HDel(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		fields []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HDel(tt.args.ctx, tt.args.key, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HDel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.HDel() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HExists(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		field string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HExists(tt.args.ctx, tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.HExists() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HGet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		field string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HGet(tt.args.ctx, tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.HGet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HGetAll(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HGetAll(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HGetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.HGetAll() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HIncrBy(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		field string
		incr  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HIncrBy(tt.args.ctx, tt.args.key, tt.args.field, tt.args.incr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HIncrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.HIncrBy() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HIncrByFloat(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		field string
		incr  float64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HIncrByFloat(tt.args.ctx, tt.args.key, tt.args.field, tt.args.incr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HIncrByFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.HIncrByFloat() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HKeys(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HKeys(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.HKeys() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HLen(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HLen(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.HLen() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HMGet(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		fields []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HMGet(tt.args.ctx, tt.args.key, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HMGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.HMGet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HSet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HSet(tt.args.ctx, tt.args.key, tt.args.value...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.HSet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HMSet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HMSet(tt.args.ctx, tt.args.key, tt.args.value...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HMSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.HMSet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HSetNX(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		field string
		value any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HSetNX(tt.args.ctx, tt.args.key, tt.args.field, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HSetNX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.HSetNX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HVals(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HVals(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HVals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.HVals() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HRandField(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HRandField(tt.args.ctx, tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HRandField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.HRandField() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_HRandFieldWithValues(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.KeyValue
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.HRandFieldWithValues(tt.args.ctx, tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HRandFieldWithValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.HRandFieldWithValues() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BLPop(t *testing.T) {
	type args struct {
		ctx     context.Context
		timeout time.Duration
		keys    []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BLPop(tt.args.ctx, tt.args.timeout, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BLPop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.BLPop() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BRPop(t *testing.T) {
	type args struct {
		ctx     context.Context
		timeout time.Duration
		keys    []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BRPop(tt.args.ctx, tt.args.timeout, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BRPop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.BRPop() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BRPopLPush(t *testing.T) {
	type args struct {
		ctx         context.Context
		source      string
		destination string
		timeout     time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BRPopLPush(tt.args.ctx, tt.args.source, tt.args.destination, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BRPopLPush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BRPopLPush() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LIndex(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		index int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LIndex(tt.args.ctx, tt.args.key, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LIndex() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LInsert(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		op    string
		pivot any
		value any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LInsert(tt.args.ctx, tt.args.key, tt.args.op, tt.args.pivot, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LInsert() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LInsertBefore(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		pivot any
		value any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LInsertBefore(tt.args.ctx, tt.args.key, tt.args.pivot, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LInsertBefore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LInsertBefore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LInsertAfter(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		pivot any
		value any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LInsertAfter(tt.args.ctx, tt.args.key, tt.args.pivot, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LInsertAfter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LInsertAfter() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LLen(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LLen(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LLen() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LPop(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LPop(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LPop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LPop() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LPopCount(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LPopCount(tt.args.ctx, tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LPopCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.LPopCount() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LPos(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value string
		args  redis.LPosArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LPos(tt.args.ctx, tt.args.key, tt.args.value, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LPos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LPos() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LPosCount(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value string
		count int64
		args  redis.LPosArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LPosCount(tt.args.ctx, tt.args.key, tt.args.value, tt.args.count, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LPosCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.LPosCount() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LPush(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		values []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LPush(tt.args.ctx, tt.args.key, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LPush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LPush() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LPushX(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LPushX(tt.args.ctx, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LPushX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LPushX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LRange(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.LRange() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LRem(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int64
		value any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LRem(tt.args.ctx, tt.args.key, tt.args.count, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LRem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LRem() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LSet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		index int64
		value any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LSet(tt.args.ctx, tt.args.key, tt.args.index, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LSet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LTrim(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LTrim(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LTrim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LTrim() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_RPop(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.RPop(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RPop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.RPop() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_RPopCount(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.RPopCount(tt.args.ctx, tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RPopCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.RPopCount() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_RPopLPush(t *testing.T) {
	type args struct {
		ctx         context.Context
		source      string
		destination string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.RPopLPush(tt.args.ctx, tt.args.source, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RPopLPush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.RPopLPush() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_RPush(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		values []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.RPush(tt.args.ctx, tt.args.key, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RPush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.RPush() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_RPushX(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.RPushX(tt.args.ctx, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RPushX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.RPushX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LMove(t *testing.T) {
	type args struct {
		ctx         context.Context
		source      string
		destination string
		srcpos      string
		destpos     string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LMove(tt.args.ctx, tt.args.source, tt.args.destination, tt.args.srcpos, tt.args.destpos)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LMove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LMove() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BLMove(t *testing.T) {
	type args struct {
		ctx         context.Context
		source      string
		destination string
		srcpos      string
		destpos     string
		ts          time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BLMove(tt.args.ctx, tt.args.source, tt.args.destination, tt.args.srcpos, tt.args.destpos, tt.args.ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BLMove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BLMove() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SAdd(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SAdd(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SAdd() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SCard(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SCard(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SCard() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SDiff(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SDiff(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SDiff() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SDiffStore(t *testing.T) {
	type args struct {
		ctx         context.Context
		destination string
		keys        []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SDiffStore(tt.args.ctx, tt.args.destination, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SDiffStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SDiffStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SInter(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SInter(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SInter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SInter() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SInterCard(t *testing.T) {
	type args struct {
		ctx   context.Context
		limit int64
		keys  []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SInterCard(tt.args.ctx, tt.args.limit, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SInterCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SInterCard() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SInterStore(t *testing.T) {
	type args struct {
		ctx         context.Context
		destination string
		keys        []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SInterStore(tt.args.ctx, tt.args.destination, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SInterStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SInterStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SIsMember(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		member any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SIsMember(tt.args.ctx, tt.args.key, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SIsMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SIsMember() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SMembers(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SMembers(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SMembers() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SMembersMap(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal map[string]struct{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SMembersMap(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SMembersMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SMembersMap() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SMove(t *testing.T) {
	type args struct {
		ctx         context.Context
		source      string
		destination string
		member      any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SMove(tt.args.ctx, tt.args.source, tt.args.destination, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SMove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SMove() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SPop(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SPop(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SPop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SPop() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SPopN(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SPopN(tt.args.ctx, tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SPopN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SPopN() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SRandMember(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SRandMember(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SRandMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SRandMember() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SRandMemberN(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SRandMemberN(tt.args.ctx, tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SRandMemberN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SRandMemberN() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SRem(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SRem(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SRem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SRem() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SUnion(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SUnion(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SUnion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SUnion() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SUnionStore(t *testing.T) {
	type args struct {
		ctx         context.Context
		destination string
		keys        []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SUnionStore(tt.args.ctx, tt.args.destination, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SUnionStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SUnionStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XAdd(t *testing.T) {
	type args struct {
		ctx context.Context
		a   *redis.XAddArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XAdd(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XAdd() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XDel(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		ids    []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XDel(tt.args.ctx, tt.args.stream, tt.args.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XDel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XDel() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XLen(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XLen(tt.args.ctx, tt.args.stream)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XLen() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XRange(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		start  string
		stop   string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XRange(tt.args.ctx, tt.args.stream, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XRange() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XRangeN(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		start  string
		stop   string
		count  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XRangeN(tt.args.ctx, tt.args.stream, tt.args.start, tt.args.stop, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XRangeN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XRangeN() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XRevRange(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		start  string
		stop   string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XRevRange(tt.args.ctx, tt.args.stream, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XRevRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XRevRange() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XRevRangeN(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		start  string
		stop   string
		count  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XRevRangeN(tt.args.ctx, tt.args.stream, tt.args.start, tt.args.stop, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XRevRangeN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XRevRangeN() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XRead(t *testing.T) {
	type args struct {
		ctx context.Context
		a   *redis.XReadArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XStream
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XRead(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XRead() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XReadStreams(t *testing.T) {
	type args struct {
		ctx     context.Context
		streams []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XStream
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XReadStreams(tt.args.ctx, tt.args.streams...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XReadStreams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XReadStreams() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XGroupCreate(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		group  string
		start  string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XGroupCreate(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XGroupCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XGroupCreate() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XGroupCreateMkStream(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		group  string
		start  string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XGroupCreateMkStream(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XGroupCreateMkStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XGroupCreateMkStream() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XGroupSetID(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		group  string
		start  string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XGroupSetID(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XGroupSetID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XGroupSetID() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XGroupDestroy(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		group  string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XGroupDestroy(tt.args.ctx, tt.args.stream, tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XGroupDestroy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XGroupDestroy() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XGroupDelConsumer(t *testing.T) {
	type args struct {
		ctx      context.Context
		stream   string
		group    string
		consumer string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XGroupDelConsumer(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.consumer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XGroupDelConsumer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XGroupDelConsumer() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XReadGroup(t *testing.T) {
	type args struct {
		ctx context.Context
		a   *redis.XReadGroupArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XStream
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XReadGroup(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XReadGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XReadGroup() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XAck(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		group  string
		ids    []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XAck(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XAck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XAck() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XPending(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream string
		group  string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal *redis.XPending
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XPending(tt.args.ctx, tt.args.stream, tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XPending() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XPending() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XPendingExt(t *testing.T) {
	type args struct {
		ctx context.Context
		a   *redis.XPendingExtArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XPendingExt
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XPendingExt(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XPendingExt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XPendingExt() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XClaim(t *testing.T) {
	type args struct {
		ctx context.Context
		a   *redis.XClaimArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XClaim(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XClaim() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XClaimJustID(t *testing.T) {
	type args struct {
		ctx context.Context
		a   *redis.XClaimArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XClaimJustID(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XClaimJustID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XClaimJustID() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XAutoClaim(t *testing.T) {
	type args struct {
		ctx context.Context
		a   *redis.XAutoClaimArgs
	}
	tests := []struct {
		name      string
		r         *Client
		args      args
		wantVal   []redis.XMessage
		wantStart string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotStart, err := tt.r.XAutoClaim(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XAutoClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XAutoClaim() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotStart != tt.wantStart {
				t.Errorf("Client.XAutoClaim() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
		})
	}
}

func TestClient_XAutoClaimJustID(t *testing.T) {
	type args struct {
		ctx context.Context
		a   *redis.XAutoClaimArgs
	}
	tests := []struct {
		name      string
		r         *Client
		args      args
		wantVal   []string
		wantStart string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotStart, err := tt.r.XAutoClaimJustID(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XAutoClaimJustID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XAutoClaimJustID() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotStart != tt.wantStart {
				t.Errorf("Client.XAutoClaimJustID() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
		})
	}
}

func TestClient_XTrimMaxLen(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		maxLen int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XTrimMaxLen(tt.args.ctx, tt.args.key, tt.args.maxLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XTrimMaxLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XTrimMaxLen() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XTrimMaxLenApprox(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		maxLen int64
		limit  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XTrimMaxLenApprox(tt.args.ctx, tt.args.key, tt.args.maxLen, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XTrimMaxLenApprox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XTrimMaxLenApprox() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XTrimMinID(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		minID string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XTrimMinID(tt.args.ctx, tt.args.key, tt.args.minID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XTrimMinID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XTrimMinID() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XTrimMinIDApprox(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		minID string
		limit int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XTrimMinIDApprox(tt.args.ctx, tt.args.key, tt.args.minID, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XTrimMinIDApprox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.XTrimMinIDApprox() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XInfoGroups(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XInfoGroup
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XInfoGroups(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XInfoGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XInfoGroups() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XInfoStream(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal *redis.XInfoStream
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XInfoStream(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XInfoStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XInfoStream() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XInfoStreamFull(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal *redis.XInfoStreamFull
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XInfoStreamFull(tt.args.ctx, tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XInfoStreamFull() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XInfoStreamFull() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_XInfoConsumers(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		group string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.XInfoConsumer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.XInfoConsumers(tt.args.ctx, tt.args.key, tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.XInfoConsumers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.XInfoConsumers() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BZPopMax(t *testing.T) {
	type args struct {
		ctx     context.Context
		timeout time.Duration
		keys    []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal *redis.ZWithKey
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BZPopMax(tt.args.ctx, tt.args.timeout, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BZPopMax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.BZPopMax() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BZPopMin(t *testing.T) {
	type args struct {
		ctx     context.Context
		timeout time.Duration
		keys    []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal *redis.ZWithKey
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BZPopMin(tt.args.ctx, tt.args.timeout, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BZPopMin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.BZPopMin() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZAdd(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []redis.Z
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZAdd(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZAdd() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZAddNX(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []redis.Z
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZAddNX(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZAddNX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZAddNX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZAddXX(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []redis.Z
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZAddXX(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZAddXX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZAddXX() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZAddArgs(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		args redis.ZAddArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZAddArgs(tt.args.ctx, tt.args.key, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZAddArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZAddArgs() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZAddArgsIncr(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		args redis.ZAddArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZAddArgsIncr(tt.args.ctx, tt.args.key, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZAddArgsIncr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZAddArgsIncr() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZCard(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZCard(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZCard() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZCount(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		min string
		max string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZCount(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZCount() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZLexCount(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		min string
		max string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZLexCount(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZLexCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZLexCount() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZIncrBy(t *testing.T) {
	type args struct {
		ctx       context.Context
		key       string
		increment float64
		member    string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZIncrBy(tt.args.ctx, tt.args.key, tt.args.increment, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZIncrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZIncrBy() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZInter(t *testing.T) {
	type args struct {
		ctx   context.Context
		store *redis.ZStore
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZInter(tt.args.ctx, tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZInter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZInter() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZInterWithScores(t *testing.T) {
	type args struct {
		ctx   context.Context
		store *redis.ZStore
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZInterWithScores(tt.args.ctx, tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZInterWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZInterWithScores() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZInterCard(t *testing.T) {
	type args struct {
		ctx   context.Context
		limit int64
		keys  []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZInterCard(tt.args.ctx, tt.args.limit, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZInterCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZInterCard() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZInterStore(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		store *redis.ZStore
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZInterStore(tt.args.ctx, tt.args.key, tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZInterStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZInterStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZMScore(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZMScore(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZMScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZMScore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZPopMax(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count []int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZPopMax(tt.args.ctx, tt.args.key, tt.args.count...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZPopMax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZPopMax() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZPopMin(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count []int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZPopMin(tt.args.ctx, tt.args.key, tt.args.count...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZPopMin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZPopMin() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRange(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRange() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRangeWithScores(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRangeWithScores(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRangeWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRangeWithScores() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRangeByScore(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		opt *redis.ZRangeBy
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRangeByScore(tt.args.ctx, tt.args.key, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRangeByScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRangeByScore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRangeByLex(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		opt *redis.ZRangeBy
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRangeByLex(tt.args.ctx, tt.args.key, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRangeByLex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRangeByLex() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRangeByScoreWithScores(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		opt *redis.ZRangeBy
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRangeByScoreWithScores(tt.args.ctx, tt.args.key, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRangeByScoreWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRangeByScoreWithScores() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRangeArgs(t *testing.T) {
	type args struct {
		ctx context.Context
		z   redis.ZRangeArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRangeArgs(tt.args.ctx, tt.args.z)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRangeArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRangeArgs() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRangeArgsWithScores(t *testing.T) {
	type args struct {
		ctx context.Context
		z   redis.ZRangeArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRangeArgsWithScores(tt.args.ctx, tt.args.z)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRangeArgsWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRangeArgsWithScores() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRangeStore(t *testing.T) {
	type args struct {
		ctx context.Context
		dst string
		z   redis.ZRangeArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRangeStore(tt.args.ctx, tt.args.dst, tt.args.z)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRangeStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZRangeStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRank(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		member string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRank(tt.args.ctx, tt.args.key, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRank() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZRank() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRem(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRem(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZRem() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRemRangeByRank(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRemRangeByRank(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRemRangeByRank() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZRemRangeByRank() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRemRangeByScore(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		min string
		max string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRemRangeByScore(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRemRangeByScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZRemRangeByScore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRemRangeByLex(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		min string
		max string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRemRangeByLex(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRemRangeByLex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZRemRangeByLex() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRevRange(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRevRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRevRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRevRange() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRevRangeWithScores(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRevRangeWithScores(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRevRangeWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRevRangeWithScores() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRevRangeByScore(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		opt *redis.ZRangeBy
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRevRangeByScore(tt.args.ctx, tt.args.key, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRevRangeByScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRevRangeByScore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRevRangeByLex(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		opt *redis.ZRangeBy
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRevRangeByLex(tt.args.ctx, tt.args.key, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRevRangeByLex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRevRangeByLex() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRevRangeByScoreWithScores(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		opt *redis.ZRangeBy
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRevRangeByScoreWithScores(tt.args.ctx, tt.args.key, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRevRangeByScoreWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRevRangeByScoreWithScores() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRevRank(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		member string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRevRank(tt.args.ctx, tt.args.key, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRevRank() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZRevRank() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZScore(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		member string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZScore(tt.args.ctx, tt.args.key, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZScore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZUnionStore(t *testing.T) {
	type args struct {
		ctx   context.Context
		dest  string
		store *redis.ZStore
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZUnionStore(tt.args.ctx, tt.args.dest, tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZUnionStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZUnionStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRandMember(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRandMember(tt.args.ctx, tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRandMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRandMember() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZRandMemberWithScores(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZRandMemberWithScores(tt.args.ctx, tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZRandMemberWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZRandMemberWithScores() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZUnion(t *testing.T) {
	type args struct {
		ctx   context.Context
		store redis.ZStore
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZUnion(tt.args.ctx, tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZUnion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZUnion() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZUnionWithScores(t *testing.T) {
	type args struct {
		ctx   context.Context
		store redis.ZStore
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZUnionWithScores(tt.args.ctx, tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZUnionWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZUnionWithScores() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZDiff(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZDiff(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZDiff() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZDiffWithScores(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.Z
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZDiffWithScores(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZDiffWithScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ZDiffWithScores() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ZDiffStore(t *testing.T) {
	type args struct {
		ctx         context.Context
		destination string
		keys        []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ZDiffStore(tt.args.ctx, tt.args.destination, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ZDiffStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ZDiffStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PFAdd(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		els []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PFAdd(tt.args.ctx, tt.args.key, tt.args.els...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PFAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.PFAdd() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PFCount(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PFCount(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PFCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.PFCount() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PFMerge(t *testing.T) {
	type args struct {
		ctx  context.Context
		dest string
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PFMerge(tt.args.ctx, tt.args.dest, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PFMerge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.PFMerge() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ConfigGet(t *testing.T) {
	type args struct {
		ctx       context.Context
		parameter string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ConfigGet(tt.args.ctx, tt.args.parameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ConfigGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ConfigGet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ConfigResetStat(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ConfigResetStat(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ConfigResetStat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ConfigResetStat() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ConfigSet(t *testing.T) {
	type args struct {
		ctx       context.Context
		parameter string
		value     string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ConfigSet(tt.args.ctx, tt.args.parameter, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ConfigSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ConfigSet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ConfigRewrite(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ConfigRewrite(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ConfigRewrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ConfigRewrite() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BgRewriteAOF(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BgRewriteAOF(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BgRewriteAOF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BgRewriteAOF() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_BgSave(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.BgSave(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BgSave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.BgSave() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClientKill(t *testing.T) {
	type args struct {
		ctx    context.Context
		ipPort string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClientKill(tt.args.ctx, tt.args.ipPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClientKill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClientKill() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClientKillByFilter(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClientKillByFilter(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClientKillByFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClientKillByFilter() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClientList(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClientList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClientList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClientList() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClientPause(t *testing.T) {
	type args struct {
		ctx context.Context
		dur time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClientPause(tt.args.ctx, tt.args.dur)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClientPause() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClientPause() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClientUnpause(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClientUnpause(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClientUnpause() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClientUnpause() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClientID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClientID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClientID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClientID() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClientUnblock(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClientUnblock(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClientUnblock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClientUnblock() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClientUnblockWithError(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClientUnblockWithError(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClientUnblockWithError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClientUnblockWithError() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_DBSize(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.DBSize(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DBSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.DBSize() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_FlushAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.FlushAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.FlushAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.FlushAll() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_FlushAllAsync(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.FlushAllAsync(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.FlushAllAsync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.FlushAllAsync() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_FlushDB(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.FlushDB(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.FlushDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.FlushDB() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_FlushDBAsync(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.FlushDBAsync(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.FlushDBAsync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.FlushDBAsync() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Info(t *testing.T) {
	type args struct {
		ctx     context.Context
		section []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Info(tt.args.ctx, tt.args.section...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Info() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_LastSave(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.LastSave(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LastSave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.LastSave() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Save(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Save(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Save() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Shutdown(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Shutdown(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Shutdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Shutdown() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ShutdownSave(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ShutdownSave(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ShutdownSave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ShutdownSave() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ShutdownNoSave(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ShutdownNoSave(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ShutdownNoSave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ShutdownNoSave() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SlaveOf(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
		port string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SlaveOf(tt.args.ctx, tt.args.host, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SlaveOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SlaveOf() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SlowLogGet(t *testing.T) {
	type args struct {
		ctx context.Context
		num int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.SlowLog
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SlowLogGet(tt.args.ctx, tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SlowLogGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.SlowLogGet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Time(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Time(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Time() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.Time() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_DebugObject(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.DebugObject(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DebugObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.DebugObject() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ReadOnly(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ReadOnly(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ReadOnly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ReadOnly() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ReadWrite(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ReadWrite(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ReadWrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ReadWrite() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_MemoryUsage(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		samples []int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.MemoryUsage(tt.args.ctx, tt.args.key, tt.args.samples...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.MemoryUsage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.MemoryUsage() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Eval(t *testing.T) {
	type args struct {
		ctx    context.Context
		script string
		keys   []string
		args   []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Eval(tt.args.ctx, tt.args.script, tt.args.keys, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.Eval() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_EvalSha(t *testing.T) {
	type args struct {
		ctx  context.Context
		sha1 string
		keys []string
		args []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.EvalSha(tt.args.ctx, tt.args.sha1, tt.args.keys, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.EvalSha() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.EvalSha() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_EvalRO(t *testing.T) {
	type args struct {
		ctx    context.Context
		script string
		keys   []string
		args   []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.EvalRO(tt.args.ctx, tt.args.script, tt.args.keys, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.EvalRO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.EvalRO() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_EvalShaRO(t *testing.T) {
	type args struct {
		ctx  context.Context
		sha1 string
		keys []string
		args []any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.EvalShaRO(tt.args.ctx, tt.args.sha1, tt.args.keys, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.EvalShaRO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.EvalShaRO() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ScriptExists(t *testing.T) {
	type args struct {
		ctx    context.Context
		hashes []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ScriptExists(tt.args.ctx, tt.args.hashes...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ScriptExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ScriptExists() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ScriptFlush(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ScriptFlush(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ScriptFlush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ScriptFlush() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ScriptKill(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ScriptKill(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ScriptKill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ScriptKill() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ScriptLoad(t *testing.T) {
	type args struct {
		ctx    context.Context
		script string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ScriptLoad(tt.args.ctx, tt.args.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ScriptLoad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ScriptLoad() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_Publish(t *testing.T) {
	type args struct {
		ctx     context.Context
		channel string
		message any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.Publish(tt.args.ctx, tt.args.channel, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.Publish() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_SPublish(t *testing.T) {
	type args struct {
		ctx     context.Context
		channel string
		message any
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.SPublish(tt.args.ctx, tt.args.channel, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SPublish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.SPublish() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PubSubChannels(t *testing.T) {
	type args struct {
		ctx     context.Context
		pattern string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PubSubChannels(tt.args.ctx, tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PubSubChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.PubSubChannels() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PubSubNumSub(t *testing.T) {
	type args struct {
		ctx      context.Context
		channels []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal map[string]int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PubSubNumSub(tt.args.ctx, tt.args.channels...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PubSubNumSub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.PubSubNumSub() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PubSubNumPat(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PubSubNumPat(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PubSubNumPat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.PubSubNumPat() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PubSubShardChannels(t *testing.T) {
	type args struct {
		ctx     context.Context
		pattern string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PubSubShardChannels(tt.args.ctx, tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PubSubShardChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.PubSubShardChannels() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_PubSubShardNumSub(t *testing.T) {
	type args struct {
		ctx      context.Context
		channels []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal map[string]int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.PubSubShardNumSub(tt.args.ctx, tt.args.channels...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PubSubShardNumSub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.PubSubShardNumSub() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterSlots(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.ClusterSlot
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterSlots(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterSlots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ClusterSlots() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterNodes(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterNodes(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterNodes() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterMeet(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
		port string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterMeet(tt.args.ctx, tt.args.host, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterMeet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterMeet() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterForget(t *testing.T) {
	type args struct {
		ctx    context.Context
		nodeID string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterForget(tt.args.ctx, tt.args.nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterForget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterForget() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterReplicate(t *testing.T) {
	type args struct {
		ctx    context.Context
		nodeID string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterReplicate(tt.args.ctx, tt.args.nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterReplicate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterReplicate() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterResetSoft(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterResetSoft(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterResetSoft() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterResetSoft() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterResetHard(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterResetHard(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterResetHard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterResetHard() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterInfo(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterInfo(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterInfo() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterKeySlot(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterKeySlot(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterKeySlot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterKeySlot() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterGetKeysInSlot(t *testing.T) {
	type args struct {
		ctx   context.Context
		slot  int
		count int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterGetKeysInSlot(tt.args.ctx, tt.args.slot, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterGetKeysInSlot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ClusterGetKeysInSlot() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterCountFailureReports(t *testing.T) {
	type args struct {
		ctx    context.Context
		nodeID string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterCountFailureReports(tt.args.ctx, tt.args.nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterCountFailureReports() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterCountFailureReports() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterCountKeysInSlot(t *testing.T) {
	type args struct {
		ctx  context.Context
		slot int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterCountKeysInSlot(tt.args.ctx, tt.args.slot)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterCountKeysInSlot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterCountKeysInSlot() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterDelSlots(t *testing.T) {
	type args struct {
		ctx   context.Context
		slots []int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterDelSlots(tt.args.ctx, tt.args.slots...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterDelSlots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterDelSlots() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterDelSlotsRange(t *testing.T) {
	type args struct {
		ctx context.Context
		min int
		max int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterDelSlotsRange(tt.args.ctx, tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterDelSlotsRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterDelSlotsRange() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterSaveConfig(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterSaveConfig(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterSaveConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterSaveConfig() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterSlaves(t *testing.T) {
	type args struct {
		ctx    context.Context
		nodeID string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterSlaves(tt.args.ctx, tt.args.nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterSlaves() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.ClusterSlaves() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterFailover(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterFailover(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterFailover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterFailover() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterAddSlots(t *testing.T) {
	type args struct {
		ctx   context.Context
		slots []int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterAddSlots(tt.args.ctx, tt.args.slots...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterAddSlots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterAddSlots() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_ClusterAddSlotsRange(t *testing.T) {
	type args struct {
		ctx context.Context
		min int
		max int
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.ClusterAddSlotsRange(tt.args.ctx, tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ClusterAddSlotsRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.ClusterAddSlotsRange() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoAdd(t *testing.T) {
	type args struct {
		ctx         context.Context
		key         string
		geoLocation []*redis.GeoLocation
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoAdd(tt.args.ctx, tt.args.key, tt.args.geoLocation...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GeoAdd() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoPos(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []*redis.GeoPos
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoPos(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoPos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.GeoPos() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoRadius(t *testing.T) {
	type args struct {
		ctx       context.Context
		key       string
		longitude float64
		latitude  float64
		query     *redis.GeoRadiusQuery
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.GeoLocation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoRadius(tt.args.ctx, tt.args.key, tt.args.longitude, tt.args.latitude, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoRadius() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.GeoRadius() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoRadiusStore(t *testing.T) {
	type args struct {
		ctx       context.Context
		key       string
		longitude float64
		latitude  float64
		query     *redis.GeoRadiusQuery
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoRadiusStore(tt.args.ctx, tt.args.key, tt.args.longitude, tt.args.latitude, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoRadiusStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GeoRadiusStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoRadiusByMember(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		member string
		query  *redis.GeoRadiusQuery
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.GeoLocation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoRadiusByMember(tt.args.ctx, tt.args.key, tt.args.member, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoRadiusByMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.GeoRadiusByMember() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoRadiusByMemberStore(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		member string
		query  *redis.GeoRadiusQuery
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoRadiusByMemberStore(tt.args.ctx, tt.args.key, tt.args.member, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoRadiusByMemberStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GeoRadiusByMemberStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoSearch(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		q   *redis.GeoSearchQuery
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoSearch(tt.args.ctx, tt.args.key, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.GeoSearch() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoSearchLocation(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		q   *redis.GeoSearchLocationQuery
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []redis.GeoLocation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoSearchLocation(tt.args.ctx, tt.args.key, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoSearchLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.GeoSearchLocation() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoSearchStore(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		store string
		q     *redis.GeoSearchStoreQuery
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoSearchStore(tt.args.ctx, tt.args.key, tt.args.store, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoSearchStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GeoSearchStore() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoDist(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		member1 string
		member2 string
		unit    string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoDist(tt.args.ctx, tt.args.key, tt.args.member1, tt.args.member2, tt.args.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoDist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Client.GeoDist() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestClient_GeoHash(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.r.GeoHash(tt.args.ctx, tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GeoHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Client.GeoHash() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
