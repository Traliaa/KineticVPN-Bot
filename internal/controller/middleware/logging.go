package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Traliaa/KineticVPN-Bot/pkg/logger"
	"github.com/Traliaa/KineticVPN-Bot/pkg/tracing"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)

		logger.Info("method: %s, url: %s, time: %s, trace_id: %s, span_id: %s", req.Method, req.RequestURI, time.Since(start), getTraceID(req.Context()), getSpanID(req.Context()))
	})
}

func getTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(tracing.TraceIDKey).(string); ok {
		return traceID
	}

	return ""
}

func getSpanID(ctx context.Context) string {
	if spanID, ok := ctx.Value(tracing.SpanIDKey).(string); ok {
		return spanID
	}

	return ""
}
