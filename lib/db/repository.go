package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
)

var (
	ErrAlreadyExists = errors.New("resource already exist")
	ErrNotFound      = errors.New("no rows in result set")
)

type CreateOrUpdateRotaParams struct {
	RotaID    string
	TeamID    string
	ChannelID string
	Name      string
	Metadata  RotaMetadata
}

func CreateOrUpdateRota(ctx context.Context, tx pgx.Tx, p CreateOrUpdateRotaParams) (string, error) {
	l := zapctx.Logger(ctx)
	client := New(tx)
	var rotaId string
	var err error
	if p.RotaID != "" {
		rotaId, err = client.UpdateRota(ctx, UpdateRotaParams{
			ID:       p.RotaID,
			Name:     p.Name,
			Metadata: p.Metadata,
		})
	} else {
		rotaId, err = client.SaveRota(ctx, SaveRotaParams{
			Name:      p.Name,
			TeamID:    p.TeamID,
			ChannelID: p.ChannelID,
			Metadata:  p.Metadata,
		})
	}
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505":
				return "", ErrAlreadyExists
			}
		}
		l.Error("failed to save rota", zap.Error(err))
		return "", err
	}
	return rotaId, nil
}
