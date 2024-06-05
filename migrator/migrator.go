package migrator

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"sso-service/config"
	"sso-service/internal/domain/service"
)

var (
	ErrorNilConnection = errors.New("connection is nil")
	ErrorSetDialect    = errors.New("ошибка установки диалекта")
)

type Migrator struct {
	config config.DatabaseConfig
	db     *pgxpool.Pool
	logger service.Logger
}

func NewMigrator(config config.DatabaseConfig, db *pgxpool.Pool, logger service.Logger) Migrator {
	return Migrator{
		config: config,
		db:     db,
		logger: logger,
	}
}

func (migrator *Migrator) Migrate() error {
	connection := stdlib.OpenDBFromPool(migrator.db)
	if connection == nil {
		return ErrorNilConnection
	}

	if err := goose.SetDialect(migrator.config.DriverName); err != nil {
		migrator.logger.Error(fmt.Sprintf("[Migrator] ошибка установки диалекта %s", err.Error()))
		return err
	}

	if err := goose.Up(connection, migrator.config.MigrationDir); err != nil {
		migrator.logger.Error(fmt.Sprintf("[Migrator] ошибка применения миграций %s", err.Error()))
		return err
	}

	return nil
}

func (migrator *Migrator) Seed(ctx context.Context) error {

	return nil
}
