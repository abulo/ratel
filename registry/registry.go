package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/abulo/ratel/server"
)

type Registry interface {
	RegisterService(context.Context, *server.ServiceInfo) error
	UnregisterService(context.Context, *server.ServiceInfo) error
	ListServices(context.Context, string, string) ([]*server.ServiceInfo, error)
	WatchServices(context.Context, string, string) (chan Endpoints, error)
	io.Closer
}

//GetServiceKey ..
func GetServiceKey(prefix string, s *server.ServiceInfo) string {
	return fmt.Sprintf("/%s/%s/%s://%s", prefix, s.Name, s.Scheme, s.Address)
}

//GetServiceValue ..
func GetServiceValue(s *server.ServiceInfo) string {
	val, _ := json.Marshal(s)
	return string(val)
}

//GetService ..
func GetService(s string) *server.ServiceInfo {
	var si server.ServiceInfo
	json.Unmarshal([]byte(s), &si)
	return &si
}

// Nop registry, used for local development/debugging
type Nop struct{}

// ListServices ...
func (n Nop) ListServices(ctx context.Context, s string, s2 string) ([]*server.ServiceInfo, error) {
	panic("implement me")
}

// WatchServices ...
func (n Nop) WatchServices(ctx context.Context, s string, s2 string) (chan Endpoints, error) {
	panic("implement me")
}

// RegisterService ...
func (n Nop) RegisterService(context.Context, *server.ServiceInfo) error { return nil }

// UnregisterService ...
func (n Nop) UnregisterService(context.Context, *server.ServiceInfo) error { return nil }

// Close ...
func (n Nop) Close() error { return nil }
