package values

import (
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/course/domain/entity"
)

type CourseResponse struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	TeacherID   uuid.UUID     `json:"teacher_id"`
	Status      entity.Status `json:"status"`
}

func NewCourseResponse(course *entity.Course) *CourseResponse {

	return &CourseResponse{

		ID:          course.GetID(),
		Title:       course.GetTitle(),
		Description: course.GetDescription(),
		TeacherID:   course.GetTeacherID(),
		Status:      course.GetStatus(),
	}

}
