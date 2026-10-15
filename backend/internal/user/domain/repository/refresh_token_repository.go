package repository

import (
	"context"

	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken entity.RefreshToken) *errorhandler.Response
	FindByHash(ctx context.Context, hash string) (*entity.RefreshToken, *errorhandler.Response)
	Revoke(ctx context.Context,hash string) *errorhandler.Response
}
