package application

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"sso-service/config"
	repositoryImpl "sso-service/internal/repository"
	serviceImpl "sso-service/internal/service"
	"sso-service/internal/service/implementation"
	"sso-service/migrator"
	"time"
)

type App struct {
	service    serviceImpl.ServiceManager
	repository repositoryImpl.RepositoryManager
	cache      *cache.Cache

	config config.Config

	conn *pgxpool.Pool
}

func NewApp(config config.Config) App {
	return App{
		config: config,
	}
}

func (app *App) GetServiceManager() serviceImpl.ServiceManager {
	return app.service
}

func (app *App) Init(ctx context.Context) error {
	logger, err := implementation.NewLogger()

	if err != nil {
		return err
	}

	conn, err := pgxpool.New(ctx, app.config.Db.ConnectionString)
	if err != nil {
		return err
	}

	app.conn = conn

	cache := cache.New(
		time.Minute*time.Duration(app.config.Cache.CodeExpirationTimeInMinutes),
		time.Minute*time.Duration(app.config.Cache.CleanupIntervalInMinutes),
	)

	repository := repositoryImpl.NewRepositoryManager(conn, cache, app.config.Auth.CodeSalt)
	service := serviceImpl.NewServiceManager(logger, app.config, repository)

	app.repository = repository
	app.service = service
	app.cache = cache

	migrator := migrator.NewMigrator(app.config.Db, conn, app.service.Logger)

	err = migrator.Migrate()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Shutdown() {
	app.conn.Close()
}
