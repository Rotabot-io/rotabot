package views

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/rotabot-io/rotabot/slack/slackclient"

	gen "github.com/rotabot-io/rotabot/gen/slack"
	"go.uber.org/zap"

	"github.com/rotabot-io/rotabot/slack/block"

	"github.com/rotabot-io/rotabot/lib/db"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"github.com/slack-go/slack"
)

type SaveRota struct {
	Queries *db.Queries
	State   *SaveRotaState
}

type SaveRotaState struct {
	TriggerID      string
	ChannelID      string
	TeamID         string
	rotaID         string
	rotaName       string
	frequency      db.RotaFrequency
	schedulingType db.RotaSchedule
	externalID     string
	previousViewID string
}

type SaveRotaProps struct {
	title  *slack.TextBlockObject
	submit *slack.TextBlockObject
	close  *slack.TextBlockObject
	blocks slack.Blocks
}

func (v SaveRota) CallbackID() ViewType {
	return VTSaveRota
}

func (v SaveRota) DefaultState() interface{} {
	return &SaveRotaState{
		frequency:      db.RFWeekly,
		schedulingType: db.RSCreated,
	}
}

func (v SaveRota) BuildProps(ctx context.Context) (interface{}, error) {
	var title *slack.TextBlockObject
	var submit *slack.TextBlockObject
	var rotaName string
	var frequency db.RotaFrequency
	var schedulingType db.RotaSchedule

	l := zapctx.Logger(ctx)
	if v.State.rotaID != "" {
		rota, err := v.Queries.FindRotaByID(ctx, v.State.rotaID)
		if err != nil {
			l.Error("failed_to_find", zap.Error(err))
			return nil, err
		}
		rotaName = rota.Name
		frequency = rota.Metadata.Frequency
		schedulingType = rota.Metadata.SchedulingType
		title = block.NewDefaultText("Update Rota")
		submit = block.NewDefaultText("Update")
	} else {
		title = block.NewDefaultText("Create Rota")
		submit = block.NewDefaultText("Create")
		rotaName = v.State.rotaName
		frequency = v.State.frequency
		schedulingType = v.State.schedulingType
	}

	blocks := []slack.Block{
		block.NewTextInput(block.TextInput{
			BlockID: "ROTA_NAME",
			Label:   "Name:",
			Hint:    "e.g. 'On Call'",
			Value:   rotaName,
		}),
		block.NewStaticSelect(block.StaticSelect{
			BlockID:       "ROTA_FREQUENCY",
			Label:         "Frequency:",
			InitialOption: block.StaticSelectOption{Text: string(frequency)},
			Options: []block.StaticSelectOption{
				{Text: string(db.RFDaily)},
				{Text: string(db.RFWeekly)},
				{Text: string(db.RFMonthly)},
			},
		}),
		block.NewStaticSelect(block.StaticSelect{
			BlockID:       "ROTA_TYPE",
			Label:         "Scheduling Type:",
			InitialOption: block.StaticSelectOption{Text: string(schedulingType)},
			Options: []block.StaticSelectOption{
				{Text: string(db.RSCreated)},
				{Text: string(db.RSRandom)},
			},
		}),
	}
	return &SaveRotaProps{
		title:  title,
		submit: submit,
		close:  block.NewDefaultText("Cancel"),
		blocks: slack.Blocks{BlockSet: blocks},
	}, nil
}

func (v SaveRota) OnAction(ctx context.Context) (*gen.ActionResponse, error) {
	zapctx.Logger(ctx).Debug("action_view")
	return &gen.ActionResponse{}, nil
}

func (v SaveRota) OnClose(ctx context.Context) (*gen.ActionResponse, error) {
	zapctx.Logger(ctx).Debug("closing_view")
	return &gen.ActionResponse{}, nil
}

