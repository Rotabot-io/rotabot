package codecs

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"

	stdHttp "net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"goa.design/goa/v3/http"
)

type badData struct{}

func (b badData) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("err")
}

var _ = Describe("Encoder", func() {
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

	Describe("RequestEncoderWithLogs", func() {
		It("should encode valid json without errors", func() {
			req := httptest.NewRequest(stdHttp.MethodGet, "/hello", strings.NewReader(`{"foo": "bar"}`))
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

		It("should fail to encode with an invalid json", func() {
			req := httptest.NewRequest(stdHttp.MethodGet, "/hello", strings.NewReader(`{"foo"`))
			req = req.WithContext(ctx)
			rd := RequestDecoderWithLogs(req)

			var target map[string]interface{}
			err := rd.Decode(&target)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unexpected EOF"))

			Expect(observedLogger.All()).To(HaveLen(1))
		})
	})

	Describe("ResponseEncoderWithLogs", func() {
		var res *httptest.ResponseRecorder
		var rd http.Encoder

		BeforeEach(func() {
			res = httptest.NewRecorder()
			rd = ResponseEncoderWithLogs(ctx, res)
		})

		It("should encode valid json without errors", func() {
			data := map[string]interface{}{"foo": "bar"}

			err := rd.Encode(data)
			Expect(err).ToNot(HaveOccurred())

			Expect(res.Body.String()).To(Equal("{\"foo\":\"bar\"}\n"))
			Expect(observedLogger.All()).To(HaveLen(0))
		})

		It("should fail to encode with an invalid json", func() {
			err := rd.Encode(badData{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("json: error calling MarshalJSON for type codecs.badData: err"))

			Expect(observedLogger.All()).To(HaveLen(1))
		})
	})
})
