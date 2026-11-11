package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	enrollmentEntity "github.com/lms-project/backend/internal/enrollment/domain/entity"
	"github.com/lms-project/backend/internal/session/domain/entity"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
	"github.com/lms-project/backend/pkg/pagination"
)

func TestSessionListByCourse(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)
	sessions := []*entity.Session{newTestSession(courseID, entity.StatusScheduled), newTestSession(courseID, entity.StatusEnded)}
	studentID := uuid.New()
	enrollment := enrollmentEntity.NewEnrollment(uuid.New(), studentID, courseID, time.Now())

	tests := []struct {
		name        string
		requesterID uuid.UUID
		role        userEntity.Role
		enrollment  *enrollmentEntity.Enrollment
		wantStatus  int // 0 means allowed
	}{
		{name: "owner teacher", requesterID: teacherID, role: userEntity.RoleTeacher},
		{name: "admin", requesterID: uuid.New(), role: userEntity.RoleAdmin},
		{name: "enrolled student", requesterID: studentID, role: userEntity.RoleStudent, enrollment: enrollment},
		{name: "student who is not enrolled", requesterID: uuid.New(), role: userEntity.RoleStudent, wantStatus: 403},
		{name: "another teacher", requesterID: uuid.New(), role: userEntity.RoleTeacher, wantStatus: 403},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionRepo := &fakeSessionRepo{list: sessions, total: 2}
			svc := NewSessionQueryServices(sessionRepo, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{existing: tt.enrollment})

			res, resErr := svc.ListByCourse(context.Background(), courseID, pagination.Params{Page: 1, Limit: 10}, tt.requesterID, tt.role)

			if tt.wantStatus == 0 {
				if resErr != nil {
					t.Fatalf("expected success, got %+v", resErr)
				}
				if len(res.Items) != 2 || res.Total != 2 {
					t.Errorf("got %d items and total %d, want 2 and 2", len(res.Items), res.Total)
				}
				return
			}

			if resErr == nil || resErr.StatusCode != tt.wantStatus || resErr.Code != errorhandler.ForbiddenErrorCode {
				t.Fatalf("want 403 forbidden, got %+v", resErr)
			}
			if sessionRepo.listCalls != 0 {
				t.Error("the sessions must not be queried when access is denied")
			}
		})
	}

	t.Run("pagination values are passed to the repository", func(t *testing.T) {
		sessionRepo := &fakeSessionRepo{list: sessions, total: 25}
		svc := NewSessionQueryServices(sessionRepo, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		res, resErr := svc.ListByCourse(context.Background(), courseID, pagination.Params{Page: 3, Limit: 5}, teacherID, userEntity.RoleTeacher)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if sessionRepo.gotLimit != 5 || sessionRepo.gotOffset != 10 {
			t.Errorf("got limit=%d offset=%d, want 5 and 10", sessionRepo.gotLimit, sessionRepo.gotOffset)
		}
		if res.Page != 3 || res.Limit != 5 || res.Total != 25 {
			t.Errorf("got page=%d limit=%d total=%d, want 3, 5, 25", res.Page, res.Limit, res.Total)
		}
	})

	t.Run("no sessions returns an empty list, not null", func(t *testing.T) {
		svc := NewSessionQueryServices(&fakeSessionRepo{}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		res, resErr := svc.ListByCourse(context.Background(), courseID, pagination.Params{Page: 1, Limit: 10}, teacherID, userEntity.RoleTeacher)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res.Items == nil || len(res.Items) != 0 {
			t.Errorf("want an empty non-nil slice, got %v", res.Items)
		}
	})

	t.Run("course does not exist", func(t *testing.T) {
		svc := NewSessionQueryServices(&fakeSessionRepo{}, &fakeCourseRepo{course: nil}, &fakeEnrollmentRepo{})

		_, resErr := svc.ListByCourse(context.Background(), courseID, pagination.Params{Page: 1, Limit: 10}, teacherID, userEntity.RoleTeacher)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.CourseNotFoundErrorCode {
			t.Fatalf("want 404 course not found, got %+v", resErr)
		}
	})

	t.Run("enrollment lookup error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewSessionQueryServices(&fakeSessionRepo{}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{err: dbErr})

		_, resErr := svc.ListByCourse(context.Background(), courseID, pagination.Params{Page: 1, Limit: 10}, uuid.New(), userEntity.RoleStudent)

		if resErr != dbErr {
			t.Errorf("want the database error, got %+v", resErr)
		}
	})
}

func TestSessionGetByID(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)
	session := newTestSession(courseID, entity.StatusScheduled)
	studentID := uuid.New()
	enrollment := enrollmentEntity.NewEnrollment(uuid.New(), studentID, courseID, time.Now())

	tests := []struct {
		name        string
		requesterID uuid.UUID
		role        userEntity.Role
		enrollment  *enrollmentEntity.Enrollment
		wantStatus  int // 0 means allowed
	}{
		{name: "owner teacher", requesterID: teacherID, role: userEntity.RoleTeacher},
		{name: "admin", requesterID: uuid.New(), role: userEntity.RoleAdmin},
		{name: "enrolled student", requesterID: studentID, role: userEntity.RoleStudent, enrollment: enrollment},
		{name: "student who is not enrolled", requesterID: uuid.New(), role: userEntity.RoleStudent, wantStatus: 403},
		{name: "another teacher", requesterID: uuid.New(), role: userEntity.RoleTeacher, wantStatus: 403},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSessionQueryServices(&fakeSessionRepo{session: session}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{existing: tt.enrollment})

			res, resErr := svc.GetByID(context.Background(), session.GetID(), tt.requesterID, tt.role)

			if tt.wantStatus == 0 {
				if resErr != nil {
					t.Fatalf("expected success, got %+v", resErr)
				}
				if res.ID != session.GetID() {
					t.Errorf("got session %v, want %v", res.ID, session.GetID())
				}
				return
			}

			if resErr == nil || resErr.StatusCode != tt.wantStatus || resErr.Code != errorhandler.ForbiddenErrorCode {
				t.Fatalf("want 403 forbidden, got %+v", resErr)
			}
		})
	}

	t.Run("session does not exist", func(t *testing.T) {
		svc := NewSessionQueryServices(&fakeSessionRepo{}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		_, resErr := svc.GetByID(context.Background(), uuid.New(), teacherID, userEntity.RoleTeacher)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.SessionNotFoundErrorCode {
			t.Fatalf("want 404 session not found, got %+v", resErr)
		}
	})

	t.Run("session repository error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewSessionQueryServices(&fakeSessionRepo{getErr: dbErr}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		_, resErr := svc.GetByID(context.Background(), uuid.New(), teacherID, userEntity.RoleTeacher)

		if resErr != dbErr {
			t.Errorf("want the database error, got %+v", resErr)
		}
	})
}
