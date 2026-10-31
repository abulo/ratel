package mongodb

import (
	"context"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/core/call"
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
		// tracer 在集合句柄构造期创建一次，追踪开启后各操作复用（PERF-2-10 [P2]）。
		tracer *globalTrace.Tracer
		// dbName、collName 在构造期缓存，避免每个操作重复解析 Database() 与 Collection 名称（PERF-2-11 [P2]）。
		dbName   string
		collName string
		filter   bson.D
		limit    int64
		skip     int64
		sort     bson.D
		fields   bson.M
	}

	// index 为索引文档的解码结构，仅在 ListIndexes 内部使用。
	index struct {
		Key  bson.D
		Name string
	}
)

// emptyMetadataMongo 为追踪开启但无上游 metadata 时复用的只读空实例，
// 避免每个操作都 metadata.New(nil) 分配空 map（PERF-2-10 [P2]）。
var emptyMetadataMongo = metadata.MD{}

// mongoBaseAttrs 为各操作共享的静态基础属性（系统名），仅在包初始化时构造一次；
// 每次 startTrace 复制出局部切片后再追加调用方位置属性（PERF-2-10 [P2]）。
var mongoBaseAttrs = []attribute.KeyValue{semconv.RPCSystemKey.String("mongo")}

// newCollection 包装原生集合并注入指标、追踪开关与空的条件状态。
func newCollection(collection *mongo.Collection, brk resource.Breaker, EnableMetric, EnableTrace bool) *DecoratedCollection {
	c := &DecoratedCollection{
		Collection:   collection,
		brk:          brk,
		EnableMetric: EnableMetric,
		EnableTrace:  EnableTrace,
		dbName:       collection.Database().Name(), // PERF-2-11 [P2] 构造期缓存库名
		collName:     collection.Name(),            // PERF-2-11 [P2] 构造期缓存集合名
		filter:       make(bson.D, 0),
		sort:         make(bson.D, 0),
	}
	// PERF-2-10 [P2] tracer 在构造期创建一次，各操作复用，避免每次 NewTracer 查询 provider。
	if EnableTrace {
		c.tracer = globalTrace.NewTracer(trace.SpanKindClient)
	}
	return c
}

// startTrace 复用构造期 tracer 开启客户端 Span，附加调用方位置与集合属性；
// 返回携带 Span 的新上下文与 Span，调用方负责 defer span.End()。
// PERF-2-10 [P2] 公共追踪构造抽取到集合句柄方法内；Caller 跳过层数较原内联实现多 1 层（本方法自身帧），保持上报位置不变。
func (c *DecoratedCollection) startTrace(ctx context.Context, op string) (context.Context, trace.Span) {
	fn, file, line := call.Caller(7)
	// PERF-2-10 [P2] 静态系统名属性只复制一份基切片，调用方位置属性追加到局部切片。
	attrs := make([]attribute.KeyValue, 0, len(mongoBaseAttrs)+3)
	attrs = append(attrs, mongoBaseAttrs...)
	attrs = append(attrs,
		semconv.CodeFunction(fn),
		semconv.CodeFilepath(file),
		semconv.CodeLineNumber(line),
	)
	newCtx, span := c.tracer.Start(ctx, op, propagation.HeaderCarrier(emptyMetadataMongo), trace.WithAttributes(attrs...))
	span.SetAttributes(semconv.DBMongoDBCollection(c.dbName))
	return newCtx, span
}

