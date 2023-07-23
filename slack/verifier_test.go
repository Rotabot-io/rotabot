package slack

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Verifier", func() {
	var ctx context.Context
	var verifier RequestVerifierMiddleware

	BeforeEach(func() {
		ctx = context.Background()
		verifier = RequestVerifierMiddleware{}
	})

	It("Verifier ignores requests for other paths", func() {
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

		handlerToTest := verifier.SlackSignatureVerifyHandler(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			},
		)

		res := httptest.NewRecorder()
		handlerToTest.ServeHTTP(res, req.WithContext(ctx))

		Expect(res.Code).To(Equal(http.StatusOK))
	})

	It("Verifier checks requests under slack path", func() {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/slack/commands", nil)

		handlerToTest := verifier.SlackSignatureVerifyHandler(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			},
		)

		res := httptest.NewRecorder()
		handlerToTest.ServeHTTP(res, req.WithContext(ctx))

		Expect(res.Code).To(Equal(http.StatusUnauthorized))
	})

})
