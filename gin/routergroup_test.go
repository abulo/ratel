// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"testing"

	"github.com/abulo/ratel/v3/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	SetMode(TestMode)
}

func TestRouterGroupBasic(t *testing.T) {
	router := New()
	group := router.Group("/hola", func(c *Context) {})
	group.Use(func(c *Context) {})

	assert.Len(t, group.Handlers, 2)
	assert.Equal(t, "/hola", group.BasePath())
	assert.Equal(t, router, group.engine)

	group2 := group.Group("manu")
	group2.Use(func(c *Context) {}, func(c *Context) {})

	assert.Len(t, group2.Handlers, 4)
	assert.Equal(t, "/hola/manu", group2.BasePath())
	assert.Equal(t, router, group2.engine)
}

func TestRouterGroupBasicHandle(t *testing.T) {
	performRequestInGroup(t, http.MethodGet)
	performRequestInGroup(t, http.MethodPost)
	performRequestInGroup(t, http.MethodPut)
	performRequestInGroup(t, http.MethodPatch)
	performRequestInGroup(t, http.MethodDelete)
	performRequestInGroup(t, http.MethodHead)
	performRequestInGroup(t, http.MethodOptions)
}

func performRequestInGroup(t *testing.T, method string) {
	router := New()
	v1 := router.Group("v1", func(c *Context) {})
	assert.Equal(t, "/v1", v1.BasePath())

	login := v1.Group("/login/", func(c *Context) {}, func(c *Context) {})
	assert.Equal(t, "/v1/login/", login.BasePath())

	handler := func(c *Context) {
		c.String(http.StatusBadRequest, "the method was %s and index %d", c.Request.Method, c.index)
	}

	switch method {
	case http.MethodGet:
		v1.GET("/test", util.Uniqid(""), handler)
		login.GET("/test", util.Uniqid(""), handler)
	case http.MethodPost:
		v1.POST("/test", util.Uniqid(""), handler)
		login.POST("/test", util.Uniqid(""), handler)
	case http.MethodPut:
		v1.PUT("/test", util.Uniqid(""), handler)
		login.PUT("/test", util.Uniqid(""), handler)
	case http.MethodPatch:
		v1.PATCH("/test", util.Uniqid(""), handler)
		login.PATCH("/test", util.Uniqid(""), handler)
	case http.MethodDelete:
		v1.DELETE("/test", util.Uniqid(""), handler)
		login.DELETE("/test", util.Uniqid(""), handler)
	case http.MethodHead:
		v1.HEAD("/test", util.Uniqid(""), handler)
		login.HEAD("/test", util.Uniqid(""), handler)
	case http.MethodOptions:
		v1.OPTIONS("/test", util.Uniqid(""), handler)
		login.OPTIONS("/test", util.Uniqid(""), handler)
	default:
		panic("unknown method")
	}

	w := PerformRequest(router, method, "/v1/login/test")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "the method was "+method+" and index 3", w.Body.String())

	w = PerformRequest(router, method, "/v1/test")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "the method was "+method+" and index 1", w.Body.String())
}

func TestRouterGroupInvalidStatic(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.Static("/path/:param", "/")
	})

	assert.Panics(t, func() {
		router.Static("/path/*param", "/")
	})
}

func TestRouterGroupInvalidStaticFile(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.StaticFile("/path/:param", "favicon.ico")
	})

	assert.Panics(t, func() {
		router.StaticFile("/path/*param", "favicon.ico")
	})
}

func TestRouterGroupInvalidStaticFileFS(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.StaticFileFS("/path/:param", "favicon.ico", Dir(".", false))
	})

	assert.Panics(t, func() {
		router.StaticFileFS("/path/*param", "favicon.ico", Dir(".", false))
	})
}

func TestRouterGroupTooManyHandlers(t *testing.T) {
	const (
		panicValue = "too many handlers"
		maximumCnt = abortIndex
	)
	router := New()
	handlers1 := make([]HandlerFunc, maximumCnt-1)
	router.Use(handlers1...)

	handlers2 := make([]HandlerFunc, maximumCnt+1)
	assert.PanicsWithValue(t, panicValue, func() {
		router.Use(handlers2...)
	})
	assert.PanicsWithValue(t, panicValue, func() {
		router.GET("/", util.Uniqid(""), handlers2...)
	})
}

func TestRouterGroupBadMethod(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.Handle(http.MethodGet, "/", util.Uniqid(""))
	})
	assert.Panics(t, func() {
		router.Handle(" GET", "/", util.Uniqid(""))
	})
	assert.Panics(t, func() {
		router.Handle("GET ", "/", util.Uniqid(""))
	})
	assert.Panics(t, func() {
		router.Handle("", "/", util.Uniqid(""))
	})
	assert.Panics(t, func() {
		router.Handle("PO ST", "/", util.Uniqid(""))
	})
	assert.Panics(t, func() {
		router.Handle("1GET", "/", util.Uniqid(""))
	})
	assert.Panics(t, func() {
		router.Handle("PATCh", "/", util.Uniqid(""))
	})
}

func TestRouterGroupPipeline(t *testing.T) {
	router := New()
	testRoutesInterface(t, router)

	v1 := router.Group("/v1")
	testRoutesInterface(t, v1)
}

func testRoutesInterface(t *testing.T, r IRoutes) {
	handler := func(c *Context) {}
	assert.Equal(t, r, r.Use(handler))

	assert.Equal(t, r, r.Handle(http.MethodGet, "/handler", util.Uniqid(""), handler))
	assert.Equal(t, r, r.Any("/any", util.Uniqid(""), handler))
	assert.Equal(t, r, r.GET("/", util.Uniqid(""), handler))
	assert.Equal(t, r, r.POST("/", util.Uniqid(""), handler))
	assert.Equal(t, r, r.DELETE("/", util.Uniqid(""), handler))
	assert.Equal(t, r, r.PATCH("/", util.Uniqid(""), handler))
	assert.Equal(t, r, r.PUT("/", util.Uniqid(""), handler))
	assert.Equal(t, r, r.OPTIONS("/", util.Uniqid(""), handler))
	assert.Equal(t, r, r.HEAD("/", util.Uniqid(""), handler))

	// assert.Equal(t, r, r.StaticFile("/file", "."))
	// assert.Equal(t, r, r.StaticFileFS("/static2", ".", Dir(".", false)))
	// assert.Equal(t, r, r.Static("/static", "."))
	// assert.Equal(t, r, r.StaticFS("/static2", Dir(".", false)))
}
