package repository

import (
	"context"
	"errors"

	"inventory-system/internal/model"

	"github.com/jackc/pgx/v5"
)

// UserRepository defines the contract for user database operations.
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// Add other methods like Create, FindByID, etc., later as needed
}

// userRepository is the concrete implementation of UserRepository.
type userRepository struct {
	db PgxIface // Injecting the interface we created in Epic 1!
}

// NewUserRepository creates and returns a new UserRepository instance.
func NewUserRepository(db PgxIface) UserRepository {
	return &userRepository{db: db}
}

// FindByEmail searches for an active user by their email address.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	// The query strictly ignores soft-deleted users (deleted_at IS NULL)
	query := `
		SELECT id, name, email, password_hash, role, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user model.User

	// Execute the query and scan the result into the user struct
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		// If no matching row is found, return a clear error
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		// Return any other database errors (e.g., connection lost)
		return nil, err
	}

	return &user, nil
}
