package service

import (
	"context"
	"errors"
	"time"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/dto/response"
	"inventory-system/internal/model"
	"inventory-system/internal/repository"
	"inventory-system/pkg/utils"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuthService defines the business logic contract for authentication.
type AuthService interface {
	// Login verifies user credentials and generates an authentication token.
	Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error)
}

// authService is the concrete implementation of AuthService.
type authService struct {
	repo   *repository.Repository
	logger *zap.Logger
}

// NewAuthService creates and returns a new instance of AuthService.
func NewAuthService(repo *repository.Repository, logger *zap.Logger) AuthService {
	return &authService{
		repo:   repo,
		logger: logger,
	}
}

// Login handles the core authentication workflow.
// It checks if the user exists, verifies the password, generates a JWT,
// and returns the sanitized user data.
func (s *authService) Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error) {
	// Extract the Request ID from the context for distributed tracing in logs.
	reqID := middleware.GetReqID(ctx)

	s.logger.Info("Attempting login", zap.String("request_id", reqID), zap.String("email", req.Email))

	// 1. Check if a user with the provided email exists in the database.
	user, err := s.repo.User.FindByEmail(ctx, req.Email)
	if err != nil {
		// Log the warning but return a generic error message to prevent email enumeration attacks.
		s.logger.Warn("Login failed: user not found", zap.String("request_id", reqID), zap.String("email", req.Email))
		return nil, errors.New("invalid email or password")
	}

	// 2. Verify if the provided plaintext password matches the hashed password in the database.
	isValid := utils.CheckPasswordHash(req.Password, user.PasswordHash)
	if !isValid {
		s.logger.Warn("Login failed: invalid password", zap.String("request_id", reqID), zap.String("email", req.Email))
		return nil, errors.New("invalid email or password")
	}

	// 3. Generate a Stateful UUID Session for the authenticated user.
	sessionID := uuid.New()
	expiredAt := time.Now().Add(24 * time.Hour)

	newSession := &model.Session{
		BaseSimple: model.BaseSimple{ID: sessionID},
		UserID:     user.ID,
		Role:       user.Role,
		ExpiredAt:  expiredAt,
	}

	err = s.repo.Session.CreateSession(ctx, newSession)
	if err != nil {
		s.logger.Error("System Error: Failed to save session to DB",
			zap.String("request_id", reqID),
			zap.Error(err),
		)
		return nil, errors.New("failed to generate authentication token")
	}

	// 4. Map the database User model to the safe UserResponse DTO.
	// This ensures sensitive data like PasswordHash and DeletedAt are not exposed to the client.
	userResponse := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}

	// 5. Construct and return the final AuthResponse containing the token and user data.
	return &response.AuthResponse{
		AccessToken: sessionID.String(),
		User:        userResponse,
	}, nil
}
