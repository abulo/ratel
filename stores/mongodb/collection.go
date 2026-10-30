package mongodb

import (
	"context"
	"reflect"
	"runtime"
	"time"

	"github.com/abulo/ratel/v3/core/call"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/metric"
	"github.com/abulo/ratel/v3/core/resource"
	globalTrace "github.com/abulo/ratel/v3/core/trace"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

type (
	// DecoratedCollection 在原生 mongo.Collection 之上提供链式条件构造、指标采集与
	// 链路追踪增强。链式条件在单次数据操作结束后自动清空，不会残留到下一次操作。
	//
	// 实例由 newCollection 内部创建，不对外构造；同一实例上的链式条件方法不保证并发安全，
	// 应在单次操作内顺序调用。
	DecoratedCollection struct {
		*mongo.Collection
		brk resource.Breaker
		// EnableMetric 是否采集 MongoDB 操作指标。
		EnableMetric bool
		// EnableTrace 是否挂载 MongoDB 链路追踪。
		EnableTrace bool
		filter      bson.D
		limit       int64
		skip        int64
		sort        bson.D
		fields      bson.M
	}

	// index 为索引文档的解码结构，仅在 ListIndexes 内部使用。
	index struct {
		Key  bson.D
		Name string
	}
)

// newCollection 包装原生集合并注入指标、追踪开关与空的条件状态。
func newCollection(collection *mongo.Collection, brk resource.Breaker, EnableMetric, EnableTrace bool) *DecoratedCollection {
	return &DecoratedCollection{
		Collection:   collection,
		brk:          brk,
		EnableMetric: EnableMetric,
		EnableTrace:  EnableTrace,
		filter:       make(bson.D, 0),
		sort:         make(bson.D, 0),
	}
}

// Where 设置查询匹配条件，覆盖此前累积的 filter，返回集合自身以便链式调用。
func (c *DecoratedCollection) Where(m bson.D) *DecoratedCollection {
	c.filter = m
	return c
}

// Limit 设置返回记录条数上限，返回集合自身以便链式调用。
func (c *DecoratedCollection) Limit(n int64) *DecoratedCollection {
	c.limit = n
	return c
}

// Skip 设置跳过的记录条数，返回集合自身以便链式调用。
func (c *DecoratedCollection) Skip(n int64) *DecoratedCollection {
	c.skip = n
	return c
}

// Sort 设置排序条件，返回集合自身以便链式调用。
func (c *DecoratedCollection) Sort(sorts bson.D) *DecoratedCollection {
	c.sort = sorts
	return c
}

// Fields 设置返回文档的投影字段，返回集合自身以便链式调用。
func (c *DecoratedCollection) Fields(fields bson.M) *DecoratedCollection {
	c.fields = fields
	return c
}

// reset 清空链式累积的条件状态，保证每次操作的条件互不残留。
func (c *DecoratedCollection) reset() {
	c.filter = nil
	c.limit = 0
	c.skip = 0
	c.sort = nil
	c.fields = nil
}

// getCtx 将 nil 或已取消的上下文替换为 context.TODO()，避免下游空上下文报错。
func getCtx(ctx context.Context) context.Context {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return ctx
}

// CreateIndex 在当前集合上创建单个索引，返回索引名称；key 为索引键，op 为索引选项，
// 可为 nil。操作按开关上报追踪与指标。
func (c *DecoratedCollection) CreateIndex(ctx context.Context, key bson.D, op *options.IndexOptions) (res string, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "CreateIndex", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		indexView := c.Collection.Indexes()
		indexModel := mongo.IndexModel{Keys: key, Options: op}
		res, err = indexView.CreateOne(ctx, indexModel)
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// ListIndexes 返回当前集合全部索引的名称列表；opts 为列表选项，可为 nil。
// 游标读取失败时清空条件并返回错误。
func (c *DecoratedCollection) ListIndexes(ctx context.Context, opts *options.ListIndexesOptions) (val []string, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "ListIndexes", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		indexView := c.Collection.Indexes()
		cursor, err := indexView.List(ctx, opts)
		if err != nil {
			c.reset()
			return err
		}
		for cursor.Next(ctx) {
			var idx index
			if err := cursor.Decode(&idx); err == nil {
				val = append(val, idx.Name)
			}
		}
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// DropIndex 按名称删除当前集合上的指定索引；opts 为删除选项，可为 nil。
func (c *DecoratedCollection) DropIndex(ctx context.Context, name string, opts *options.DropIndexesOptions) (err error) {
	start := time.Now()
	indexView := c.Collection.Indexes()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "DropIndex", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		_, err = indexView.DropOne(ctx, name, opts)
		if err != nil {
			c.reset()
			return err
		}
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// InsertOne 向当前集合写入单条文档；document 为结构体、bson.M 或原生文档，
// 写入前由 BeforeCreate 完成转换与 _id 生成，返回插入结果。
func (c *DecoratedCollection) InsertOne(ctx context.Context, document any) (val *mongo.InsertOneResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		data := BeforeCreate(document)
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "InsertOne", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		val, err = c.Collection.InsertOne(ctx, data)
		c.reset()

		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// InsertMany 批量写入文档；documents 须为结构体切片或 bson 文档切片，
// 写入前由 BeforeCreate 逐项转换并生成 _id。返回批量插入结果。
func (c *DecoratedCollection) InsertMany(ctx context.Context, documents any) (val *mongo.InsertManyResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		data := BeforeCreate(documents).([]any)
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "InsertMany", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		val, err = c.Collection.InsertMany(ctx, data)
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// Aggregate 执行聚合管道，结果解码到 result；pipeline 为聚合阶段数组，
// result 须为可写入的结果接收对象。游标错误会被清空条件后返回。
func (c *DecoratedCollection) Aggregate(ctx context.Context, pipeline any, result any) (err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "Aggregate", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		cursor, err := c.Collection.Aggregate(ctx, pipeline)
		if err != nil {
			c.reset()
			return err
		}
		err = cursor.All(ctx, result)
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return

}

// UpdateOrInsert 在 Where 设置的 filter 上执行批量更新，匹配不到文档时插入；
// documents 中的文档须携带 _id。返回更新结果。
func (c *DecoratedCollection) UpdateOrInsert(ctx context.Context, documents []any) (val *mongo.UpdateResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "UpdateOrInsert", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		var upsert = true
		val, err = c.Collection.UpdateMany(ctx, c.filter, documents, &options.UpdateOptions{Upsert: &upsert})
		c.reset()

		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// UpdateOne 在 Where 设置的 filter 上更新单个文档，更新字段由 BeforeUpdate 提取后包装为
// $set；返回更新结果。
func (c *DecoratedCollection) UpdateOne(ctx context.Context, document any) (val *mongo.UpdateResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		update := bson.M{"$set": BeforeUpdate(document)}
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "UpdateOne", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		val, err = c.Collection.UpdateOne(ctx, c.filter, update)
		c.reset()

		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb").Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// UpdateOneRaw 在 Where 设置的 filter 上更新单个文档，document 作为原生更新表达式透传，
// 不做 $set 包装；opt 为可选更新选项。返回更新结果。
func (c *DecoratedCollection) UpdateOneRaw(ctx context.Context, document any, opt ...*options.UpdateOptions) (val *mongo.UpdateResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "UpdateOneRaw", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		val, err = c.Collection.UpdateOne(ctx, c.filter, document, opt...)
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// UpdateMany 在 Where 设置的 filter 上批量更新文档，更新字段由 BeforeUpdate 提取后包装为
// $set；返回更新结果。
func (c *DecoratedCollection) UpdateMany(ctx context.Context, document any) (val *mongo.UpdateResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		update := bson.M{"$set": BeforeUpdate(document)}
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "UpdateMany", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		val, err = c.Collection.UpdateMany(ctx, c.filter, update)
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// FindOne 按累积的 filter、skip、sort、fields 查询单条文档并解码到 document；
// document 须为指向目标结构的指针。文档不存在时返回 mongo.ErrNoDocuments。
func (c *DecoratedCollection) FindOne(ctx context.Context, document any) (err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "FindOne", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		result := c.Collection.FindOne(ctx, c.filter, &options.FindOneOptions{
			Skip:       &c.skip,
			Sort:       c.sort,
			Projection: c.fields,
		})
		err = result.Decode(document)
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// FindMany 按累积的 filter、skip、limit、sort、fields 查询多条文档并解码到 documents；
// documents 必须为指向切片的指针，否则返回错误。游标在读取结束后关闭。
func (c *DecoratedCollection) FindMany(ctx context.Context, documents any) (err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "FindMany", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		result, err := c.Collection.Find(ctx, c.filter, &options.FindOptions{
			Skip:       &c.skip,
			Limit:      &c.limit,
			Sort:       c.sort,
			Projection: c.fields,
		})
		if err != nil {
			c.reset()
			return err
		}
		defer func(ctx context.Context) {
			if err = result.Close(ctx); err != nil {
				logger.Logger.Error("Error closing result: ", err)
				return
			}
		}(ctx)
		val := reflect.ValueOf(documents)
		if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Slice {
			err = errors.New("result argument must be a slice address")
			c.reset()
			return err
		}
		slice := reflect.MakeSlice(val.Elem().Type(), 0, 0)
		itemTyp := val.Elem().Type().Elem()
		for result.Next(ctx) {
			item := reflect.New(itemTyp)
			err = result.Decode(item.Interface())
			if err != nil {
				err = errors.New("result argument must be a slice address")
				c.reset()
				return err
			}
			slice = reflect.Append(slice, reflect.Indirect(item))
		}
		val.Elem().Set(slice)
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return nil
	}, acceptable)
	return
}

// Delete 按 Where 设置的 filter 删除匹配文档并返回删除数量；filter 为空时拒绝执行，
// 避免误删整集合。操作附带 5 秒超时。
func (c *DecoratedCollection) Delete(ctx context.Context) (count int64, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "Delete", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		if len(c.filter) == 0 {
			err = errors.New("you can't delete all documents, it's very dangerous")
			c.reset()
			return err
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		result, err := c.Collection.DeleteMany(ctx, c.filter)
		if err != nil {
			c.reset()
			return err
		}
		count = result.DeletedCount
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return nil
	}, acceptable)
	return
}

// Drop 删除当前集合本身，操作附带 5 秒超时；调用方须确认该集合可被整体删除。
func (c *DecoratedCollection) Drop(ctx context.Context) (err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "Drop", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = c.Collection.Drop(ctx)
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// Count 统计匹配 Where filter 的文档数量，操作附带 5 秒超时。
func (c *DecoratedCollection) Count(ctx context.Context) (result int64, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if c.EnableTrace {
			tracer := globalTrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("mongo"),
			}
			fn, file, line := call.Caller(6)
			attrs = append(attrs,
				semconv.CodeFunction(fn),
				semconv.CodeFilepath(file),
				semconv.CodeLineNumber(line),
			)
			md := metadata.New(nil)
			newCtx, span := tracer.Start(ctx, "Count", propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

			span.SetAttributes(
				semconv.DBMongoDBCollection(c.Database().Name()),
			)
			defer span.End()
			ctx = newCtx
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		result, err = c.Collection.CountDocuments(ctx, c.filter)
		c.reset()
		if c.EnableMetric {
			cost := time.Since(start)
			if err != nil {
				metric.LibHandleCounter.WithLabelValues("mongodb", c.Database().Name(), c.Name(), "ERR").Inc()
			} else {
				metric.LibHandleCounter.Inc("mongodb", c.Database().Name(), c.Name(), "OK")
			}
			metric.LibHandleHistogram.WithLabelValues("mongodb", c.Database().Name(), c.Name()).Observe(cost.Seconds())
		}
		return err
	}, acceptable)
	return
}

// BeforeCreate 在写入前将输入转换为 BSON 文档：指针自动解引用，数组/切片递归逐项转换，
// 结构体按 bson 标签转为 bson.M，并为 ObjectID 类型或空 Id 字段补充 _id。
// 输入类型与转换假设不匹配时会触发类型断言 panic，调用方须传入结构体、切片或 bson.M。
func BeforeCreate(document any) any {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)

	switch typ.Kind() {
	case reflect.Pointer:
		return BeforeCreate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]any, val.Len(), val.Cap())
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
		if val.FieldByName("Id").Type() == reflect.TypeFor[primitive.ObjectID]() {
			dataVal.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID()))
		}

		if val.FieldByName("Id").Interface() == "" {
			dataVal.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID().String()))
		}
		return dataVal.Interface()

	default:
		if val.Type() == reflect.TypeFor[bson.M]() {
			if !val.MapIndex(reflect.ValueOf("_id")).IsValid() {
				val.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID()))
			}
		}
		return val.Interface()
	}
}

