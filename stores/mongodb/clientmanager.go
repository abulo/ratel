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
	defaultTimeout       = time.Second * 3
	defaultSlowThreshold = time.Millisecond * 500
)

var clientManager = resource.NewResourceManager()
var slowThreshold = resource.ForAtomicDuration(defaultSlowThreshold)

// MongoDB wraps *mongo.Client and provides a Close method.
type (
	MongoDB struct {
		*mongo.Client
		DisableMetric bool   // 关闭指标采集
		DisableTrace  bool   // 关闭链路追踪
		Name          string // 数据库名称
		Uri           string
	}
	// Option defines the method to customize a mongo model.
	Option func(opts *options.ClientOptions)

	Model struct {
		Collection *DecoratedCollection
		// name       string
		cli *mongo.Client
		brk resource.Breaker
	}
)

// Close disconnects the underlying *mongo.Client.
func (cs *MongoDB) Close() error {
	return cs.Client.Disconnect(context.Background())
}

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

func defaultTimeoutOption() Option {
	return func(opts *options.ClientOptions) {
		opts.SetTimeout(defaultTimeout)
	}
}

// WithTimeout set the mon client operation timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options.ClientOptions) {
		opts.SetTimeout(timeout)
	}
}

// WithMaxConnIdleTime set the mon client max connection idle time.
func WithMaxConnIdleTime(d time.Duration) Option {
	return func(opts *options.ClientOptions) {
		opts.SetMaxConnIdleTime(d)
	}
}

// WithMaxPoolSize set the mon client max pool size.
func WithMaxPoolSize(size uint64) Option {
	return func(opts *options.ClientOptions) {
		opts.SetMaxPoolSize(size)
	}
}

// WithMinPoolSize set the mon client min pool size.
func WithMinPoolSize(size uint64) Option {
	return func(opts *options.ClientOptions) {
		opts.SetMinPoolSize(size)
	}
}

// NewMongoDBClient
func NewMongoDBClient(uri string, opts ...Option) (*MongoDB, error) {
	val, err := clientManager.GetResource(uri, func() (io.Closer, error) {
		o := options.Client().ApplyURI(uri)
		opts = append([]Option{defaultTimeoutOption()}, opts...)
		for _, opt := range opts {
			opt(o)
		}
		//解析URL
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

func (m *MongoDB) SetDisableMetric(disableMetric bool) *MongoDB {
	m.DisableMetric = disableMetric
	return m
}

func (m *MongoDB) SetDisableTrace(disableTrace bool) *MongoDB {
	m.DisableTrace = disableTrace
	return m
}

func (m *MongoDB) Ping() error {
	return m.Client.Ping(context.Background(), nil)
}

func (m *MongoDB) NewModel(collection string) (*Model, error) {
	err := m.Ping()
	if err != nil {
		return nil, err
	}
	brk := resource.GetBreaker(m.Uri)
	coll := newCollection(m.Client.Database(m.Name).Collection(collection), brk, m.DisableMetric, m.DisableTrace)
	return &Model{
		Collection: coll,
		cli:        m.Client,
		brk:        brk,
	}, nil
}

// MustNewModel returns a Model, exits on errors.
func (m *MongoDB) MustNewModel(collection string) *Model {
	model, err := m.NewModel(collection)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	return model
}

func (m *MongoDB) NewCollection(collection string) (*DecoratedCollection, error) {
	err := m.Ping()
	if err != nil {
		return nil, err
	}
	brk := resource.GetBreaker(m.Uri)
	return newCollection(m.Client.Database(m.Name).Collection(collection), brk, m.DisableMetric, m.DisableTrace), nil
}

func (m *MongoDB) MustNewCollection(collection string) *DecoratedCollection {
	err := m.Ping()
	if err != nil {
		logger.Logger.Fatal(err)
	}
	brk := resource.GetBreaker(m.Uri)
	return newCollection(m.Client.Database(m.Name).Collection(collection), brk, m.DisableMetric, m.DisableTrace)
}
