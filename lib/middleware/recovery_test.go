package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

var _ = Describe("Recovery", func() {
	It("should recover from a panic and log message", func() {
		ctx := context.Background()

		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		observedLogger := zap.New(observedZapCore)
		ctx = zapctx.WithLogger(ctx, observedLogger)

		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

		handlerToTest := RecoveryHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("test")
			}),
		)

		res := httptest.NewRecorder()
		handlerToTest.ServeHTTP(res, req.WithContext(ctx))

		Expect(res.Code).To(Equal(http.StatusInternalServerError))

		Expect(observedLogs.Len()).To(Equal(1))

		loggedEntry := observedLogs.AllUntimed()[0]
		Expect(loggedEntry.Message).To(Equal("request_panic"))
		Expect(loggedEntry.Context).To(HaveLen(2))
		Expect(loggedEntry.Context[0].Key).To(Equal("stacktrace"))
	})

	It("should recover from a panic and log error", func() {
		ctx := context.Background()

		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		observedLogger := zap.New(observedZapCore)
		ctx = zapctx.WithLogger(ctx, observedLogger)

		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

		handlerToTest := RecoveryHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(errors.New("hello darkness my old friend"))
			}),
		)

		res := httptest.NewRecorder()
		handlerToTest.ServeHTTP(res, req.WithContext(ctx))

		Expect(res.Code).To(Equal(http.StatusInternalServerError))

		Expect(observedLogs.Len()).To(Equal(1))

		loggedEntry := observedLogs.AllUntimed()[0]
		Expect(loggedEntry.Message).To(Equal("request_panic"))
		Expect(loggedEntry.Context).To(HaveLen(2))
		Expect(loggedEntry.Context[0].Key).To(Equal("stacktrace"))
		Expect(loggedEntry.Context[1].Key).To(Equal("error"))
	})
})
