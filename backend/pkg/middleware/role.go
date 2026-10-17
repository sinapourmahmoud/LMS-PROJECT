package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

func RequireRole(roles ...entity.Role) gin.HandlerFunc {

	return func(c *gin.Context) {

		role, exists := c.Get("role")

		if !exists {

			c.JSON(401, &errorhandler.Response{Code: errorhandler.InvalidAccessTokenErrorCode, Message: errorhandler.InvalidAccessTokenMessage, StatusCode: 401})
			c.Abort()
			return

		}

		userRole, ok := role.(entity.Role)

		if !ok {
			c.JSON(500, &errorhandler.Response{Code: errorhandler.InternalSystemErrorCode, Message: errorhandler.InternalSystemErrorMessage, StatusCode: 500})
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(403, &errorhandler.Response{Code: errorhandler.ForbiddenErrorCode, Message: errorhandler.ForbiddenMessage, StatusCode: 403})
		c.Abort()

	}
}
