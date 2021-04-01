package main

import (
	"github.com/abulo/ratel"
	"github.com/abulo/ratel/gin"
	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/server/http"
)

type Engine struct {
	ratel.Ratel
}

func main() {
	eng := NewEngine()

	if err := eng.Run(); err != nil {
		logger.Logger.Panic(err)
	}

}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.serveHTTP,
	); err != nil {
		logger.Logger.Panic("startup", err)
	}
	return eng
}
func (eng *Engine) serveHTTP() error {
	config := &http.Config{
		Host:    "127.0.0.1",
		Port:    7777,
		Mode:    gin.DebugMode,
		Name:    "admin",
		Network: "tcp4",
	}
	server := config.Build()

	server.Use(gin.Logger(), gin.Recovery())
	server.GET("/ping", "ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "7777",
		})
	})

	return eng.Serve(server)
}
