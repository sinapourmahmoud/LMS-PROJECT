package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	id uuid.UUID

	userID    uuid.UUID
	email     string
	tokenHash string
	expiresAt time.Time
	createdAt time.Time
	revokedAt *time.Time
}

func (r *RefreshToken) GetID() uuid.UUID {
	return r.id
}
func (r *RefreshToken) GetUserID() uuid.UUID {
	return r.userID
}
func (r *RefreshToken) GetEmail() string {
	return r.email
}
func (r *RefreshToken) GetTokenHash() string {
	return r.tokenHash
}
func (r *RefreshToken) GetExpiresAt() time.Time {
	return r.expiresAt
}
func (r *RefreshToken) GetCreatedAt() time.Time {
	return r.createdAt
}
func (r *RefreshToken) GetRevokedAt() *time.Time {
	return r.revokedAt
}

func (r *RefreshToken) IsRevoked() bool {
	return r.revokedAt != nil
}

func (r *RefreshToken) IsExpired(now time.Time) bool {
	return now.After(r.expiresAt)
}

func NewRefreshToken(id uuid.UUID, userID uuid.UUID, tokenHash string, expiresAt time.Time, revokedAt *time.Time, createdAt time.Time) *RefreshToken {
	return &RefreshToken{
		id:        id,
		userID:    userID,
		tokenHash: tokenHash,
		expiresAt: expiresAt,
		revokedAt: revokedAt,
		createdAt: createdAt,
	}
}
