package main

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

var _ = Describe("Server", func() {
	It("should execute middlewares in the correct order", func() {
		// We want to ensure that the order of the middlewares is the correct one.
		// This ensures that if a panic happens then we get the access logs afterward
		ctx := context.Background()

		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		observedLogger := zap.New(observedZapCore)
		ctx = zapctx.WithLogger(ctx, observedLogger.With(zap.String("Test", "flag")))

		req := httptest.NewRequest(http.MethodGet, "http://testing/foo", nil)

		h := wireUpMiddlewares(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				zapctx.Logger(r.Context()).Info("About to panic!")
				panic("Test")
			}),
		)

		res := httptest.NewRecorder()
		h.ServeHTTP(res, req.WithContext(ctx))

		Expect(res.Code).To(Equal(http.StatusInternalServerError))

		Expect(observedLogs.Len()).To(Equal(4))
		Expect(observedLogs.AllUntimed()[0].Message).To(Equal("request.start"))

		// This ensures that our baseContext is injected properly
		Expect(observedLogs.AllUntimed()[0].Context[0].Key).To(Equal("Test"))
		Expect(observedLogs.AllUntimed()[0].Context[0].String).To(Equal("flag"))

		Expect(observedLogs.AllUntimed()[1].Message).To(Equal("About to panic!"))
		Expect(observedLogs.AllUntimed()[2].Message).To(Equal("request_panic"))
		Expect(observedLogs.AllUntimed()[3].Message).To(Equal("request.finish"))
	})
})
