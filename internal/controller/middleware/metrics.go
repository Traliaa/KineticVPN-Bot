package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
	"time"
)

var (
	totalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_total_requests",
			Help: "Total number of HTTP requests to all endpoints.",
		},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request latencies.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "cleaned_url"},
	)
)

func init() {
	prometheus.MustRegister(
		totalRequests,
		requestDuration,
	)
}

func observeTotalRequests() {
	totalRequests.Inc()
}

func observeRequestDuration(statusCode int, cleanedURL string, duration float64) {
	requestDuration.WithLabelValues(http.StatusText(statusCode), cleanedURL).Observe(duration)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		observeTotalRequests()
		start := time.Now()
		statusCode := http.StatusOK

		defer func() {
			duration := time.Since(start).Seconds()
			cleanedURL := cleanURL(r.URL.Path)
			observeRequestDuration(statusCode, cleanedURL, duration)
		}()

		rw := &statusResponseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		statusCode = rw.statusCode
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func cleanURL(url string) string {
	return strings.TrimSuffix(url, "/")
}
