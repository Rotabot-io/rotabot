package db

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/rotabot-io/rotabot/assets"
	"github.com/rotabot-io/rotabot/lib/zapctx"
	"go.uber.org/zap"
)

func Migrate(ctx context.Context, dsn string) error {
	logger := zapctx.Logger(ctx)
	iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
	if err != nil {
		logger.Error("pulling_db_migrations", zap.Error(err))
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, dsn)
	if err != nil {
		logger.Error("connecting_for_migrations", zap.Error(err))
		return err
	}

	err = migrator.Up()
	switch {
	case errors.Is(err, migrate.ErrNoChange):
		logger.Info("migrations_up_to_date")
	case err != nil:
		logger.Error("failing_to_migrate", zap.Error(err))
		return err
	}
	return nil
}
