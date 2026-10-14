package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) *errorhandler.Response
	GetByEmail(ctx context.Context, email string) (*entity.User, *errorhandler.Response)
	GetByUsername(ctx context.Context, username string) (*entity.User, *errorhandler.Response)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, *errorhandler.Response)
}
