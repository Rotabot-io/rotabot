package db

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
)

type CreateOrUpdateRotaParams struct {
	RotaID    string
	TeamID    string
	ChannelID string
	Name      string
	Metadata  RotaMetadata
	Members   []Member
}

func (q *Queries) CreateOrUpdateRota(ctx context.Context, p CreateOrUpdateRotaParams) (string, error) {
	l := zapctx.Logger(ctx)
	var rotaId string
	var err error
	if p.RotaID != "" {
		rotaId, err = q.updateRota(ctx, updateRotaParams{
			ID:       p.RotaID,
			Name:     p.Name,
			Metadata: p.Metadata,
		})
	} else {
		rotaId, err = q.saveRota(ctx, saveRotaParams{
			Name:      p.Name,
			TeamID:    p.TeamID,
			ChannelID: p.ChannelID,
			Metadata:  p.Metadata,
		})
	}
	if err != nil {
		err = mapError(err)
		l.Error("failed to save rota", zap.Error(err))
		return "", err
	}
	err = q.updateMembersList(ctx, rotaId, p.Members)
	if err != nil {
		err = mapError(err)
		l.Error("failed to save rota members", zap.Error(err))
		return "", err
	}
	return rotaId, nil
}

func (q *Queries) updateMembersList(ctx context.Context, rotaId string, members []Member) error {
	l := zapctx.Logger(ctx)
	e, err := q.ListUserIDsByRotaID(ctx, rotaId)
	if err != nil {
		l.Error("unable_to_fetch_existing_members", zap.Error(err))
		return err
	}
	for _, userId := range e {
		inx := slices.IndexFunc(members, func(member Member) bool { return member.UserID == userId })
		if inx == -1 {
			err = q.deleteMember(ctx, userId)
			if err != nil {
				l.Error("unable_to_fetch_existing_members", zap.Error(err))
				return err
			}
		}
	}
	for _, m := range members {
		_, err = q.saveMember(ctx, saveMemberParams{
			RotaID:   rotaId,
			UserID:   m.UserID,
			Metadata: m.Metadata,
		})
		if err != nil {
			err = mapError(err)
			if errors.Is(err, ErrAlreadyExists) {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func mapError(err error) error {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch pgError.Code {
		case "23505":
			return ErrAlreadyExists
		}
	}
	return err
}
