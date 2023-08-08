package views

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"

	"github.com/getsentry/sentry-go"
	"github.com/rotabot-io/rotabot/lib/db"

	"github.com/slack-go/slack"
)

type ResolverParams struct {
	Action slack.InteractionCallback
}

func Resolve(ctx context.Context, p ResolverParams) (View, error) {
	switch p.Action.View.CallbackID {
	case string(VTHome):
		return resolveHomeView(ctx, p)
	case string(VTSaveRota):
		return resolveSaveRota(ctx, p)
	default:
		zapctx.Logger(ctx).Warn("unknown_callback_id", zap.String("callback_id", p.Action.View.CallbackID))
		sentry.CaptureMessage(fmt.Sprintf("unknown_callback_id: %s", p.Action.View.CallbackID))
		return nil, errors.New("unknown_callback_id")
	}
}

func resolveHomeView(ctx context.Context, p ResolverParams) (View, error) {
	m, err := unMarshallMetadata(ctx, p.Action.View.PrivateMetadata)
	if err != nil {
		return nil, err
	}

	view := &Home{}
	view.State = view.DefaultState().(*HomeState)
	view.State.TriggerID = p.Action.TriggerID
	view.State.TeamID = p.Action.Team.ID
	view.State.ChannelID = m.ChannelID

	if p.Action.ActionCallback.BlockActions != nil {
		blockAction := p.Action.ActionCallback.BlockActions[0]
		if blockAction.ActionID == string(HSRota) {
			view.State.action = HomeAction(blockAction.SelectedOption.Value)
			view.State.rotaID = blockAction.BlockID
		} else {
			view.State.action = HomeAction(p.Action.ActionCallback.BlockActions[0].ActionID)
		}
	}

	return view, nil
}

func resolveSaveRota(ctx context.Context, p ResolverParams) (View, error) {
	m, err := unMarshallMetadata(ctx, p.Action.View.PrivateMetadata)
	if err != nil {
		return nil, err
	}

	view := &SaveRota{}
	view.State = view.DefaultState().(*SaveRotaState)
	view.State.TriggerID = p.Action.TriggerID
	view.State.rotaID = m.RotaID
	view.State.ChannelID = m.ChannelID
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
}

func unMarshallMetadata(ctx context.Context, metadata string) (Metadata, error) {
	var m Metadata
	err := json.Unmarshal([]byte(metadata), &m)
	if err != nil {
		zapctx.Logger(ctx).Error("unmarshall_metadata", zap.Error(err))
		return m, err
	}
	return m, nil
}
