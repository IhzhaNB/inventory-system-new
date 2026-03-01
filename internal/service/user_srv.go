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
	CreateUser(ctx context.Context, req request.CreateUserRequest, requesterRole string) (*response.UserResponse, error)
	GetUsers(ctx context.Context, req request.PaginationQuery) (*response.PaginatedResponse[response.UserResponse], error)
	UpdateUser(ctx context.Context, id uuid.UUID, req request.UpdateUserRequest, requesterRole string) (*response.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID, requesterRole string) error
}

type userService struct {
	repo   *repository.Repository
	logger *zap.Logger
}

func NewUserService(repo *repository.Repository, logger *zap.Logger) UserService {
	return &userService{repo: repo, logger: logger}
}

// CreateUser handles the business logic for registering a new user.
func (s *userService) CreateUser(ctx context.Context, req request.CreateUserRequest, requesterRole string) (*response.UserResponse, error) {
	if requesterRole == string(model.RoleAdmin) && req.Role == string(model.RoleSuperAdmin) {
		s.logger.Warn("Admin attempted to create a super_admin", zap.String("requester_role", requesterRole))
		return nil, errors.New("forbidden: admin cannot create a super_admin")
	}

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
	totalItems, err := s.repo.User.Count(ctx, req.Search)
	if err != nil {
		return nil, errors.New("failed to count users")
	}

	// 4. Query Repo: "Fetch [Limit] records starting from the [Offset] position"
	users, err := s.repo.User.FindAll(ctx, req.Limit, offset, req.Search)
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

// UpdateUser handles the business logic for updating a user's details.
func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req request.UpdateUserRequest, requesterRole string) (*response.UserResponse, error) {
	// 1. Check if the user exists
	user, err := s.repo.User.FindByID(ctx, id)
	if err != nil {
		s.logger.Warn("Attempted to update a non-existent user", zap.String("user_id", id.String()), zap.Error(err))
		return nil, errors.New("user not found")
	}

	// 🛡️ GUARD 1: Admin tidak boleh mengedit data Super Admin
	if requesterRole == string(model.RoleAdmin) && user.Role == model.RoleSuperAdmin {
		s.logger.Warn("Admin attempted to modify a super_admin", zap.String("target_user_id", id.String()))
		return nil, errors.New("forbidden: admin cannot modify a super_admin")
	}

	// 🛡️ GUARD 2: Admin tidak boleh me-naikkan jabatan seseorang menjadi Super Admin
	if requesterRole == string(model.RoleAdmin) && req.Role == string(model.RoleSuperAdmin) {
		s.logger.Warn("Admin attempted to promote someone to super_admin", zap.String("target_user_id", id.String()))
		return nil, errors.New("forbidden: admin cannot promote a user to super_admin")
	}

	// 2. Update the user model with new data
	user.Name = req.Name
	user.Role = model.UserRole(req.Role)

	// 3. Save the changes to the database
	if err := s.repo.User.Update(ctx, user); err != nil {
		s.logger.Error("Database error while updating user", zap.String("user_id", id.String()), zap.Error(err))
		return nil, errors.New("internal server error")
	}

	s.logger.Info("User updated successfully", zap.String("user_id", id.String()))

	// 4. Map the updated model to a safe response DTO
	res := response.ToUserResponse(user)
	return &res, nil
}

// DeleteUser handles the business logic for removing a user.
func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID, requesterRole string) error {
	// 1. Ensure the user exists before attempting to delete
	user, err := s.repo.User.FindByID(ctx, id)
	if err != nil {
		s.logger.Warn("Attempted to delete a non-existent user", zap.String("user_id", id.String()), zap.Error(err))
		return errors.New("user not found")
	}

	// 🛡️ GUARD: Admin tidak boleh menghapus Super Admin
	if requesterRole == string(model.RoleAdmin) && user.Role == model.RoleSuperAdmin {
		s.logger.Warn("Admin attempted to delete a super_admin", zap.String("target_user_id", id.String()))
		return errors.New("forbidden: admin cannot delete a super_admin")
	}

	// 2. Execute the deletion
	if err := s.repo.User.Delete(ctx, id); err != nil {
		s.logger.Error("Database error while deleting user", zap.String("user_id", id.String()), zap.Error(err))
		return errors.New("internal server error")
	}

	s.logger.Info("User deleted successfully", zap.String("user_id", id.String()))

	return nil
}
