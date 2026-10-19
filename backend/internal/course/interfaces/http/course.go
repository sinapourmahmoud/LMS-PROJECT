package http

import (
	"reflect"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/course/application/dto"
	"github.com/lms-project/backend/internal/course/domain/services"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
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
	courseCommandServices services.CourseCommandDomainServices
	courseQueryServices   services.CourseQueryDomainServices
}

func NewCourseHandler(courseCommandServices services.CourseCommandDomainServices, courseQueryDomainServices services.CourseQueryDomainServices) *handler {
	return &handler{
		courseCommandServices: courseCommandServices,
		courseQueryServices:   courseQueryDomainServices,
	}
}

func (h *handler) CreateCourse() gin.HandlerFunc {

	return func(c *gin.Context) {

		var request dto.CreateCourseReqDto

		if err := c.ShouldBindJSON(&request); err != nil {

			c.JSON(400, &errorhandler.Response{Code: errorhandler.RequestBodyParseErrorCode, Message: errorhandler.RequestBodyParseErrorMessage, StatusCode: 400})
			return

		}

		if err := validate.Struct(request); err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: errorhandler.FormatValidationErrors(err), StatusCode: 400})
			return
		}

		teacherID := c.MustGet("userID").(uuid.UUID)

		res, resErr := h.courseCommandServices.Create(c, request.Title, request.Description, teacherID)

		if resErr != nil {

			c.JSON(resErr.StatusCode, resErr)
			return

		}
		c.JSON(201, res)

	}

}

func (h *handler) ApproveCourse() gin.HandlerFunc {
	return func(c *gin.Context) {

		courseID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		resErr := h.courseCommandServices.Approve(c, courseID)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, gin.H{"message": "course approved"})

	}
}

func (h *handler) DeleteCourse() gin.HandlerFunc {

	return func(c *gin.Context) {

		courseID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		requesterID := c.MustGet("userID").(uuid.UUID)
		requesterRole := c.MustGet("role").(userEntity.Role)

		resErr := h.courseCommandServices.Delete(c, courseID, requesterID, requesterRole)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, gin.H{"message": "course deleted"})

	}
}

func (h *handler) GetCourse() gin.HandlerFunc {

	return func(c *gin.Context) {

		courseID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		res, resErr := h.courseQueryServices.GetByID(c, courseID)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)

	}
}

func (h *handler) ListPublishedCourses() gin.HandlerFunc {

	return func(c *gin.Context) {

		res, resErr := h.courseQueryServices.ListPublished(c)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)

	}
}

func (h *handler) ListMyCourses() gin.HandlerFunc {

	return func(c *gin.Context) {

		teacherID := c.MustGet("userID").(uuid.UUID)

		res, resErr := h.courseQueryServices.ListByTeacher(c, teacherID)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)

	}
}

func (h *handler) ListPendingCourses() gin.HandlerFunc {

	return func(c *gin.Context) {

		res, resErr := h.courseQueryServices.ListPending(c)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, res)

	}
}
