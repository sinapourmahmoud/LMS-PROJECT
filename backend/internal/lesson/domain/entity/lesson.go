package entity

import (
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	id        uuid.UUID
	courseId  uuid.UUID
	title     string
	content   string
	position  int
	createdAt time.Time
	updatedAt *time.Time
}

func (l *Lesson) GetID() uuid.UUID {
	return l.id
}

func (l *Lesson) GetCourseID() uuid.UUID {
	return l.courseId
}

func (l *Lesson) GetTitle() string {
	return l.title
}

func (l *Lesson) GetContent() string {
	return l.content
}
func (l *Lesson) GetPosition() int {
	return l.position
}
func (l *Lesson) GetCreatedAt() time.Time {
	return l.createdAt
}
func (l *Lesson) GetUpdatedAt() *time.Time {
	return l.updatedAt
}

func NewLesson(id uuid.UUID, courseId uuid.UUID, title string, content string, position int, createdAt time.Time, updatedAt *time.Time) *Lesson {
	return &Lesson{
		id:        id,
		courseId:  courseId,
		title:     title,
		content:   content,
		position:  position,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// for updating the process
func (l *Lesson) Update(title, content string, position int) {
	now := time.Now()
	l.title = title
	l.content = content
	l.position = position
	l.updatedAt = &now
}
