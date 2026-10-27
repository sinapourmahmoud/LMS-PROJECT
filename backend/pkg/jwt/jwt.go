package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/user/domain/entity"
)

type Claims struct {
	SessionID uuid.UUID   `json:"session_id"`
	UserID    uuid.UUID   `json:"user_id"`
	Role      entity.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(sessionID uuid.UUID, userID uuid.UUID, role entity.Role, secret string, ttl time.Duration) (string, error) {

	claims := Claims{
		SessionID: sessionID,
		UserID:    userID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}

func GenerateRefreshToken(sessionID uuid.UUID, userID uuid.UUID, secret string, ttl time.Duration) (string, error) {
	claims := Claims{
		SessionID: sessionID,

		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])

		}
		return []byte(secret), nil

	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {

		return nil, fmt.Errorf("invalid token")

	}
	return claims, nil

}
