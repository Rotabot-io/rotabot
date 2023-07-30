package slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/rotabot-io/rotabot/slack/slackclient"
	"github.com/slack-go/slack/slackevents"

	"github.com/slack-go/slack"

	"github.com/rotabot-io/rotabot/lib/goaerrors"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"github.com/rotabot-io/rotabot/slack/views"
	"go.uber.org/zap"

	gen "github.com/rotabot-io/rotabot/gen/slack"
	"github.com/rotabot-io/rotabot/lib/db"
)

func New(q *db.Queries) gen.Service {
	return &svc{
		queries: q,
	}
}

type svc struct {
	queries *db.Queries
}

func (s svc) Commands(ctx context.Context, c *gen.Command) error {
	ctx = zapctx.WithLogger(ctx, zapctx.Logger(ctx).
		With(zap.String("cmd", c.Command)).
		With(zap.String("user_id", c.UserID)).
		With(zap.String("channel_id", c.ChannelID)).
		With(zap.String("team_id", c.TeamID)).
		With(zap.String("trigger_id", c.TriggerID)))
	l := zapctx.Logger(ctx)

	client, err := slackclient.ClientFor(ctx, c.TeamID)
	if err != nil {
		l.Error("failed to get slack client", zap.Error(err))
		return goaerrors.NewInternalError()
	}
	ctx = slackclient.WithClient(ctx, client)

	view := views.Home{
		Queries: s.queries,
		State: &views.HomeState{
			TriggerID: c.TriggerID,
			ChannelID: c.ChannelID,
			TeamID:    c.TeamID,
		},
	}

	p, err := view.BuildProps(ctx)
	if err != nil {
		l.Error("failed to build props", zap.Error(err))
		return goaerrors.NewInternalError()
	}
	props, ok := p.(*views.HomeProps)
	if !ok {
		l.Error("received_invalid_props")
		return errors.New("received invalid props")
	}

	return view.Render(ctx, props)
}

func (s svc) Events(_ context.Context, event *gen.Event) (*gen.EventResponse, error) {
	if event.Type == slackevents.URLVerification {
		return &gen.EventResponse{Challenge: event.Challenge}, nil
	}
	return &gen.EventResponse{}, nil
}

func (s svc) MessageActions(ctx context.Context, event *gen.Action) (*gen.ActionResponse, error) {
	action, err := marshallCallback(ctx, event)
	if err != nil {
		return nil, err
	}
	ctx = zapctx.WithLogger(ctx, zapctx.Logger(ctx).
		With(zap.String("type", string(action.Type))).
		With(zap.String("user_id", action.User.ID)).
		With(zap.String("team_id", action.Team.ID)).
		With(zap.String("callback_id", action.View.CallbackID)))
	l := zapctx.Logger(ctx)

	client, err := slackclient.ClientFor(ctx, action.Team.ID)
	if err != nil {
		l.Error("failed to get slack client", zap.Error(err))
		return nil, goaerrors.NewInternalError()
	}
	ctx = slackclient.WithClient(ctx, client)

	view, err := views.Resolve(ctx, views.ResolverParams{
		Action:  action,
		Queries: s.queries,
	})
	if err != nil {
		l.Error("failed to resolve view", zap.Error(err))
		return nil, goaerrors.NewInternalError()
	}
	// We only use views for now so these are the only events that make sense for us to handle.
	switch action.Type { // nolint:exhaustive
	case slack.InteractionTypeBlockActions:
		return view.OnAction(ctx)
	case slack.InteractionTypeViewSubmission:
		return view.OnSubmit(ctx)
	case slack.InteractionTypeViewClosed:
		return view.OnClose(ctx)
	default:
		sentry.CaptureMessage(fmt.Sprintf("unknown_action_type: %s", action.Type))
		l.Warn("unknown_action_type", zap.String("type", string(action.Type)))
		response := string(slack.RAClear)
		return &gen.ActionResponse{ResponseAction: &response}, nil
	}
}

func marshallCallback(ctx context.Context, event *gen.Action) (slack.InteractionCallback, error) {
	l := zapctx.Logger(ctx)
	var action slack.InteractionCallback
	err := json.Unmarshal(event.Payload, &action)
	if err != nil {
		l.Error("failed to unmarshal action", zap.Error(err))
		return slack.InteractionCallback{}, goaerrors.NewInternalError()
	}
	return action, nil
}

var _ gen.Service = &svc{}
