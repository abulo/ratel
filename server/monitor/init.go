package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
)

var (
	// DefaultServeMux ...
	DefaultServeMux = http.NewServeMux()
	routes          = []string{}
)

func (s *Server) InitHandle() {
	// 获取全部治理路由
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

// HandleFunc ...
func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	DefaultServeMux.HandleFunc(pattern, handler)
	routes = append(routes, pattern)
}
