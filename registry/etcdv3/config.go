package etcdv3

import (
	"time"

	"github.com/abulo/ratel/v3/client/etcdv3"
	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/registry"
)

// Config ...
type Config struct {
	*etcdv3.Config
	ReadTimeout time.Duration
	// ConfigKey   string
	Prefix     string
	ServiceTTL time.Duration
}

func New() *Config {
	return &Config{
		Config:      etcdv3.New(),
		ReadTimeout: time.Second * 3,
		Prefix:      "ratel",
		ServiceTTL:  0,
	}
}

// Build ...
func (config Config) Build() (registry.Registry, error) {
	return newETCDRegistry(&config)
}

func (config Config) MustBuild() registry.Registry {
	reg, err := config.Build()
	if err != nil {
		logger.Logger.Panicf("build registry failed: %v", err)
	}
	return reg
}