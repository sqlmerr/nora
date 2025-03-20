package main

import (
	"context"
	"fmt"
	"nora/internal/config"
	"nora/internal/handler"
	"nora/internal/repository"
	"nora/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	dbPool, err := pgxpool.New(context.Background(), cfg.ConnectionString())
	userRepo := repository.NewUserRepository(dbPool)

	s := service.New(userRepo, cfg)
	h := handler.New(logger, s, cfg)

	var serverAddr = fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	h.Listen(serverAddr)
}
