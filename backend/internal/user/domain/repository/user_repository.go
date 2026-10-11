package repository

import (
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/internal/user/domain/values"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type UserRepository interface {
	Create(user *entity.User) *errorhandler.Response
	GetByEmail(email string) (*entity.User,*values.UserResponse)
	GetByUsername(username string) (*entity.User,*values.UserResponse)
}