// reportMetric 按统一四标签上报操作计数与耗时；库名、集合名使用构造期缓存值，
// 避免每个操作重复调用 Database().Name() 与 Name()（PERF-2-11 [P2]）。
func (c *DecoratedCollection) reportMetric(start time.Time, err error) {
	if !c.EnableMetric {
		return
	}
	cost := time.Since(start)
	status := "OK"
	if err != nil {
		status = "ERR"
	}
	metric.LibHandleCounter.WithLabelValues("mongodb", c.dbName, c.collName, status).Inc()
	metric.LibHandleHistogram.WithLabelValues("mongodb", c.dbName, c.collName).Observe(cost.Seconds())
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata，不再每操作 NewTracer/metadata.New。
			ctx, span = c.startTrace(ctx, "CreateIndex")
			defer span.End()
		}
		indexView := c.Collection.Indexes()
		indexModel := mongo.IndexModel{Keys: key, Options: op}
		res, err = indexView.CreateOne(ctx, indexModel)
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "ListIndexes")
			defer span.End()
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
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "DropIndex")
			defer span.End()
		}
		_, err = indexView.DropOne(ctx, name, opts)
		if err != nil {
			c.reset()
			return err
		}
		c.reset()
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "InsertOne")
			defer span.End()
		}
		val, err = c.Collection.InsertOne(ctx, data)
		c.reset()

		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "InsertMany")
			defer span.End()
		}
		val, err = c.Collection.InsertMany(ctx, data)
		c.reset()
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "Aggregate")
			defer span.End()
		}
		cursor, err := c.Collection.Aggregate(ctx, pipeline)
		if err != nil {
			c.reset()
			return err
		}
		err = cursor.All(ctx, result)
		c.reset()
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "UpdateOrInsert")
			defer span.End()
		}
		var upsert = true
		val, err = c.Collection.UpdateMany(ctx, c.filter, documents, &options.UpdateOptions{Upsert: &upsert})
		c.reset()

		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "UpdateOne")
			defer span.End()
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "UpdateOneRaw")
			defer span.End()
		}
		val, err = c.Collection.UpdateOne(ctx, c.filter, document, opt...)
		c.reset()
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "UpdateMany")
			defer span.End()
		}
		val, err = c.Collection.UpdateMany(ctx, c.filter, update)
		c.reset()
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "FindOne")
			defer span.End()
		}
		// PERF-2-12 [P2] skip 为 0 时不传指针，避免零值指针参与选项解析。
		findOneOpts := &options.FindOneOptions{Sort: c.sort, Projection: c.fields}
		if c.skip > 0 {
			findOneOpts.Skip = &c.skip
		}
		result := c.Collection.FindOne(ctx, c.filter, findOneOpts)
		err = result.Decode(document)
		c.reset()
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "FindMany")
			defer span.End()
		}
		// PERF-2-12 [P2] skip/limit 为 0 时不传指针，避免零值指针参与选项解析。
		findOpts := &options.FindOptions{Sort: c.sort, Projection: c.fields}
		if c.skip > 0 {
			findOpts.Skip = &c.skip
		}
		if c.limit > 0 {
			findOpts.Limit = &c.limit
		}
		result, err := c.Collection.Find(ctx, c.filter, findOpts)
		if err != nil {
			c.reset()
			return err
		}
		// PERF-2-8 [P1] 保留与原实现一致的“指向切片的指针”类型校验；随后由驱动 cursor.All
		// 直接批量解码到目标切片，省去逐文档 reflect.New/reflect.Append 与切片反复扩容，
		// 文档顺序与字段映射与原手写反射实现保持一致。
		val := reflect.ValueOf(documents)
		if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Slice {
			err = errors.New("result argument must be a slice address")
			c.reset()
			return err
		}
		if err = result.All(ctx, documents); err != nil {
			c.reset()
			return err
		}
		c.reset()
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "Delete")
			defer span.End()
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
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "Drop")
			defer span.End()
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = c.Collection.Drop(ctx)
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
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
			var span trace.Span
			// PERF-2-10 [P2] 复用构造期 tracer 与共享空 metadata。
			ctx, span = c.startTrace(ctx, "Count")
			defer span.End()
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		result, err = c.Collection.CountDocuments(ctx, c.filter)
		c.reset()
		// PERF-2-11 [P2] 指标标签使用构造期缓存的库名与集合名。
		c.reportMetric(start, err)
		return err
	}, acceptable)
	return
}

// mongoStructFieldMeta 缓存单个结构体字段的反射元信息，避免每次写入重复遍历字段与读取 tag。
type mongoStructFieldMeta struct {
	index    int
	bsonKey  string
	skipOver bool
}

// mongoStructMetaCache 按结构体 reflect.Type 缓存字段元信息；类型集合在运行期稳定，
// 缓存命中后反射遍历与 tag 读取只在首次发生（PERF-2-9 [P1]）。
var mongoStructMetaCache sync.Map // reflect.Type -> []mongoStructFieldMeta

// structFieldMeta 返回指定结构体类型缓存后的字段元信息；首次访问时构建并放入缓存。
func structFieldMeta(typ reflect.Type) []mongoStructFieldMeta {
	if v, ok := mongoStructMetaCache.Load(typ); ok {
		return v.([]mongoStructFieldMeta)
	}
	fields := make([]mongoStructFieldMeta, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		_, over := f.Tag.Lookup("over")
		fields = append(fields, mongoStructFieldMeta{index: i, bsonKey: f.Tag.Get("bson"), skipOver: over})
	}
	mongoStructMetaCache.Store(typ, fields)
	return fields
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
		// PERF-2-9 [P1] 字段元信息按类型缓存，消除每次写入的字段遍历与重复 tag 读取。
		data := make(bson.M)
		for _, f := range structFieldMeta(typ) {
			data[f.bsonKey] = val.Field(f.index).Interface()
		}
		// PERF-2-9 [P1] 合并原两次 FieldByName("Id") 重复反射查找；字段不存在时跳过，避免零值反射恐慌。
		if idField, ok := typ.FieldByName("Id"); ok {
			idVal := val.FieldByName("Id")
			if idField.Type == reflect.TypeFor[primitive.ObjectID]() {
				data["_id"] = primitive.NewObjectID()
			} else if idVal.Interface() == "" {
				data["_id"] = primitive.NewObjectID().String()
			}
		}
		return data

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
		// PERF-2-9 [P1] 字段元信息按类型缓存；over 标签在缓存期一次性判定。
		data := make(bson.M)
		for _, f := range structFieldMeta(typ) {
			if f.skipOver {
				continue
			}
			data[f.bsonKey] = val.Field(f.index).Interface()
		}
		return data

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
