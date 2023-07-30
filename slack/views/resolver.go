package views

import (
	"context"
	"errors"
	"fmt"

	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"

	"github.com/getsentry/sentry-go"
	"github.com/rotabot-io/rotabot/lib/db"

	"github.com/slack-go/slack"
)

type ResolverParams struct {
	Action  slack.InteractionCallback
	Queries *db.Queries
}

func Resolve(ctx context.Context, p ResolverParams) (View, error) {
	switch p.Action.View.CallbackID {
	case string(VTHome):
		view := &Home{
			Queries: p.Queries,
		}
		view.State = view.DefaultState().(*HomeState)
		view.State.TriggerID = p.Action.TriggerID
		view.State.ChannelID = p.Action.View.PrivateMetadata
		view.State.TeamID = p.Action.Team.ID

		if p.Action.ActionCallback.BlockActions != nil {
			view.State.Action = HomeAction(p.Action.ActionCallback.BlockActions[0].ActionID)
		}

		return view, nil
	case string(VTAddRota):
		view := &AddRota{
			Queries: p.Queries,
		}
		view.State = view.DefaultState().(*AddRotaState)
		view.State.TriggerID = p.Action.TriggerID
		view.State.ChannelID = p.Action.View.PrivateMetadata
		view.State.TeamID = p.Action.Team.ID
		view.State.previousViewID = p.Action.View.PreviousViewID
		view.State.externalID = p.Action.View.ExternalID

		values := p.Action.View.State.Values
		if values != nil {
			view.State.rotaName = values["ROTA_NAME"]["ROTA_NAME"].Value
			view.State.frequency = db.RotaFrequency(values["ROTA_FREQUENCY"]["ROTA_FREQUENCY"].SelectedOption.Value)
			view.State.schedulingType = db.RotaSchedule(values["ROTA_TYPE"]["ROTA_TYPE"].SelectedOption.Value)
		}

		return view, nil
	default:
		zapctx.Logger(ctx).Warn("unknown_callback_id", zap.String("callback_id", p.Action.View.CallbackID))
		sentry.CaptureMessage(fmt.Sprintf("unknown_callback_id: %s", p.Action.View.CallbackID))
		return nil, errors.New("unknown_callback_id")
	}
}
