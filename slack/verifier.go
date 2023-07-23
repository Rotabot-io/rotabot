package slack

import (
	"bytes"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/rotabot-io/rotabot/lib/zapctx"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type RequestVerifierMiddleware struct {
	SigningSecret string
}

func (sv *RequestVerifierMiddleware) SlackSignatureVerifyHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(path.Clean(r.URL.EscapedPath()), "/slack") {
			// We don't want to verify requests that aren't from slack
			next.ServeHTTP(w, r)
		} else {
			l := zapctx.Logger(r.Context())
			verifier, err := slack.NewSecretsVerifier(r.Header, sv.SigningSecret)
			if err != nil {
				l.Error("failed to create secrets verifier", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Ensure the body can be read again
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			if _, err = verifier.Write(body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err = verifier.Ensure(); err != nil {
				l.Error("failed to verify request signature", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
