package main

import (
	"net"

	"github.com/rotabot-io/rotabot/slack"

	"github.com/rotabot-io/rotabot/lib/db"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var rotabotCommand = &cli.Command{
	Name:  "serve",
	Usage: "Starts the rotabot server and its dependencies",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "dev",
			Usage: "This will run all rotabot's dependency in docker containers with rotabot, so  avoid using this in production",
			Value: false,
		},
		&cli.StringFlag{
			Name:        "server.addr",
			Usage:       "Port for the app to listen on",
			DefaultText: ":8080",
			Value:       ":8080",
		},
		&cli.StringFlag{
			Name:        "metrics.addr",
			Usage:       "Port for the app to listen on",
			DefaultText: ":8081",
			Value:       ":8081",
		},
		&cli.StringFlag{
			Name:        "database.url",
			Usage:       "Host on which the database is running",
			DefaultText: "localhost",
			Required:    false,
		},
		&cli.BoolFlag{
			Name:  "migrate",
			Usage: "This run the db migrations automatically",
			Value: true,
		},
		&cli.StringFlag{
			Name:     "slack.signing_secret",
			Usage:    "Secret that ensures the requests from slack are real",
			Required: true,
			EnvVars:  []string{"SLACK_SIGNING_SECRET"},
		},
		&cli.StringFlag{
			Name:     "slack.client_secret",
			Usage:    "Secret that allows app to access the slack API, format: `xoxb-*`",
			Required: true,
			EnvVars:  []string{"SLACK_CLIENT_SECRET"},
		},
	},
	Action: commandAction(),
}

func commandAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		logger := zapctx.Logger(c.Context)
		err := provideSentry(c.Context, c)
		if err != nil {
			logger.Error("unable to setup sentry", zap.Error(err))
			return err
		}

		dbUrl, err := provideConnString(c)
		if err != nil {
			logger.Error("unable to fetch provideConnString from command or ")
			return err
		}

		if err = db.Migrate(c.Context, dbUrl); err != nil {
			logger.Error("unable to run database migrations", zap.Error(err))
			return err
		}

		queries, err := provideQueries(c.Context, dbUrl)
		if err != nil {
			logger.Error("failed to connect to database", zap.Error(err))
			return err
		}

		httpListener, err := net.Listen("tcp", c.String("server.addr"))
		if err != nil {
			logger.Error("failed to start http listener", zap.Error(err))
			return err
		}

		metricListener, err := net.Listen("tcp", c.String("metrics.addr"))
		if err != nil {
			logger.Error("failed to start metrics listener", zap.Error(err))
			return err
		}

		params := &ServerParams{
			BaseContext:      c.Context,
			AppComponent:     "backend",
			MetricsComponent: "metrics",

			SlackConfig: &slack.Config{
				ClientSecret:  c.String("slack.client_secret"),
				SigningSecret: c.String("slack.signing_secret"),
			},

			SlackService: provideSlackService(c.Context, queries, c),

			HttpListener:    httpListener,
			MetricsListener: metricListener,
		}
		if err = NewServer(params).Run(); err != nil {
			logger.Error("failed to run server", zap.Error(err))
			return err
		}
		return nil
	}
}
