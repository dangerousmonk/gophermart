package service

import (
	"github.com/dangerousmonk/gophermart/cmd/config"
	"github.com/dangerousmonk/gophermart/internal/repository"
)

type GophermartService struct {
	Repo repository.Repository
	Cfg  *config.Config
}

func NewGophermartService(r repository.Repository, cfg *config.Config) *GophermartService {
	s := GophermartService{Repo: r, Cfg: cfg}
	return &s
}
