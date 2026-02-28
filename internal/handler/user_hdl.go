package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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
// @Description  **Required Roles:** `super_admin`, `admin`
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
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

// GetUsers godoc
// @Summary      Get all users
// @Description  Retrieve a paginated list of users with optional search filtering.
// @Description  **Required Roles:** `super_admin`, `admin`
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page    query     int     false  "Page number for pagination (default: 1)"
// @Param        limit   query     int     false  "Number of items per page (default: 10)"
// @Param        search  query     string  false  "Search filter for user name or email"
// @Success 200 {object} utils.Response{data=response.UserPaginatedResponse} "Users retrieved successfully"
// @Failure      401  {object}  utils.Response "Unauthorized - Invalid or expired session"
// @Failure      500  {object}  utils.Response "Internal server error"
// @Router       /api/v1/users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// 1. Extract values from URL (they are still strings at this stage)
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

	// 2. Convert Strings to Integers (Int)
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// 3. Wrap into a Pagination Request DTO
	query := request.PaginationQuery{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	// 4. Pass the request to the Service layer
	result, err := h.userService.GetUsers(r.Context(), query)
	if err != nil {
		utils.Error(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 5. Return the response to the Client
	utils.Success(w, r, http.StatusOK, "Users retrieved successfully", result)
}
