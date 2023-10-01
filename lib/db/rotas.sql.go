// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: rotas.sql

package db

import (
	"context"
)

const findRotaByID = `-- name: FindRotaByID :one
SELECT rotas.id, rotas.team_id, rotas.channel_id, rotas.name, rotas.metadata, rotas.created_at, rotas.updated_at
FROM ROTAS
WHERE ID = $1
`

func (q *Queries) FindRotaByID(ctx context.Context, id string) (Rota, error) {
	row := q.db.QueryRow(ctx, findRotaByID, id)
	var i Rota
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.ChannelID,
		&i.Name,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listRotasByChannel = `-- name: ListRotasByChannel :many
SELECT rotas.id, rotas.team_id, rotas.channel_id, rotas.name, rotas.metadata, rotas.created_at, rotas.updated_at
from ROTAS
WHERE ROTAS.CHANNEL_ID = $1
  AND ROTAS.TEAM_ID = $2
`

type ListRotasByChannelParams struct {
	ChannelID string `json:"channel_id"`
	TeamID    string `json:"team_id"`
}

func (q *Queries) ListRotasByChannel(ctx context.Context, arg ListRotasByChannelParams) ([]Rota, error) {
	rows, err := q.db.Query(ctx, listRotasByChannel, arg.ChannelID, arg.TeamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Rota
	for rows.Next() {
		var i Rota
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.ChannelID,
			&i.Name,
			&i.Metadata,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveRota = `-- name: SaveRota :one
INSERT INTO ROTAS (TEAM_ID, CHANNEL_ID, NAME, METADATA)
VALUES ($1, $2, $3, $4) RETURNING ID
`

type SaveRotaParams struct {
	TeamID    string       `json:"team_id"`
	ChannelID string       `json:"channel_id"`
	Name      string       `json:"name"`
	Metadata  RotaMetadata `json:"metadata"`
}

func (q *Queries) SaveRota(ctx context.Context, arg SaveRotaParams) (string, error) {
	row := q.db.QueryRow(ctx, saveRota,
		arg.TeamID,
		arg.ChannelID,
		arg.Name,
		arg.Metadata,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const updateRota = `-- name: UpdateRota :one
UPDATE ROTAS
SET NAME       = $1,
    METADATA   = $2
WHERE ID = $3
RETURNING ID
`

type UpdateRotaParams struct {
	Name     string       `json:"name"`
	Metadata RotaMetadata `json:"metadata"`
	ID       string       `json:"id"`
}

func (q *Queries) UpdateRota(ctx context.Context, arg UpdateRotaParams) (string, error) {
	row := q.db.QueryRow(ctx, updateRota, arg.Name, arg.Metadata, arg.ID)
	var id string
	err := row.Scan(&id)
	return id, err
}
