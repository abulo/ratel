package initial

import (
	"github.com/abulo/ratel/v3/core/trace"
	"github.com/abulo/ratel/v3/core/trace/jaeger"
	"github.com/spf13/cast"
)

// InitTrace ...
func (initial *Initial) InitTrace() {
	opt := jaeger.NewJaeger()
	conf := initial.Config.Get("trace")
	res := conf.(map[string]interface{})
	opt.EnableRPCMetrics = cast.ToBool(res["EnableRPCMetrics"])
	opt.LocalAgentHostPort = cast.ToString(res["LocalAgentHostPort"])
	opt.LogSpans = cast.ToBool(res["LogSpans"])
	opt.Param = cast.ToFloat64(res["Param"])
	opt.PanicOnError = cast.ToBool(res["PanicOnError"])
	client := opt.Build().Build()
	trace.SetGlobalTracer(client)
}
