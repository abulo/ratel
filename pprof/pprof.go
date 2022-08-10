package pprof

import (
	"net/http"
	"net/http/pprof"

	"github.com/abulo/ratel/v3/gin"
)

// DefaultPrefix ...
const (
	// DefaultPrefix url prefix of pprof
	DefaultPrefix = "/debug/pprof"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

// Register the standard HandlerFuncs from the net/http/pprof package with
// the provided gin.Engine. prefixOptions is a optional. If not prefixOptions,
// the default path prefix is used, otherwise first prefixOptions will be path prefix.
func Register(r *gin.Engine, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)

	prefixRouter := r.Group(prefix)
	{
		prefixRouter.GET("/", "debug_index", pprofHandler(pprof.Index))
		prefixRouter.GET("/cmdline", "debug_cmdline_get", pprofHandler(pprof.Cmdline))
		prefixRouter.GET("/profile", "debug_profile_get", pprofHandler(pprof.Profile))
		prefixRouter.POST("/symbol", "debug_symbol_post", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/symbol", "debug_symbol_get", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/trace", "debug_trace_get", pprofHandler(pprof.Trace))
		prefixRouter.GET("/allocs", "debug_allocs_get", pprofHandler(pprof.Handler("allocs").ServeHTTP))
		prefixRouter.GET("/block", "debug_block_get", pprofHandler(pprof.Handler("block").ServeHTTP))
		prefixRouter.GET("/goroutine", "debug_goroutine_get", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		prefixRouter.GET("/heap", "debug_heap_get", pprofHandler(pprof.Handler("heap").ServeHTTP))
		prefixRouter.GET("/mutex", "debug_mutex_get", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		prefixRouter.GET("/threadcreate", "debug_threadcreate_get", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	}
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
