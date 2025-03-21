package service

import (
	"nora/internal/config"
	"nora/internal/repository"
)

type Service struct {
	u   repository.UserRepo
	t   repository.TaskRepo
	tg  repository.TaskGroupRepo
	sp  repository.SpaceRepo
	usp repository.UserSpaceRepo
	cfg *config.Config
}

func New(u repository.UserRepo, t repository.TaskRepo, tg repository.TaskGroupRepo, sp repository.SpaceRepo, usp repository.UserSpaceRepo, config *config.Config) *Service {
	return &Service{u: u, t: t, tg: tg, sp: sp, usp: usp, cfg: config}
}
