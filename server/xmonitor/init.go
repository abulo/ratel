package xmonitor

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
)

var (
	// DefaultServeMux 为监控服务端使用的默认路由表，所有治理接口均注册到该 mux。
	DefaultServeMux = http.NewServeMux()
	routes          = []string{}
)

// InitHandle 注册监控服务内置的治理接口：已注册路由清单、pprof 性能分析接口，
// 以及返回当前模块构建信息的 /modInfo 接口。
func (s *Server) InitHandle() {
	// /routes 接口返回当前已注册的全部路由模式清单。
	s.HandleFunc("/routes", func(resp http.ResponseWriter, req *http.Request) {
		_ = json.NewEncoder(resp).Encode(routes)
	})
	s.HandleFunc("/debug/pprof/", pprof.Index)
	s.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.HandleFunc("/debug/pprof/trace", pprof.Trace)

	if info, ok := debug.ReadBuildInfo(); ok {
		s.HandleFunc("/modInfo", func(w http.ResponseWriter, r *http.Request) {
			encoder := json.NewEncoder(w)
			if r.URL.Query().Get("pretty") == "true" {
				encoder.SetIndent("", "    ")
			}
			_ = encoder.Encode(info)
		})
	}
}

// HandleFunc 在 DefaultServeMux 上注册路由处理器，并将路由模式记入路由清单供 /routes 接口返回。
func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	DefaultServeMux.HandleFunc(pattern, handler)
	routes = append(routes, pattern)
}
