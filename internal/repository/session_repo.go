package repository

import (
	"context"
	"inventory-system/internal/model"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *model.Session) error
}

type sessionRepository struct {
	db PgxIface
}

func NewSessionRepository(db PgxIface) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, role, expired_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.Role,
		session.ExpiredAt,
	)
	return err
}
