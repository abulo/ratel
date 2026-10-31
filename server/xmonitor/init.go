package xmonitor

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"sync"
)

var (
	// DefaultServeMux 为监控服务端使用的默认路由表，所有治理接口均注册到该 mux。
	DefaultServeMux = http.NewServeMux()
	routes          = []string{}

	// PERF-6-21 [P2] 路由清单在启动注册完成后保持不变，/routes 响应 JSON 一次性编码并缓存，
	// 避免每次探测都重新走 json 编码；后续新增路由时重置缓存以保持语义正确。
	routesOnce  sync.Once
	routesCache []byte
)

// InitHandle 注册监控服务内置的治理接口：已注册路由清单、pprof 性能分析接口，
// 以及返回当前模块构建信息的 /modInfo 接口。
func (s *Server) InitHandle() {
	// /routes 接口返回当前已注册的全部路由模式清单。
	s.HandleFunc("/routes", func(resp http.ResponseWriter, req *http.Request) {
		routesOnce.Do(func() {
			// 与原 json.Encoder.Encode 行为一致：编码 JSON 后追加一个换行。
			if b, err := json.Marshal(routes); err == nil {
				routesCache = append(b, '\n')
			} else {
				routesCache = []byte("[]\n")
			}
		})
		_, _ = resp.Write(routesCache)
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
	// 路由清单变化后使已缓存的 /routes 响应失效，下次请求重新编码。
	routesOnce = sync.Once{}
	routesCache = nil
}
