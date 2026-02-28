package repository

import (
	"context"
	"inventory-system/internal/model"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	Revoke(ctx context.Context, sessionID uuid.UUID) error
	GetValid(ctx context.Context, sessionID uuid.UUID) (*model.Session, error)
}

type sessionRepository struct {
	db PgxIface
}

func NewSessionRepository(db PgxIface) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *model.Session) error {
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

func (r *sessionRepository) Revoke(ctx context.Context, sessionID uuid.UUID) error {
	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, sessionID)
	return err
}

func (r *sessionRepository) GetValid(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	query := `
		SELECT id, user_id, role, expired_at, revoked_at, created_at
		FROM sessions
		WHERE id = $1
		  AND expired_at > NOW()
		  AND revoked_at IS NULL
	`

	session := &model.Session{}
	err := r.db.QueryRow(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.Role,
		&session.ExpiredAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return session, nil
}
