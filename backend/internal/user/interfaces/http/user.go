package http

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lms-project/backend/cmd/app/config"
	"github.com/lms-project/backend/internal/user/application/dto"
	"github.com/lms-project/backend/internal/user/domain/services"
	"github.com/lms-project/backend/pkg/errorhandler"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

type handler struct {
	userCommandServices     services.UserCommandDomainServices
	userQueryDomainServices services.UserQueryDomainServices
}

func NewUserHandler(userCommandServices services.UserCommandDomainServices, userQueryDomainServices services.UserQueryDomainServices) *handler {
	return &handler{
		userCommandServices:     userCommandServices,
		userQueryDomainServices: userQueryDomainServices,
	}
}

// sign up handler

func (h *handler) SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {

		var request dto.SignUpReqDto

		if err := c.ShouldBindJSON(&request); err != nil {

			c.JSON(400, &errorhandler.Response{Code: errorhandler.RequestBodyParseErrorCode, Message: errorhandler.RequestBodyParseErrorMessage, StatusCode: 400})
			return

		}

		if err := validate.Struct(request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: errorhandler.FormatValidationErrors(err), StatusCode: 400})
			return
		}

		res, resErr := h.userCommandServices.SignUp(c, &request)

		if resErr != nil {

			c.JSON(resErr.StatusCode, resErr)
			return

		}
		c.JSON(201, res)

	}
}

func (h *handler) Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		var request dto.LoginReqDto

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.RequestBodyParseErrorCode, Message: errorhandler.RequestBodyParseErrorMessage, StatusCode: 400})
			return
		}

		if err := validate.Struct(request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: errorhandler.FormatValidationErrors(err), StatusCode: 400})
			return
		}

		res, resErr := h.userQueryDomainServices.Login(c, &request)
		if resErr != nil {

			c.JSON(resErr.StatusCode, resErr)
			return

		}

		http.SetCookie(c.Writer, &http.Cookie{

			Name:     "refresh_token",
			Value:    res.RefreshToken,
			Path:     "/refresh",
			HttpOnly: true,
			MaxAge:   int(time.Duration(config.Config.Auth.RefreshTokenTTLHr) * time.Hour / time.Second),
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		c.JSON(200, gin.H{"access_token": res.AccessToken})

	}

}

func (h *handler) Refresh() gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := c.Cookie("refresh_token")
		if err != nil {

			c.JSON(401, &errorhandler.Response{Code: errorhandler.InvalidRefreshTokenErrorCode, Message: errorhandler.InvalidRefreshTokenMessage, StatusCode: 401})
			return

		}

		res, resErr := h.userQueryDomainServices.Refresh(c, token)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return

		}

		c.JSON(200, res)

	}
}

func (h *handler) Logout() gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := c.Cookie("refresh_token")

		if err != nil {
			c.JSON(401, &errorhandler.Response{Code: errorhandler.InvalidRefreshTokenErrorCode, Message: errorhandler.InvalidRefreshTokenMessage, StatusCode: 401})
			return
		}

		resErr := h.userQueryDomainServices.Logout(c, token)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/refresh",
			HttpOnly: true,
			MaxAge:   -1,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		c.JSON(200, gin.H{"message": "Logged out"})

	}
}
