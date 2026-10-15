package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/internal/user/domain/repository"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type refreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) repository.RefreshTokenRepository {

	return &refreshTokenRepository{
		db: db,
	}

}

func (r *refreshTokenRepository) Create(ctx context.Context, refreshToken *entity.RefreshToken) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO refresh_tokens(id,user_id,token_hash,expires_at,revoked_at,created_at) VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.db.ExecContext(ctx, query,
		refreshToken.GetID(),
		refreshToken.GetUserID(),
		refreshToken.GetTokenHash(),
		refreshToken.GetExpiresAt(),
		refreshToken.GetRevokedAt(),
		refreshToken.GetCreatedAt(),
	)

	if err != nil {
		return &errorhandler.Response{
			Message:    errorhandler.DatabaseErrorMessage,
			Code:       errorhandler.DatabaseErrorCode,
			StatusCode: 500,
		}
	}

	return nil
}
func (r *refreshTokenRepository) FindByHash(ctx context.Context, hash string) (*entity.RefreshToken, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
	SELECT id,user_id,token_hash,expires_at,revoked_at,created_at FROM refresh_tokens WHERE token_hash=$1
	`

	var (
		id        uuid.UUID
		userID    uuid.UUID
		tokenHash string
		expiresAt time.Time
		createdAt time.Time
		revokedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query,hash).Scan(&id, &userID, &tokenHash, &expiresAt, &revokedAt, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}

	}

	var revokedAtPtr *time.Time
	if revokedAt.Valid {
		revokedAtPtr = &revokedAt.Time
	}

	return entity.NewRefreshToken(id, userID, tokenHash, expiresAt, revokedAtPtr, createdAt), nil

}
func (r *refreshTokenRepository) Revoke(ctx context.Context, hash string) *errorhandler.Response {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
	UPDATE  refresh_tokens SET revoked_at = NOW() WHERE token_hash=$1 
	`

	_, err := r.db.ExecContext(ctx, query, hash)

	if err != nil {

		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}

	}

	return nil
}
