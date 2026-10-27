package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/user/application/dto"
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/internal/user/domain/repository"
	"github.com/lms-project/backend/internal/user/domain/values"
	"github.com/lms-project/backend/pkg/errorhandler"
	"golang.org/x/crypto/bcrypt"
)

type UserCommandDomainServices interface {
	SignUp(ctx context.Context, req *dto.SignUpReqDto) (*values.UserResponse, *errorhandler.Response)
}

type userCommandDomainServices struct {
	userRepository repository.UserRepository
}

func NewUserCommandDomainServices(userRepository repository.UserRepository) UserCommandDomainServices {

	return &userCommandDomainServices{
		userRepository: userRepository,
	}

}

func (r *userCommandDomainServices) SignUp(ctx context.Context, req *dto.SignUpReqDto) (*values.UserResponse, *errorhandler.Response) {

	user, err := r.userRepository.GetByEmail(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.EmailAlreadyExistsErrorCode,
			Message:    errorhandler.EmailAlreadyExistsMessage,
			StatusCode: 409,
		}

	}

	user, err = r.userRepository.GetByUsername(ctx, req.Username)

	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.UsernameAlreadyExistsErrorCode,
			Message:    errorhandler.UsernameAlreadyExistsMessage,
			StatusCode: 409,
		}

	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.InternalSystemErrorCode,
			Message:    errorhandler.InternalSystemErrorMessage,
			StatusCode: 500,
		}
	}

	user = entity.NewUser(uuid.New(), req.Username, req.Email, string(hashedPassword), entity.RoleStudent, time.Now(), nil)
	err = r.userRepository.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return values.NewUserResponse(user), nil

}
