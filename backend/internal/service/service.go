package service

import (
	"nora/internal/config"
	"nora/internal/repository"
)

type Service struct {
	u   repository.UserRepo
	cfg *config.Config
}

func New(u repository.UserRepo, config *config.Config) *Service {
	return &Service{u: u, cfg: config}
}
