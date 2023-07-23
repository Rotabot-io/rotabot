package cli

import (
	"fmt"
	"runtime"

	"github.com/rotabot-io/rotabot/lib/zapctx"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type Params struct {
	AppName   string
	Usage     string
	BuildDate string
	Sha       string
	Command   *cli.Command
}

func New(p *Params) *cli.App {
	cli.VersionPrinter = func(c *cli.Context) {
		_, _ = fmt.Printf(
			"Application: %s\nSha: %s\nGo Version: %v\nGo OS/Arch: %v/%v\nBuilt at: %v\n",
			p.AppName, p.Sha, runtime.Version(), runtime.GOOS, runtime.GOARCH, p.BuildDate,
		)
	}

	return &cli.App{
		EnableBashCompletion: true,
		Name:                 p.AppName,
		Usage:                p.Usage,
		Version:              p.Sha,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "This will enable all possible logging, this means that sensitive information will be logged, so  avoid using this in production",
				Value: false,
			},
			&cli.StringFlag{
				Name:        "log-format",
				Usage:       "[json | pretty | logfmt] change the output format of log entries (pretty and logfmt loggers are slower than json)",
				EnvVars:     []string{"ROTABOT_LOG_FORMAT"},
				DefaultText: "json",
			},
			&cli.BoolFlag{
				Name:  "sentry",
				Usage: "Whether to configure sentry or not",
				Value: false,
			},
			&cli.StringFlag{
				Name:     "sentry.dsn",
				Usage:    "Secret that allows rotabot to connect to sentry",
				Required: false,
				EnvVars:  []string{"ROTABOT_SENTRY_DSN"},
			},
		},
		Commands: []*cli.Command{p.Command},
		Before: func(ctx *cli.Context) error {
			cfg := zapctx.DefaultLoggerConfig()
			switch ctx.String("log-format") {
			case "pretty":
				cfg.Encoding = "console"
			case "logfmt":
				cfg.Encoding = "logfmt"
			default:
				// default encoding in the config is json
			}

			zl, err := cfg.Build()
			if err != nil {
				return err
			}

			zl = zl.With(
				zap.String("application", p.AppName),
				zap.String("sha", p.Sha),
			)

			zap.RedirectStdLog(zl)
			zap.ReplaceGlobals(zl)

			zl.Info("Logger initialized", zap.String("level", cfg.Level.String()))
			ctx.Context = zapctx.WithLogger(ctx.Context, zl)
			return nil
		},
	}
}
