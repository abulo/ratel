package mongodb

import (
	"context"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/abulo/ratel/core/logger"
	"github.com/abulo/ratel/core/metric"
	"github.com/abulo/ratel/core/trace"
	"github.com/abulo/ratel/util"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB 连接
type MongoDB struct {
	Client        *mongo.Client
	Name          string
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
}

// collection *mongo.Client
type collection struct {
	Database      *mongo.Database
	Table         *mongo.Collection
	filter        bson.D
	limit         int64
	skip          int64
	sort          bson.D
	fields        bson.M
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
}

// Config 配置
type Config struct {
	URI             string
	MaxConnIdleTime time.Duration
	MaxPoolSize     uint64
	MinPoolSize     uint64
	DisableMetric   bool // 关闭指标采集
	DisableTrace    bool // 关闭链路追踪
}

type index struct {
	Key  bson.D
	Name string
}

// NewClient New 数据库连接
func NewClient(config *Config) *MongoDB {
	//数据库连接
	mongoOptions := options.Client()
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
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Logger.Panic("MongoDB连接失败->", err)
		return nil
	}
	return &MongoDB{Client: client, Name: name, DisableMetric: config.DisableMetric, DisableTrace: config.DisableTrace}
}

func (collection *collection) reset() {
	collection.filter = nil
	collection.limit = 0
	collection.skip = 0
	collection.sort = nil
	collection.fields = nil
	collection.Table = nil
	collection.DisableMetric = true
	collection.DisableTrace = true
}

// Collection 得到一个mongo操作对象
func (client *MongoDB) Collection(table string) *collection {
	database := client.Client.Database(client.Name)
	return &collection{
		Database:      database,
		Table:         database.Collection(table),
		filter:        make(bson.D, 0),
		sort:          make(bson.D, 0),
		DisableMetric: client.DisableMetric,
		DisableTrace:  client.DisableTrace,
	}
}

// Where 条件查询, bson.M{"field": "value"}
func (collection *collection) Where(m bson.D) *collection {
	collection.filter = m
	return collection
}

// Limit 限制条数
func (collection *collection) Limit(n int64) *collection {
	collection.limit = n
	return collection
}

// Skip 跳过条数
func (collection *collection) Skip(n int64) *collection {
	collection.skip = n
	return collection
}

// Sort 排序 bson.M{"created_at":-1}
func (collection *collection) Sort(sorts bson.D) *collection {
	collection.sort = sorts
	return collection
}

// Fields 指定查询字段
func (collection *collection) Fields(fields bson.M) *collection {
	collection.fields = fields
	return collection
}

// CreateIndex CreateOneIndex 创建单个普通索引
func (collection *collection) CreateIndex(ctx context.Context, key bson.D, op *options.IndexOptions) (res string, err error) {
	start := time.Now()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
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
	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}
	return
}

// ListIndexes 获取所有所有
func (collection *collection) ListIndexes(ctx context.Context, opts *options.ListIndexesOptions) (interface{}, error) {
	start := time.Now()
	var results []string
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
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
	for cursor.Next(ctx) {
		var idx index
		if err := cursor.Decode(&idx); err == nil {
			results = append(results, idx.Name)
		}
	}
	collection.reset()

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}

	return results, nil
}

// DropIndex 删除索引
func (collection *collection) DropIndex(ctx context.Context, name string, opts *options.DropIndexesOptions) error {
	start := time.Now()
	indexView := collection.Table.Indexes()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
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

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}
	return nil
}

// InsertOne 写入单条数据
func (collection *collection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	start := time.Now()
	var data interface{}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	data = BeforeCreate(document)
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("method", "InsertOne")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("document", data))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	result, err := collection.Table.InsertOne(ctx, data)
	collection.reset()

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}

	return result, err
}

// InsertMany 写入多条数据
func (collection *collection) InsertMany(ctx context.Context, documents interface{}) (*mongo.InsertManyResult, error) {
	start := time.Now()
	var data []interface{}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	data = BeforeCreate(documents).([]interface{})
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("method", "InsertMany")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("documents", documents))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	result, err := collection.Table.InsertMany(ctx, data)
	collection.reset()

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}

	return result, err
}

