package entity

import (
	"github.com/google/uuid"
	"time"
)

type Enrollment struct {
	id         uuid.UUID
	studentID  uuid.UUID
	courseID   uuid.UUID
	enrolledAt time.Time
}

func (e *Enrollment) GetID() uuid.UUID {
	return e.id
}

func (e *Enrollment) GetStudentID() uuid.UUID {
	return e.studentID
}

func (e *Enrollment) GetCourseID() uuid.UUID {
	return e.courseID
}

func (e *Enrollment) GetEnrolledAt() time.Time {
	return e.enrolledAt
}

func NewEnrollment(id uuid.UUID, studentId uuid.UUID, courseId uuid.UUID, enrolledAt time.Time) *Enrollment {

	return &Enrollment{
		id:         id,
		studentID:  studentId,
		courseID:   courseId,
		enrolledAt: enrolledAt,
	}
}
