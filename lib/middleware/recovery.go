package middleware

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rotabot-io/rotabot/lib/metrics"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
)

func RecoveryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rawErr := recover(); rawErr != nil {
				l := zapctx.Logger(r.Context())
				sentry.CurrentHub().RecoverWithContext(r.Context(), rawErr)
				sentry.Flush(time.Second * 5)

				var err error
				switch e := rawErr.(type) {
				case error:
					err = e
				default:
					err = fmt.Errorf("panic: %v", rawErr)
				}

				w.WriteHeader(http.StatusInternalServerError)
				endpoint := metrics.Endpoint(requestMethod(r), requestPath(r))
				metrics.PanicsTotal.With(prometheus.Labels{"endpoint": endpoint}).Inc()
				l.Error("request_panic", zap.Stack("stacktrace"), zap.Error(err))
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
