package xmonitor

import (
	"fmt"
)

// ModName ..
const ModName = "server.monitor"

// Config HTTP config
type Config struct {
	Host           string
	Port           int
	Deployment     string
	Mode           string
	ServiceAddress string // ServiceAddress service address in registry info, default to 'Host:Port'
}

func New() *Config {
	return &Config{}
}

// WithHost ...
func (config *Config) WithHost(host string) *Config {
	config.Host = host
	return config
}

// WithPort ...
func (config *Config) WithPort(port int) *Config {
	config.Port = port
	return config
}

// WithDeployment ...
func (config *Config) WithDeployment(deployment string) *Config {
	config.Deployment = deployment
	return config
}

// WithMode ...
func (config *Config) WithMode(mode string) *Config {
	config.Mode = mode
	return config
}

// WithServiceAddress ...
func (config *Config) WithServiceAddress(serviceAddress string) *Config {
	config.ServiceAddress = serviceAddress
	return config
}

// Build ...
func (config *Config) Build() *Server {
	return newServer(config)
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
