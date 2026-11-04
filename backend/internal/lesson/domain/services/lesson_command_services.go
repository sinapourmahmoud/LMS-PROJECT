package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	courseRepository "github.com/lms-project/backend/internal/course/domain/repository"
	"github.com/lms-project/backend/internal/lesson/domain/entity"
	"github.com/lms-project/backend/internal/lesson/domain/repository"
	"github.com/lms-project/backend/internal/lesson/domain/values"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type LessonCommandServices interface {
	Create(ctx context.Context, creatorID uuid.UUID, courseID uuid.UUID, title string, content string, position int) (*values.LessonResponse, *errorhandler.Response)
	Delete(ctx context.Context, lessonID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) *errorhandler.Response
	Update(ctx context.Context, lessonID uuid.UUID, title, content string, position int, requesterID uuid.UUID) (*values.LessonResponse, *errorhandler.Response)
}

type lessonsCommandServices struct {
	lessonRepository repository.LessonRepository
	courseRepository courseRepository.CourseRepository
}

func NewLessonCommandServices(lessonRepository repository.LessonRepository, courseRepository courseRepository.CourseRepository) LessonCommandServices {
	return &lessonsCommandServices{
		lessonRepository: lessonRepository,
		courseRepository: courseRepository,
	}
}

func (r *lessonsCommandServices) Create(ctx context.Context, creatorID uuid.UUID, courseID uuid.UUID, title string, content string, position int) (*values.LessonResponse, *errorhandler.Response) {

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

	if course.GetTeacherID() != creatorID {
		return nil, &errorhandler.Response{
			Code:       errorhandler.ForbiddenErrorCode,
			Message:    errorhandler.ForbiddenMessage,
			StatusCode: 403,
		}
	}

	lesson := entity.NewLesson(uuid.New(), courseID, title, content, position, time.Now(), nil)
	err = r.lessonRepository.Create(ctx, lesson)

	if err != nil {
		return nil, err
	}

	return values.NewLessonResponse(lesson), nil

}
func (r *lessonsCommandServices) Delete(ctx context.Context, lessonID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) *errorhandler.Response {

	lesson, err := r.lessonRepository.GetByID(ctx, lessonID)

	if err != nil {
		return err
	}

	if lesson == nil {
		return &errorhandler.Response{
			Code:       errorhandler.LessonNotFoundErrorCode,
			Message:    errorhandler.LessonNotFoundMessage,
			StatusCode: 404,
		}

	}

	course, err := r.courseRepository.GetByID(ctx, lesson.GetCourseID())

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

	return r.lessonRepository.Delete(ctx, lessonID)

}
func (r *lessonsCommandServices) Update(ctx context.Context, lessonID uuid.UUID, title, content string, position int, requesterID uuid.UUID) (*values.LessonResponse, *errorhandler.Response) {

	lesson, err := r.lessonRepository.GetByID(ctx, lessonID)

	if err != nil {
		return nil, err
	}

	if lesson == nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.LessonNotFoundErrorCode,
			Message:    errorhandler.LessonNotFoundMessage,
			StatusCode: 404,
		}
	}

	course, err := r.courseRepository.GetByID(ctx, lesson.GetCourseID())

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

	if course.GetTeacherID() != requesterID {
		return nil, &errorhandler.Response{
			Code:       errorhandler.ForbiddenErrorCode,
			Message:    errorhandler.ForbiddenMessage,
			StatusCode: 403,
		}
	}

	lesson.Update(title, content, position)

	err = r.lessonRepository.Update(ctx, lesson)

	if err != nil {
		return nil, err
	}

	return values.NewLessonResponse(lesson), nil
}
