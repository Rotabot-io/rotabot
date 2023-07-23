package middleware

import (
	"net/http"
	"net/http/httptest"

	uuidGen "github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestIdHandler", func() {
	It("should add a unique request ID to the request context", func() {
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

		handlerToTest := RequestIdHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(w.Header().Get(string(RequestIdKey))).ToNot(BeEmpty())
				_, err := uuidGen.Parse(w.Header().Get(string(RequestIdKey)))
				Expect(err).ToNot(HaveOccurred())
				Expect(r.Context().Value(RequestIdKey)).ToNot(BeEmpty())

				Expect(r.Context().Value(RequestIdKey)).To(Equal(w.Header().Get(string(RequestIdKey))))
			}),
		)

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	})

	It("should reuse UUID if present on request", func() {
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
		req.Header.Set(string(RequestIdKey), "123")

		handlerToTest := RequestIdHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(w.Header().Get(string(RequestIdKey))).ToNot(BeEmpty())
				Expect(r.Context().Value(RequestIdKey)).ToNot(BeEmpty())

				Expect(r.Context().Value(RequestIdKey)).To(Equal(w.Header().Get(string(RequestIdKey))))

				Expect(r.Context().Value(RequestIdKey)).To(Equal("123"))
			}),
		)

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	})
})