func (v SaveRota) OnSubmit(ctx context.Context) (*gen.ActionResponse, error) {
	l := zapctx.Logger(ctx)
	id, err := saveOrUpdate(ctx, v)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505":
				// TODO: handle this error in a better way
				// duplicate key value violates unique constraint
				response := string(slack.RAErrors)
				return &gen.ActionResponse{
					ResponseAction: &response,
					Errors: map[string]string{
						"ROTA_NAME": "A rota with this name already exists in this channel.",
					},
				}, nil
			}
		}
		l.Error("failed to save rota", zap.Error(err))
		return nil, err
	}
	l.Info("saved_rota", zap.String("id", id))

	h := Home{
		Queries: v.Queries,
		State: &HomeState{
			TriggerID: v.State.TriggerID,
			ChannelID: v.State.ChannelID,
			TeamID:    v.State.TeamID,
		},
	}
	p, err := h.BuildProps(ctx)
	if err != nil {
		l.Error("failed_to_build_home_props", zap.Error(err))
		return nil, err
	}
	props, ok := p.(*HomeProps)
	if !ok {
		l.Error("received_invalid_props")
		return nil, errors.New("received invalid props")
	}

	// Slack does not recommend using the update view API when a modal has been submitted but in our case
	// it's the only way to go back to the home view after creating the new rota.
	// See https://slack.dev/java-slack-sdk/guides/modals
	client, err := slackclient.ClientFor(ctx, v.State.TeamID)
	if err != nil {
		l.Error("failed_to_get_client", zap.Error(err))
		sentry.CaptureException(err)
		return nil, err
	}

	bytes, err := json.Marshal(Metadata{RotaID: id, ChannelID: v.State.ChannelID})
	if err != nil {
		l.Error("failed_to_marshal_metadata", zap.Error(err))
		return nil, err
	}

	r := slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           props.title,
		Blocks:          props.blocks,
		CallbackID:      string(h.CallbackID()),
		NotifyOnClose:   true,
		ClearOnClose:    true,
		PrivateMetadata: string(bytes),
	}
	emptyHash := "" // This is empty to avoid slack thinking this view is outdated (and fail with a hash_conflict error)
	if _, err = client.UpdateViewContext(ctx, r, v.State.externalID, emptyHash, v.State.previousViewID); err != nil {
		l.Error("failed_to_home_view", zap.Error(err))
		return nil, err
	}
	return &gen.ActionResponse{}, nil
}

func (v SaveRota) Render(ctx context.Context, p interface{}) error {
	l := zapctx.Logger(ctx)
	props, ok := p.(*SaveRotaProps)
	if !ok {
		return errors.New("received invalid props")
	}

	bytes, err := json.Marshal(Metadata{ChannelID: v.State.ChannelID})
	if err != nil {
		l.Error("failed_to_marshal_metadata", zap.Error(err))
		return err
	}

	view := slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           props.title,
		Submit:          props.submit,
		Close:           props.close,
		Blocks:          props.blocks,
		CallbackID:      string(v.CallbackID()),
		NotifyOnClose:   true,
		ClearOnClose:    true,
		PrivateMetadata: string(bytes),
	}
	client, err := slackclient.ClientFor(ctx, v.State.TeamID)
	if err != nil {
		l.Error("failed_to_get_client", zap.Error(err))
		sentry.CaptureException(err)
		return err
	}
	_, err = client.OpenViewContext(ctx, v.State.TriggerID, view)
	if err != nil {
		l.Error("failed_to_open_view", zap.Error(err))
		return err
	}
	return nil
}

func saveOrUpdate(ctx context.Context, v SaveRota) (string, error) {
	if v.State.rotaID != "" {
		_, err := v.Queries.UpdateRota(ctx, db.UpdateRotaParams{
			ID:   v.State.rotaID,
			Name: v.State.rotaName,
			Metadata: db.RotaMetadata{
				Frequency:      v.State.frequency,
				SchedulingType: v.State.schedulingType,
			},
		})
		if err != nil {
			return "", err
		}
		return v.State.rotaID, nil
	} else {
		id, err := v.Queries.SaveRota(ctx, db.SaveRotaParams{
			TeamID:    v.State.TeamID,
			ChannelID: v.State.ChannelID,
			Name:      v.State.rotaName,
			Metadata: db.RotaMetadata{
				Frequency:      v.State.frequency,
				SchedulingType: v.State.schedulingType,
			},
		})
		if err != nil {
			return "", err
		}
		return id, nil
	}
}
