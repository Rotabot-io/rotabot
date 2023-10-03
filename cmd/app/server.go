package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rotabot-io/rotabot/slack"

	"go.uber.org/zap/zapcore"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapio"

	genSlack "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/lib/middleware"
	"github.com/rotabot-io/rotabot/lib/zapctx"
)

type ServerParams struct {
	BaseContext context.Context

	AppComponent     string
	MetricsComponent string

	SlackSigningSecret string
	SlackService       genSlack.Service

	HttpListener    net.Listener
	MetricsListener net.Listener
}

type Server struct {
	group   *run.Group
	ctx     context.Context
	running int32

	Server        *http.Server
	MetricsServer *http.Server
}

func NewServer(params *ServerParams) *Server {
	var group run.Group
	return &Server{
		group:         &group,
		ctx:           params.BaseContext,
		Server:        initHttpServer(params, &group),
		MetricsServer: initMetricsServer(params, &group),
	}
}

func (s *Server) Run() error {
	swapped := atomic.CompareAndSwapInt32(&s.running, 0, 1)
	if !swapped {
		return errors.New("server already running")
	}

	s.group.Add(func() error {
		<-s.ctx.Done()
		return nil
	}, func(error) {
		// noop
	})

	return s.group.Run()
}

func initHttpServer(p *ServerParams, rg *run.Group) *http.Server {
	ctx, cancel := context.WithCancel(p.BaseContext)
	logger := zapctx.Logger(ctx).With(zap.String("component", p.AppComponent))
	ctx = zapctx.WithLogger(ctx, logger)

	mux := provideGoaMux(ctx, p.SlackService)
	srv := &http.Server{
		Handler: wireUpMiddlewares(p, http.Handler(mux)),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		ErrorLog:          zapToStdLog(zapctx.Logger(ctx)),
		ReadTimeout:       100 * time.Millisecond,
		ReadHeaderTimeout: 100 * time.Millisecond,
	}

	rg.Add(func() error {
		logger.Info("starting server", zap.Stringer("address", p.HttpListener.Addr()))
		return srv.Serve(p.HttpListener)
	}, func(error) {
		logger.Info("stopping server")
		cancel()

		if err := srv.Close(); err != nil {
			logger.Error("failed to close server", zap.Error(err))
		}
	})

	return srv
}

func wireUpMiddlewares(p *ServerParams, handler http.Handler) http.Handler {
	handler = slack.RequestVerifier(handler, p.SlackSigningSecret)
	handler = middleware.RecoveryHandler(handler)
	handler = middleware.RequestAccessLogHandler(handler)
	handler = middleware.LoggerInjectionHandler(handler)
	handler = middleware.RequestIdHandler(handler)
	return handler
}

func initMetricsServer(params *ServerParams, rg *run.Group) *http.Server {
	ctx, cancel := context.WithCancel(params.BaseContext)
	logger := zapctx.Logger(ctx).With(zap.String("component", params.MetricsComponent))
	ctx = zapctx.WithLogger(ctx, logger)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		ErrorLog:          zapToStdLog(zapctx.Logger(ctx)),
		ReadTimeout:       100 * time.Millisecond,
		ReadHeaderTimeout: 100 * time.Millisecond,
	}

	rg.Add(func() error {
		logger.Info("starting server", zap.Stringer("address", params.MetricsListener.Addr()))
		return srv.Serve(params.MetricsListener)
	}, func(error) {
		logger.Info("stopping server")
		cancel()

		if err := srv.Close(); err != nil {
			logger.Error("failed to close server", zap.Error(err))
		}
	})

	return srv
}

func zapToStdLog(l *zap.Logger) *log.Logger {
	return log.New(&zapio.Writer{Log: l, Level: zapcore.ErrorLevel}, "", 0)
}
