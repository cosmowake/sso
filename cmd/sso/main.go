package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

//go run cmd/sso/main.go --config="./config/local.yaml"

func main() {
	cfg := config.MustLoad()
	fmt.Println("SSO server starting with configuration:", cfg)

	log := setupLogger(cfg.Env)
	log.Info("start application")

	application := app.New(
		log,
		cfg.GRPC.Port,
		cfg.Mongo.Address,
		cfg.TokenTTL,
	)
	log.Info("application initialized")

	log.Info("GRPCServer run")
	go application.GRPCServer.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}
