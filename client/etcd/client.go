package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/abulo/ratel/logger"
	"github.com/coreos/etcd/clientv3"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// Config ...
type Config struct {
	Endpoints        []string      `json:"endpoints"`
	CertFile         string        `json:"certFile"`
	KeyFile          string        `json:"keyFile"`
	CaCert           string        `json:"caCert"`
	BasicAuth        bool          `json:"basicAuth"`
	UserName         string        `json:"userName"`
	Password         string        `json:"-"`
	ConnectTimeout   time.Duration `json:"connectTimeout"` // 连接超时时间
	Secure           bool          `json:"secure"`
	AutoSyncInterval time.Duration `json:"autoAsyncInterval"` // 自动同步member list的间隔
	TTL              int           // 单位：s
}

// Client ...
type Client struct {
	*clientv3.Client
	config *Config
}

func (config *Config) Build() *Client {
	cc := newClient(config)
	return cc
}

// New ...
func newClient(config *Config) *Client {
	conf := clientv3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          config.ConnectTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		},
		AutoSyncInterval: config.AutoSyncInterval,
	}
	if config.Endpoints == nil {
		logger.Panic("client etcd endpoints empty", config)
	}
	if !config.Secure {
		conf.DialOptions = append(conf.DialOptions, grpc.WithInsecure())
	}
	if config.BasicAuth {
		conf.Username = config.UserName
		conf.Password = config.Password
	}
	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}
	if config.CaCert != "" {
		certBytes, err := ioutil.ReadFile(config.CaCert)
		if err != nil {
			logger.Panic("parse CaCert failed", err)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
		tlsEnabled = true
	}
	if config.CertFile != "" && config.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			logger.Panic("load CertFile or KeyFile failed", err)
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}
	if tlsEnabled {
		conf.TLS = tlsConfig
	}

	client, err := clientv3.New(conf)

	if err != nil {
		logger.Panic("client etcd start panic", err)
	}

	cc := &Client{
		Client: client,
		config: config,
	}
	logger.Info("dial etcd server")
	return cc
}
