package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

var _ = Describe("RequestAccess", func() {
	It("should log request access information", func() {
		ctx := context.Background()

		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		observedLogger := zap.New(observedZapCore)
		ctx = zapctx.WithLogger(ctx, observedLogger)

		req := httptest.NewRequest(http.MethodGet, "http://testing/foo", nil)

		handlerToTest := RequestAccessLogHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}),
		)

		res := httptest.NewRecorder()
		handlerToTest.ServeHTTP(res, req.WithContext(ctx))

		Expect(res.Code).To(Equal(http.StatusBadRequest))

		Expect(observedLogs.Len()).To(Equal(2))

		Expect(observedLogs.AllUntimed()[0].Message).To(Equal("request.start"))
		Expect(observedLogs.AllUntimed()[1].Message).To(Equal("request.finish"))

		loggedEntry := observedLogs.AllUntimed()[1]
		Expect(len(loggedEntry.Context)).To(Equal(2))
		Expect(loggedEntry.Context[0].Key).To(Equal("duration"))
		Expect(loggedEntry.Context[1].Key).To(Equal("status"))
		Expect(loggedEntry.Context[1].Integer).To(Equal(int64(http.StatusBadRequest)))
	})
})
