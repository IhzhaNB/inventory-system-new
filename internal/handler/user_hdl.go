package handler

import (
	"encoding/json"
	"net/http"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/service"
	"inventory-system/pkg/utils" // Adjust package name as needed

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.UserService
	logger      *zap.Logger // <-- Added Logger here
}

// NewUserHandler initializes the UserHandler with necessary dependencies.
func NewUserHandler(userService service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger, // <-- Injected here
	}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Register a new user (Admin or Staff) into the system. Requires authentication.
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Insert your token with format: Bearer <token>"
// @Param        request body request.CreateUserRequest true "User data payload"
// @Success      201  {object}  utils.Response{data=response.UserResponse} "User created successfully"
// @Failure      400  {object}  utils.Response "Invalid request payload"
// @Failure      401  {object}  utils.Response "Unauthorized"
// @Failure      500  {object}  utils.Response "Internal server error"
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Extract request ID for distributed tracing in logs
	reqID := middleware.GetReqID(r.Context())

	var req request.CreateUserRequest

	// 1. Decode JSON payload into the request DTO.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode JSON payload", zap.String("request_id", reqID), zap.Error(err))
		utils.Error(w, r, http.StatusBadRequest, "Invalid request payload format", nil)
		return
	}

	// Optional: Add request validation here using go-playground/validator if you have it setup.

	// 2. Pass the parsed request to the Service layer for business logic execution.
	userRes, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		// Log the error returned by the service layer
		h.logger.Error("Service failed to create user", zap.String("request_id", reqID), zap.Error(err))

		// Differentiating client errors (400) vs server errors (500)
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid user role. Must be super_admin, admin, or staff" || err.Error() == "email already exists or database error occurred" {
			statusCode = http.StatusBadRequest
		}

		utils.Error(w, r, statusCode, err.Error(), nil)
		return
	}

	h.logger.Info("User created successfully", zap.String("request_id", reqID), zap.String("email", req.Email))

	// 3. Return a successful 201 Created response.
	utils.Success(w, r, http.StatusCreated, "User created successfully", userRes)
}
