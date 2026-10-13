package services

import (
	"context"
	"time"

	"github.com/lms-project/backend/cmd/app/config"
	"github.com/lms-project/backend/internal/user/application/dto"
	"github.com/lms-project/backend/internal/user/domain/repository"
	"github.com/lms-project/backend/internal/user/domain/values"
	"github.com/lms-project/backend/pkg/errorhandler"
	"github.com/lms-project/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserQueryDomainServices interface {
	Login(ctx context.Context, req *dto.LoginReqDto) (*values.TokenResponse, *errorhandler.Response)
}

type userQueryDomainServices struct {
	userRepository repository.UserRepository
}

func NewUserQueryServices(userRepository repository.UserRepository) UserQueryDomainServices {

	return &userQueryDomainServices{
		userRepository: userRepository,
	}

}

func (r *userQueryDomainServices) Login(ctx context.Context, req *dto.LoginReqDto) (*values.TokenResponse, *errorhandler.Response) {

	user, err := r.userRepository.GetByEmail(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.InvalidCredentialsErrorCode,
			Message:    errorhandler.InvalidCredentialsMessage,
			StatusCode: 401,
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.GetHashedPassword()), []byte(req.Password)); err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.InvalidCredentialsErrorCode,
			Message:    errorhandler.InvalidCredentialsMessage,
			StatusCode: 401,
		}
	}

	accessToken, tokenErr := jwt.GenerateAccessToken(user.GetID(), user.GetRole(), config.Config.Auth.JWTAccessSecret, time.Duration(config.Config.Auth.AccessTokenTTLMin)*time.Minute)
	if tokenErr != nil {
		return nil, &errorhandler.Response{Code: errorhandler.InternalSystemErrorCode, Message: errorhandler.InternalSystemErrorMessage, StatusCode: 500}
	}

	refreshToken, tokenErr := jwt.GenerateRefreshToken(user.GetID(), config.Config.Auth.JWTRefreshSecret, time.Duration(config.Config.Auth.RefreshTokenTTLHr)*time.Hour)
	if tokenErr != nil {
		return nil, &errorhandler.Response{Code: errorhandler.InternalSystemErrorCode, Message: errorhandler.InternalSystemErrorMessage, StatusCode: 500}
	}

	return values.NewTokenResponse(accessToken, refreshToken), nil
}
