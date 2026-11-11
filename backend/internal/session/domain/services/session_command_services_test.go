package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/session/domain/entity"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

func TestSessionCreate(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)
	future := time.Now().Add(24 * time.Hour)

	t.Run("owner teacher schedules a session", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		res, resErr := svc.Create(context.Background(), teacherID, courseID, "Q&A", future)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res.Title != "Q&A" || res.CourseID != courseID || res.Status != entity.StatusScheduled {
			t.Errorf("unexpected response: %+v", res)
		}
		if sessionRepo.createCalls != 1 || sessionRepo.created.GetCourseID() != courseID {
			t.Error("the session must be saved once, in the right course")
		}
	})

	t.Run("another teacher is forbidden", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		res, resErr := svc.Create(context.Background(), uuid.New(), courseID, "Q&A", future)

		if res != nil || resErr == nil || resErr.StatusCode != 403 || resErr.Code != errorhandler.ForbiddenErrorCode {
			t.Fatalf("want 403 forbidden, got res=%v err=%+v", res, resErr)
		}
		if sessionRepo.createCalls != 0 {
			t.Error("nothing must be saved when access is denied")
		}
	})

	t.Run("start time in the past is rejected", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		_, resErr := svc.Create(context.Background(), teacherID, courseID, "Q&A", time.Now().Add(-time.Hour))

		if resErr == nil || resErr.StatusCode != 400 || resErr.Code != errorhandler.ValidationErrorCode {
			t.Fatalf("want 400 validation error, got %+v", resErr)
		}
		if sessionRepo.createCalls != 0 {
			t.Error("nothing must be saved when the time is in the past")
		}
	})

	t.Run("course does not exist", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: nil})

		_, resErr := svc.Create(context.Background(), teacherID, courseID, "Q&A", future)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.CourseNotFoundErrorCode {
			t.Fatalf("want 404 course not found, got %+v", resErr)
		}
		if sessionRepo.createCalls != 0 {
			t.Error("nothing must be saved when the course does not exist")
		}
	})

	t.Run("course repository error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewSessionCommandServices(&fakeSessionRepo{}, &fakeCourseRepo{err: dbErr})

		_, resErr := svc.Create(context.Background(), teacherID, courseID, "Q&A", future)

		if resErr != dbErr {
			t.Errorf("want the database error, got %+v", resErr)
		}
	})

	t.Run("session repository error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewSessionCommandServices(&fakeSessionRepo{createErr: dbErr}, &fakeCourseRepo{course: course})

		res, resErr := svc.Create(context.Background(), teacherID, courseID, "Q&A", future)

		if res != nil || resErr != dbErr {
			t.Errorf("want the database error and no response, got res=%v err=%+v", res, resErr)
		}
	})
}

func TestSessionUpdate(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)
	future := time.Now().Add(48 * time.Hour)

	t.Run("owner teacher updates a scheduled session", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusScheduled)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		res, resErr := svc.Update(context.Background(), uuid.New(), "New title", future, teacherID)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res.Title != "New title" || res.UpdatedAt == nil {
			t.Errorf("unexpected response: %+v", res)
		}
		if sessionRepo.updateCalls != 1 {
			t.Error("the session must be saved once")
		}
	})

	t.Run("another teacher is forbidden", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusScheduled)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		_, resErr := svc.Update(context.Background(), uuid.New(), "New title", future, uuid.New())

		if resErr == nil || resErr.StatusCode != 403 {
			t.Fatalf("want 403, got %+v", resErr)
		}
		if sessionRepo.updateCalls != 0 {
			t.Error("nothing must be saved when access is denied")
		}
	})

	t.Run("session does not exist", func(t *testing.T) {
		svc := NewSessionCommandServices(&fakeSessionRepo{}, &fakeCourseRepo{course: course})

		_, resErr := svc.Update(context.Background(), uuid.New(), "New title", future, teacherID)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.SessionNotFoundErrorCode {
			t.Fatalf("want 404 session not found, got %+v", resErr)
		}
	})

	t.Run("start time in the past is rejected", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusScheduled)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		_, resErr := svc.Update(context.Background(), uuid.New(), "New title", time.Now().Add(-time.Hour), teacherID)

		if resErr == nil || resErr.StatusCode != 400 {
			t.Fatalf("want 400, got %+v", resErr)
		}
		if sessionRepo.updateCalls != 0 {
			t.Error("nothing must be saved")
		}
	})

	t.Run("a live session cannot be updated", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusLive)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		_, resErr := svc.Update(context.Background(), uuid.New(), "New title", future, teacherID)

		if resErr == nil || resErr.StatusCode != 409 || resErr.Code != errorhandler.SessionInvalidStateErrorCode {
			t.Fatalf("want 409 invalid state, got %+v", resErr)
		}
		if sessionRepo.updateCalls != 0 {
			t.Error("nothing must be saved")
		}
	})
}

