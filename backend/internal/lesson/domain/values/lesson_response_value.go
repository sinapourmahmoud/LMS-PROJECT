package values

import (
	"time"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/lesson/domain/entity"
)

type LessonResponse struct {
	ID        uuid.UUID  `json:"id"`
	CourseID  uuid.UUID  `json:"course_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Position  int        `json:"position"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewLessonResponse(l *entity.Lesson) *LessonResponse {
	return &LessonResponse{
		ID:        l.GetID(),
		CourseID:  l.GetCourseID(),
		Title:     l.GetTitle(),
		Content:   l.GetContent(),
		Position:  l.GetPosition(),
		CreatedAt: l.GetCreatedAt(),
		UpdatedAt: l.GetUpdatedAt(),
	}
}