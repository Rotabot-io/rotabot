package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/getsentry/sentry-go"
	"github.com/rotabot-io/rotabot/lib/db"
	"github.com/urfave/cli/v2"

	"github.com/prometheus/client_golang/prometheus"
	httpSlack "github.com/rotabot-io/rotabot/gen/http/slack/server"
	genSlack "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/lib/metrics"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"github.com/rotabot-io/rotabot/slack"
	"go.uber.org/zap"
	goahttp "goa.design/goa/v3/http"
)

func provideGoaMux(
	ctx context.Context,
	slackService genSlack.Service,
) goahttp.Muxer {
	mux := goahttp.NewMuxer()
	l := zapctx.Logger(ctx)

	slackSvr := slack.NewServer(mux, slackService)
	httpSlack.Mount(mux, slackSvr)
	for _, m := range slackSvr.Mounts {
		initMetrics(m.Verb, m.Pattern)
		l.Info("mounts",
			zap.String("verb", m.Verb),
			zap.String("path", m.Pattern),
			zap.String("method", m.Method),
		)
	}

	return mux
}

func initMetrics(verb, path string) {
	endpoint := metrics.Endpoint(verb, path)
	metrics.RequestsTotal.With(prometheus.Labels{"endpoint": endpoint}).Add(0)

	metrics.RequestDuration.With(prometheus.Labels{"endpoint": endpoint, "status": "200"}).Observe(0)
	metrics.RequestDuration.With(prometheus.Labels{"endpoint": endpoint, "status": "200"}).Observe(0)

	metrics.ResponsesTotal.With(prometheus.Labels{"endpoint": endpoint, "status": "200"}).Add(0)
	metrics.ResponsesTotal.With(prometheus.Labels{"endpoint": endpoint, "status": "200"}).Add(0)

	metrics.PanicsTotal.With(prometheus.Labels{"endpoint": endpoint}).Add(0)

	metrics.AppTotal.With(prometheus.Labels{"app_name": AppName, "sha": Sha}).Add(0)
}

func provideConnString(c *cli.Context) (string, error) {
	if c.Bool("dev") {
		container, err := postgres.RunContainer(c.Context,
			testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		)
		if err != nil {
			return "", err
		}

		dbUrl, err := container.ConnectionString(c.Context, "sslmode=disable")
		if err != nil {
			return "", err
		}

		return dbUrl, nil
	}
	if dsn := c.String("database.url"); dsn != "" {
		return dsn, nil
	}
	return "", errors.New("provideConnString not found")
}

func provideQueries(ctx context.Context, dsn string) (*db.Queries, error) {
	logger := zapctx.Logger(ctx)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("unable to connect to database", zap.Error(err))
		return nil, err
	}
	return db.New(pool), nil
}

func provideSlackService(_ context.Context, q *db.Queries) genSlack.Service {
	return slack.New(q)
}

func provideSentry(ctx context.Context, c *cli.Context) error {
	if c.Bool("sentry") {
		logger := zapctx.Logger(ctx)
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              c.String("sentry.dsn"),
			TracesSampleRate: 1.0,
			Debug:            c.Bool("verbose"),
			Release:          fmt.Sprintf("%s-%s", AppName, Sha),
		})
		if err != nil {
			logger.Error("unable to configure sentry", zap.Error(err))
			return err
		}
		logger.Debug("connected to sentry")

		// Flush buffered events before the program terminates.
		defer sentry.Flush(2 * time.Second)
	}
	return nil
}
