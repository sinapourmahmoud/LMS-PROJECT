package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lms-project/backend/cmd/app/config"
	"github.com/lms-project/backend/pkg/errorhandler"
	"github.com/lms-project/backend/pkg/jwt"
)

type SessionChecker interface {
	IsSessionValid(ctx context.Context, sessionID uuid.UUID) *errorhandler.Response
}

func AuthMiddleware(sessionChecker SessionChecker) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, &errorhandler.Response{
				Code:       errorhandler.InvalidAccessTokenErrorCode,
				Message:    errorhandler.InvalidAccessTokenMessage,
				StatusCode: 401,
			})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1], config.Config.Auth.JWTAccessSecret)
		if err != nil {

			if ve, ok := err.(*gojwt.ValidationError); ok && ve.Errors&gojwt.ValidationErrorExpired != 0 {
				c.JSON(401, &errorhandler.Response{Code: errorhandler.ExpiredAccessTokenErrorCode, Message: errorhandler.ExpiredAccessTokenMessage, StatusCode: 401})
				c.Abort()
				return
			}

			c.JSON(401, &errorhandler.Response{Code: errorhandler.InvalidAccessTokenErrorCode, Message: errorhandler.InvalidAccessTokenMessage, StatusCode: 401})
			c.Abort()
			return
		}

		if resErr := sessionChecker.IsSessionValid(c, claims.SessionID); resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()

	}
}
