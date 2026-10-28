package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	courseEntity "github.com/lms-project/backend/internal/course/domain/entity"
	"github.com/lms-project/backend/internal/enrollment/domain/entity"
	"github.com/lms-project/backend/internal/enrollment/domain/values"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

func TestListByCourse(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := courseEntity.NewCourse(courseID, "title", "desc", teacherID, courseEntity.StatusPublished, time.Now(), nil)
	roster := []*values.RosterEntry{
		{StudentID: uuid.New(), Username: "ali", Email: "ali@test.com", EnrolledAt: time.Now()},
		{StudentID: uuid.New(), Username: "sara", Email: "sara@test.com", EnrolledAt: time.Now()},
	}

	tests := []struct {
		name          string
		course        *courseEntity.Course
		requesterID   uuid.UUID
		requesterRole userEntity.Role
		wantStatus    int // 0 means allowed
		wantCode      int
		wantLen       int
	}{
		{name: "teacher of the course", course: course, requesterID: teacherID, requesterRole: userEntity.RoleTeacher, wantLen: 2},
		{name: "admin", course: course, requesterID: uuid.New(), requesterRole: userEntity.RoleAdmin, wantLen: 2},
		{name: "another teacher", course: course, requesterID: uuid.New(), requesterRole: userEntity.RoleTeacher, wantStatus: 403, wantCode: errorhandler.ForbiddenErrorCode},
		{name: "student", course: course, requesterID: uuid.New(), requesterRole: userEntity.RoleStudent, wantStatus: 403, wantCode: errorhandler.ForbiddenErrorCode},
		{name: "course does not exist", course: nil, requesterID: teacherID, requesterRole: userEntity.RoleTeacher, wantStatus: 404, wantCode: errorhandler.CourseNotFoundErrorCode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enrollmentRepo := &fakeEnrollmentRepo{listRoster: roster}
			svc := NewEnrollmentQueryServices(enrollmentRepo, &fakeCourseRepo{course: tt.course})

			res, resErr := svc.ListByCourse(context.Background(), courseID, tt.requesterID, tt.requesterRole)

			if tt.wantStatus == 0 {
				if resErr != nil {
					t.Fatalf("expected success, got error: %+v", resErr)
				}
				if len(res) != tt.wantLen {
					t.Errorf("got %d enrollments, want %d", len(res), tt.wantLen)
				}
				return
			}

			if resErr == nil {
				t.Fatalf("expected status %d, got success", tt.wantStatus)
			}
			if resErr.StatusCode != tt.wantStatus || resErr.Code != tt.wantCode {
				t.Errorf("got status %d code %d, want status %d code %d", resErr.StatusCode, resErr.Code, tt.wantStatus, tt.wantCode)
			}
			if enrollmentRepo.listCalls != 0 {
				t.Error("the repository must not be queried when access is denied")
			}
		})
	}
}

func TestListByCourseRoster(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := courseEntity.NewCourse(courseID, "title", "desc", teacherID, courseEntity.StatusPublished, time.Now(), nil)

	t.Run("returns the roster from the repository", func(t *testing.T) {
		want := []*values.RosterEntry{{StudentID: uuid.New(), Username: "ali", Email: "ali@test.com", EnrolledAt: time.Now()}}
		svc := NewEnrollmentQueryServices(&fakeEnrollmentRepo{listRoster: want}, &fakeCourseRepo{course: course})

		res, resErr := svc.ListByCourse(context.Background(), courseID, teacherID, userEntity.RoleTeacher)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if len(res) != 1 || res[0].Username != "ali" || res[0].Email != "ali@test.com" {
			t.Errorf("unexpected roster: %+v", res)
		}
	})

	t.Run("repository error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewEnrollmentQueryServices(&fakeEnrollmentRepo{rosterErr: dbErr}, &fakeCourseRepo{course: course})

		res, resErr := svc.ListByCourse(context.Background(), courseID, teacherID, userEntity.RoleTeacher)

		if res != nil || resErr == nil || resErr.StatusCode != 500 {
			t.Errorf("want the 500 error, got res=%v err=%+v", res, resErr)
		}
	})
}

func TestListByStudent(t *testing.T) {
	studentID := uuid.New()

	t.Run("returns the student's enrollments", func(t *testing.T) {
		repo := &fakeEnrollmentRepo{list: []*entity.Enrollment{entity.NewEnrollment(uuid.New(), studentID, uuid.New(), time.Now())}}
		svc := NewEnrollmentQueryServices(repo, &fakeCourseRepo{})

		res, resErr := svc.ListByStudent(context.Background(), studentID)

		if resErr != nil || len(res) != 1 {
			t.Fatalf("got %v, %+v; want 1 enrollment", res, resErr)
		}
	})

	t.Run("no enrollments gives an empty slice, not nil", func(t *testing.T) {
		svc := NewEnrollmentQueryServices(&fakeEnrollmentRepo{}, &fakeCourseRepo{})

		res, resErr := svc.ListByStudent(context.Background(), studentID)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res == nil || len(res) != 0 {
			t.Errorf("want empty non-nil slice (JSON [] not null), got %#v", res)
		}
	})

	t.Run("repository error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewEnrollmentQueryServices(&fakeEnrollmentRepo{listErr: dbErr}, &fakeCourseRepo{})

		_, resErr := svc.ListByStudent(context.Background(), studentID)

		if resErr == nil || resErr.StatusCode != 500 {
			t.Errorf("want the 500 error, got %+v", resErr)
		}
	})
}
