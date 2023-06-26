package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rotabot-io/rotabot/lib/metrics"

	"github.com/rotabot-io/rotabot/lib/cli"

	// Automatically set GOMEMLIMIT to match Linux cgroups(7) memory limit.
	// This will only take effect on Linux environments.
	_ "github.com/KimMachineGun/automemlimit"
	// Automatically set GOMAXPROCS to match Linux container CPU quota.
	// This will only take effect on Linux environments.
	_ "go.uber.org/automaxprocs"
)

var (
	AppName = "rotabot"
	Sha     = "unknown"
	Date    = "unknown"
)

func main() {
	params := &cli.Params{
		Usage:     "SlackApp that makes team rotations easy",
		Command:   rotabotCommand,
		AppName:   AppName,
		Sha:       Sha,
		BuildDate: Date,
	}
	cliApp := cli.New(params)

	rootCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()

	metrics.AppTotal.With(prometheus.Labels{"app_name": AppName, "sha": Sha}).Add(1)
	if err := cliApp.RunContext(rootCtx, os.Args); err != nil {
		log.Fatalln("app errored: ", err.Error())
	}
}
