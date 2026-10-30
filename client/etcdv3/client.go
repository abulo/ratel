package etcdv3

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client 封装 go.etcd.io/etcd/client/v3 客户端及其构建配置，键值、锁、监听等能力均通过该对象发起。
//
// 零值不可直接使用，须由 NewClient 或 Config.Build 构建；内嵌的 clientv3.Client 方法可直接调用。
type Client struct {
	*clientv3.Client
	config *Config
}

// NewClient 依据 config 构建 etcd 客户端并装配指标与（可选的）链路追踪拦截器。
//
// Endpoints 为 nil 时返回错误；读取 CA 证书或客户端证书/私钥失败时会直接 panic，调用方须保证
// 配置的证书文件路径存在且可读；BasicAuth 为 true 时使用配置的用户名密码认证。
func NewClient(config *Config) (*Client, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(grpcprom.UnaryClientInterceptor),
		grpc.WithChainStreamInterceptor(grpcprom.StreamClientInterceptor),
	}

	if config.EnableTrace {
		dialOptions = append(dialOptions,
			grpc.WithChainUnaryInterceptor(traceUnaryClientInterceptor()),
			grpc.WithChainStreamInterceptor(traceStreamClientInterceptor()),
		)
	}

	conf := clientv3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          config.ConnectTimeout,
		DialKeepAliveTime:    5 * time.Minute,
		DialKeepAliveTimeout: 20 * time.Second,
		DialOptions:          dialOptions,
		AutoSyncInterval:     config.AutoSyncInterval,
	}

	if config.Endpoints == nil {
		return nil, fmt.Errorf("client etcd endpoints empty, empty endpoints")
	}

	if !config.Secure {
		conf.DialOptions = append(conf.DialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		certBytes, err := os.ReadFile(config.CaCert)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Panic("parse CaCert failed")
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
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Panic("load CertFile or KeyFile failed")
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}

	if tlsEnabled {
		conf.TLS = tlsConfig
	}

	client, err := clientv3.New(conf)

	if err != nil {
		return nil, fmt.Errorf("client etcd start failed: %v", err)
	}

	cc := &Client{
		Client: client,
		config: config,
	}
	logger.Logger.Info("dial etcd server")
	return cc, nil
}

// GetKeyValue 查询单个键对应的 mvccpb.KeyValue；键不存在时返回 (nil, nil)，查询失败时返回错误。
func (client *Client) GetKeyValue(ctx context.Context, key string) (kv *mvccpb.KeyValue, err error) {
	rp, err := client.Client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(rp.Kvs) > 0 {
		return rp.Kvs[0], nil
	}

	return
}

// GetPrefix 查询以 prefix 为前缀的全部键值，返回键到值的映射；查询失败时返回已收集的空映射与错误。
func (client *Client) GetPrefix(ctx context.Context, prefix string) (map[string]string, error) {
	var (
		vars = make(map[string]string)
	)

	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return vars, err
	}

	for _, kv := range resp.Kvs {
		vars[string(kv.Key)] = string(kv.Value)
	}

	return vars, nil
}

// DelPrefix 删除以 prefix 为前缀的全部键，返回被删除的键数量；删除失败时返回错误。
func (client *Client) DelPrefix(ctx context.Context, prefix string) (deleted int64, err error) {
	resp, err := client.Delete(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}
	return resp.Deleted, err
}

// GetValues 按前缀批量查询多组键值并汇总为键到值的映射；任一批次查询失败时返回已收集的映射与错误。
func (client *Client) GetValues(ctx context.Context, keys ...string) (map[string]string, error) {
	var (
		firstRevision = int64(0)
		vars          = make(map[string]string)
		maxTxnOps     = 128
	)

	// 每批最多 128 个键；首批以最新修订读取，之后各批沿用首个响应的修订版本，保证跨批读取的
	// 修订基线一致。
	doTxn := func(ops []string) error {
		txnOps := make([]clientv3.Op, 0, maxTxnOps)

		for _, k := range ops {
			txnOps = append(txnOps, clientv3.OpGet(k,
				clientv3.WithPrefix(),
				clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend),
				clientv3.WithRev(firstRevision)))
		}

		result, err := client.Txn(ctx).Then(txnOps...).Commit()
		if err != nil {
			return err
		}
		for i, r := range result.Responses {
			originKey := ops[i]
			originKeyFixed := originKey
			if !strings.HasSuffix(originKeyFixed, "/") {
				originKeyFixed = originKey + "/"
			}
			for _, ev := range r.GetResponseRange().Kvs {
				k := string(ev.Key)
				if k == originKey || strings.HasPrefix(k, originKeyFixed) {
					vars[string(ev.Key)] = string(ev.Value)
				}
			}
		}
		if firstRevision == 0 {
			firstRevision = result.Header.GetRevision()
		}
		return nil
	}
	cnt := len(keys) / maxTxnOps
	for i := 0; i <= cnt; i++ {
		switch temp := (i == cnt); temp {
		case false:
			if err := doTxn(keys[i*maxTxnOps : (i+1)*maxTxnOps]); err != nil {
				return vars, err
			}
		case true:
			if err := doTxn(keys[i*maxTxnOps:]); err != nil {
				return vars, err
			}
		}
	}
	return vars, nil
}

// GetLeaseSession 基于当前客户端创建租约会话，可传入 SessionOption 调整 TTL 等参数；
// 会话创建失败时返回错误。
func (client *Client) GetLeaseSession(ctx context.Context, opts ...concurrency.SessionOption) (leaseSession *concurrency.Session, err error) {
	return concurrency.NewSession(client.Client, opts...)
}
