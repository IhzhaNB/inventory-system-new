package handler

import (
	"inventory-system/internal/service"

	"go.uber.org/zap"
)

type Handler struct {
	Auth AuthHandler
	User UserHandler
}

func NewHandler(service *service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		Auth: *NewAuthHandler(service.Auth, logger),
		User: *NewUserHandler(service.User, logger),
	}
}
