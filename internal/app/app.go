package app

import (
	"errors"
	grpcapp "github.com/cuongpiger/sso/internal/app/grpc"
	"github.com/cuongpiger/sso/internal/config"
	"github.com/cuongpiger/sso/internal/lib/migrator"
	"github.com/cuongpiger/sso/internal/services/apps"
	"github.com/cuongpiger/sso/internal/services/auth"
	perm "github.com/cuongpiger/sso/internal/services/permissions"
	"github.com/cuongpiger/sso/internal/storage/postgres"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
)

type App struct {
	GRPCSrv *grpcapp.App
	Storage *postgres.Storage
}

func New(log *slog.Logger, cfg *config.Config) *App {
	storage, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}
	err = migrator.Migrate(cfg)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Debug("no migrations to apply")
		} else {
			panic(err)
		}
	}
	log.Debug("migrations applied successfully")

	authServer := auth.New(log, storage, storage, storage, cfg.TokenTTL)
	permServer := perm.New(log, storage, storage)
	appsServer := apps.New(log, storage, storage, storage, storage)

	grpcApp := grpcapp.New(log, authServer, permServer, appsServer, cfg.GRPC.Host, cfg.GRPC.Port)
	return &App{GRPCSrv: grpcApp, Storage: storage}
}
