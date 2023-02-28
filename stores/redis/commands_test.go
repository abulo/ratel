package redis

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestClient_Pipeline(t *testing.T) {
	tests := []struct {
		name    string
		r       *Client
		wantRes redis.Pipeliner
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Pipeline(); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Pipeline() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes []redis.Cmder
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.r.Pipelined(tt.args.ctx, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Pipelined() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Pipelined() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes []redis.Cmder
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.r.TxPipelined(tt.args.ctx, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TxPipelined() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.TxPipelined() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_TxPipeline(t *testing.T) {
	tests := []struct {
		name    string
		r       *Client
		wantRes redis.Pipeliner
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.TxPipeline(); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.TxPipeline() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.CommandsInfoCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Command(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Command() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClientGetName(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClientGetName() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_Echo(t *testing.T) {
	type args struct {
		ctx     context.Context
		message interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Echo(tt.args.ctx, tt.args.message); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Echo() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Ping(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Ping() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Quit(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Quit() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Del(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Del() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Unlink(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Unlink() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Dump(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Dump() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Exists(tt.args.ctx, tt.args.key...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Exists() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Expire(tt.args.ctx, tt.args.key, tt.args.expiration); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Expire() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ExpireAt(tt.args.ctx, tt.args.key, tt.args.tm); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ExpireAt() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ExpireNX(tt.args.ctx, tt.args.key, tt.args.tm); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ExpireNX() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ExpireXX(tt.args.ctx, tt.args.key, tt.args.tm); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ExpireXX() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ExpireGT(tt.args.ctx, tt.args.key, tt.args.tm); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ExpireGT() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ExpireLT(tt.args.ctx, tt.args.key, tt.args.tm); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ExpireLT() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Keys(tt.args.ctx, tt.args.pattern); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Keys() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Migrate(tt.args.ctx, tt.args.host, tt.args.port, tt.args.key, tt.args.db, tt.args.timeout); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Migrate() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Move(tt.args.ctx, tt.args.key, tt.args.db); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Move() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ObjectRefCount(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ObjectRefCount() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ObjectEncoding(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ObjectEncoding() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.DurationCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ObjectIdleTime(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ObjectIdleTime() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Persist(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Persist() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PExpire(tt.args.ctx, tt.args.key, tt.args.expiration); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PExpire() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PExpireAt(tt.args.ctx, tt.args.key, tt.args.tm); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PExpireAt() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.DurationCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PTTL(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PTTL() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.RandomKey(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.RandomKey() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Rename(tt.args.ctx, tt.args.key, tt.args.newkey); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Rename() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.RenameNX(tt.args.ctx, tt.args.key, tt.args.newkey); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.RenameNX() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Restore(tt.args.ctx, tt.args.key, tt.args.ttl, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Restore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.RestoreReplace(tt.args.ctx, tt.args.key, tt.args.ttl, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.RestoreReplace() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Sort(tt.args.ctx, tt.args.key, tt.args.sort); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Sort() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SortRO(tt.args.ctx, tt.args.key, tt.args.sort); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SortRO() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SortStore(tt.args.ctx, tt.args.key, tt.args.store, tt.args.sort); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SortStore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.SliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SortInterfaces(tt.args.ctx, tt.args.key, tt.args.sort); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SortInterfaces() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Touch(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Touch() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.DurationCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.TTL(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.TTL() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Type(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Type() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Append(tt.args.ctx, tt.args.key, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Append() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Decr(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Decr() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.DecrBy(tt.args.ctx, tt.args.key, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.DecrBy() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Get(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Get() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GetRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.end); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GetRange() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_GetSet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GetSet(tt.args.ctx, tt.args.key, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GetSet() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GetEx(tt.args.ctx, tt.args.key, tt.args.ts); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GetEx() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GetDel(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GetDel() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Incr(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Incr() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.IncrBy(tt.args.ctx, tt.args.key, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.IncrBy() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.FloatCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.IncrByFloat(tt.args.ctx, tt.args.key, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.IncrByFloat() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.SliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.MGet(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.MGet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_MSet(t *testing.T) {
	type args struct {
		ctx    context.Context
		values []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.MSet(tt.args.ctx, tt.args.values...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.MSet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_MSetNX(t *testing.T) {
	type args struct {
		ctx    context.Context
		values []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.MSetNX(tt.args.ctx, tt.args.values...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.MSetNX() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_Set(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		value      interface{}
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Set() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SetArgs(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value interface{}
		a     redis.SetArgs
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SetArgs(tt.args.ctx, tt.args.key, tt.args.value, tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SetArgs() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SetEx(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		value      interface{}
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SetEx(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SetEx() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SetNX(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		value      interface{}
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SetNX(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SetNX() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SetXX(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		value      interface{}
		expiration time.Duration
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SetXX(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SetXX() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SetRange(tt.args.ctx, tt.args.key, tt.args.offset, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SetRange() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.StrLen(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.StrLen() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Copy(tt.args.ctx, tt.args.sourceKey, tt.args.destKey, tt.args.db, tt.args.replace); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Copy() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GetBit(tt.args.ctx, tt.args.key, tt.args.offset); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GetBit() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SetBit(tt.args.ctx, tt.args.key, tt.args.offset, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SetBit() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BitCount(tt.args.ctx, tt.args.key, tt.args.bitCount); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BitCount() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BitOpAnd(tt.args.ctx, tt.args.destKey, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BitOpAnd() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BitOpOr(tt.args.ctx, tt.args.destKey, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BitOpOr() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BitOpXor(tt.args.ctx, tt.args.destKey, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BitOpXor() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BitOpNot(tt.args.ctx, tt.args.destKey, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BitOpNot() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BitPos(tt.args.ctx, tt.args.key, tt.args.bit, tt.args.pos...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BitPos() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_BitField(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		args []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BitField(tt.args.ctx, tt.args.key, tt.args.args...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BitField() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_Scan(t *testing.T) {
	type args struct {
		ctx    context.Context
		cursor uint64
		match  string
		count  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.ScanCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Scan(tt.args.ctx, tt.args.cursor, tt.args.match, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Scan() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_ScanType(t *testing.T) {
	type args struct {
		ctx     context.Context
		cursor  uint64
		match   string
		count   int64
		keyType string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.ScanCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ScanType(tt.args.ctx, tt.args.cursor, tt.args.match, tt.args.count, tt.args.keyType); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ScanType() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SScan(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		cursor uint64
		match  string
		count  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.ScanCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SScan(tt.args.ctx, tt.args.key, tt.args.cursor, tt.args.match, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SScan() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_HScan(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		cursor uint64
		match  string
		count  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.ScanCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HScan(tt.args.ctx, tt.args.key, tt.args.cursor, tt.args.match, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HScan() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_ZScan(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		cursor uint64
		match  string
		count  int64
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.ScanCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZScan(tt.args.ctx, tt.args.key, tt.args.cursor, tt.args.match, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZScan() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HDel(tt.args.ctx, tt.args.key, tt.args.fields...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HDel() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HExists(tt.args.ctx, tt.args.key, tt.args.field); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HExists() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HGet(tt.args.ctx, tt.args.key, tt.args.field); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HGet() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.MapStringStringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HGetAll(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HGetAll() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HIncrBy(tt.args.ctx, tt.args.key, tt.args.field, tt.args.incr); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HIncrBy() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.FloatCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HIncrByFloat(tt.args.ctx, tt.args.key, tt.args.field, tt.args.incr); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HIncrByFloat() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HKeys(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HKeys() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HLen(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HLen() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.SliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HMGet(tt.args.ctx, tt.args.key, tt.args.fields...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HMGet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_HSet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HSet(tt.args.ctx, tt.args.key, tt.args.value...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HSet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_HMSet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HMSet(tt.args.ctx, tt.args.key, tt.args.value...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HMSet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_HSetNX(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		field string
		value interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HSetNX(tt.args.ctx, tt.args.key, tt.args.field, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HSetNX() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HVals(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HVals() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HRandField(tt.args.ctx, tt.args.key, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HRandField() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.KeyValueSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.HRandFieldWithValues(tt.args.ctx, tt.args.key, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.HRandFieldWithValues() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BLPop(tt.args.ctx, tt.args.timeout, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BLPop() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BRPop(tt.args.ctx, tt.args.timeout, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BRPop() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BRPopLPush(tt.args.ctx, tt.args.source, tt.args.destination, tt.args.timeout); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BRPopLPush() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LIndex(tt.args.ctx, tt.args.key, tt.args.index); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LIndex() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_LInsert(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		op    string
		pivot interface{}
		value interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LInsert(tt.args.ctx, tt.args.key, tt.args.op, tt.args.pivot, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LInsert() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_LInsertBefore(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		pivot interface{}
		value interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LInsertBefore(tt.args.ctx, tt.args.key, tt.args.pivot, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LInsertBefore() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_LInsertAfter(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		pivot interface{}
		value interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LInsertAfter(tt.args.ctx, tt.args.key, tt.args.pivot, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LInsertAfter() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LLen(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LLen() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LPop(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LPop() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LPopCount(tt.args.ctx, tt.args.key, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LPopCount() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LPos(tt.args.ctx, tt.args.key, tt.args.value, tt.args.args); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LPos() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LPosCount(tt.args.ctx, tt.args.key, tt.args.value, tt.args.count, tt.args.args); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LPosCount() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_LPush(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		values []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LPush(tt.args.ctx, tt.args.key, tt.args.values...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LPush() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_LPushX(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LPushX(tt.args.ctx, tt.args.key, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LPushX() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LRange() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_LRem(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int64
		value interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LRem(tt.args.ctx, tt.args.key, tt.args.count, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LRem() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_LSet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		index int64
		value interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LSet(tt.args.ctx, tt.args.key, tt.args.index, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LSet() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LTrim(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LTrim() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.RPop(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.RPop() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.RPopCount(tt.args.ctx, tt.args.key, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.RPopCount() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.RPopLPush(tt.args.ctx, tt.args.source, tt.args.destination); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.RPopLPush() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_RPush(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		values []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.RPush(tt.args.ctx, tt.args.key, tt.args.values...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.RPush() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_RPushX(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.RPushX(tt.args.ctx, tt.args.key, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.RPushX() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LMove(tt.args.ctx, tt.args.source, tt.args.destination, tt.args.srcpos, tt.args.destpos); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LMove() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BLMove(tt.args.ctx, tt.args.source, tt.args.destination, tt.args.srcpos, tt.args.destpos, tt.args.ts); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BLMove() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SAdd(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SAdd(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SAdd() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SCard(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SCard() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SDiff(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SDiff() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SDiffStore(tt.args.ctx, tt.args.destination, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SDiffStore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SInter(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SInter() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SInterCard(tt.args.ctx, tt.args.limit, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SInterCard() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SInterStore(tt.args.ctx, tt.args.destination, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SInterStore() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SIsMember(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		member interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SIsMember(tt.args.ctx, tt.args.key, tt.args.member); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SIsMember() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SMembers(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SMembers() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringStructMapCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SMembersMap(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SMembersMap() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SMove(t *testing.T) {
	type args struct {
		ctx         context.Context
		source      string
		destination string
		member      interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SMove(tt.args.ctx, tt.args.source, tt.args.destination, tt.args.member); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SMove() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SPop(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SPop() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SPopN(tt.args.ctx, tt.args.key, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SPopN() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SRandMember(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SRandMember() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SRandMemberN(tt.args.ctx, tt.args.key, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SRandMemberN() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SRem(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SRem(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SRem() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SUnion(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SUnion() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SUnionStore(tt.args.ctx, tt.args.destination, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SUnionStore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XAdd(tt.args.ctx, tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XAdd() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XDel(tt.args.ctx, tt.args.stream, tt.args.ids...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XDel() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XLen(tt.args.ctx, tt.args.stream); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XLen() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XMessageSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XRange(tt.args.ctx, tt.args.stream, tt.args.start, tt.args.stop); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XRange() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XMessageSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XRangeN(tt.args.ctx, tt.args.stream, tt.args.start, tt.args.stop, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XRangeN() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XMessageSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XRevRange(tt.args.ctx, tt.args.stream, tt.args.start, tt.args.stop); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XRevRange() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XMessageSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XRevRangeN(tt.args.ctx, tt.args.stream, tt.args.start, tt.args.stop, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XRevRangeN() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XStreamSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XRead(tt.args.ctx, tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XRead() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XStreamSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XReadStreams(tt.args.ctx, tt.args.streams...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XReadStreams() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XGroupCreate(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.start); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XGroupCreate() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XGroupCreateMkStream(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.start); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XGroupCreateMkStream() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XGroupSetID(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.start); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XGroupSetID() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XGroupDestroy(tt.args.ctx, tt.args.stream, tt.args.group); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XGroupDestroy() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XGroupDelConsumer(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.consumer); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XGroupDelConsumer() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XStreamSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XReadGroup(tt.args.ctx, tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XReadGroup() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XAck(tt.args.ctx, tt.args.stream, tt.args.group, tt.args.ids...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XAck() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XPendingCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XPending(tt.args.ctx, tt.args.stream, tt.args.group); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XPending() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XPendingExtCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XPendingExt(tt.args.ctx, tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XPendingExt() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XMessageSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XClaim(tt.args.ctx, tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XClaim() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XClaimJustID(tt.args.ctx, tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XClaimJustID() = %v, want %v", gotRes, tt.wantRes)
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
		name    string
		r       *Client
		args    args
		wantRes *redis.XAutoClaimCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XAutoClaim(tt.args.ctx, tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XAutoClaim() = %v, want %v", gotRes, tt.wantRes)
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
		name    string
		r       *Client
		args    args
		wantRes *redis.XAutoClaimJustIDCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XAutoClaimJustID(tt.args.ctx, tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XAutoClaimJustID() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XTrimMaxLen(tt.args.ctx, tt.args.key, tt.args.maxLen); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XTrimMaxLen() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XTrimMaxLenApprox(tt.args.ctx, tt.args.key, tt.args.maxLen, tt.args.limit); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XTrimMaxLenApprox() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XTrimMinID(tt.args.ctx, tt.args.key, tt.args.minID); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XTrimMinID() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XTrimMinIDApprox(tt.args.ctx, tt.args.key, tt.args.minID, tt.args.limit); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XTrimMinIDApprox() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XInfoGroupsCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XInfoGroups(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XInfoGroups() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XInfoStreamCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XInfoStream(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XInfoStream() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XInfoStreamFullCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XInfoStreamFull(tt.args.ctx, tt.args.key, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XInfoStreamFull() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.XInfoConsumersCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.XInfoConsumers(tt.args.ctx, tt.args.key, tt.args.group); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.XInfoConsumers() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZWithKeyCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BZPopMax(tt.args.ctx, tt.args.timeout, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BZPopMax() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZWithKeyCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BZPopMin(tt.args.ctx, tt.args.timeout, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BZPopMin() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZAdd(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZAdd() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZAddNX(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZAddNX() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZAddXX(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZAddXX() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZAddArgs(tt.args.ctx, tt.args.key, tt.args.args); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZAddArgs() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.FloatCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZAddArgsIncr(tt.args.ctx, tt.args.key, tt.args.args); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZAddArgsIncr() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZCard(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZCard() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZCount(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZCount() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZLexCount(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZLexCount() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.FloatCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZIncrBy(tt.args.ctx, tt.args.key, tt.args.increment, tt.args.member); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZIncrBy() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZInter(tt.args.ctx, tt.args.store); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZInter() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZInterWithScores(tt.args.ctx, tt.args.store); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZInterWithScores() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZInterCard(tt.args.ctx, tt.args.limit, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZInterCard() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZInterStore(tt.args.ctx, tt.args.key, tt.args.store); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZInterStore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.FloatSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZMScore(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZMScore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZPopMax(tt.args.ctx, tt.args.key, tt.args.count...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZPopMax() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZPopMin(tt.args.ctx, tt.args.key, tt.args.count...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZPopMin() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRange() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRangeWithScores(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRangeWithScores() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRangeByScore(tt.args.ctx, tt.args.key, tt.args.opt); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRangeByScore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRangeByLex(tt.args.ctx, tt.args.key, tt.args.opt); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRangeByLex() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRangeByScoreWithScores(tt.args.ctx, tt.args.key, tt.args.opt); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRangeByScoreWithScores() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRangeArgs(tt.args.ctx, tt.args.z); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRangeArgs() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRangeArgsWithScores(tt.args.ctx, tt.args.z); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRangeArgsWithScores() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRangeStore(tt.args.ctx, tt.args.dst, tt.args.z); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRangeStore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRank(tt.args.ctx, tt.args.key, tt.args.member); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRank() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_ZRem(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRem(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRem() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRemRangeByRank(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRemRangeByRank() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRemRangeByScore(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRemRangeByScore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRemRangeByLex(tt.args.ctx, tt.args.key, tt.args.min, tt.args.max); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRemRangeByLex() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRevRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRevRange() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRevRangeWithScores(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRevRangeWithScores() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRevRangeByScore(tt.args.ctx, tt.args.key, tt.args.opt); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRevRangeByScore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRevRangeByLex(tt.args.ctx, tt.args.key, tt.args.opt); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRevRangeByLex() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRevRangeByScoreWithScores(tt.args.ctx, tt.args.key, tt.args.opt); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRevRangeByScoreWithScores() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRevRank(tt.args.ctx, tt.args.key, tt.args.member); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRevRank() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.FloatCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZScore(tt.args.ctx, tt.args.key, tt.args.member); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZScore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZUnionStore(tt.args.ctx, tt.args.dest, tt.args.store); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZUnionStore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRandMember(tt.args.ctx, tt.args.key, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRandMember() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZRandMemberWithScores(tt.args.ctx, tt.args.key, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZRandMemberWithScores() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZUnion(tt.args.ctx, tt.args.store); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZUnion() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZUnionWithScores(tt.args.ctx, tt.args.store); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZUnionWithScores() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZDiff(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZDiff() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ZSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZDiffWithScores(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZDiffWithScores() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ZDiffStore(tt.args.ctx, tt.args.destination, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ZDiffStore() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_PFAdd(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		els []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PFAdd(tt.args.ctx, tt.args.key, tt.args.els...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PFAdd() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PFCount(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PFCount() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PFMerge(tt.args.ctx, tt.args.dest, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PFMerge() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.MapStringStringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ConfigGet(tt.args.ctx, tt.args.parameter); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ConfigGet() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ConfigResetStat(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ConfigResetStat() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ConfigSet(tt.args.ctx, tt.args.parameter, tt.args.value); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ConfigSet() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ConfigRewrite(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ConfigRewrite() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BgRewriteAOF(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BgRewriteAOF() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.BgSave(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.BgSave() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClientKill(tt.args.ctx, tt.args.ipPort); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClientKill() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClientKillByFilter(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClientKillByFilter() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClientList(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClientList() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClientPause(tt.args.ctx, tt.args.dur); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClientPause() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClientUnpause(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClientUnpause() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClientID(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClientID() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClientUnblock(tt.args.ctx, tt.args.id); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClientUnblock() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClientUnblockWithError(tt.args.ctx, tt.args.id); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClientUnblockWithError() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.DBSize(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.DBSize() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.FlushAll(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.FlushAll() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.FlushAllAsync(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.FlushAllAsync() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.FlushDB(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.FlushDB() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.FlushDBAsync(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.FlushDBAsync() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Info(tt.args.ctx, tt.args.section...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Info() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.LastSave(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.LastSave() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Save(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Save() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Shutdown(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Shutdown() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ShutdownSave(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ShutdownSave() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ShutdownNoSave(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ShutdownNoSave() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SlaveOf(tt.args.ctx, tt.args.host, tt.args.port); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SlaveOf() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.SlowLogCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SlowLogGet(tt.args.ctx, tt.args.num); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SlowLogGet() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.TimeCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Time(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Time() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.DebugObject(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.DebugObject() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ReadOnly(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ReadOnly() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ReadWrite(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ReadWrite() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.MemoryUsage(tt.args.ctx, tt.args.key, tt.args.samples...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.MemoryUsage() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_Eval(t *testing.T) {
	type args struct {
		ctx    context.Context
		script string
		keys   []string
		args   []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Eval(tt.args.ctx, tt.args.script, tt.args.keys, tt.args.args...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Eval() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_EvalSha(t *testing.T) {
	type args struct {
		ctx  context.Context
		sha1 string
		keys []string
		args []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.EvalSha(tt.args.ctx, tt.args.sha1, tt.args.keys, tt.args.args...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.EvalSha() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_EvalRO(t *testing.T) {
	type args struct {
		ctx    context.Context
		script string
		keys   []string
		args   []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.EvalRO(tt.args.ctx, tt.args.script, tt.args.keys, tt.args.args...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.EvalRO() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_EvalShaRO(t *testing.T) {
	type args struct {
		ctx  context.Context
		sha1 string
		keys []string
		args []interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.EvalShaRO(tt.args.ctx, tt.args.sha1, tt.args.keys, tt.args.args...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.EvalShaRO() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.BoolSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ScriptExists(tt.args.ctx, tt.args.hashes...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ScriptExists() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ScriptFlush(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ScriptFlush() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ScriptKill(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ScriptKill() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ScriptLoad(tt.args.ctx, tt.args.script); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ScriptLoad() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_Publish(t *testing.T) {
	type args struct {
		ctx     context.Context
		channel string
		message interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.Publish(tt.args.ctx, tt.args.channel, tt.args.message); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.Publish() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_SPublish(t *testing.T) {
	type args struct {
		ctx     context.Context
		channel string
		message interface{}
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.SPublish(tt.args.ctx, tt.args.channel, tt.args.message); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.SPublish() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PubSubChannels(tt.args.ctx, tt.args.pattern); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PubSubChannels() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.MapStringIntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PubSubNumSub(tt.args.ctx, tt.args.channels...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PubSubNumSub() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PubSubNumPat(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PubSubNumPat() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PubSubShardChannels(tt.args.ctx, tt.args.pattern); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PubSubShardChannels() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.MapStringIntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.PubSubShardNumSub(tt.args.ctx, tt.args.channels...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.PubSubShardNumSub() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.ClusterSlotsCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterSlots(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterSlots() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterNodes(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterNodes() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterMeet(tt.args.ctx, tt.args.host, tt.args.port); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterMeet() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterForget(tt.args.ctx, tt.args.nodeID); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterForget() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterReplicate(tt.args.ctx, tt.args.nodeID); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterReplicate() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterResetSoft(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterResetSoft() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterResetHard(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterResetHard() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterInfo(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterInfo() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterKeySlot(tt.args.ctx, tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterKeySlot() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterGetKeysInSlot(tt.args.ctx, tt.args.slot, tt.args.count); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterGetKeysInSlot() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterCountFailureReports(tt.args.ctx, tt.args.nodeID); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterCountFailureReports() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterCountKeysInSlot(tt.args.ctx, tt.args.slot); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterCountKeysInSlot() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterDelSlots(tt.args.ctx, tt.args.slots...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterDelSlots() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterDelSlotsRange(tt.args.ctx, tt.args.min, tt.args.max); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterDelSlotsRange() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterSaveConfig(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterSaveConfig() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterSlaves(tt.args.ctx, tt.args.nodeID); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterSlaves() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterFailover(tt.args.ctx); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterFailover() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterAddSlots(tt.args.ctx, tt.args.slots...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterAddSlots() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StatusCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.ClusterAddSlotsRange(tt.args.ctx, tt.args.min, tt.args.max); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.ClusterAddSlotsRange() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoAdd(tt.args.ctx, tt.args.key, tt.args.geoLocation...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoAdd() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.GeoPosCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoPos(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoPos() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.GeoLocationCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoRadius(tt.args.ctx, tt.args.key, tt.args.longitude, tt.args.latitude, tt.args.query); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoRadius() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoRadiusStore(tt.args.ctx, tt.args.key, tt.args.longitude, tt.args.latitude, tt.args.query); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoRadiusStore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.GeoLocationCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoRadiusByMember(tt.args.ctx, tt.args.key, tt.args.member, tt.args.query); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoRadiusByMember() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoRadiusByMemberStore(tt.args.ctx, tt.args.key, tt.args.member, tt.args.query); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoRadiusByMemberStore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoSearch(tt.args.ctx, tt.args.key, tt.args.q); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoSearch() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.GeoSearchLocationCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoSearchLocation(tt.args.ctx, tt.args.key, tt.args.q); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoSearchLocation() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoSearchStore(tt.args.ctx, tt.args.key, tt.args.store, tt.args.q); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoSearchStore() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.FloatCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoDist(tt.args.ctx, tt.args.key, tt.args.member1, tt.args.member2, tt.args.unit); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoDist() = %v, want %v", gotRes, tt.wantRes)
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
		wantRes *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GeoHash(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GeoHash() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
