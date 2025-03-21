package main

import (
	"context"
	"fmt"
	"nora/internal/config"
	"nora/internal/handler"
	"nora/internal/repository"
	"nora/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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
	defer dbPool.Close()

	db := stdlib.OpenDB(*dbPool.Config().ConnConfig)
	defer db.Close()
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepository(dbPool)
	taskRepo := repository.NewTaskRepository(dbPool)
	taskGroupRepo := repository.NewTaskGroupRepository(dbPool)
	spaceRepo := repository.NewSpaceRepository(dbPool)
	userSpaceRepo := repository.NewUserSpaceRepository(dbPool)

	s := service.New(userRepo, taskRepo, taskGroupRepo, spaceRepo, userSpaceRepo, cfg)
	h := handler.New(logger, s, cfg)

	var serverAddr = fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	h.Listen(serverAddr)
}
