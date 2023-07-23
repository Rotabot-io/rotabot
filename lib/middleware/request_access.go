package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rotabot-io/rotabot/lib/metrics"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
)

type responseCapture struct {
	http.ResponseWriter
	statusCode int
}

func (rc *responseCapture) WriteHeader(statusCode int) {
	rc.statusCode = statusCode
	rc.ResponseWriter.WriteHeader(statusCode)
}

// RequestAccessLogHandler waits for the request to complete logging the outcome of it.
func RequestAccessLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		l := zapctx.Logger(r.Context())
		endpoint := metrics.Endpoint(requestMethod(r), requestPath(r))
		metrics.RequestsTotal.With(prometheus.Labels{
			"endpoint": endpoint,
		}).Inc()

		l.Info("request.start")
		capture := &responseCapture{ResponseWriter: w}

		defer func() {
			duration := time.Since(start).Seconds()
			metrics.RequestDuration.With(prometheus.Labels{
				"endpoint": endpoint,
				"status":   fmt.Sprintf("%v", capture.statusCode),
			}).Observe(duration)
			metrics.ResponsesTotal.With(prometheus.Labels{
				"endpoint": endpoint,
				"status":   fmt.Sprintf("%v", capture.statusCode),
			}).Inc()
			l.Info("request.finish",
				zap.Float64("duration", duration),
				zap.Int("status", capture.statusCode),
			)
		}()

		next.ServeHTTP(capture, r)
	})
}
