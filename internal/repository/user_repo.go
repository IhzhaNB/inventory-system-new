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
	Create(ctx context.Context, user *model.User) error
	CountUsers(ctx context.Context, search string) (int64, error)
	FindAllUsers(ctx context.Context, limit, offset int, search string) ([]*model.User, error)
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

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
	)
	return err
}

func (r *userRepository) CountUsers(ctx context.Context, search string) (int64, error) {
	query := `SELECT COUNT(id) FROM users WHERE name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%'`
	var total int64
	err := r.db.QueryRow(ctx, query, search).Scan(&total)
	return total, err
}

func (r *userRepository) FindAllUsers(ctx context.Context, limit, offset int, search string) ([]*model.User, error) {
	query := `
		SELECT id, name, email, role
		FROM users
		WHERE name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%'
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}
