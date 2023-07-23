package middleware

import (
	"fmt"
	"net/http"
	"path"

	"github.com/rotabot-io/rotabot/lib/zapctx"

	"go.uber.org/zap"
)

// LoggerInjectionHandler Injects the logger into the context with shared info of the request
// i.e. the request ID
func LoggerInjectionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := zapctx.Logger(ctx).With(
			zap.String("method", r.Method),
			zap.String("path", path.Clean(r.URL.EscapedPath())),
			zap.String("request_id", fmt.Sprintf("%s", ctx.Value(RequestIdKey))),
		)
		next.ServeHTTP(w, r.WithContext(zapctx.WithLogger(ctx, l)))
	})
}
