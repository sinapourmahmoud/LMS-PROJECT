package entity

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPublished Status = "published"
)

type Course struct {
	id          uuid.UUID
	title       string
	description string
	teacherID   uuid.UUID
	status      Status
	createdAt   time.Time
	updatedAt   *time.Time
}

func (c *Course) GetID() uuid.UUID {
	return c.id
}
func (c *Course) GetTitle() string {
	return c.title
}

func (c *Course) GetDescription() string {
	return c.description
}
func (c *Course) GetTeacherID() uuid.UUID {
	return c.teacherID
}
func (c *Course) GetStatus() Status {
	return c.status
}
func (c *Course) GetCreatedAt() time.Time {
	return c.createdAt
}
func (c *Course) GetUpdatedAt() *time.Time {
	return c.updatedAt
}

func NewCourse(
	id uuid.UUID, title string, description string,
	teacherID uuid.UUID,
	status Status,
	createdAt time.Time,
	updatedAt *time.Time) *Course {
	return &Course{
		id: id, title: title, description: description, teacherID: teacherID, status: status, createdAt: createdAt, updatedAt: updatedAt,
	}
}
