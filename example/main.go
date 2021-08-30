package main

import (
	"os"

	"github.com/abulo/ratel"
	ccc "github.com/abulo/ratel/config"
	"github.com/abulo/ratel/gin"

	// "github.com/abulo/ratel/gin/multitemplate"
	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/server/http"
	"github.com/abulo/ratel/util"
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

	//辅助函数
	server.InitFuncMap()
	server.AddFuncMap("config", ccc.String)
	server.AddFuncMap("marshalHtml", util.MarshalHTML)
	server.AddFuncMap("marshalJs", util.MarshalJS)
	server.AddFuncMap("static", util.Static)
	server.AddFuncMap("js", util.JS)
	server.AddFuncMap("formatDate", util.FormatDate)
	server.AddFuncMap("formatDateTime", util.FormatDateTime)
	server.AddFuncMap("inArray", util.InArray)
	server.AddFuncMap("multiArray", util.MultiArray)
	server.AddFuncMap("empty", util.Empty)
	server.AddFuncMap("divide", util.Divide)
	server.AddFuncMap("add", util.Add)
	server.AddFuncMap("strReplace", util.StrReplace)
	server.AddFuncMap("debugFormat", util.DebugFormat)
	server.GET("/ping", "ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": os.Getpid(),
		})
	})

	return eng.Serve(server)
}
