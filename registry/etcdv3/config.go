package etcdv3

import (
	"time"

	"github.com/abulo/ratel/client/etcdv3"
	"github.com/abulo/ratel/core/constant"
	"github.com/abulo/ratel/core/logger"
	"github.com/abulo/ratel/core/singleton"
	"github.com/abulo/ratel/registry"
)

// Config ...
type Config struct {
	*etcdv3.Config
	ReadTimeout time.Duration
	Node        string
	Prefix      string
	ServiceTTL  time.Duration
}

// New ...
func New() *Config {
	return &Config{
		Config:      etcdv3.New(),
		ReadTimeout: time.Second * 3,
		Prefix:      "ratel",
		ServiceTTL:  0,
	}
}

// SetNode ...
func (config *Config) SetNode(prefix string) *Config {
	config.Node = prefix
	return config
}

func (config *Config) GetNode() string {
	return config.Node
}

// SetPrefix ...
func (config *Config) SetPrefix(prefix string) *Config {
	config.Prefix = prefix
	return config
}

// Build ...
func (config *Config) Build() (registry.Registry, error) {
	return newETCDRegistry(config)
}

// MustBuild ...
func (config *Config) MustBuild() registry.Registry {
	reg, err := config.Build()
	if err != nil {
		logger.Logger.Panicf("build registry failed: %v", err)
	}
	return reg
}

func (config *Config) Singleton() (registry.Registry, error) {
	if val, ok := singleton.Load(constant.ModuleClientEtcd, config.Node); ok {
		return val.(registry.Registry), nil
	}

	reg, err := config.Build()
	if err != nil {
		return nil, err
	}

	singleton.Store(constant.ModuleClientEtcd, config.Node, reg)

	return reg, nil
}

func (config *Config) MustSingleton() registry.Registry {
	reg, err := config.Singleton()
	if err != nil {
		logger.Logger.Panicf("build registry failed: %v", err)
	}

	return reg
}
