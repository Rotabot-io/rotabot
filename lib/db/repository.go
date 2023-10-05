//go:generate mockgen -package mock_db -destination=mock_db/db.go . Repository
package db

import (
	"context"
	"errors"
)

var (
	ErrAlreadyExists = errors.New("resource already exist")
	ErrNotFound      = errors.New("no rows in result set")
)

// RotaSchedule is the type that defines how the members of a rota are scheduled
type RotaSchedule string

// RotaFrequency is the type that defines how long a rota lasts
type RotaFrequency string

const (
	RFDaily   = RotaFrequency("Daily")
	RFWeekly  = RotaFrequency("Weekly")
	RFMonthly = RotaFrequency("Monthly")

	RSCreated = RotaSchedule("Created At")
	RSRandom  = RotaSchedule("Randomly")
)

type RotaMetadata struct {
	Frequency      RotaFrequency `json:"frequency"`
	SchedulingType RotaSchedule  `json:"scheduling_type"`
}

type MemberMetadata struct{}

type Repository interface {
	CreateOrUpdateRota(ctx context.Context, p CreateOrUpdateRotaParams) (string, error)
	FindRotaByID(ctx context.Context, id string) (Rota, error)
	ListRotasByChannel(ctx context.Context, args ListRotasByChannelParams) ([]Rota, error)
	ListUserIDsByRotaID(ctx context.Context, rotaID string) ([]string, error)
}
