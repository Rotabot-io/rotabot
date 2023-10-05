// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Member struct {
	ID        string           `json:"id"`
	RotaID    string           `json:"rota_id"`
	UserID    string           `json:"user_id"`
	Metadata  MemberMetadata   `json:"metadata"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type Rota struct {
	ID        string           `json:"id"`
	TeamID    string           `json:"team_id"`
	ChannelID string           `json:"channel_id"`
	Name      string           `json:"name"`
	Metadata  RotaMetadata     `json:"metadata"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}
