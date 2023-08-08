package views

import (
	"context"
	"encoding/json"

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

func (m Metadata) ToJson() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (m Metadata) FromJson(payload string) error {
	bytes := []byte(payload)
	return json.Unmarshal(bytes, &m)
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
