package lib

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"log"
	"net/http"
)

//InitializeJaeger initializes a new Jaeger instance.
func InitializeJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("error: failed to read config from env vars: %v\n", err)
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatalf("error: cannot init Jaeger: %v\n", err)
	}
	return tracer, closer
}


// Inject req with Active un-finished Span.
func Inject(span opentracing.Span, request *http.Request) error {
	carrier := opentracing.HTTPHeadersCarrier(request.Header)
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		carrier)
}

// Extract span from ctx
func Extract(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}

//StartSpanFromRequest starts the span
func StartSpanFromRequest(tracer opentracing.Tracer, r *http.Request) opentracing.Span {
	spanCtx, _ := Extract(tracer, r)
	return tracer.StartSpan("ping-receive", ext.RPCServerOption(spanCtx))
}
