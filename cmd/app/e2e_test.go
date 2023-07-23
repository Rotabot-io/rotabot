package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotabot-io/rotabot/internal"
	"github.com/rotabot-io/rotabot/lib/db"
	"github.com/rotabot-io/rotabot/slack"
)

var _ = Describe("E2E", func() {
	var ctx context.Context
	var cancel context.CancelFunc

	var server *Server

	var httpPort string

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		httpListener, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())

		metricListener, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())

		container, err := internal.RunContainer(ctx)
		Expect(err).ToNot(HaveOccurred())

		err = db.Migrate(ctx, container.ConnectionString())
		Expect(err).ToNot(HaveOccurred())

		conn, err := pgx.Connect(ctx, container.ConnectionString())
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			cancel()
			httpListener.Close()
			metricListener.Close()
			conn.Close(ctx)
		})

		server = NewServer(&ServerParams{
			BaseContext:      ctx,
			AppComponent:     "backend",
			MetricsComponent: "metrics",

			SlackService: slack.New(
				&slack.Config{
					ClientSecret:  "TEST",
					SigningSecret: "TEST",
				},
				db.New(conn),
			),

			HttpListener:    httpListener,
			MetricsListener: metricListener,
		})

		errc := make(chan error, 1)
		go func() {
			errc <- server.Run()
		}()
		Eventually(errc).ShouldNot(Receive())

		httpPort = httpListener.Addr().String()
	})

	It("Healthcheck should return 200", func() {
		Eventually(func() interface{} {
			u := url.URL{
				Scheme: "http",
				Host:   httpPort,
				Path:   "/api/health_check",
			}

			res, err := http.Get(u.String())
			Expect(err).NotTo(HaveOccurred())

			return res.StatusCode
		}).Should(Equal(404)) // TODO: Fix this when we have a healthcheck
	})

	It("Running app twice fails", func() {
		errc := make(chan error, 1)
		go func() {
			errc <- server.Run()
		}()

		go func() {
			errc <- server.Run()
		}()

		err := errors.New("server already running")
		Eventually(errc).Should(Receive(MatchError(err)))
	})
})
