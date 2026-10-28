package values

import (
	"time"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/enrollment/domain/entity"
)

type EnrollmentResponse struct {
	ID         uuid.UUID `json:"id"`
	StudentID  uuid.UUID `json:"student_id"`
	CourseID   uuid.UUID `json:"course_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
}



func NewEnrollmentResponse(e *entity.Enrollment) *EnrollmentResponse {
	return &EnrollmentResponse{
		ID:         e.GetID(),
		StudentID:  e.GetStudentID(),
		CourseID:   e.GetCourseID(),
		EnrolledAt: e.GetEnrolledAt(),
	}
}
