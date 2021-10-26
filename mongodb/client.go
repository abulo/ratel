package mongodb

import (
	"context"
	"net/url"
	"reflect"
	"time"

	"github.com/abulo/ratel/v1/logger"
	"github.com/abulo/ratel/v1/util"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoDB 连接
type MongoDB struct {
	Client *mongo.Client
	Name   string
}

//collection *mongo.Client
type collection struct {
	Database *mongo.Database
	Table    *mongo.Collection
	filter   bson.D
	limit    int64
	skip     int64
	sort     bson.D
	fields   bson.M
}

//Config 配置
type Config struct {
	URI string
	// AppName string
	// ConnectTimeout  time.Duration
	MaxConnIdleTime time.Duration
	MaxPoolSize     uint64
	MinPoolSize     uint64
}

//New 数据库连接
func New(config *Config) *MongoDB {
	//数据库连接
	mongoOptions := options.Client()
	// mongoOptions.SetAppName(config.AppName)
	// mongoOptions.SetConnectTimeout(config.ConnectTimeout)
	mongoOptions.SetMaxConnIdleTime(config.MaxConnIdleTime)
	mongoOptions.SetMaxPoolSize(config.MaxPoolSize)
	mongoOptions.SetMinPoolSize(config.MinPoolSize)
	client, err := mongo.NewClient(mongoOptions.ApplyURI(config.URI))
	if err != nil {
		logger.Logger.Panic(err)
		return nil
	}
	//解析URL
	u, err := url.Parse(config.URI)
	if err != nil {
		logger.Logger.Panic(err)
		return nil
	}
	if util.Empty(u.Path) || util.Empty(u.Path[1:]) {
		logger.Logger.Panic(errors.New("no database"))
		return nil
	}
	name := u.Path[1:]
	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Logger.Panic("MongoDB连接失败->", err)
		return nil
	}
	return &MongoDB{Client: client, Name: name}
}

func (collection *collection) reset() {
	collection.filter = nil
	collection.limit = 0
	collection.skip = 0
	collection.sort = nil
	collection.fields = nil
	collection.Table = nil
}

// Collection 得到一个mongo操作对象
func (client *MongoDB) Collection(table string) *collection {
	database := client.Client.Database(client.Name)
	return &collection{
		Database: database,
		Table:    database.Collection(table),
		filter:   make(bson.D, 0),
		sort:     make(bson.D, 0),
	}
}

// 条件查询, bson.M{"field": "value"}
func (collection *collection) Where(m bson.D) *collection {
	collection.filter = m
	return collection
}

// 限制条数
func (collection *collection) Limit(n int64) *collection {
	collection.limit = n
	return collection
}

// 跳过条数
func (collection *collection) Skip(n int64) *collection {
	collection.skip = n
	return collection
}

// 排序 bson.M{"created_at":-1}
func (collection *collection) Sort(sorts bson.D) *collection {
	collection.sort = sorts
	return collection
}

// 指定查询字段
func (collection *collection) Fields(fields bson.M) *collection {
	collection.fields = fields
	return collection
}

//CreateOneIndex 创建单个普通索引
func (collection *collection) CreateIndex(ctx context.Context, key bson.D, op *options.IndexOptions) (res string, err error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "CreateIndex")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("keys", key))
			span.LogFields(log.Object("options", op))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	indexView := collection.Table.Indexes()
	indexModel := mongo.IndexModel{Keys: key, Options: op}
	res, err = indexView.CreateOne(ctx, indexModel)
	return
}

//ListIndexes 获取所有所有
func (collection *collection) ListIndexes(ctx context.Context, opts *options.ListIndexesOptions) (interface{}, error) {
	var results interface{}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "ListIndexes")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("options", opts))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	indexView := collection.Table.Indexes()
	cursor, err := indexView.List(ctx, opts)
	if err != nil {
		collection.reset()
		return nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		collection.reset()
		return nil, err
	}
	collection.reset()
	return results, nil
}

//DropIndex 删除索引
func (collection *collection) DropIndex(ctx context.Context, name string, opts *options.DropIndexesOptions) error {
	indexView := collection.Table.Indexes()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "DropIndex")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.String("indexname", name))
			span.LogFields(log.Object("options", opts))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	_, err := indexView.DropOne(ctx, name, opts)
	if err != nil {
		collection.reset()
		return err
	}
	collection.reset()
	return nil
}

// 写入单条数据
func (collection *collection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	var data interface{}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	data = BeforeCreate(document)
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "InsertOne")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("document", data))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	result, err := collection.Table.InsertOne(ctx, data)
	collection.reset()
	return result, err
}

// 写入多条数据
func (collection *collection) InsertMany(ctx context.Context, documents interface{}) (*mongo.InsertManyResult, error) {
	var data []interface{}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	data = BeforeCreate(documents).([]interface{})
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "InsertMany")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("documents", documents))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	result, err := collection.Table.InsertMany(ctx, data)
	collection.reset()
	return result, err
}

func (collection *collection) Aggregate(ctx context.Context, pipeline interface{}, result interface{}) (err error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "Aggregate")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("pipeline", pipeline))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	cursor, err := collection.Table.Aggregate(ctx, pipeline)
	if err != nil {
		collection.reset()
		return
	}
	err = cursor.All(ctx, result)
	if err != nil {
		collection.reset()
		return
	}
	collection.reset()
	return
}

// 存在更新,不存在写入, documents 里边的文档需要有 _id 的存在
func (collection *collection) UpdateOrInsert(ctx context.Context, documents []interface{}) (*mongo.UpdateResult, error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "UpdateOrInsert")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Object("documents", documents))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	var upsert = true
	result, err := collection.Table.UpdateMany(ctx, collection.filter, documents, &options.UpdateOptions{Upsert: &upsert})
	collection.reset()
	return result, err
}

