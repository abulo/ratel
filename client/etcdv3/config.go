package etcdv3

import (
	"time"

	"github.com/abulo/ratel/v3/core/logger"
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

// New 新建连接
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

// MustBuild ...
func (config *Config) MustBuild() *Client {
	client, err := config.Build()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("build etcd client failed")
	}
	return client
}

// SetEndpoints ...
func (config *Config) SetEndpoints(endpoint []string) *Config {
	for _, item := range endpoint {
		if !util.InArray(item, config.Endpoints) {
			config.Endpoints = append(config.Endpoints, item)
		}
	}
	return config
}

// SetCertFile ...
func (config *Config) SetCertFile(cert string) *Config {
	config.CertFile = cert
	return config
}

// SetKeyFile ...
func (config *Config) SetKeyFile(key string) *Config {
	config.KeyFile = key
	return config
}

// SetCaCert ...
func (config *Config) SetCaCert(ca string) *Config {
	config.CaCert = ca
	return config
}

// SetBasicAuth ...
func (config *Config) SetBasicAuth(auth bool) *Config {
	config.BasicAuth = auth
	return config
}

// SetUserName ...
func (config *Config) SetUserNam(userName string) *Config {
	config.UserName = userName
	return config
}

// SetPassword ...
func (config *Config) SetPassword(pwd string) *Config {
	config.Password = pwd
	return config
}

// SetConnectTimeout ...
func (config *Config) SetConnectTimeout(timeout time.Duration) *Config {
	config.ConnectTimeout = timeout
	return config
}

// SetSecure ...
func (config *Config) SetSecure(secure bool) *Config {
	config.Secure = secure
	return config
}

// SetAutoSyncInterval ...
func (config *Config) SetAutoSyncInterval(autoSyncInterval time.Duration) *Config {
	config.AutoSyncInterval = autoSyncInterval
	return config
}
