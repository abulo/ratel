package trace

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const MaxContentLength = 1 << 16

var elasticComponent = opentracing.Tag{string(ext.Component), "elastic"}

type ESTracedTransport struct {
	*http.Transport
}

func (t *ESTracedTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {

	ctx := r.Context()
	var span opentracing.Span
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if parentSpan := opentracing.SpanFromContext(r.Context()); parentSpan != nil {
		parentCtx := parentSpan.Context()
		span = opentracing.StartSpan("elastic", opentracing.ChildOf(parentCtx))
		ext.SpanKindRPCClient.Set(span)
		ext.PeerService.Set(span, "elastic")

		span.SetTag(string(ext.DBType), "elastic")
		span.SetTag(string(ext.DBInstance), r.URL.Host)
		span.SetTag("elastic.method", r.Method)
		span.SetTag("elastic.url", r.URL.Path)
		span.SetTag("elastic.params", r.URL.Query().Encode())
		// defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	r = r.WithContext(ctx)

	// span, ctx := StartSpanFromContext(
	// 	r.Context(),
	// 	"elastic",
	// 	TagComponent("http"),
	// 	TagSpanKind("server"),
	// 	HeaderExtractor(r.Header),
	// 	CustomTag("http.url", r.URL.Path),
	// 	CustomTag("http.method", r.Method),
	// )
	r = r.WithContext(ctx)
	defer func() {
		if err != nil {
			span.SetTag("elastic.error", err.Error())
			span.SetTag(string(ext.Error), true)
		}
		span.Finish()
	}()

	// span.SetTag(string(ext.DBType), "elastic")
	// span.SetTag(string(ext.DBInstance), r.URL.Host)
	// span.SetTag("elastic.method", r.Method)
	// span.SetTag("elastic.url", r.URL.Path)
	// span.SetTag("elastic.params", r.URL.Query().Encode())

	contentLength, _ := strconv.Atoi(r.Header.Get("Content-Length"))

	if r.Body != nil && contentLength < MaxContentLength {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		span.SetTag(string(ext.DBStatement), string(buf))
		span.LogFields(log.String("params", string(buf)))
		r.Body = io.NopCloser(bytes.NewBuffer(buf))
	}

	resp, err = t.Transport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	span.SetTag("elastic.status_code", resp.StatusCode)
	return resp, err
}

func ESTraceServerInterceptor() *http.Client {
	return &http.Client{
		Transport: &ESTracedTransport{&http.Transport{}},
	}
}