// BeforeUpdate 在更新前将输入转换为 BSON 文档：指针解引用、数组递归转换，结构体按 bson
// 标签提取字段；标记 over 标签的字段被跳过，不进入更新内容。
func BeforeUpdate(document any) any {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)

	switch typ.Kind() {
	case reflect.Pointer:
		return BeforeUpdate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]any, val.Len(), val.Cap())
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
		return val.Interface()
	}
}

// acceptable 判断错误是否属于可忽略的常规情况：文档不存在、空值文档、空游标、空切片，
// 以及各类会话事务状态错误；其余错误视为需要上报的失败。
func acceptable(err error) bool {
	return err == nil || err == mongo.ErrNoDocuments || err == mongo.ErrNilValue ||
		err == mongo.ErrNilDocument || err == mongo.ErrNilCursor || err == mongo.ErrEmptySlice ||
		// 会话与事务状态错误
		err == session.ErrSessionEnded || err == session.ErrNoTransactStarted ||
		err == session.ErrTransactInProgress || err == session.ErrAbortAfterCommit ||
		err == session.ErrAbortTwice || err == session.ErrCommitAfterAbort ||
		err == session.ErrUnackWCUnsupported || err == session.ErrSnapshotTransaction
}

// Caller 返回跳过 skip 层调用栈后的调用方文件行号与函数名映射，供日志与追踪定位使用；
// skip 超出调用栈深度时返回零值字段。
func Caller(skip int) map[string]string {
	pc, file, lineNo, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	return map[string]string{
		"path": file + ":" + cast.ToString(lineNo),
		"func": name,
	}
}
