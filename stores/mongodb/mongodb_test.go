package mongodb

import (
	"context"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewClient(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want *MongoDB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_reset(t *testing.T) {
	tests := []struct {
		name       string
		collection *collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.collection.reset()
		})
	}
}

func TestMongoDB_Collection(t *testing.T) {
	type args struct {
		table string
	}
	tests := []struct {
		name   string
		client *MongoDB
		args   args
		want   *collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.client.Collection(tt.args.table); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MongoDB.Collection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_Where(t *testing.T) {
	type args struct {
		m bson.D
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.collection.Where(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.Where() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_Limit(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.collection.Limit(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_Skip(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.collection.Skip(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.Skip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_Sort(t *testing.T) {
	type args struct {
		sorts bson.D
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.collection.Sort(tt.args.sorts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_Fields(t *testing.T) {
	type args struct {
		fields bson.M
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.collection.Fields(tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_CreateIndex(t *testing.T) {
	type args struct {
		ctx context.Context
		key bson.D
		op  *options.IndexOptions
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		wantRes    string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.collection.CreateIndex(tt.args.ctx, tt.args.key, tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("collection.CreateIndex() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_collection_ListIndexes(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts *options.ListIndexesOptions
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       interface{}
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.collection.ListIndexes(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.ListIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.ListIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_DropIndex(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		opts *options.DropIndexesOptions
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.collection.DropIndex(tt.args.ctx, tt.args.name, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("collection.DropIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_collection_InsertOne(t *testing.T) {
	type args struct {
		ctx      context.Context
		document interface{}
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *mongo.InsertOneResult
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.collection.InsertOne(tt.args.ctx, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.InsertOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.InsertOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_InsertMany(t *testing.T) {
	type args struct {
		ctx       context.Context
		documents interface{}
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *mongo.InsertManyResult
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.collection.InsertMany(tt.args.ctx, tt.args.documents)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.InsertMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.InsertMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_Aggregate(t *testing.T) {
	type args struct {
		ctx      context.Context
		pipeline interface{}
		result   interface{}
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.collection.Aggregate(tt.args.ctx, tt.args.pipeline, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("collection.Aggregate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_collection_UpdateOrInsert(t *testing.T) {
	type args struct {
		ctx       context.Context
		documents []interface{}
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *mongo.UpdateResult
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.collection.UpdateOrInsert(tt.args.ctx, tt.args.documents)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.UpdateOrInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.UpdateOrInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_UpdateOne(t *testing.T) {
	type args struct {
		ctx      context.Context
		document interface{}
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *mongo.UpdateResult
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.collection.UpdateOne(tt.args.ctx, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.UpdateOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.UpdateOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_UpdateOneRaw(t *testing.T) {
	type args struct {
		ctx      context.Context
		document interface{}
		opt      []*options.UpdateOptions
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *mongo.UpdateResult
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.collection.UpdateOneRaw(tt.args.ctx, tt.args.document, tt.args.opt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.UpdateOneRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.UpdateOneRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_UpdateMany(t *testing.T) {
	type args struct {
		ctx      context.Context
		document interface{}
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		want       *mongo.UpdateResult
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.collection.UpdateMany(tt.args.ctx, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.UpdateMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collection.UpdateMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collection_FindOne(t *testing.T) {
	type args struct {
		ctx      context.Context
		document interface{}
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.collection.FindOne(tt.args.ctx, tt.args.document); (err != nil) != tt.wantErr {
				t.Errorf("collection.FindOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_collection_FindMany(t *testing.T) {
	type args struct {
		ctx       context.Context
		documents interface{}
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.collection.FindMany(tt.args.ctx, tt.args.documents); (err != nil) != tt.wantErr {
				t.Errorf("collection.FindMany() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_collection_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		wantCount  int64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := tt.collection.Delete(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("collection.Delete() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func Test_collection_Drop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.collection.Drop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("collection.Drop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_collection_Count(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		collection *collection
		args       args
		wantResult int64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.collection.Count(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("collection.Count() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBeforeCreate(t *testing.T) {
	type args struct {
		document interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
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
		document interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
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
