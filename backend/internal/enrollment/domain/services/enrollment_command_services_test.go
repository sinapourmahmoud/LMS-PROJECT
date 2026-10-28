package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	courseEntity "github.com/lms-project/backend/internal/course/domain/entity"
	"github.com/lms-project/backend/internal/enrollment/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

func newCourse(status courseEntity.Status) *courseEntity.Course {
	return courseEntity.NewCourse(uuid.New(), "title", "desc", uuid.New(), status, time.Now(), nil)
}

func TestEnroll(t *testing.T) {
	studentID := uuid.New()
	courseID := uuid.New()
	dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}

	tests := []struct {
		name            string
		courseRepo      *fakeCourseRepo
		enrollmentRepo  *fakeEnrollmentRepo
		wantStatus      int // 0 means "no error expected"
		wantCode        int
		wantCreateCalls int
	}{
		{
			name:            "success",
			courseRepo:      &fakeCourseRepo{course: newCourse(courseEntity.StatusPublished)},
			enrollmentRepo:  &fakeEnrollmentRepo{},
			wantCreateCalls: 1,
		},
		{
			name:           "course does not exist",
			courseRepo:     &fakeCourseRepo{course: nil},
			enrollmentRepo: &fakeEnrollmentRepo{},
			wantStatus:     404,
			wantCode:       errorhandler.CourseNotFoundErrorCode,
		},
		{
			name:           "course is still pending",
			courseRepo:     &fakeCourseRepo{course: newCourse(courseEntity.StatusPending)},
			enrollmentRepo: &fakeEnrollmentRepo{},
			wantStatus:     404,
			wantCode:       errorhandler.CourseNotFoundErrorCode,
		},
		{
			name:       "already enrolled",
			courseRepo: &fakeCourseRepo{course: newCourse(courseEntity.StatusPublished)},
			enrollmentRepo: &fakeEnrollmentRepo{
				existing: entity.NewEnrollment(uuid.New(), studentID, courseID, time.Now()),
			},
			wantStatus: 409,
			wantCode:   errorhandler.AlreadyEnrolledErrorCode,
		},
		{
			name:           "course repository fails: error passes through, not a 404",
			courseRepo:     &fakeCourseRepo{err: dbErr},
			enrollmentRepo: &fakeEnrollmentRepo{},
			wantStatus:     500,
			wantCode:       errorhandler.DatabaseErrorCode,
		},
		{
			name:           "duplicate check fails",
			courseRepo:     &fakeCourseRepo{course: newCourse(courseEntity.StatusPublished)},
			enrollmentRepo: &fakeEnrollmentRepo{findErr: dbErr},
			wantStatus:     500,
			wantCode:       errorhandler.DatabaseErrorCode,
		},
		{
			name:            "create fails",
			courseRepo:      &fakeCourseRepo{course: newCourse(courseEntity.StatusPublished)},
			enrollmentRepo:  &fakeEnrollmentRepo{createErr: dbErr},
			wantStatus:      500,
			wantCode:        errorhandler.DatabaseErrorCode,
			wantCreateCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEnrollmentCommandServices(tt.enrollmentRepo, tt.courseRepo)

			res, resErr := svc.Enroll(context.Background(), studentID, courseID)

			if tt.wantStatus == 0 {
				if resErr != nil {
					t.Fatalf("expected success, got error: %+v", resErr)
				}
				if res == nil {
					t.Fatal("expected a response, got nil")
				}
				if res.StudentID != studentID || res.CourseID != courseID {
					t.Errorf("response has wrong ids: %+v", res)
				}
			} else {
				if resErr == nil {
					t.Fatalf("expected status %d, got success", tt.wantStatus)
				}
				if resErr.StatusCode != tt.wantStatus {
					t.Errorf("status = %d, want %d", resErr.StatusCode, tt.wantStatus)
				}
				if resErr.Code != tt.wantCode {
					t.Errorf("code = %d, want %d", resErr.Code, tt.wantCode)
				}
				if res != nil {
					t.Errorf("expected nil response on error, got %+v", res)
				}
			}

			if tt.enrollmentRepo.createCalls != tt.wantCreateCalls {
				t.Errorf("Create called %d times, want %d", tt.enrollmentRepo.createCalls, tt.wantCreateCalls)
			}
		})
	}
}

func TestUnEnroll(t *testing.T) {
	studentID := uuid.New()
	courseID := uuid.New()
	dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
	enrolled := entity.NewEnrollment(uuid.New(), studentID, courseID, time.Now())

	tests := []struct {
		name            string
		enrollmentRepo  *fakeEnrollmentRepo
		wantStatus      int // 0 means success
		wantCode        int
		wantDeleteCalls int
	}{
		{
			name:            "success",
			enrollmentRepo:  &fakeEnrollmentRepo{existing: enrolled},
			wantDeleteCalls: 1,
		},
		{
			name:           "not enrolled",
			enrollmentRepo: &fakeEnrollmentRepo{existing: nil},
			wantStatus:     404,
			wantCode:       errorhandler.NotEnrolledErrorCode,
		},
		{
			name:           "find fails",
			enrollmentRepo: &fakeEnrollmentRepo{findErr: dbErr},
			wantStatus:     500,
			wantCode:       errorhandler.DatabaseErrorCode,
		},
		{
			name:            "delete fails",
			enrollmentRepo:  &fakeEnrollmentRepo{existing: enrolled, deleteErr: dbErr},
			wantStatus:      500,
			wantCode:        errorhandler.DatabaseErrorCode,
			wantDeleteCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEnrollmentCommandServices(tt.enrollmentRepo, &fakeCourseRepo{})

			resErr := svc.UnEnroll(context.Background(), studentID, courseID)

			if tt.wantStatus == 0 {
				if resErr != nil {
					t.Fatalf("expected success, got error: %+v", resErr)
				}
			} else {
				if resErr == nil {
					t.Fatalf("expected status %d, got success", tt.wantStatus)
				}
				if resErr.StatusCode != tt.wantStatus {
					t.Errorf("status = %d, want %d", resErr.StatusCode, tt.wantStatus)
				}
				if resErr.Code != tt.wantCode {
					t.Errorf("code = %d, want %d", resErr.Code, tt.wantCode)
				}
			}

			if tt.enrollmentRepo.deleteCalls != tt.wantDeleteCalls {
				t.Errorf("Delete called %d times, want %d", tt.enrollmentRepo.deleteCalls, tt.wantDeleteCalls)
			}
		})
	}
}
