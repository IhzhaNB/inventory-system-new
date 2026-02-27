package service

import (
	"inventory-system/internal/config"
	"inventory-system/internal/repository"

	"go.uber.org/zap"
)

type Service struct {
	Auth AuthService
}

func NewService(repo *repository.Repository, config *config.Config, logger *zap.Logger) *Service {
	return &Service{
		Auth: NewAuthService(repo, config.JWT, logger),
	}
}
