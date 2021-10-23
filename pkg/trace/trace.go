package trace

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/xqk/good/pkg/ilog"
)

var (
	// String ...
	String = log.String
)

// SetGlobalTracer ...
func SetGlobalTracer(tracer opentracing.Tracer) {
	ilog.Info("set global tracer", ilog.FieldMod("trace"))
	opentracing.SetGlobalTracer(tracer)
}

// Start ...
func StartSpanFromContext(ctx context.Context, op string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, op, opts...)
}

// SpanFromContext ...
func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}
