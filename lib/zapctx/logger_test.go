package zapctx

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var _ = Describe("Logger", func() {
	It("should panic when nil context is passed", func() {
		Expect(func() {
			Logger(nil) //nolint:staticcheck // This is a test, we want to panic
		}).To(Panic())
	})

	It("should carry the logger through the context", func() {
		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		ctx := WithLogger(
			context.Background(),
			zap.New(observedZapCore).
				With(zap.String("key", "value")),
		)

		l := Logger(ctx)

		l.Info("Send Random Test")

		loggedEntry := observedLogs.AllUntimed()[0]
		Expect(len(loggedEntry.Context)).To(Equal(1))
		Expect(loggedEntry.Context[0].Key).To(Equal("key"))
		Expect(loggedEntry.Context[0].String).To(Equal("value"))
	})

	It("should return a default logger when no logger is set", func() {
		l := Logger(context.Background())

		Expect(l).ToNot(BeNil())
		Expect(l.Core().Enabled(zapcore.InfoLevel)).To(BeTrue())
	})
})
