package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	gojwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lms-project/backend/cmd/app/config"
	"github.com/lms-project/backend/internal/user/application/dto"
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/internal/user/domain/repository"
	"github.com/lms-project/backend/internal/user/domain/values"
	"github.com/lms-project/backend/pkg/errorhandler"
	"github.com/lms-project/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserQueryDomainServices interface {
	Login(ctx context.Context, req *dto.LoginReqDto) (*values.TokenResponse, *errorhandler.Response)
	Refresh(ctx context.Context, refreshToken string) (*values.AccessTokenResponse, *errorhandler.Response)
	Logout(ctx context.Context, refreshToken string) *errorhandler.Response
	IsSessionValid(ctx context.Context, sessionID uuid.UUID) (entity.Role, *errorhandler.Response)
}

type userQueryDomainServices struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
}

func NewUserQueryServices(userRepository repository.UserRepository, refreshTokenRepository repository.RefreshTokenRepository) UserQueryDomainServices {

	return &userQueryDomainServices{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
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

	sessionID := uuid.New()

	accessToken, tokenErr := jwt.GenerateAccessToken(sessionID, user.GetID(), user.GetRole(), config.Config.Auth.JWTAccessSecret, time.Duration(config.Config.Auth.AccessTokenTTLMin)*time.Minute)
	if tokenErr != nil {
		return nil, &errorhandler.Response{Code: errorhandler.InternalSystemErrorCode, Message: errorhandler.InternalSystemErrorMessage, StatusCode: 500}
	}

	refreshToken, tokenErr := jwt.GenerateRefreshToken(sessionID, user.GetID(), config.Config.Auth.JWTRefreshSecret, time.Duration(config.Config.Auth.RefreshTokenTTLHr)*time.Hour)
	if tokenErr != nil {
		return nil, &errorhandler.Response{Code: errorhandler.InternalSystemErrorCode, Message: errorhandler.InternalSystemErrorMessage, StatusCode: 500}
	}

	hashToken := hashRefreshToken(refreshToken)

	tokenEntity := entity.NewRefreshToken(sessionID, user.GetID(), hashToken, time.Now().Add(time.Duration(config.Config.Auth.RefreshTokenTTLHr)*time.Hour), nil, time.Now())

	if repoErr := r.refreshTokenRepository.Create(ctx, tokenEntity); repoErr != nil {

		return nil, repoErr

	}

	return values.NewTokenResponse(accessToken, refreshToken), nil
}

func (r *userQueryDomainServices) Refresh(ctx context.Context, refreshToken string) (*values.AccessTokenResponse, *errorhandler.Response) {

	claims, err := jwt.ParseToken(refreshToken, config.Config.Auth.JWTRefreshSecret)
	if err != nil {
		if ve, ok := err.(*gojwt.ValidationError); ok && ve.Errors&gojwt.ValidationErrorExpired != 0 {
			return nil, &errorhandler.Response{Code: errorhandler.ExpiredRefreshTokenErrorCode, Message: errorhandler.ExpiredRefreshTokenMessage, StatusCode: 401}
		}
		return nil, &errorhandler.Response{Code: errorhandler.InvalidRefreshTokenErrorCode, Message: errorhandler.InvalidRefreshTokenMessage, StatusCode: 401}
	}

	hash := hashRefreshToken(refreshToken)

	tokenRecord, repoErr := r.refreshTokenRepository.FindByHash(ctx, hash)
	if repoErr != nil {
		return nil, repoErr
	}
	if tokenRecord == nil || tokenRecord.IsRevoked() || tokenRecord.IsExpired(time.Now()) {
		return nil, &errorhandler.Response{Code: errorhandler.InvalidRefreshTokenErrorCode, Message: errorhandler.InvalidRefreshTokenMessage, StatusCode: 401}
	}

	user, repoErr := r.userRepository.GetByID(ctx, claims.UserID)

	if repoErr != nil {
		return nil, repoErr
	}
	if user == nil {
		return nil, &errorhandler.Response{Code: errorhandler.UserNotFoundErrorCode, Message: errorhandler.UserNotFoundMessage, StatusCode: 404}
	}

	accessToken, tokenErr := jwt.GenerateAccessToken(tokenRecord.GetID(), user.GetID(), user.GetRole(), config.Config.Auth.JWTAccessSecret, time.Duration(config.Config.Auth.AccessTokenTTLMin)*time.Minute)
	if tokenErr != nil {
		return nil, &errorhandler.Response{Code: errorhandler.InternalSystemErrorCode, Message: errorhandler.InternalSystemErrorMessage, StatusCode: 500}
	}

	return values.NewAccessTokenResponse(accessToken), nil

}

func (r *userQueryDomainServices) Logout(ctx context.Context, refreshToken string) *errorhandler.Response {
	hash := hashRefreshToken(refreshToken)

	if repoErr := r.refreshTokenRepository.Revoke(ctx, hash); repoErr != nil {
		return repoErr
	}

	return nil
}

func (r *userQueryDomainServices) IsSessionValid(ctx context.Context, sessionID uuid.UUID) (entity.Role, *errorhandler.Response) {

	tokenRecord, repoErr := r.refreshTokenRepository.FindByID(ctx, sessionID)

	if repoErr != nil {
		return "", repoErr
	}

	if tokenRecord == nil || tokenRecord.IsRevoked() || tokenRecord.IsExpired(time.Now()) {

		return "", &errorhandler.Response{Code: errorhandler.InvalidAccessTokenErrorCode, Message: errorhandler.InvalidAccessTokenMessage, StatusCode: 401}
	}

	user, repoErr := r.userRepository.GetByID(ctx, tokenRecord.GetUserID())
	if repoErr != nil {
		return "", repoErr
	}
	if user == nil {
		return "", &errorhandler.Response{Code: errorhandler.UserNotFoundErrorCode, Message: errorhandler.UserNotFoundMessage, StatusCode: 404}
	}

	return user.GetRole(), nil
}

// hash function
func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
