package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/course/domain/entity"
	"github.com/lms-project/backend/internal/course/domain/repository"
	"github.com/lms-project/backend/internal/course/domain/values"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	userRepository "github.com/lms-project/backend/internal/user/domain/repository"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type CourseCommandDomainServices interface {
	Create(ctx context.Context, title string, description string, teacherID uuid.UUID) (*values.CourseResponse, *errorhandler.Response)
	Approve(ctx context.Context, courseID uuid.UUID) *errorhandler.Response
	Delete(ctx context.Context, courseID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) *errorhandler.Response
}

type courseCommandServices struct {
	courseRepository repository.CourseRepository
	userRepository   userRepository.UserRepository
}

func NewCourseCommandServices(courseRepository repository.CourseRepository, userRepository userRepository.UserRepository) CourseCommandDomainServices {
	return &courseCommandServices{
		courseRepository: courseRepository,
		userRepository:   userRepository,
	}
}

func (r *courseCommandServices) Create(ctx context.Context, title string, description string, teacherID uuid.UUID) (*values.CourseResponse, *errorhandler.Response) {

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

	course := entity.NewCourse(uuid.New(), title, description, user.GetID(), entity.StatusPending, time.Now(), nil)

	if resErr := r.courseRepository.Create(ctx, course); resErr != nil {
		return nil, resErr
	}

	return &values.CourseResponse{ID: course.GetID(), Title: course.GetTitle(), Description: course.GetDescription(), TeacherID: user.GetID(), Status: course.GetStatus()}, nil
}
func (r *courseCommandServices) Approve(ctx context.Context, courseID uuid.UUID) *errorhandler.Response {

	course, err := r.courseRepository.GetByID(ctx, courseID)

	if err != nil {
		return err
	}

	if course == nil {
		return &errorhandler.Response{
			Code:       errorhandler.CourseNotFoundErrorCode,
			Message:    errorhandler.CourseNotFoundMessage,
			StatusCode: 404,
		}
	}

	err = r.courseRepository.UpdateStatus(ctx, courseID, entity.StatusPublished)

	if err != nil {
		return err
	}

	return nil

}

func (r *courseCommandServices) Delete(ctx context.Context, courseID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) *errorhandler.Response {

	course, err := r.courseRepository.GetByID(ctx, courseID)
	if err != nil {
		return err
	}

	if course == nil {
		return &errorhandler.Response{
			Code:       errorhandler.CourseNotFoundErrorCode,
			Message:    errorhandler.CourseNotFoundMessage,
			StatusCode: 404,
		}
	}

	if requesterRole != userEntity.RoleAdmin && course.GetTeacherID() != requesterID {
		return &errorhandler.Response{
			Code:       errorhandler.ForbiddenErrorCode,
			Message:    errorhandler.ForbiddenMessage,
			StatusCode: 403,
		}
	}

	if err := r.courseRepository.Delete(ctx, courseID); err != nil {
		return err
	}

	return nil
}
