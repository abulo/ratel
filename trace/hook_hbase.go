package trace

import (
	"context"
	"io"

	"github.com/abulo/ratel/hbase"
	"github.com/tsuna/gohbase/hrpc"
)

//HBaseTrace ...
func HBaseTrace(component, instance string) hbase.HookFunc {
	return func(ctx context.Context, call hrpc.Call, customName string) func(err error) {
		if customName == "" {
			customName = call.Name()
		}
		statement := string(call.Table()) + " " + string(call.Key())
		span, ctx := StartSpanFromContext(
			ctx,
			"Hbase:"+customName,
			TagComponent("Hbase"),
			TagSpanKind("client"),
			CustomTag("statement", statement),
		)
		return func(err error) {
			if err == io.EOF {
				err = nil
			}
			defer span.Finish()
		}
	}
}
