package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	BaseSimple
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Role      UserRole   `json:"role" db:"role"`
	ExpiredAt time.Time  `json:"expired_at" db:"expired_at"`
	RevokedAt *time.Time `json:"revoked_at" db:"revoked_at"`
}
