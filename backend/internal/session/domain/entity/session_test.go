package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestSession(status Status) *Session {
	return NewSession(uuid.New(), uuid.New(), "Live Q&A", time.Now().Add(time.Hour), status, nil, nil, time.Now(), nil)
}

func TestSessionStart(t *testing.T) {
	t.Run("a scheduled session becomes live", func(t *testing.T) {
		s := newTestSession(StatusScheduled)

		if err := s.Start(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.GetStatus() != StatusLive || s.GetStartedAt() == nil || s.GetUpdatedAt() == nil {
			t.Errorf("want live with startedAt and updatedAt set, got %+v", s)
		}
	})

	for _, status := range []Status{StatusLive, StatusEnded} {
		t.Run("cannot start a "+string(status)+" session", func(t *testing.T) {
			s := newTestSession(status)

			if err := s.Start(); !errors.Is(err, ErrInvalidState) {
				t.Errorf("want ErrInvalidState, got %v", err)
			}
			if s.GetStatus() != status {
				t.Error("the status must not change")
			}
		})
	}
}

func TestSessionEnd(t *testing.T) {
	t.Run("a live session becomes ended", func(t *testing.T) {
		s := newTestSession(StatusLive)

		if err := s.End(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.GetStatus() != StatusEnded || s.GetEndedAt() == nil {
			t.Errorf("want ended with endedAt set, got %+v", s)
		}
	})

	for _, status := range []Status{StatusScheduled, StatusEnded} {
		t.Run("cannot end a "+string(status)+" session", func(t *testing.T) {
			s := newTestSession(status)

			if err := s.End(); !errors.Is(err, ErrInvalidState) {
				t.Errorf("want ErrInvalidState, got %v", err)
			}
			if s.GetStatus() != status {
				t.Error("the status must not change")
			}
		})
	}
}

func TestSessionUpdate(t *testing.T) {
	newTime := time.Now().Add(48 * time.Hour)

	t.Run("a scheduled session can be changed", func(t *testing.T) {
		s := newTestSession(StatusScheduled)

		if err := s.Update("New title", newTime); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.GetTitle() != "New title" || !s.GetStartsAt().Equal(newTime) || s.GetUpdatedAt() == nil {
			t.Errorf("session was not updated: %+v", s)
		}
	})

	for _, status := range []Status{StatusLive, StatusEnded} {
		t.Run("a "+string(status)+" session cannot be changed", func(t *testing.T) {
			s := newTestSession(status)

			if err := s.Update("New title", newTime); !errors.Is(err, ErrInvalidState) {
				t.Errorf("want ErrInvalidState, got %v", err)
			}
			if s.GetTitle() != "Live Q&A" {
				t.Error("the title must not change")
			}
		})
	}
}
