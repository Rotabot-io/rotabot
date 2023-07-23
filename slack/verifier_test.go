package slack

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verifier", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	It("Verifier ignores requests for other paths", func() {
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

		handlerToTest := RequestVerifier(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}),
			"SECRET",
		)

		res := httptest.NewRecorder()
		handlerToTest.ServeHTTP(res, req.WithContext(ctx))

		Expect(res.Code).To(Equal(http.StatusOK))
	})

	It("Verifier checks requests under slack path", func() {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/slack/commands", nil)

		handlerToTest := RequestVerifier(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}),
			"SECRET",
		)

		res := httptest.NewRecorder()
		handlerToTest.ServeHTTP(res, req.WithContext(ctx))

		Expect(res.Code).To(Equal(http.StatusUnauthorized))
	})
})
