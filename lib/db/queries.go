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
		// Check existing member needs to be deleted
		if inx := slices.IndexFunc(members, func(member Member) bool { return member.UserID == userId }); inx == -1 {
			err = q.deleteMember(ctx, userId)
			if err != nil {
				l.Error("unable_to_delete_member",
					zap.Error(err),
					zap.String("user_id", userId),
					zap.Strings("existing", e),
				)
				return err
			}
		}
	}
	for _, m := range members {
		// Check if desired member already exist before adding
		if inx := slices.IndexFunc(e, func(userId string) bool { return m.UserID == userId }); inx == -1 {
			_, err = q.saveMember(ctx, saveMemberParams{
				RotaID:   rotaId,
				UserID:   m.UserID,
				Metadata: m.Metadata,
			})
		}
		if err != nil {
			err = mapError(err)
			l.Error("unable_to_add_member",
				zap.Error(err),
				zap.String("user_id", m.UserID),
				zap.Strings("existing", e),
			)
			return err
		}
		// Avoid trying to add the same user_id twice within the same request
		e = append(e, m.UserID)
	}
	return nil
}

func mapError(err error) error {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		// The magic list of errors can be found here
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		switch pgError.Code {
		case "23505":
			return ErrAlreadyExists
		}
	}
	return err
}
