package views

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/rotabot-io/rotabot/slack/slackclient"

	"github.com/getsentry/sentry-go"

	gen "github.com/rotabot-io/rotabot/gen/slack"

	"github.com/rotabot-io/rotabot/lib/db"
	"github.com/rotabot-io/rotabot/lib/goaerrors"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"github.com/rotabot-io/rotabot/slack/block"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

// HomeAction defines the list of possible actions that be taken on the home view
type HomeAction string

// HomeSection defines the list of possible sections that can be interacted with on the home view
type HomeSection string

const (
	HASaveRota = HomeAction("HOME_SAVE_ROTA")

	HSHomeActions = HomeSection("HOME_ACTIONS")
	HSRota        = HomeSection("ROTA_ELEMENT")
)

type Home struct {
	Queries *db.Queries
	State   *HomeState
}

type HomeState struct {
	TriggerID string
	ChannelID string
	TeamID    string
	action    HomeAction
	rotaID    string
}

type HomeProps struct {
	title  *slack.TextBlockObject
	blocks slack.Blocks
}

func (v Home) CallbackID() ViewType {
	return VTHome
}

func (v Home) DefaultState() interface{} {
	return &HomeState{}
}

func (v Home) BuildProps(ctx context.Context) (interface{}, error) {
	l := zapctx.Logger(ctx)
	rotas, err := v.Queries.ListRotasByChannel(ctx, db.ListRotasByChannelParams{ChannelID: v.State.ChannelID, TeamID: v.State.TeamID})
	if err != nil {
		l.Error("failed to list rotas", zap.Error(err))
		return nil, goaerrors.NewInternalError()
	}

	blocks := []slack.Block{
		slack.NewActionBlock(
			string(HSHomeActions),
			block.NewButton(block.Button{Text: "Add Rota :heavy_plus_sign:", ActionID: string(HASaveRota)}),
		),
		block.NewHeader("Active Rotas:"),
	}
	for _, rota := range rotas {
		blocks = append(
			blocks, block.NewOverflowSectionElement(
				block.OverflowSection{
					ElementID:   rota.ID,
					ElementName: rota.Name,
					SectionName: string(HSRota),
					Actions: []block.OverflowAction{
						{Name: ":spiral_note_pad: Edit Rota", Action: string(HASaveRota)},
					},
				},
			),
		)
	}
	return &HomeProps{
		title:  block.NewDefaultText("Rotabot Home"),
		blocks: slack.Blocks{BlockSet: blocks},
	}, nil
}

func (v Home) OnAction(ctx context.Context) (*gen.ActionResponse, error) {
	switch v.State.action {
	case HASaveRota:
		return handleAddRotaAction(ctx, v)
	default:
		zapctx.Logger(ctx).Warn("unknown_action", zap.String("action", string(v.State.action)))
		sentry.CaptureMessage("unknown_action")
		return nil, errors.New("unknown_action")
	}
}

func (v Home) OnClose(ctx context.Context) (*gen.ActionResponse, error) {
	zapctx.Logger(ctx).Debug("closing_home_view")
	return &gen.ActionResponse{}, nil
}

func (v Home) OnSubmit(ctx context.Context) (*gen.ActionResponse, error) {
	zapctx.Logger(ctx).Error("submitting_home_view")
	return nil, goaerrors.NewInternalError()
}

func (v Home) Render(ctx context.Context, p interface{}) error {
	l := zapctx.Logger(ctx)
	props, ok := p.(*HomeProps)
	if !ok {
		zapctx.Logger(ctx).Error("received_invalid_props")
		return errors.New("received invalid props")
	}

	client, err := slackclient.ClientFor(ctx, v.State.TeamID)
	if err != nil {
		l.Error("failed_to_get_client", zap.Error(err))
		sentry.CaptureException(err)
		return err
	}

	bytes, err := json.Marshal(Metadata{RotaID: v.State.rotaID, ChannelID: v.State.ChannelID})
	if err != nil {
		l.Error("failed_to_marshal_metadata", zap.Error(err))
		return err
	}

	_, err = client.OpenViewContext(
		ctx,
		v.State.TriggerID,
		slack.ModalViewRequest{
			Type:            slack.VTModal,
			Title:           props.title,
			Blocks:          props.blocks,
			CallbackID:      string(v.CallbackID()),
			NotifyOnClose:   true,
			ClearOnClose:    true,
			PrivateMetadata: string(bytes),
		},
	)
	return err
}

func handleAddRotaAction(ctx context.Context, v Home) (*gen.ActionResponse, error) {
	l := zapctx.Logger(ctx)
	view := SaveRota{
		Queries: v.Queries,
	}
	view.State = view.DefaultState().(*SaveRotaState)
	view.State.ChannelID = v.State.ChannelID
	view.State.TeamID = v.State.TeamID
	view.State.rotaID = v.State.rotaID

	p, err := view.BuildProps(ctx)
	if err != nil {
		l.Error("failed to build props", zap.Error(err))
		return nil, errors.New("failed to build add rota props")
	}
	props, ok := p.(*SaveRotaProps)
	if !ok {
		l.Error("received_invalid_props")
		return nil, errors.New("received invalid props")
	}

	client, err := slackclient.ClientFor(ctx, v.State.TeamID)
	if err != nil {
		l.Error("failed_to_get_client", zap.Error(err))
		sentry.CaptureException(err)
		return nil, err
	}

	bytes, err := json.Marshal(Metadata{RotaID: v.State.rotaID, ChannelID: v.State.ChannelID})
	if err != nil {
		l.Error("failed_to_marshal_metadata", zap.Error(err))
		return nil, err
	}

	_, err = client.PushViewContext(ctx, v.State.TriggerID, slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           props.title,
		Close:           props.close,
		Submit:          props.submit,
		Blocks:          props.blocks,
		CallbackID:      string(view.CallbackID()),
		NotifyOnClose:   true,
		ClearOnClose:    true,
		PrivateMetadata: string(bytes),
	})
	response := string(slack.RAClear)
	return &gen.ActionResponse{ResponseAction: &response}, err
}
