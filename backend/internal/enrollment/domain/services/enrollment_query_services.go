package services

import (
	"context"

	"github.com/google/uuid"
	courseRepository "github.com/lms-project/backend/internal/course/domain/repository"
	"github.com/lms-project/backend/internal/enrollment/domain/repository"
	"github.com/lms-project/backend/internal/enrollment/domain/values"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type EnrollmentQueryDomainServices interface {
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*values.EnrollmentResponse, *errorhandler.Response)
	ListByCourse(ctx context.Context, courseID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) ([]*values.EnrollmentResponse, *errorhandler.Response)
}

type enrollmentQueryServices struct {
	enrollmentRepository repository.EnrollmentRepository
	courseRepository     courseRepository.CourseRepository
}

func NewEnrollmentQueryServices(enrollmentRepository repository.EnrollmentRepository, courseRepository courseRepository.CourseRepository) EnrollmentQueryDomainServices {
	return &enrollmentQueryServices{
		enrollmentRepository: enrollmentRepository,
		courseRepository:     courseRepository,
	}
}

func (r *enrollmentQueryServices) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*values.EnrollmentResponse, *errorhandler.Response) {

	entityEnrollments, err := r.enrollmentRepository.ListByStudent(ctx, studentID)

	if err != nil {
		return nil, err
	}

	enrollments := make([]*values.EnrollmentResponse, 0, len(entityEnrollments))

	for _, enrollment := range entityEnrollments {
		enrollments = append(enrollments, values.NewEnrollmentResponse(enrollment))

	}
	return enrollments, nil

}

func (r *enrollmentQueryServices) ListByCourse(ctx context.Context, courseID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) ([]*values.RosterEntry, *errorhandler.Response) {

	course, err := r.courseRepository.GetByID(ctx, courseID)

	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.CourseNotFoundErrorCode,
			Message:    errorhandler.CourseNotFoundMessage,
			StatusCode: 404,
		}
	}

	if requesterRole != userEntity.RoleAdmin && course.GetTeacherID() != requesterID {
		return nil, &errorhandler.Response{
			Code:       errorhandler.ForbiddenErrorCode,
			Message:    errorhandler.ForbiddenMessage,
			StatusCode: 403,
		}
	}

	roster, err := r.enrollmentRepository.ListRosterByCourse(ctx, courseID)

	if err != nil {
		return nil, err
	}

	return roster, nil
}