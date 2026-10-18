package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/course/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entity.Course) *errorhandler.Response
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, *errorhandler.Response)
	ListPublished(ctx context.Context) ([]*entity.Course, *errorhandler.Response)
	ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]*entity.Course, *errorhandler.Response)
	ListPending(ctx context.Context) ([]*entity.Course, *errorhandler.Response)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.Status) *errorhandler.Response
	Delete(ctx context.Context, id uuid.UUID) *errorhandler.Response
}