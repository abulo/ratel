package mongodb

import (
	"context"
	"errors"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/metric"
	"github.com/abulo/ratel/v3/core/resource"
	"github.com/abulo/ratel/v3/core/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

type (
	DecoratedCollection struct {
		*mongo.Collection
		// db            string
		// name          string
		brk           resource.Breaker
		DisableMetric bool // 关闭指标采集
		DisableTrace  bool // 关闭链路追踪
		filter        bson.D
		limit         int64
		skip          int64
		sort          bson.D
		fields        bson.M
	}

	index struct {
		Key  bson.D
		Name string
	}
)

func newCollection(collection *mongo.Collection, brk resource.Breaker, DisableMetric, DisableTrace bool) *DecoratedCollection {
	return &DecoratedCollection{
		Collection: collection,
		// name:          c.Name(),
		// db:            c.Database().Name(),
		brk:           brk,
		DisableMetric: DisableMetric,
		DisableTrace:  DisableTrace,
		filter:        make(bson.D, 0),
		sort:          make(bson.D, 0),
	}
}

// Where 条件查询, bson.M{"field": "value"}
func (c *DecoratedCollection) Where(m bson.D) *DecoratedCollection {
	c.filter = m
	return c
}

// Limit 限制条数
func (c *DecoratedCollection) Limit(n int64) *DecoratedCollection {
	c.limit = n
	return c
}

// Skip 跳过条数
func (c *DecoratedCollection) Skip(n int64) *DecoratedCollection {
	c.skip = n
	return c
}

// Sort 排序 bson.M{"created_at":-1}
func (c *DecoratedCollection) Sort(sorts bson.D) *DecoratedCollection {
	c.sort = sorts
	return c
}

// Fields 指定查询字段
func (c *DecoratedCollection) Fields(fields bson.M) *DecoratedCollection {
	c.fields = fields
	return c
}

func (c *DecoratedCollection) reset() {
	c.filter = nil
	c.limit = 0
	c.skip = 0
	c.sort = nil
	c.fields = nil
}

func getCtx(ctx context.Context) context.Context {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return ctx
}

// CreateIndex CreateOneIndex 创建单个普通索引
func (c *DecoratedCollection) CreateIndex(ctx context.Context, key bson.D, op *options.IndexOptions) (res string, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("keys", key))
				span.LogFields(log.Object("options", op))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		indexView := c.Collection.Indexes()
		indexModel := mongo.IndexModel{Keys: key, Options: op}
		res, err = indexView.CreateOne(ctx, indexModel)
		if !c.DisableMetric {
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

// ListIndexes 获取所有所有
func (c *DecoratedCollection) ListIndexes(ctx context.Context, opts *options.ListIndexesOptions) (val []string, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("options", opts))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
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
		if !c.DisableMetric {
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

// DropIndex 删除索引
func (c *DecoratedCollection) DropIndex(ctx context.Context, name string, opts *options.DropIndexesOptions) (err error) {
	start := time.Now()
	indexView := c.Collection.Indexes()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.String("indexname", name))
				span.LogFields(log.Object("options", opts))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		_, err = indexView.DropOne(ctx, name, opts)
		if err != nil {
			c.reset()
			return err
		}
		c.reset()
		if !c.DisableMetric {
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

// InsertOne 写入单条数据
func (c *DecoratedCollection) InsertOne(ctx context.Context, document any) (val *mongo.InsertOneResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		// var data any
		data := BeforeCreate(document)
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("document", data))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		val, err = c.Collection.InsertOne(ctx, data)
		c.reset()

		if !c.DisableMetric {
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

// InsertMany 写入多条数据
func (c *DecoratedCollection) InsertMany(ctx context.Context, documents any) (val *mongo.InsertManyResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		// var data []any
		data := BeforeCreate(documents).([]any)
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("documents", documents))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		val, err = c.Collection.InsertMany(ctx, data)
		c.reset()
		if !c.DisableMetric {
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

// Aggregate ...
func (c *DecoratedCollection) Aggregate(ctx context.Context, pipeline any, result any) (err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("pipeline", pipeline))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		cursor, err := c.Collection.Aggregate(ctx, pipeline)
		if err != nil {
			c.reset()
			return err
		}
		err = cursor.All(ctx, result)
		c.reset()
		if !c.DisableMetric {
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

// UpdateOrInsert 存在更新,不存在写入, documents 里边的文档需要有 _id 的存在
func (c *DecoratedCollection) UpdateOrInsert(ctx context.Context, documents []any) (val *mongo.UpdateResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("filter", c.filter))
				span.LogFields(log.Object("documents", documents))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		var upsert = true
		val, err = c.Collection.UpdateMany(ctx, c.filter, documents, &options.UpdateOptions{Upsert: &upsert})
		c.reset()

		if !c.DisableMetric {
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

// UpdateOne ...
func (c *DecoratedCollection) UpdateOne(ctx context.Context, document any) (val *mongo.UpdateResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		update := bson.M{"$set": BeforeUpdate(document)}
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("filter", c.filter))
				span.LogFields(log.Object("update", update))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		val, err = c.Collection.UpdateOne(ctx, c.filter, update)
		c.reset()

		if !c.DisableMetric {
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

// UpdateOneRaw 原生update
func (c *DecoratedCollection) UpdateOneRaw(ctx context.Context, document any, opt ...*options.UpdateOptions) (val *mongo.UpdateResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("filter", c.filter))
				span.LogFields(log.Object("update", document))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		val, err = c.Collection.UpdateOne(ctx, c.filter, document, opt...)
		c.reset()
		if !c.DisableMetric {
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

// UpdateMany ...
func (c *DecoratedCollection) UpdateMany(ctx context.Context, document any) (val *mongo.UpdateResult, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		update := bson.M{"$set": BeforeUpdate(document)}
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("filter", c.filter))
				span.LogFields(log.Object("update", update))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		val, err = c.Collection.UpdateMany(ctx, c.filter, update)
		c.reset()
		if !c.DisableMetric {
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

// FindOne 查询一条数据
func (c *DecoratedCollection) FindOne(ctx context.Context, document any) (err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("filter", c.filter))
				span.LogFields(log.Int64("skip", c.skip))
				span.LogFields(log.Object("sort", c.sort))
				span.LogFields(log.Object("fields", c.fields))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		result := c.Collection.FindOne(ctx, c.filter, &options.FindOneOptions{
			Skip:       &c.skip,
			Sort:       c.sort,
			Projection: c.fields,
		})
		err = result.Decode(document)
		c.reset()
		if !c.DisableMetric {
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

// FindMany 查询多条数据
func (c *DecoratedCollection) FindMany(ctx context.Context, documents any) (err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("filter", c.filter))
				span.LogFields(log.Int64("skip", c.skip))
				span.LogFields(log.Int64("limit", c.limit))
				span.LogFields(log.Object("sort", c.sort))
				span.LogFields(log.Object("fields", c.fields))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
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
		if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
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
		if !c.DisableMetric {
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

// Delete 删除数据,并返回删除成功的数量
func (c *DecoratedCollection) Delete(ctx context.Context) (count int64, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("filter", c.filter))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		if c.filter == nil || len(c.filter) == 0 {
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
		if !c.DisableMetric {
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

// Drop ...
func (c *DecoratedCollection) Drop(ctx context.Context) (err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = c.Collection.Drop(ctx)
		if !c.DisableMetric {
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

// Count ...
func (c *DecoratedCollection) Count(ctx context.Context) (result int64, err error) {
	start := time.Now()
	ctx = getCtx(ctx)
	err = c.brk.DoWithAcceptable(func() error {
		if !c.DisableTrace {
			call := Caller(6)
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
				span.LogFields(log.String("database", c.Database().Name()))
				span.LogFields(log.String("table", c.Name()))
				span.LogFields(log.Object("filter", c.filter))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		result, err = c.Collection.CountDocuments(ctx, c.filter)
		c.reset()
		if !c.DisableMetric {
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

// BeforeCreate ...
func BeforeCreate(document any) any {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)

	switch typ.Kind() {
	case reflect.Ptr:
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
func BeforeUpdate(document any) any {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)

	switch typ.Kind() {
	case reflect.Ptr:
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

func acceptable(err error) bool {
	return err == nil || err == mongo.ErrNoDocuments || err == mongo.ErrNilValue ||
		err == mongo.ErrNilDocument || err == mongo.ErrNilCursor || err == mongo.ErrEmptySlice ||
		// session errors
		err == session.ErrSessionEnded || err == session.ErrNoTransactStarted ||
		err == session.ErrTransactInProgress || err == session.ErrAbortAfterCommit ||
		err == session.ErrAbortTwice || err == session.ErrCommitAfterAbort ||
		err == session.ErrUnackWCUnsupported || err == session.ErrSnapshotTransaction
}

func Caller(skip int) map[string]string {
	pc, file, lineNo, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	return map[string]string{
		"path": file + ":" + cast.ToString(lineNo),
		"func": name,
	}
}
