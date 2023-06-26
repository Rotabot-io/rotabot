package slack

import (
	"context"

	gen "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/lib/db"
)

type Config struct {
	ClientSecret  string
	SigningSecret string
}

func New(c *Config, q *db.Queries) gen.Service {
	return &svc{
		config:  c,
		queries: q,
	}
}

type svc struct {
	queries *db.Queries
	config  *Config
}

func (s svc) Commands(ctx context.Context, c *gen.Command) error {
	return nil
}

func (s svc) Events(ctx context.Context, event *gen.Event) (*gen.EventResponse, error) {
	if *event.Type == "url_verification" {
		return &gen.EventResponse{
			Challenge: event.Challenge,
		}, nil
	}
	return nil, nil
}

var _ gen.Service = &svc{}
