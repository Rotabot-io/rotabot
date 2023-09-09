package views

import (
	"context"

	"github.com/jackc/pgx/v5"

	gen "github.com/rotabot-io/rotabot/gen/slack"
)

type ViewType string

const (
	VTHome     = ViewType("Home")
	VTSaveRota = ViewType("SaveRota")
)

type Metadata struct {
	RotaID    string `json:"rota_id"`
	ChannelID string `json:"channel_id"`
}

type View interface {
	CallbackID() ViewType
	DefaultState() interface{}
	BuildProps(ctx context.Context, tx pgx.Tx) (interface{}, error)
	OnAction(ctx context.Context, tx pgx.Tx) (*gen.ActionResponse, error)
	OnClose(ctx context.Context, tx pgx.Tx) (*gen.ActionResponse, error)
	OnSubmit(ctx context.Context, tx pgx.Tx) (*gen.ActionResponse, error)
	Render(ctx context.Context, props interface{}) error
}
