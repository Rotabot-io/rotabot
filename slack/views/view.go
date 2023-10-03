package views

import (
	"context"

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
	BuildProps(ctx context.Context) (interface{}, error)
	OnAction(ctx context.Context) (*gen.ActionResponse, error)
	OnClose(ctx context.Context) (*gen.ActionResponse, error)
	OnSubmit(ctx context.Context) (*gen.ActionResponse, error)
	Render(ctx context.Context, props interface{}) error
}
