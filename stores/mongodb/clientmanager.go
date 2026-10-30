package mongodb

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/resource"
	"github.com/abulo/ratel/v3/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// defaultTimeout 为客户端操作默认超时时间。
	defaultTimeout = time.Second * 3
	// defaultSlowThreshold 为慢操作阈值的初始值。
	defaultSlowThreshold = time.Millisecond * 500
)

var clientManager = resource.NewResourceManager()
var slowThreshold = resource.ForAtomicDuration(defaultSlowThreshold)

type (
	// MongoDB 包装原生 mongo.Client，附加指标、追踪开关、库名与 URI 等元信息，
	// 并提供断开连接、Ping 校验与集合模型创建方法。
	//
	// 实例由 NewMongoDBClient 创建并按 URI 缓存复用；并发使用安全由底层 mongo.Client 保证。
	MongoDB struct {
		*mongo.Client
		// EnableMetric 是否采集 MongoDB 操作指标。
		EnableMetric bool
		// EnableTrace 是否挂载 MongoDB 链路追踪。
		EnableTrace bool
		// Name 为连接的数据库名称，取自 URI 路径。
		Name string
		// Uri 为原始连接 URI，作为连接缓存的去重标识。
		Uri string
	}

	// Option 是 mongo.ClientOptions 的配置选项函数，在 NewMongoDBClient 中依次装配。
	Option func(opts *options.ClientOptions)

	// Model 绑定一个集合的增强操作句柄，对外提供链式条件构造与数据操作方法。
	//
	// 实例由 MongoDB.NewModel 创建；Collection 对外暴露全部增强操作能力。
	Model struct {
		// Collection 为包装了熔断、指标与追踪的增强集合。
		Collection *DecoratedCollection
		cli        *mongo.Client
		brk        resource.Breaker
	}
)

// Close 断开底层 mongo.Client 与服务端的连接，返回 Disconnect 的错误。
func (cs *MongoDB) Close() error {
	return cs.Client.Disconnect(context.Background())
}

// SetSlowThreshold 全局设置 MongoDB 慢操作判定阈值。
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

// defaultTimeoutOption 注入默认的客户端操作超时时间，调用方 Option 可在其后覆盖。
func defaultTimeoutOption() Option {
	return func(opts *options.ClientOptions) {
		opts.SetTimeout(defaultTimeout)
	}
}

// WithTimeout 设置客户端单次操作超时时间。
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options.ClientOptions) {
		opts.SetTimeout(timeout)
	}
}

// WithMaxConnIdleTime 设置连接池中连接的最大空闲时长。
func WithMaxConnIdleTime(d time.Duration) Option {
	return func(opts *options.ClientOptions) {
		opts.SetMaxConnIdleTime(d)
	}
}

// WithMaxPoolSize 设置连接池最大连接数。
func WithMaxPoolSize(size uint64) Option {
	return func(opts *options.ClientOptions) {
		opts.SetMaxPoolSize(size)
	}
}

// WithMinPoolSize 设置连接池最小保持的连接数。
func WithMinPoolSize(size uint64) Option {
	return func(opts *options.ClientOptions) {
		opts.SetMinPoolSize(size)
	}
}

// NewMongoDBClient 依据 URI 创建并缓存 MongoDB 客户端。
//
// 同一 URI 复用已建立的连接；URI 缺少数据库名、URI 解析失败、连接或 Ping 失败时，
// 记录日志并 panic。建立连接的上下文超时为 5 秒，默认操作超时 3 秒，可由传入的 Option 覆盖。
func NewMongoDBClient(uri string, opts ...Option) (*MongoDB, error) {
	val, err := clientManager.GetResource(uri, func() (io.Closer, error) {
		o := options.Client().ApplyURI(uri)
		opts = append([]Option{defaultTimeoutOption()}, opts...)
		for _, opt := range opts {
			opt(o)
		}
		u, err := url.Parse(uri)
		if err != nil {
			logger.Logger.Panic(err)
			return nil, err
		}
		if util.Empty(u.Path) || util.Empty(u.Path[1:]) {
			logger.Logger.Panic(errors.New("no database"))
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cli, err := mongo.Connect(ctx, o)
		if err != nil {
			logger.Logger.Panic(err)
			return nil, err
		}

		err = cli.Ping(ctx, nil)
		if err != nil {
			logger.Logger.Panic(err)
			return nil, err
		}
		name := u.Path[1:]
		conn := &MongoDB{
			Client: cli,
			Name:   name,
			Uri:    uri,
		}
		return conn, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*MongoDB), nil
}

// SetEnableMetric 设置指标采集开关并返回客户端自身，便于链式调用。
func (m *MongoDB) SetEnableMetric(disableMetric bool) *MongoDB {
	m.EnableMetric = disableMetric
	return m
}

// SetEnableTrace 设置链路追踪开关并返回客户端自身，便于链式调用。
func (m *MongoDB) SetEnableTrace(disableTrace bool) *MongoDB {
	m.EnableTrace = disableTrace
	return m
}

// Ping 校验与 MongoDB 服务端的连通性，返回 Ping 的错误。
func (m *MongoDB) Ping() error {
	return m.Client.Ping(context.Background(), nil)
}

// NewModel 校验连通性后为指定集合创建 Model，继承当前客户端的指标与追踪开关；
// Ping 失败时返回错误。collection 为集合名称，调用方须保证非空。
func (m *MongoDB) NewModel(collection string) (*Model, error) {
	err := m.Ping()
	if err != nil {
		return nil, err
	}
	brk := resource.GetBreaker(m.Uri)
	coll := newCollection(m.Client.Database(m.Name).Collection(collection), brk, m.EnableMetric, m.EnableTrace)
	return &Model{
		Collection: coll,
		cli:        m.Client,
		brk:        brk,
	}, nil
}

// MustNewModel 为指定集合创建 Model；Ping 或创建失败时记录日志并终止程序，适用于启动期初始化。
func (m *MongoDB) MustNewModel(collection string) *Model {
	model, err := m.NewModel(collection)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	return model
}

// NewCollection 校验连通性后返回指定集合的 DecoratedCollection；Ping 失败时返回错误。
func (m *MongoDB) NewCollection(collection string) (*DecoratedCollection, error) {
	err := m.Ping()
	if err != nil {
		return nil, err
	}
	brk := resource.GetBreaker(m.Uri)
	return newCollection(m.Client.Database(m.Name).Collection(collection), brk, m.EnableMetric, m.EnableTrace), nil
}

// MustNewCollection 返回指定集合的 DecoratedCollection；Ping 或创建失败时记录日志并终止程序。
func (m *MongoDB) MustNewCollection(collection string) *DecoratedCollection {
	err := m.Ping()
	if err != nil {
		logger.Logger.Fatal(err)
	}
	brk := resource.GetBreaker(m.Uri)
	return newCollection(m.Client.Database(m.Name).Collection(collection), brk, m.EnableMetric, m.EnableTrace)
}
