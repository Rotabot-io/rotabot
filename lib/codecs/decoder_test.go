package codecs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var _ = Describe("Decoder", func() {
	var observedZapCore zapcore.Core
	var observedLogger *observer.ObservedLogs
	var ctx context.Context

	BeforeEach(func() {
		observedZapCore, observedLogger = observer.New(zapcore.InfoLevel)
		ctx = zapctx.WithLogger(
			context.Background(),
			zap.New(observedZapCore),
		)
	})

	It("should decode valid json without errors", func() {
		req := httptest.NewRequest(http.MethodGet, "/hello", strings.NewReader(`{"foo": "bar"}`))
		req = req.WithContext(ctx)
		rd := RequestDecoderWithLogs(req)

		var target map[string]interface{}
		err := rd.Decode(&target)
		Expect(err).ToNot(HaveOccurred())

		Expect(target).To(Equal(map[string]interface{}{
			"foo": "bar",
		}))

		Expect(observedLogger.All()).To(HaveLen(0))
	})

	It("should decode valid FormRequest without errors", func() {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("a=1&b=2"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(ctx)
		rd := RequestDecoderWithLogs(req)

		var target map[string]interface{}
		err := rd.Decode(&target)
		Expect(err).ToNot(HaveOccurred())

		Expect(target).To(Equal(map[string]interface{}{
			"b": "2",
			"a": "1",
		}))

		Expect(observedLogger.All()).To(HaveLen(0))
	})

	It("should fail to decode invalid json", func() {
		req := httptest.NewRequest(http.MethodGet, "/hello", strings.NewReader(`{"foo"`))
		req = req.WithContext(ctx)
		rd := RequestDecoderWithLogs(req)

		var target map[string]interface{}
		err := rd.Decode(&target)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unexpected EOF"))

		Expect(observedLogger.All()).To(HaveLen(1))
	})
})
