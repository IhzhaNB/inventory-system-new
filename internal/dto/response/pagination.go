package response

import "math"

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func NewPaginatedResponse[T any](data []T, page, limit int, totalItems int64) PaginatedResponse[T] {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	totalPages = max(1, totalPages)

	return PaginatedResponse[T]{
		Data: data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: int(totalItems),
			TotalPages: totalPages,
		},
	}
}

// UserPaginatedResponse is a concrete type for Swagger documentation.
// This helps 'swag' parser find the definition easily.
type UserPaginatedResponse PaginatedResponse[UserResponse]
