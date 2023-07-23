package slack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

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

	It("Verifier reads headers and body", func() {
		validSigningSecret := "e6b19c573432dcc6b075501d51b51bb8"
		validBody := `{"token":"aF5ynEYQH0dFN9imlgcADxDB","team_id":"XXXXXXXXX","api_app_id":"YYYYYYYYY","event":{"type":"app_mention","user":"AAAAAAAAA","text":"<@EEEEEEEEE> hello world","client_msg_id":"477cc591-ch73-a14z-4db8-g0cd76321bec","ts":"1531431954.000073","channel":"TTTTTTTTT","event_ts":"1531431954.000073"},"type":"event_callback","event_id":"TvBP7LRED7","event_time":1531431954,"authed_users":["EEEEEEEEE"]}`

		req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/slack/commands", strings.NewReader(validBody))
		req.Header.Set("X-Slack-Signature", "v0=adada4ed31709aef585c2580ca3267678c6a8eaeb7e0c1aca3ee57b656886b2c")
		req.Header.Set("X-Slack-Request-Timestamp", "1531431954")

		handlerToTest := RequestVerifier(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}),
			validSigningSecret,
		)

		res := httptest.NewRecorder()
		handlerToTest.ServeHTTP(res, req.WithContext(ctx))

		Expect(res.Code).To(Equal(http.StatusUnauthorized))
	})
})
