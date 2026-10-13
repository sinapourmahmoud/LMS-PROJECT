package http

import (
	"reflect"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	userCommandServices services.UserCommandDomainServices
}

func NewUserHandler(userCommandServices services.UserCommandDomainServices) *handler {
	return &handler{
		userCommandServices: userCommandServices,
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
