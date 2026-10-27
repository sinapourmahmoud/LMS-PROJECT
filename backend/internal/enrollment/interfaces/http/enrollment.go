package http

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/enrollment/domain/services"
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
	enrollmentCommandServices services.EnrollmentCommandDomainServices
	enrollmentQueryServices   services.EnrollmentQueryDomainServices
}

func NewEnrollmentHandler(enrollmentCommandServices services.EnrollmentCommandDomainServices, enrollmentQueryServices services.EnrollmentQueryDomainServices) *handler {

	return &handler{
		enrollmentCommandServices: enrollmentCommandServices,
		enrollmentQueryServices:   enrollmentQueryServices,
	}
}

func (h *handler) Enroll() gin.HandlerFunc {

	return func(c *gin.Context) {

		studentID := c.MustGet("userID").(uuid.UUID)
		courseID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		enrollmentRes, resErr := h.enrollmentCommandServices.Enroll(c, studentID, courseID)

		if resErr != nil {

			c.JSON(resErr.StatusCode, resErr)
			return

		}

		c.JSON(201, enrollmentRes)

	}
}

func (h *handler) ListMyEnrollments() gin.HandlerFunc {

	return func(c *gin.Context) {

		userID := c.MustGet("userID").(uuid.UUID)

		enrollmentRes, resErr := h.enrollmentQueryServices.ListByStudent(c, userID)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, enrollmentRes)

	}
}

func (h *handler) ListCourseEnrollments() gin.HandlerFunc {

	return func(c *gin.Context) {

		courseID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		teacherID := c.MustGet("userID").(uuid.UUID)
		role := c.MustGet("role").(userEntity.Role)

		enrollmentRes, resErr := h.enrollmentQueryServices.ListByCourse(c, courseID, teacherID, role)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, enrollmentRes)

	}

}

func (h *handler) Unenroll() gin.HandlerFunc {
	return func(c *gin.Context) {

		courseID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(400, &errorhandler.Response{Code: errorhandler.ValidationErrorCode, Message: "invalid course id", StatusCode: 400})
			return
		}

		studentID := c.MustGet("userID").(uuid.UUID)

		resErr := h.enrollmentCommandServices.UnEnroll(c, studentID, courseID)

		if resErr != nil {
			c.JSON(resErr.StatusCode, resErr)
			return
		}

		c.JSON(200, gin.H{"Message": "Unenrolled"})

	}
}
