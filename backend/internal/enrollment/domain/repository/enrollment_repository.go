package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/enrollment/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *entity.Enrollment) *errorhandler.Response
	FindByStudentAndCourse(ctx context.Context, studentID uuid.UUID, courseID uuid.UUID) (*entity.Enrollment, *errorhandler.Response)
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Enrollment, *errorhandler.Response)
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.Enrollment, *errorhandler.Response)
}
