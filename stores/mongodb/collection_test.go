package mongodb

import (
	"context"
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/core/resource"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_newCollection(t *testing.T) {
	type args struct {
		collection    *mongo.Collection
		brk           resource.Breaker
		DisableMetric bool
		DisableTrace  bool
	}
	tests := []struct {
		name string
		args args
		want *DecoratedCollection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCollection(tt.args.collection, tt.args.brk, tt.args.DisableMetric, tt.args.DisableTrace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedCollection_Where(t *testing.T) {
	type args struct {
		m bson.D
	}
	tests := []struct {
		name string
		c    *DecoratedCollection
		args args
		want *DecoratedCollection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Where(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecoratedCollection.Where() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedCollection_Limit(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name string
		c    *DecoratedCollection
		args args
		want *DecoratedCollection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Limit(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecoratedCollection.Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedCollection_Skip(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name string
		c    *DecoratedCollection
		args args
		want *DecoratedCollection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Skip(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecoratedCollection.Skip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedCollection_Sort(t *testing.T) {
	type args struct {
		sorts bson.D
	}
	tests := []struct {
		name string
		c    *DecoratedCollection
		args args
		want *DecoratedCollection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Sort(tt.args.sorts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecoratedCollection.Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedCollection_Fields(t *testing.T) {
	type args struct {
		fields bson.M
	}
	tests := []struct {
		name string
		c    *DecoratedCollection
		args args
		want *DecoratedCollection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Fields(tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecoratedCollection.Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedCollection_reset(t *testing.T) {
	tests := []struct {
		name string
		c    *DecoratedCollection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.reset()
		})
	}
}

func TestDecoratedCollection_getCtx(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		c    *DecoratedCollection
		args args
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.getCtx(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecoratedCollection.getCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedCollection_CreateIndex(t *testing.T) {
	type args struct {
		ctx context.Context
		key bson.D
		op  *options.IndexOptions
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantRes string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.c.CreateIndex(tt.args.ctx, tt.args.key, tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("DecoratedCollection.CreateIndex() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestDecoratedCollection_ListIndexes(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts *options.ListIndexesOptions
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantVal []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.c.ListIndexes(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.ListIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("DecoratedCollection.ListIndexes() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestDecoratedCollection_DropIndex(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		opts *options.DropIndexesOptions
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.DropIndex(tt.args.ctx, tt.args.name, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.DropIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecoratedCollection_InsertOne(t *testing.T) {
	type args struct {
		ctx      context.Context
		document any
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantVal *mongo.InsertOneResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.c.InsertOne(tt.args.ctx, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.InsertOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("DecoratedCollection.InsertOne() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestDecoratedCollection_InsertMany(t *testing.T) {
	type args struct {
		ctx       context.Context
		documents any
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantVal *mongo.InsertManyResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.c.InsertMany(tt.args.ctx, tt.args.documents)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.InsertMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("DecoratedCollection.InsertMany() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestDecoratedCollection_Aggregate(t *testing.T) {
	type args struct {
		ctx      context.Context
		pipeline any
		result   any
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Aggregate(tt.args.ctx, tt.args.pipeline, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.Aggregate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecoratedCollection_UpdateOrInsert(t *testing.T) {
	type args struct {
		ctx       context.Context
		documents []any
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantVal *mongo.UpdateResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.c.UpdateOrInsert(tt.args.ctx, tt.args.documents)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.UpdateOrInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("DecoratedCollection.UpdateOrInsert() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestDecoratedCollection_UpdateOne(t *testing.T) {
	type args struct {
		ctx      context.Context
		document any
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantVal *mongo.UpdateResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.c.UpdateOne(tt.args.ctx, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.UpdateOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("DecoratedCollection.UpdateOne() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestDecoratedCollection_UpdateOneRaw(t *testing.T) {
	type args struct {
		ctx      context.Context
		document any
		opt      []*options.UpdateOptions
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantVal *mongo.UpdateResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.c.UpdateOneRaw(tt.args.ctx, tt.args.document, tt.args.opt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.UpdateOneRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("DecoratedCollection.UpdateOneRaw() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestDecoratedCollection_UpdateMany(t *testing.T) {
	type args struct {
		ctx      context.Context
		document any
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantVal *mongo.UpdateResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.c.UpdateMany(tt.args.ctx, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.UpdateMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("DecoratedCollection.UpdateMany() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestDecoratedCollection_FindOne(t *testing.T) {
	type args struct {
		ctx      context.Context
		document any
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.FindOne(tt.args.ctx, tt.args.document); (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.FindOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecoratedCollection_FindMany(t *testing.T) {
	type args struct {
		ctx       context.Context
		documents any
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.FindMany(tt.args.ctx, tt.args.documents); (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.FindMany() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecoratedCollection_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		c         *DecoratedCollection
		args      args
		wantCount int64
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := tt.c.Delete(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("DecoratedCollection.Delete() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestDecoratedCollection_Drop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *DecoratedCollection
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Drop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.Drop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecoratedCollection_Count(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		c          *DecoratedCollection
		args       args
		wantResult int64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.c.Count(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecoratedCollection.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("DecoratedCollection.Count() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBeforeCreate(t *testing.T) {
	type args struct {
		document any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeforeCreate(tt.args.document); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BeforeCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeforeUpdate(t *testing.T) {
	type args struct {
		document any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeforeUpdate(tt.args.document); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BeforeUpdate() = %v, want %v", got, tt.want)
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
