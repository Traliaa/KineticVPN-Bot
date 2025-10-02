package middleware

import (
	"context"
	"net/http"

	"github.com/Traliaa/KineticVPN-Bot/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := opentracing.GlobalTracer()
		span := tracer.StartSpan(r.URL.Path)
		defer span.Finish()

		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			ctx := context.WithValue(r.Context(), tracing.TraceIDKey, sc.TraceID().String())
			ctx = context.WithValue(ctx, tracing.SpanIDKey, sc.SpanID().String())
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
