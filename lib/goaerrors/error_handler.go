package goaerrors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
)

func ErrorHandler() func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, rawError error) {
		zapctx.Logger(ctx).Error(
			"unexpected_error",
			zap.Error(rawError),
		)
		http.Error(w, fmt.Sprintf("%s: Unexpected error occurred", "http"), http.StatusInternalServerError)
	}
}
