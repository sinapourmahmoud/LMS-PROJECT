package http

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/lesson/application/dto"
	"github.com/lms-project/backend/internal/lesson/domain/services"
	"github.com/lms-project/backend/internal/user/domain/entity"
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
	lessonCommandServices services.LessonCommandServices
	lessonQueryServices   services.LessonQueryServices
}

func NewLessonHandler(lessonCommandServices services.LessonCommandServices, lessonQueryServices services.LessonQueryServices) *handler {
	return &handler{
		lessonCommandServices: lessonCommandServices,
		lessonQueryServices:   lessonQueryServices,
	}
}

func (h *handler) CreateLesson() gin.HandlerFunc {

	return func(c *gin.Context) {

		courseID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		var request dto.LessonRequest

		if err := c.ShouldBindJSON(&request); err != nil {

			c.JSON(400, &errorhandler.Response{Code: errorhandler.RequestBodyParseErrorCode, Message: errorhandler.RequestBodyParseErrorMessage, StatusCode: 400})
			return

		}

		if err := validate.Struct(request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: errorhandler.FormatValidationErrors(err), StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)

		res, resErr := h.lessonCommandServices.Create(c, requesterID, courseID, request.Title, request.Content, request.Position)
		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}
		c.JSON(201, res)

	}

}

func (h *handler) ListLessons() gin.HandlerFunc {

	return func(c *gin.Context) {

		courseID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)
		requesterRole := c.MustGet("role").(entity.Role)

		res, resErr := h.lessonQueryServices.ListByCourse(c, courseID, requesterID, requesterRole)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)

	}

}

func (h *handler) GetLesson() gin.HandlerFunc {

	return func(c *gin.Context) {

		lessonID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid lesson id", StatusCode: 400})
			return
		}
		requesterID := c.MustGet("userID").(uuid.UUID)
		requesterRole := c.MustGet("role").(entity.Role)

		res, resErr := h.lessonQueryServices.GetByID(c, lessonID, requesterID, requesterRole)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}
		c.JSON(200, res)

	}
}

func (h *handler) UpdateLesson() gin.HandlerFunc {

	return func(c *gin.Context) {

		lessonID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid lesson id", StatusCode: 400})
			return
		}

		var request dto.LessonRequest

		if err := c.ShouldBindJSON(&request); err != nil {

			c.JSON(400, &errorhandler.Response{Code: errorhandler.RequestBodyParseErrorCode, Message: errorhandler.RequestBodyParseErrorMessage, StatusCode: 400})
			return

		}

		if err := validate.Struct(request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: errorhandler.FormatValidationErrors(err), StatusCode: 400})
			return
		}
		requesterID := c.MustGet("userID").(uuid.UUID)

		res, resErr := h.lessonCommandServices.Update(c, lessonID, request.Title, request.Content, request.Position, requesterID)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)

	}

}

func (h *handler) DeleteLesson() gin.HandlerFunc {

	return func(c *gin.Context) {
		lessonID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid lesson id", StatusCode: 400})
			return
		}
		requesterID := c.MustGet("userID").(uuid.UUID)
		requesterRole := c.MustGet("role").(entity.Role)

		resErr := h.lessonCommandServices.Delete(c, lessonID, requesterID, requesterRole)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}
		c.Status(204)

	}
}
