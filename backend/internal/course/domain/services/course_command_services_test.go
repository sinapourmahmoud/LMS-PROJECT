package services
 
import (
	"context"
	"testing"
	"time"
 
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/course/domain/entity"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)
 
var dbErr = &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
 
func newCourse(teacherID uuid.UUID, status entity.Status) *entity.Course {
	return entity.NewCourse(uuid.New(), "title", "desc", teacherID, status, time.Now(), nil)
}
 
func newUser(role userEntity.Role) *userEntity.User {
	return userEntity.NewUser(uuid.New(), "name", "name@example.com", "hash", role, time.Now(), nil)
}
 
func TestCreateCourse(t *testing.T) {
	teacher := newUser(userEntity.RoleTeacher)
 
	t.Run("success: course is created as pending", func(t *testing.T) {
		courseRepo := &fakeCourseRepo{}
		svc := NewCourseCommandServices(courseRepo, &fakeUserRepo{user: teacher})
 
		res, resErr := svc.Create(context.Background(), "Go", "learn Go", teacher.GetID())
 
		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res.Status != entity.StatusPending {
			t.Errorf("status = %q, want pending (an admin must approve it)", res.Status)
		}
		if res.TeacherID != teacher.GetID() || res.Title != "Go" {
			t.Errorf("response has wrong data: %+v", res)
		}
		if courseRepo.createCalls != 1 {
			t.Errorf("Create called %d times, want 1", courseRepo.createCalls)
		}
	})
 
	t.Run("teacher does not exist", func(t *testing.T) {
		courseRepo := &fakeCourseRepo{}
		svc := NewCourseCommandServices(courseRepo, &fakeUserRepo{user: nil})
 
		_, resErr := svc.Create(context.Background(), "Go", "learn Go", uuid.New())
 
		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.UserNotFoundErrorCode {
			t.Fatalf("want 404 user not found, got %+v", resErr)
		}
		if courseRepo.createCalls != 0 {
			t.Error("Create must not be called when the teacher does not exist")
		}
	})
 
	t.Run("user repository fails", func(t *testing.T) {
		svc := NewCourseCommandServices(&fakeCourseRepo{}, &fakeUserRepo{err: dbErr})
 
		_, resErr := svc.Create(context.Background(), "Go", "learn Go", teacher.GetID())
 
		if resErr == nil || resErr.StatusCode != 500 {
			t.Errorf("want the 500 error passed through, got %+v", resErr)
		}
	})
 
	t.Run("course repository fails on create", func(t *testing.T) {
		svc := NewCourseCommandServices(&fakeCourseRepo{err: dbErr}, &fakeUserRepo{user: teacher})
 
		res, resErr := svc.Create(context.Background(), "Go", "learn Go", teacher.GetID())
 
		if resErr == nil || resErr.StatusCode != 500 {
			t.Errorf("want the 500 error, got %+v", resErr)
		}
		if res != nil {
			t.Errorf("want nil response on error, got %+v", res)
		}
	})
}
 
func TestApproveCourse(t *testing.T) {
	t.Run("success: status becomes published", func(t *testing.T) {
		courseRepo := &fakeCourseRepo{course: newCourse(uuid.New(), entity.StatusPending)}
		svc := NewCourseCommandServices(courseRepo, &fakeUserRepo{})
 
		resErr := svc.Approve(context.Background(), uuid.New())
 
		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if courseRepo.updateCalls != 1 || courseRepo.updatedStatus != entity.StatusPublished {
			t.Errorf("UpdateStatus calls=%d status=%q, want 1 call with published", courseRepo.updateCalls, courseRepo.updatedStatus)
		}
	})
 
	t.Run("course does not exist", func(t *testing.T) {
		courseRepo := &fakeCourseRepo{}
		svc := NewCourseCommandServices(courseRepo, &fakeUserRepo{})
 
		resErr := svc.Approve(context.Background(), uuid.New())
 
		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.CourseNotFoundErrorCode {
			t.Fatalf("want 404 course not found, got %+v", resErr)
		}
		if courseRepo.updateCalls != 0 {
			t.Error("UpdateStatus must not be called for a missing course")
		}
	})
 
	t.Run("repository error is passed through", func(t *testing.T) {
		svc := NewCourseCommandServices(&fakeCourseRepo{err: dbErr}, &fakeUserRepo{})
 
		resErr := svc.Approve(context.Background(), uuid.New())
 
		if resErr == nil || resErr.StatusCode != 500 {
			t.Errorf("want the 500 error, got %+v", resErr)
		}
	})
}
 
func TestDeleteCourse(t *testing.T) {
	ownerID := uuid.New()
	course := newCourse(ownerID, entity.StatusPublished)
 
	tests := []struct {
		name        string
		course      *entity.Course
		repoErr     *errorhandler.Response
		requesterID uuid.UUID
		role        userEntity.Role
		wantStatus  int // 0 means the delete is allowed
		wantCode    int
		wantDeletes int
	}{
		{name: "owner teacher", course: course, requesterID: ownerID, role: userEntity.RoleTeacher, wantDeletes: 1},
		{name: "admin", course: course, requesterID: uuid.New(), role: userEntity.RoleAdmin, wantDeletes: 1},
		{name: "another teacher", course: course, requesterID: uuid.New(), role: userEntity.RoleTeacher, wantStatus: 403, wantCode: errorhandler.ForbiddenErrorCode},
		{name: "course does not exist", course: nil, requesterID: ownerID, role: userEntity.RoleTeacher, wantStatus: 404, wantCode: errorhandler.CourseNotFoundErrorCode},
		{name: "repository fails", course: course, repoErr: dbErr, requesterID: ownerID, role: userEntity.RoleTeacher, wantStatus: 500, wantCode: errorhandler.DatabaseErrorCode},
	}
 
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseRepo := &fakeCourseRepo{course: tt.course, err: tt.repoErr}
			svc := NewCourseCommandServices(courseRepo, &fakeUserRepo{})
 
			resErr := svc.Delete(context.Background(), uuid.New(), tt.requesterID, tt.role)
 
			if tt.wantStatus == 0 {
				if resErr != nil {
					t.Fatalf("expected success, got %+v", resErr)
				}
			} else {
				if resErr == nil {
					t.Fatalf("expected status %d, got success", tt.wantStatus)
				}
				if resErr.StatusCode != tt.wantStatus || resErr.Code != tt.wantCode {
					t.Errorf("got status %d code %d, want status %d code %d", resErr.StatusCode, resErr.Code, tt.wantStatus, tt.wantCode)
				}
			}
			if courseRepo.deleteCalls != tt.wantDeletes {
				t.Errorf("Delete called %d times, want %d", courseRepo.deleteCalls, tt.wantDeletes)
			}
		})
	}
}
 