// Aggregate ...
func (collection *collection) Aggregate(ctx context.Context, pipeline interface{}, result interface{}) (err error) {
	start := time.Now()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("method", "Aggregate")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("pipeline", pipeline))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
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

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}

	return
}

// UpdateOrInsert 存在更新,不存在写入, documents 里边的文档需要有 _id 的存在
func (collection *collection) UpdateOrInsert(ctx context.Context, documents []interface{}) (*mongo.UpdateResult, error) {
	start := time.Now()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("method", "UpdateOrInsert")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Object("documents", documents))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	var upsert = true
	result, err := collection.Table.UpdateMany(ctx, collection.filter, documents, &options.UpdateOptions{Upsert: &upsert})
	collection.reset()

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}
	return result, err
}

// UpdateOne ...
func (collection *collection) UpdateOne(ctx context.Context, document interface{}) (*mongo.UpdateResult, error) {
	start := time.Now()
	// var update bson.M
	update := bson.M{"$set": BeforeUpdate(document)}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("method", "UpdateOne")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Object("update", update))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}

	}
	result, err := collection.Table.UpdateOne(ctx, collection.filter, update)
	collection.reset()

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb").Observe(cost.Seconds())
	}
	return result, err
}

// UpdateOneRaw 原生update
func (collection *collection) UpdateOneRaw(ctx context.Context, document interface{}, opt ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	start := time.Now()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("method", "UpdateOneRaw")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Object("update", document))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	result, err := collection.Table.UpdateOne(ctx, collection.filter, document, opt...)
	collection.reset()

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}
	return result, err
}

// UpdateMany ...
func (collection *collection) UpdateMany(ctx context.Context, document interface{}) (*mongo.UpdateResult, error) {
	start := time.Now()
	// var update bson.M
	update := bson.M{"$set": BeforeUpdate(document)}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("method", "UpdateMany")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			span.LogFields(log.Object("update", update))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	result, err := collection.Table.UpdateMany(ctx, collection.filter, update)
	collection.reset()

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}
	return result, err
}

// FindOne 查询一条数据
func (collection *collection) FindOne(ctx context.Context, document interface{}) error {
	start := time.Now()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
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

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}

	return nil
}

// FindMany 查询多条数据
func (collection *collection) FindMany(ctx context.Context, documents interface{}) (err error) {
	start := time.Now()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
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

	defer func(ctx context.Context) {
		if err := result.Close(ctx); err != nil {
			logger.Logger.Error("Error closing result: ", err)
		}
	}(ctx)

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

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}
	return
}

// Delete 删除数据,并返回删除成功的数量
func (collection *collection) Delete(ctx context.Context) (count int64, err error) {
	start := time.Now()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := collection.Table.DeleteMany(ctx, collection.filter)
	if err != nil {
		collection.reset()
		return
	}
	count = result.DeletedCount
	collection.reset()

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}
	return
}

// Drop ...
func (collection *collection) Drop(ctx context.Context) error {
	start := time.Now()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("method", "Drop")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}

	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := collection.Table.Drop(ctx)

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}
	return err
}

// Count ...
func (collection *collection) Count(ctx context.Context) (result int64, err error) {
	start := time.Now()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !collection.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("mongodb", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("method", "Count")
			span.LogFields(log.String("database", collection.Database.Name()))
			span.LogFields(log.String("table", collection.Table.Name()))
			span.LogFields(log.Object("filter", collection.filter))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err = collection.Table.CountDocuments(ctx, collection.filter)
	if err != nil {
		collection.reset()
		return
	}
	collection.reset()

	if !collection.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("mongodb", collection.Database.Name(), collection.Table.Name(), "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("mongodb", collection.Database.Name(), collection.Table.Name()).Observe(cost.Seconds())
	}

	return
}

// BeforeCreate ...
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

// BeforeUpdate ...
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

// func isZero(value reflect.Value) bool {
// 	switch value.Kind() {
// 	case reflect.String:
// 		return value.Len() == 0
// 	case reflect.Bool:
// 		return !value.Bool()
// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		return value.Int() == 0
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
// 		return value.Uint() == 0
// 	case reflect.Float32, reflect.Float64:
// 		return value.Float() == 0
// 	case reflect.Interface, reflect.Ptr:
// 		return value.IsNil()
// 	}
// 	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
// }
