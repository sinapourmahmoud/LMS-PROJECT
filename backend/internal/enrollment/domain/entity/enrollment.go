package entity

import (
	"time"
	"github.com/google/uuid"
)

type Enrollment struct {
	id         uuid.UUID
	studentId  uuid.UUID
	courseId   uuid.UUID
	enrolledAt time.Time
}

func (e *Enrollment) GetID() uuid.UUID {
	return e.id
}

func (e *Enrollment) GetStudentID() uuid.UUID {
	return e.studentId
}

func (e *Enrollment) GetCourseID() uuid.UUID {
	return e.courseId
}

func (e *Enrollment) GetEnrlledAt() time.Time {
	return e.enrolledAt
}

func NewEnrollment(id uuid.UUID, studentId uuid.UUID, courseId uuid.UUID, enrolledAt time.Time) *Enrollment {

	return &Enrollment{
		id:         id,
		studentId:  studentId,
		courseId:   courseId,
		enrolledAt: enrolledAt,
	}
}