//
func (collection *collection) UpdateOne(ctx context.Context, document interface{}) (*mongo.UpdateResult, error) {
	var update bson.M
	update = bson.M{"$set": BeforeUpdate(document)}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "UpdateOne")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Object("update", update))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	result, err := collection.Table.UpdateOne(ctx, collection.filter, update)

	collection.reset()
	return result, err
}

//原生update
func (collection *collection) UpdateOneRaw(ctx context.Context, document interface{}, opt ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "UpdateOneRaw")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Object("update", document))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	result, err := collection.Table.UpdateOne(ctx, collection.filter, document, opt...)
	collection.reset()
	return result, err
}

//
func (collection *collection) UpdateMany(ctx context.Context, document interface{}) (*mongo.UpdateResult, error) {
	var update bson.M
	update = bson.M{"$set": BeforeUpdate(document)}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "UpdateMany")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Object("update", update))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	result, err := collection.Table.UpdateMany(ctx, collection.filter, update)

	collection.reset()
	return result, err
}

// 查询一条数据
func (collection *collection) FindOne(ctx context.Context, document interface{}) error {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "FindOne")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Int64("skip", collection.skip))
			span.LogFields(log.Object("sort", collection.sort))
			span.LogFields(log.Object("fields", collection.fields))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	result := collection.Table.FindOne(ctx, collection.filter, &options.FindOneOptions{
		Skip:       &collection.skip,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	err := result.Decode(document)
	if err != nil {
		collection.reset()
		return err
	}
	collection.reset()
	return nil
}

// 查询多条数据
func (collection *collection) FindMany(ctx context.Context, documents interface{}) (err error) {

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "FindMany")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Int64("skip", collection.skip))
			span.LogFields(log.Int64("limit", collection.limit))
			span.LogFields(log.Object("sort", collection.sort))
			span.LogFields(log.Object("fields", collection.fields))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	result, err := collection.Table.Find(ctx, collection.filter, &options.FindOptions{
		Skip:       &collection.skip,
		Limit:      &collection.limit,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	if err != nil {
		collection.reset()
		return
	}
	defer result.Close(ctx)
	val := reflect.ValueOf(documents)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		err = errors.New("result argument must be a slice address")
		collection.reset()
		return
	}

	slice := reflect.MakeSlice(val.Elem().Type(), 0, 0)
	itemTyp := val.Elem().Type().Elem()
	for result.Next(ctx) {
		item := reflect.New(itemTyp)
		err := result.Decode(item.Interface())
		if err != nil {
			err = errors.New("result argument must be a slice address")
			collection.reset()
			return err
		}

		slice = reflect.Append(slice, reflect.Indirect(item))
	}
	val.Elem().Set(slice)
	collection.reset()
	return
}

// 删除数据,并返回删除成功的数量
func (collection *collection) Delete(ctx context.Context) (count int64, err error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "Delete")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	if collection.filter == nil || len(collection.filter) == 0 {
		err = errors.New("you can't delete all documents, it's very dangerous")
		collection.reset()
		return
	}

	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	result, err := collection.Table.DeleteMany(ctx, collection.filter)
	if err != nil {
		collection.reset()
		return
	}
	count = result.DeletedCount
	collection.reset()
	return
}

func (collection *collection) Drop(ctx context.Context) error {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "Drop")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	err := collection.Table.Drop(ctx)
	return err
}

func (collection *collection) Count(ctx context.Context) (result int64, err error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "mongodb")
			span.SetTag("method", "Count")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	result, err = collection.Table.CountDocuments(ctx, collection.filter)
	if err != nil {
		collection.reset()
		return
	}
	collection.reset()
	return
}
func BeforeCreate(document interface{}) interface{} {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)

	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeCreate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]interface{}, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			sliceData[i] = BeforeCreate(val.Index(i).Interface()).(bson.M)
		}
		return sliceData

	case reflect.Struct:
		var data = make(bson.M)
		for i := 0; i < typ.NumField(); i++ {
			data[typ.Field(i).Tag.Get("bson")] = val.Field(i).Interface()
		}
		dataVal := reflect.ValueOf(data)
		if val.FieldByName("Id").Type() == reflect.TypeOf(primitive.ObjectID{}) {
			dataVal.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID()))
		}

		if val.FieldByName("Id").Interface() == "" {
			dataVal.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID().String()))
		}
		return dataVal.Interface()

	default:
		if val.Type() == reflect.TypeOf(bson.M{}) {
			if !val.MapIndex(reflect.ValueOf("_id")).IsValid() {
				val.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID()))
			}
		}
		return val.Interface()
	}
}

func BeforeUpdate(document interface{}) interface{} {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)

	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeUpdate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]interface{}, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			sliceData[i] = BeforeCreate(val.Index(i).Interface()).(bson.M)
		}
		return sliceData

	case reflect.Struct:
		var data = make(bson.M)
		for i := 0; i < typ.NumField(); i++ {
			_, ok := typ.Field(i).Tag.Lookup("over")
			if ok {
				continue
			}
			data[typ.Field(i).Tag.Get("bson")] = val.Field(i).Interface()
		}
		dataVal := reflect.ValueOf(data)
		return dataVal.Interface()

	default:
		// if val.Type() == reflect.TypeOf(bson.M{}) {
		// 	if !val.MapIndex(reflect.ValueOf("_id")).IsValid() {
		// 		val.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID()))
		// 	}
		// }
		return val.Interface()
	}
}
func isZero(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
