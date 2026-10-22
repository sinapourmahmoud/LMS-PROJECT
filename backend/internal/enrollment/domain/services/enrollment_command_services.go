package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	courseEntity "github.com/lms-project/backend/internal/course/domain/entity"
	courseRepository "github.com/lms-project/backend/internal/course/domain/repository"
	"github.com/lms-project/backend/internal/enrollment/domain/entity"
	"github.com/lms-project/backend/internal/enrollment/domain/repository"
	"github.com/lms-project/backend/internal/enrollment/domain/values"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type EnrollmentCommandDomainServices interface {
	Enroll(ctx context.Context, studentID uuid.UUID, courseID uuid.UUID) (*values.EnrollmentResponse, *errorhandler.Response)
}

type enrollmentCommandServices struct {
	enrollmentRepository repository.EnrollmentRepository
	courseRepository     courseRepository.CourseRepository
}

func NewEnrollmentCommandServices(enrollmentRepository repository.EnrollmentRepository, courseRepository courseRepository.CourseRepository) EnrollmentCommandDomainServices {

	return &enrollmentCommandServices{
		enrollmentRepository: enrollmentRepository,
		courseRepository:     courseRepository,
	}

}

func (r *enrollmentCommandServices) Enroll(ctx context.Context, studentID uuid.UUID, courseID uuid.UUID) (*values.EnrollmentResponse, *errorhandler.Response) {

	course, err := r.courseRepository.GetByID(ctx, courseID)

	if err != nil {
		return nil, err
	}

	if course == nil || course.GetStatus() != courseEntity.StatusPublished {
		return nil, &errorhandler.Response{
			Code:       errorhandler.CourseNotFoundErrorCode,
			Message:    errorhandler.CourseNotFoundMessage,
			StatusCode: 404,
		}
	}

	existing, err := r.enrollmentRepository.FindByStudentAndCourse(ctx, studentID, courseID)

	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.AlreadyEnrolledErrorCode,
			Message:    errorhandler.AlreadyEnrolledErrorMessage,
			StatusCode: 409,
		}
	}

	enrollment := entity.NewEnrollment(uuid.New(), studentID, courseID, time.Now())

	if resErr := r.enrollmentRepository.Create(ctx, enrollment); resErr != nil {
		return nil, resErr
	}

	return values.NewEnrollmentResponse(enrollment), nil
}
