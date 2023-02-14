package registry

import (
	"log"

	"github.com/abulo/ratel/core/logger"
)

var registryBuilder = make(map[string]Builder)

// Config ...
type Config map[string]ConfigLab

// ConfigLab ...
type ConfigLab struct {
	Kind          string `json:"kind" description:"底层注册器类型, eg: etcdv3, consul"`
	ConfigKey     string `json:"configKey" description:"底册注册器的配置键"`
	DeplaySeconds int    `json:"deplaySeconds" description:"延迟注册"`
}

// DefaultRegisterer default register
var DefaultRegisterer Registry = &Local{}

// Builder ...
type Builder func(string) Registry

// BuildFunc ...
type BuildFunc func(string) (Registry, error)

// RegisterBuilder ...
func RegisterBuilder(kind string, build Builder) {
	if _, ok := registryBuilder[kind]; ok {
		log.Panicf("duplicate register registry builder: %s", kind)
	}
	registryBuilder[kind] = build
}

// New ...
func New() Config {
	var config Config
	return config
}

// Lab ...
func (config Config) Lab(name string, lab ConfigLab) Config {
	config["ddd"] = lab
	return config
}

// InitDefaultRegister ...
func (config Config) InitDefaultRegister() {
	for name, item := range config {
		var itemKind = item.Kind
		if itemKind == "" {
			itemKind = "etcdv3"
		}
		build, ok := registryBuilder[itemKind]
		if !ok {
			logger.Logger.Printf("invalid registry kind: %s", itemKind)
			continue
		}
		DefaultRegisterer = build(item.ConfigKey)
		logger.Logger.Printf("build registry %s with config: %s", name, item.ConfigKey)
	}
}
