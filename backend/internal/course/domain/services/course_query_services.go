package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/course/domain/repository"
	"github.com/lms-project/backend/internal/course/domain/values"
	userRepository "github.com/lms-project/backend/internal/user/domain/repository"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type CourseQueryDomainServices interface {
	GetByID(ctx context.Context, id uuid.UUID) (*values.CourseResponse, *errorhandler.Response)
	ListPublished(ctx context.Context) ([]*values.CourseResponse, *errorhandler.Response)
	ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]*values.CourseResponse, *errorhandler.Response)
	ListPending(ctx context.Context) ([]*values.CourseResponse, *errorhandler.Response)
}

type courseQueryServices struct {
	courseRepository repository.CourseRepository
	userRepository   userRepository.UserRepository
}

func NewCourseQueryServices(courseRepository repository.CourseRepository, userRepository userRepository.UserRepository) CourseQueryDomainServices {
	return &courseQueryServices{
		courseRepository: courseRepository,
		userRepository:   userRepository,
	}
}

func (r *courseQueryServices) GetByID(ctx context.Context, id uuid.UUID) (*values.CourseResponse, *errorhandler.Response) {

	course, err := r.courseRepository.GetByID(ctx, id)

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

	return values.NewCourseResponse(course), nil

}

func (r *courseQueryServices) ListPublished(ctx context.Context) ([]*values.CourseResponse, *errorhandler.Response) {

	entityCourses, err := r.courseRepository.ListPublished(ctx)

	if err != nil {
		return nil, err
	}
	var courses = make([]*values.CourseResponse, 0, len(entityCourses))

	for _, course := range entityCourses {
		courses = append(courses, values.NewCourseResponse(course))
	}

	return courses, nil

}

func (r *courseQueryServices) ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]*values.CourseResponse, *errorhandler.Response) {

	user, err := r.userRepository.GetByID(ctx, teacherID)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.UserNotFoundErrorCode,
			Message:    errorhandler.UserNotFoundMessage,
			StatusCode: 404,
		}
	}

	entityCourses, err := r.courseRepository.ListByTeacher(ctx, teacherID)

	if err != nil {
		return nil, err
	}
	var courses = make([]*values.CourseResponse, 0, len(entityCourses))

	for _, course := range entityCourses {
		courses = append(courses, values.NewCourseResponse(course))
	}

	return courses, nil

}
func (r *courseQueryServices) ListPending(ctx context.Context) ([]*values.CourseResponse, *errorhandler.Response) {

	entityCourses, err := r.courseRepository.ListPending(ctx)

	if err != nil {
		return nil, err
	}
	var courses = make([]*values.CourseResponse, 0, len(entityCourses))

	for _, course := range entityCourses {
		courses = append(courses, values.NewCourseResponse(course))
	}

	return courses, nil

}
