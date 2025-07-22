package app

import (
	"context"
	"log/slog"
	"sso/internal/services/auth"
	dbmongo "sso/internal/storage/mongo"
	"time"

	grpcapp "sso/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
	db         *dbmongo.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	connectionURI string,
	tokenTTL time.Duration,
) *App {

	storage, err := dbmongo.New(context.TODO(), connectionURI)
	if err != nil {
		log.Error("mongo error", slog.Any("error", err))
		panic(err)
	}

	authService := auth.New(log, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		db:         storage,
		GRPCServer: grpcApp,
	}
}
