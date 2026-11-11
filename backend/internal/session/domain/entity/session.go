package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusScheduled Status = "scheduled"
	StatusLive      Status = "live"
	StatusEnded     Status = "ended"
)

// ErrInvalidState is returned when an action is not allowed in the current status of the session.
var ErrInvalidState = errors.New("action not allowed in the current session status")

type Session struct {
	id        uuid.UUID
	courseID  uuid.UUID
	title     string
	startsAt  time.Time
	status    Status
	startedAt *time.Time
	endedAt   *time.Time
	createdAt time.Time
	updatedAt *time.Time
}

func NewSession(id uuid.UUID, courseID uuid.UUID, title string, startsAt time.Time, status Status, startedAt *time.Time, endedAt *time.Time, createdAt time.Time, updatedAt *time.Time) *Session {
	return &Session{
		id:        id,
		courseID:  courseID,
		title:     title,
		startsAt:  startsAt,
		status:    status,
		startedAt: startedAt,
		endedAt:   endedAt,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (s *Session) GetID() uuid.UUID {
	return s.id
}

func (s *Session) GetCourseID() uuid.UUID {
	return s.courseID
}

func (s *Session) GetTitle() string {
	return s.title
}

func (s *Session) GetStartsAt() time.Time {
	return s.startsAt
}

func (s *Session) GetStatus() Status {
	return s.status
}

func (s *Session) GetStartedAt() *time.Time {
	return s.startedAt
}

func (s *Session) GetEndedAt() *time.Time {
	return s.endedAt
}

func (s *Session) GetCreatedAt() time.Time {
	return s.createdAt
}

func (s *Session) GetUpdatedAt() *time.Time {
	return s.updatedAt
}

// Update changes the title and the planned start time. Only a scheduled session can be changed.
func (s *Session) Update(title string, startsAt time.Time) error {
	if s.status != StatusScheduled {
		return ErrInvalidState
	}

	now := time.Now()
	s.title = title
	s.startsAt = startsAt
	s.updatedAt = &now
	return nil
}

// Start moves the session from scheduled to live.
func (s *Session) Start() error {
	if s.status != StatusScheduled {
		return ErrInvalidState
	}

	now := time.Now()
	s.status = StatusLive
	s.startedAt = &now
	s.updatedAt = &now
	return nil
}

// End moves the session from live to ended.
func (s *Session) End() error {
	if s.status != StatusLive {
		return ErrInvalidState
	}

	now := time.Now()
	s.status = StatusEnded
	s.endedAt = &now
	s.updatedAt = &now
	return nil
}
