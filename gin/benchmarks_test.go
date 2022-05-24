// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"html/template"
	"net/http"
	"os"
	"testing"

	"github.com/abulo/ratel/v3/util"
)

func BenchmarkOneRoute(B *testing.B) {
	router := New()
	router.GET("/ping", util.Uniqid(""), func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkRecoveryMiddleware(B *testing.B) {
	router := New()
	router.Use(Recovery())
	router.GET("/", util.Uniqid(""), func(c *Context) {})
	runRequest(B, router, "GET", "/")
}

func BenchmarkLoggerMiddleware(B *testing.B) {
	router := New()
	router.Use(LoggerWithWriter(newMockWriter()))
	router.GET("/", util.Uniqid(""), func(c *Context) {})
	runRequest(B, router, "GET", "/")
}

func BenchmarkManyHandlers(B *testing.B) {
	router := New()
	router.Use(Recovery(), LoggerWithWriter(newMockWriter()))
	router.Use(func(c *Context) {})
	router.Use(func(c *Context) {})
	router.GET("/ping", util.Uniqid(""), func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func Benchmark5Params(B *testing.B) {
	DefaultWriter = os.Stdout
	router := New()
	router.Use(func(c *Context) {})
	router.GET("/param/:param1/:params2/:param3/:param4/:param5", util.Uniqid(""), func(c *Context) {})
	runRequest(B, router, "GET", "/param/path/to/parameter/john/12345")
}

func BenchmarkOneRouteJSON(B *testing.B) {
	router := New()
	data := struct {
		Status string `json:"status"`
	}{"ok"}
	router.GET("/json", util.Uniqid(""), func(c *Context) {
		c.JSON(http.StatusOK, data)
	})
	runRequest(B, router, "GET", "/json")
}

func BenchmarkOneRouteHTML(B *testing.B) {
	router := New()
	t := template.Must(template.New("index").Parse(`
		<html><body><h1>{{.}}</h1></body></html>`))
	router.SetHTMLTemplate(t)

	router.GET("/html", util.Uniqid(""), func(c *Context) {
		c.HTML(http.StatusOK, "index", "hola")
	})
	runRequest(B, router, "GET", "/html")
}

func BenchmarkOneRouteSet(B *testing.B) {
	router := New()
	router.GET("/ping", util.Uniqid(""), func(c *Context) {
		c.Set("key", "value")
	})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkOneRouteString(B *testing.B) {
	router := New()
	router.GET("/text", util.Uniqid(""), func(c *Context) {
		c.String(http.StatusOK, "this is a plain text")
	})
	runRequest(B, router, "GET", "/text")
}

func BenchmarkManyRoutesFist(B *testing.B) {
	router := New()
	router.Any("/ping", util.Uniqid(""), func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkManyRoutesLast(B *testing.B) {
	router := New()
	router.Any("/ping", util.Uniqid(""), func(c *Context) {})
	runRequest(B, router, "OPTIONS", "/ping")
}

func Benchmark404(B *testing.B) {
	router := New()
	router.Any("/something", util.Uniqid(""), func(c *Context) {})
	router.NoRoute(func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func Benchmark404Many(B *testing.B) {
	router := New()
	router.GET("/", util.Uniqid(""), func(c *Context) {})
	router.GET("/path/to/something", util.Uniqid(""), func(c *Context) {})
	router.GET("/post/:id", util.Uniqid(""), func(c *Context) {})
	router.GET("/view/:id", util.Uniqid(""), func(c *Context) {})
	router.GET("/favicon.ico", util.Uniqid(""), func(c *Context) {})
	router.GET("/robots.txt", util.Uniqid(""), func(c *Context) {})
	router.GET("/delete/:id", util.Uniqid(""), func(c *Context) {})
	router.GET("/user/:id/:mode", util.Uniqid(""), func(c *Context) {})

	router.NoRoute(func(c *Context) {})
	runRequest(B, router, "GET", "/viewfake")
}

type mockWriter struct {
	headers http.Header
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		http.Header{},
	}
}

func (m *mockWriter) Header() (h http.Header) {
	return m.headers
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockWriter) WriteHeader(int) {}

func runRequest(B *testing.B, r *Engine, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}
