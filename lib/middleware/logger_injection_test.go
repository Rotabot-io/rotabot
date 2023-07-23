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

var _ = Describe("RequestIdHandler", func() {
	It("should inject the request information into the context", func() {
		ctx := context.Background()

		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		observedLogger := zap.New(observedZapCore)
		ctx = zapctx.WithLogger(ctx, observedLogger)

		req := httptest.NewRequest(http.MethodGet, "http://testing/foo", nil)

		handlerToTest := LoggerInjectionHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				zapctx.Logger(r.Context()).Info("hello")
			}),
		)

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req.WithContext(ctx))

		Expect(observedLogs.Len()).To(Equal(1))

		loggedEntry := observedLogs.AllUntimed()[0]
		Expect(loggedEntry.Message).To(Equal("hello"))

		Expect(loggedEntry.Context).To(HaveLen(3))
		Expect(loggedEntry.Context[0].Key).To(Equal("method"))
		Expect(loggedEntry.Context[1].Key).To(Equal("path"))
		Expect(loggedEntry.Context[2].Key).To(Equal("request_id"))
	})
})
