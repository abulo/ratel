package etcdv3

import (
	"time"

	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/util"
	"github.com/sirupsen/logrus"
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

// 新建连接
func New() *Config {
	return &Config{
		BasicAuth:      false,
		ConnectTimeout: util.Duration("5s"),
		Secure:         false,
	}
}

// Build ...
func (config *Config) Build() (*Client, error) {
	return newClient(config)
}

func (config *Config) MustBuild() *Client {
	client, err := config.Build()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("build etcd client failed")
	}
	return client
}
