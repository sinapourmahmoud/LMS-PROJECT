package http

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/session/application/dto"
	"github.com/lms-project/backend/internal/session/domain/services"
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
	"github.com/lms-project/backend/pkg/pagination"
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
	sessionCommandServices services.SessionCommandServices
	sessionQueryServices   services.SessionQueryServices
}

func NewSessionHandler(sessionCommandServices services.SessionCommandServices, sessionQueryServices services.SessionQueryServices) *handler {
	return &handler{
		sessionCommandServices: sessionCommandServices,
		sessionQueryServices:   sessionQueryServices,
	}
}

func (h *handler) CreateSession() gin.HandlerFunc {

	return func(c *gin.Context) {

		courseID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		var request dto.SessionRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.RequestBodyParseErrorCode, Message: errorhandler.RequestBodyParseErrorMessage, StatusCode: 400})
			return
		}

		if err := validate.Struct(request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: errorhandler.FormatValidationErrors(err), StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)

		res, resErr := h.sessionCommandServices.Create(c, requesterID, courseID, request.Title, request.StartsAt)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(201, res)
	}
}

func (h *handler) ListSessions() gin.HandlerFunc {

	return func(c *gin.Context) {

		courseID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		params, err := pagination.Parse(c.Query("page"), c.Query("limit"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: err.Error(), StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)
		requesterRole := c.MustGet("role").(entity.Role)

		res, resErr := h.sessionQueryServices.ListByCourse(c, courseID, params, requesterID, requesterRole)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)
	}
}

func (h *handler) GetSession() gin.HandlerFunc {

	return func(c *gin.Context) {

		sessionID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid session id", StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)
		requesterRole := c.MustGet("role").(entity.Role)

		res, resErr := h.sessionQueryServices.GetByID(c, sessionID, requesterID, requesterRole)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)
	}
}

func (h *handler) UpdateSession() gin.HandlerFunc {

	return func(c *gin.Context) {

		sessionID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid session id", StatusCode: 400})
			return
		}

		var request dto.SessionRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.RequestBodyParseErrorCode, Message: errorhandler.RequestBodyParseErrorMessage, StatusCode: 400})
			return
		}

		if err := validate.Struct(request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: errorhandler.FormatValidationErrors(err), StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)

		res, resErr := h.sessionCommandServices.Update(c, sessionID, request.Title, request.StartsAt, requesterID)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)
	}
}

func (h *handler) StartSession() gin.HandlerFunc {

	return func(c *gin.Context) {

		sessionID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid session id", StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)

		res, resErr := h.sessionCommandServices.Start(c, sessionID, requesterID)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)
	}
}

func (h *handler) EndSession() gin.HandlerFunc {

	return func(c *gin.Context) {

		sessionID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid session id", StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)

		res, resErr := h.sessionCommandServices.End(c, sessionID, requesterID)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)
	}
}

func (h *handler) DeleteSession() gin.HandlerFunc {

	return func(c *gin.Context) {

		sessionID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid session id", StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)
		requesterRole := c.MustGet("role").(entity.Role)

		resErr := h.sessionCommandServices.Delete(c, sessionID, requesterID, requesterRole)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.Status(204)
	}
}
