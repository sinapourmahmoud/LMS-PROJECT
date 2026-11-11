package values

import (
	"time"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/session/domain/entity"
)

type SessionResponse struct {
	ID        uuid.UUID     `json:"id"`
	CourseID  uuid.UUID     `json:"course_id"`
	Title     string        `json:"title"`
	StartsAt  time.Time     `json:"starts_at"`
	Status    entity.Status `json:"status"`
	StartedAt *time.Time    `json:"started_at"`
	EndedAt   *time.Time    `json:"ended_at"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at"`
}

func NewSessionResponse(s *entity.Session) *SessionResponse {
	return &SessionResponse{
		ID:        s.GetID(),
		CourseID:  s.GetCourseID(),
		Title:     s.GetTitle(),
		StartsAt:  s.GetStartsAt(),
		Status:    s.GetStatus(),
		StartedAt: s.GetStartedAt(),
		EndedAt:   s.GetEndedAt(),
		CreatedAt: s.GetCreatedAt(),
		UpdatedAt: s.GetUpdatedAt(),
	}
}