func TestSessionStartAndEnd(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)

	t.Run("owner teacher starts a scheduled session", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusScheduled)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		res, resErr := svc.Start(context.Background(), uuid.New(), teacherID)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res.Status != entity.StatusLive || res.StartedAt == nil {
			t.Errorf("want a live session with startedAt, got %+v", res)
		}
		if sessionRepo.updateCalls != 1 {
			t.Error("the session must be saved once")
		}
	})

	t.Run("starting a live session returns 409", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusLive)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		_, resErr := svc.Start(context.Background(), uuid.New(), teacherID)

		if resErr == nil || resErr.StatusCode != 409 || resErr.Code != errorhandler.SessionInvalidStateErrorCode {
			t.Fatalf("want 409 invalid state, got %+v", resErr)
		}
		if sessionRepo.updateCalls != 0 {
			t.Error("nothing must be saved")
		}
	})

	t.Run("another teacher cannot start", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusScheduled)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		_, resErr := svc.Start(context.Background(), uuid.New(), uuid.New())

		if resErr == nil || resErr.StatusCode != 403 {
			t.Fatalf("want 403, got %+v", resErr)
		}
		if sessionRepo.updateCalls != 0 {
			t.Error("nothing must be saved")
		}
	})

	t.Run("owner teacher ends a live session", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusLive)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		res, resErr := svc.End(context.Background(), uuid.New(), teacherID)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res.Status != entity.StatusEnded || res.EndedAt == nil {
			t.Errorf("want an ended session with endedAt, got %+v", res)
		}
		if sessionRepo.updateCalls != 1 {
			t.Error("the session must be saved once")
		}
	})

	t.Run("ending a scheduled session returns 409", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusScheduled)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		_, resErr := svc.End(context.Background(), uuid.New(), teacherID)

		if resErr == nil || resErr.StatusCode != 409 || resErr.Code != errorhandler.SessionInvalidStateErrorCode {
			t.Fatalf("want 409 invalid state, got %+v", resErr)
		}
		if sessionRepo.updateCalls != 0 {
			t.Error("nothing must be saved")
		}
	})

	t.Run("another teacher cannot end", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusLive)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		_, resErr := svc.End(context.Background(), uuid.New(), uuid.New())

		if resErr == nil || resErr.StatusCode != 403 {
			t.Fatalf("want 403, got %+v", resErr)
		}
	})

	t.Run("update error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusScheduled), updateErr: dbErr}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

		res, resErr := svc.Start(context.Background(), uuid.New(), teacherID)

		if res != nil || resErr != dbErr {
			t.Errorf("want the database error and no response, got res=%v err=%+v", res, resErr)
		}
	})
}

func TestSessionDelete(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)

	tests := []struct {
		name        string
		requesterID uuid.UUID
		role        userEntity.Role
		status      entity.Status
		wantStatus  int // 0 means deleted
		wantCode    int
	}{
		{name: "owner teacher deletes a scheduled session", requesterID: teacherID, role: userEntity.RoleTeacher, status: entity.StatusScheduled},
		{name: "owner teacher deletes an ended session", requesterID: teacherID, role: userEntity.RoleTeacher, status: entity.StatusEnded},
		{name: "admin deletes a session", requesterID: uuid.New(), role: userEntity.RoleAdmin, status: entity.StatusScheduled},
		{name: "another teacher is forbidden", requesterID: uuid.New(), role: userEntity.RoleTeacher, status: entity.StatusScheduled, wantStatus: 403, wantCode: errorhandler.ForbiddenErrorCode},
		{name: "a live session cannot be deleted", requesterID: teacherID, role: userEntity.RoleTeacher, status: entity.StatusLive, wantStatus: 409, wantCode: errorhandler.SessionInvalidStateErrorCode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, tt.status)}
			svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{course: course})

			resErr := svc.Delete(context.Background(), uuid.New(), tt.requesterID, tt.role)

			if tt.wantStatus == 0 {
				if resErr != nil {
					t.Fatalf("unexpected error: %+v", resErr)
				}
				if sessionRepo.deleteCalls != 1 {
					t.Error("the session must be deleted once")
				}
				return
			}

			if resErr == nil || resErr.StatusCode != tt.wantStatus || resErr.Code != tt.wantCode {
				t.Fatalf("want status %d code %d, got %+v", tt.wantStatus, tt.wantCode, resErr)
			}
			if sessionRepo.deleteCalls != 0 {
				t.Error("nothing must be deleted")
			}
		})
	}

	t.Run("session does not exist", func(t *testing.T) {
		svc := NewSessionCommandServices(&fakeSessionRepo{}, &fakeCourseRepo{course: course})

		resErr := svc.Delete(context.Background(), uuid.New(), teacherID, userEntity.RoleTeacher)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.SessionNotFoundErrorCode {
			t.Fatalf("want 404 session not found, got %+v", resErr)
		}
	})

	t.Run("a failed course lookup does not delete anything", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		sessionRepo := &fakeSessionRepo{session: newTestSession(courseID, entity.StatusScheduled)}
		svc := NewSessionCommandServices(sessionRepo, &fakeCourseRepo{err: dbErr})

		resErr := svc.Delete(context.Background(), uuid.New(), teacherID, userEntity.RoleTeacher)

		if resErr != dbErr {
			t.Errorf("want the database error, got %+v", resErr)
		}
		if sessionRepo.deleteCalls != 0 {
			t.Error("nothing must be deleted")
		}
	})
}
