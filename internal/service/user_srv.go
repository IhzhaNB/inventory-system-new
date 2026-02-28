package service

import (
	"context"
	"errors"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/dto/response"
	"inventory-system/internal/model"
	"inventory-system/internal/repository"
	"inventory-system/pkg/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(ctx context.Context, req request.CreateUserRequest) (*response.UserResponse, error)
}

type userService struct {
	repo   *repository.Repository
	logger *zap.Logger
}

func NewUserService(repo *repository.Repository, logger *zap.Logger) UserService {
	return &userService{repo: repo, logger: logger}
}

// CreateUser handles the business logic for registering a new user.
func (s *userService) CreateUser(ctx context.Context, req request.CreateUserRequest) (*response.UserResponse, error) {
	// 1. Validate the provided Role against allowed enums to ensure data integrity.
	role := model.UserRole(req.Role)
	if role != model.RoleSuperAdmin && role != model.RoleAdmin && role != model.RoleStaff {
		return nil, errors.New("invalid user role. Must be super_admin, admin, or staff")
	}

	// 2. Hash the user's plaintext password securely.
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	// 3. Construct the User model instance.
	newUser := &model.User{
		BaseModel:    model.BaseModel{ID: uuid.New()}, // Adjust based on your actual Base struct
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	// 4. Save the new user to the database via the repository layer.
	err = s.repo.User.Create(ctx, newUser)
	if err != nil {
		s.logger.Error("Failed to insert user to DB", zap.Error(err), zap.String("email", req.Email))
		return nil, errors.New("email already exists or database error occurred")
	}

	// 5. Map the saved model to a safe response DTO, omitting sensitive data.
	resp := response.ToUserResponse(newUser)
	return &resp, nil
}
