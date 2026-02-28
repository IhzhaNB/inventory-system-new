package service

import (
	"inventory-system/internal/repository"

	"go.uber.org/zap"
)

type Service struct {
	Auth AuthService
	User UserService
}

func NewService(repo *repository.Repository, logger *zap.Logger) *Service {
	return &Service{
		Auth: NewAuthService(repo, logger),
		User: NewUserService(repo, logger),
	}
}
