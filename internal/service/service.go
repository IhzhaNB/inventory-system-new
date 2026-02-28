package service

import (
	"inventory-system/internal/repository"

	"go.uber.org/zap"
)

type Service struct {
	Auth AuthService
}

func NewService(repo *repository.Repository, logger *zap.Logger) *Service {
	return &Service{
		Auth: NewAuthService(repo, logger),
	}
}
