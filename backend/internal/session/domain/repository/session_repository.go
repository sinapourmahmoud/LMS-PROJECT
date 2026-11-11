package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/session/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) *errorhandler.Response
	GetByID(ctx context.Context, sessionID uuid.UUID) (*entity.Session, *errorhandler.Response)
	Update(ctx context.Context, session *entity.Session) *errorhandler.Response
	Delete(ctx context.Context, sessionID uuid.UUID) *errorhandler.Response
	// ListByCourse returns one page of the sessions of a course (ordered by start time) and the total number of sessions.
	ListByCourse(ctx context.Context, courseID uuid.UUID, limit int, offset int) ([]*entity.Session, int, *errorhandler.Response)
}
