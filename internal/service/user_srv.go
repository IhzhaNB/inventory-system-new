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
	GetUsers(ctx context.Context, req request.PaginationQuery) (*response.PaginatedResponse[response.UserResponse], error)
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

func (s *userService) GetUsers(ctx context.Context, req request.PaginationQuery) (*response.PaginatedResponse[response.UserResponse], error) {
	// 1. Set default values if the URL does not provide page or limit
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	// 2. Offset Formula: (Page - 1) * Limit
	// Example: If Page 2 and Limit 10 are requested -> (2-1)*10 = 10. (The database skips the first 10 records)
	offset := (req.Page - 1) * req.Limit

	// 3. Query Repo: "What is the total number of records in the DB?"
	totalItems, err := s.repo.User.CountUsers(ctx, req.Search)
	if err != nil {
		return nil, errors.New("failed to count users")
	}

	// 4. Query Repo: "Fetch [Limit] records starting from the [Offset] position"
	users, err := s.repo.User.FindAllUsers(ctx, req.Limit, offset, req.Search)
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	// 5. Map Database Models to Data Transfer Objects (DTOs)
	// We use the ToUserResponse helper created previously
	var userResponses []response.UserResponse
	for _, u := range users {
		userResponses = append(userResponses, response.ToUserResponse(u))
	}

	// Ensure the slice is initialized so the JSON output is an empty array [] instead of null
	if userResponses == nil {
		userResponses = []response.UserResponse{}
	}

	// 6. Wrap data into the Paginated Response container
	// Note: Generics [response.UserResponse] will automatically adjust based on the type.
	result := response.NewPaginatedResponse(userResponses, req.Page, req.Limit, totalItems)

	return &result, nil
}
