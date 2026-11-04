package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/lesson/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type LessonRepository interface {
	Create(ctx context.Context, lessonEntity *entity.Lesson) *errorhandler.Response
	GetByID(ctx context.Context, lessonID uuid.UUID) (*entity.Lesson, *errorhandler.Response)
	Update(ctx context.Context, lesson *entity.Lesson) *errorhandler.Response
	Delete(ctx context.Context, lessonID uuid.UUID) *errorhandler.Response
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, *errorhandler.Response)
}
