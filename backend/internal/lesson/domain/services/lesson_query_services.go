package services

import (
	"context"

	"github.com/google/uuid"
	courseEntity "github.com/lms-project/backend/internal/course/domain/entity"
	courseRepository "github.com/lms-project/backend/internal/course/domain/repository"
	enrollmentRepository "github.com/lms-project/backend/internal/enrollment/domain/repository"
	"github.com/lms-project/backend/internal/lesson/domain/repository"
	"github.com/lms-project/backend/internal/lesson/domain/values"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type LessonQueryServices interface {
	ListByCourse(ctx context.Context, courseID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) ([]*values.LessonResponse, *errorhandler.Response)
	GetByID(ctx context.Context, lessonID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) (*values.LessonResponse, *errorhandler.Response)
}

type lessonQueryServices struct {
	lessonRepository     repository.LessonRepository
	courseRepository     courseRepository.CourseRepository
	enrollmentRepository enrollmentRepository.EnrollmentRepository
}

func NewLessonQueryServices(lessonRepository repository.LessonRepository, courseRepository courseRepository.CourseRepository, enrollmentRepository enrollmentRepository.EnrollmentRepository) LessonQueryServices {
	return &lessonQueryServices{
		lessonRepository:     lessonRepository,
		courseRepository:     courseRepository,
		enrollmentRepository: enrollmentRepository,
	}
}

func (r *lessonQueryServices) canView(ctx context.Context, course *courseEntity.Course, requesterID uuid.UUID, requesterRole userEntity.Role) (bool, *errorhandler.Response) {

	if requesterRole == userEntity.RoleAdmin || requesterID == course.GetTeacherID() {
		return true, nil
	}

	enrollmentEntity, err := r.enrollmentRepository.FindByStudentAndCourse(ctx, requesterID, course.GetID())

	if err != nil {
		return false, err
	}

	if enrollmentEntity == nil {
		return false, nil
	}

	return enrollmentEntity.GetStudentID() == requesterID, nil

}

func (r *lessonQueryServices) ListByCourse(ctx context.Context, courseID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) ([]*values.LessonResponse, *errorhandler.Response) {

	courseEntity, err := r.courseRepository.GetByID(ctx, courseID)

	if err != nil {
		return nil, err
	}

	if courseEntity == nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.CourseNotFoundErrorCode,
			Message:    errorhandler.CourseNotFoundMessage,
			StatusCode: 404,
		}

	}

	allowed, err := r.canView(ctx, courseEntity, requesterID, requesterRole)

	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, &errorhandler.Response{
			Code:       errorhandler.ForbiddenErrorCode,
			Message:    errorhandler.ForbiddenMessage,
			StatusCode: 403,
		}
	}

	lessonsEntity, err := r.lessonRepository.ListByCourse(ctx, courseEntity.GetID())

	if err != nil {
		return nil, err
	}
	lessons := make([]*values.LessonResponse, 0, len(lessonsEntity))

	for _, lesson := range lessonsEntity {
		lessons = append(lessons, values.NewLessonResponse(lesson))
	}

	return lessons, nil

}
func (r *lessonQueryServices) GetByID(ctx context.Context, lessonID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) (*values.LessonResponse, *errorhandler.Response) {

	lessonEntity, err := r.lessonRepository.GetByID(ctx, lessonID)

	if err != nil {
		return nil, err
	}

	if lessonEntity == nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.LessonNotFoundErrorCode,
			Message:    errorhandler.LessonNotFoundMessage,
			StatusCode: 404,
		}
	}

	courseEntity, err := r.courseRepository.GetByID(ctx, lessonEntity.GetCourseID())

	if err != nil {
		return nil, err
	}

	if courseEntity == nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.LessonNotFoundErrorCode,
			Message:    errorhandler.LessonNotFoundMessage,
			StatusCode: 404,
		}
	}

	allowed, err := r.canView(ctx, courseEntity, requesterID, requesterRole)

	if err != nil {

		return nil, err
	}

	if !allowed {
		return nil, &errorhandler.Response{
			Code:       errorhandler.ForbiddenErrorCode,
			Message:    errorhandler.ForbiddenMessage,
			StatusCode: 403,
		}
	}

	return values.NewLessonResponse(lessonEntity), nil

}
