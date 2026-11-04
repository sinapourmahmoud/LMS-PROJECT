package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

func TestLessonCreate(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)

	t.Run("owner teacher creates a lesson", func(t *testing.T) {
		lessonRepo := &fakeLessonRepo{}
		svc := NewLessonCommandServices(lessonRepo, &fakeCourseRepo{course: course})

		res, resErr := svc.Create(context.Background(), teacherID, courseID, "Intro", "Hello", 1)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res.Title != "Intro" || res.Content != "Hello" || res.Position != 1 || res.CourseID != courseID {
			t.Errorf("unexpected response: %+v", res)
		}
		if lessonRepo.createCalls != 1 || lessonRepo.created.GetCourseID() != courseID {
			t.Error("the lesson must be saved once, in the right course")
		}
	})

	t.Run("another teacher is forbidden", func(t *testing.T) {
		lessonRepo := &fakeLessonRepo{}
		svc := NewLessonCommandServices(lessonRepo, &fakeCourseRepo{course: course})

		res, resErr := svc.Create(context.Background(), uuid.New(), courseID, "Intro", "Hello", 1)

		if res != nil || resErr == nil || resErr.StatusCode != 403 || resErr.Code != errorhandler.ForbiddenErrorCode {
			t.Fatalf("want 403 forbidden, got res=%v err=%+v", res, resErr)
		}
		if lessonRepo.createCalls != 0 {
			t.Error("nothing must be saved when access is denied")
		}
	})

	t.Run("course does not exist", func(t *testing.T) {
		lessonRepo := &fakeLessonRepo{}
		svc := NewLessonCommandServices(lessonRepo, &fakeCourseRepo{course: nil})

		_, resErr := svc.Create(context.Background(), teacherID, courseID, "Intro", "Hello", 1)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.CourseNotFoundErrorCode {
			t.Fatalf("want 404 course not found, got %+v", resErr)
		}
		if lessonRepo.createCalls != 0 {
			t.Error("nothing must be saved when the course does not exist")
		}
	})

	t.Run("course repository error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewLessonCommandServices(&fakeLessonRepo{}, &fakeCourseRepo{err: dbErr})

		_, resErr := svc.Create(context.Background(), teacherID, courseID, "Intro", "Hello", 1)

		if resErr != dbErr {
			t.Errorf("want the database error, got %+v", resErr)
		}
	})

	t.Run("duplicate position is passed through as 409", func(t *testing.T) {
		taken := &errorhandler.Response{Code: errorhandler.LessonPositionTakenErrorCode, StatusCode: 409}
		svc := NewLessonCommandServices(&fakeLessonRepo{createErr: taken}, &fakeCourseRepo{course: course})

		res, resErr := svc.Create(context.Background(), teacherID, courseID, "Intro", "Hello", 1)

		if res != nil || resErr != taken {
			t.Errorf("want the 409 error and no response, got res=%v err=%+v", res, resErr)
		}
	})
}

func TestLessonUpdate(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)

	t.Run("owner teacher updates a lesson", func(t *testing.T) {
		lesson := newTestLesson(courseID, 1)
		lessonRepo := &fakeLessonRepo{lesson: lesson}
		svc := NewLessonCommandServices(lessonRepo, &fakeCourseRepo{course: course})

		res, resErr := svc.Update(context.Background(), lesson.GetID(), "New title", "New content", 3, teacherID)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res.Title != "New title" || res.Content != "New content" || res.Position != 3 {
			t.Errorf("unexpected response: %+v", res)
		}
		if res.UpdatedAt == nil {
			t.Error("updated_at must be set after an update")
		}
		if lessonRepo.updateCalls != 1 || lessonRepo.updated.GetTitle() != "New title" {
			t.Error("the changed lesson must be saved once")
		}
	})

	t.Run("another teacher is forbidden", func(t *testing.T) {
		lesson := newTestLesson(courseID, 1)
		lessonRepo := &fakeLessonRepo{lesson: lesson}
		svc := NewLessonCommandServices(lessonRepo, &fakeCourseRepo{course: course})

		_, resErr := svc.Update(context.Background(), lesson.GetID(), "New title", "New content", 3, uuid.New())

		if resErr == nil || resErr.StatusCode != 403 {
			t.Fatalf("want 403, got %+v", resErr)
		}
		if lessonRepo.updateCalls != 0 {
			t.Error("nothing must be saved when access is denied")
		}
	})

	t.Run("lesson does not exist", func(t *testing.T) {
		svc := NewLessonCommandServices(&fakeLessonRepo{lesson: nil}, &fakeCourseRepo{course: course})

		_, resErr := svc.Update(context.Background(), uuid.New(), "t", "c", 1, teacherID)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.LessonNotFoundErrorCode {
			t.Fatalf("want 404 lesson not found, got %+v", resErr)
		}
	})

	t.Run("duplicate position is passed through as 409", func(t *testing.T) {
		lesson := newTestLesson(courseID, 1)
		taken := &errorhandler.Response{Code: errorhandler.LessonPositionTakenErrorCode, StatusCode: 409}
		svc := NewLessonCommandServices(&fakeLessonRepo{lesson: lesson, updateErr: taken}, &fakeCourseRepo{course: course})

		res, resErr := svc.Update(context.Background(), lesson.GetID(), "t", "c", 2, teacherID)

		if res != nil || resErr != taken {
			t.Errorf("want the 409 error and no response, got res=%v err=%+v", res, resErr)
		}
	})
}

func TestLessonDelete(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)
	lesson := newTestLesson(courseID, 1)

	tests := []struct {
		name        string
		requesterID uuid.UUID
		role        userEntity.Role
		wantStatus  int // 0 means allowed
	}{
		{name: "owner teacher", requesterID: teacherID, role: userEntity.RoleTeacher},
		{name: "admin", requesterID: uuid.New(), role: userEntity.RoleAdmin},
		{name: "another teacher", requesterID: uuid.New(), role: userEntity.RoleTeacher, wantStatus: 403},
		{name: "student", requesterID: uuid.New(), role: userEntity.RoleStudent, wantStatus: 403},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lessonRepo := &fakeLessonRepo{lesson: lesson}
			svc := NewLessonCommandServices(lessonRepo, &fakeCourseRepo{course: course})

			resErr := svc.Delete(context.Background(), lesson.GetID(), tt.requesterID, tt.role)

			if tt.wantStatus == 0 {
				if resErr != nil {
					t.Fatalf("expected success, got %+v", resErr)
				}
				if lessonRepo.deleteCalls != 1 {
					t.Error("the lesson must be deleted once")
				}
				return
			}

			if resErr == nil || resErr.StatusCode != tt.wantStatus {
				t.Fatalf("want status %d, got %+v", tt.wantStatus, resErr)
			}
			if lessonRepo.deleteCalls != 0 {
				t.Error("nothing must be deleted when access is denied")
			}
		})
	}

	t.Run("lesson does not exist", func(t *testing.T) {
		lessonRepo := &fakeLessonRepo{lesson: nil}
		svc := NewLessonCommandServices(lessonRepo, &fakeCourseRepo{course: course})

		resErr := svc.Delete(context.Background(), uuid.New(), teacherID, userEntity.RoleTeacher)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.LessonNotFoundErrorCode {
			t.Fatalf("want 404 lesson not found, got %+v", resErr)
		}
	})

	t.Run("course lookup error is passed through, not treated as success", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		lessonRepo := &fakeLessonRepo{lesson: lesson}
		svc := NewLessonCommandServices(lessonRepo, &fakeCourseRepo{err: dbErr})

		resErr := svc.Delete(context.Background(), lesson.GetID(), teacherID, userEntity.RoleTeacher)

		if resErr != dbErr {
			t.Fatalf("want the database error, got %+v", resErr)
		}
		if lessonRepo.deleteCalls != 0 {
			t.Error("nothing must be deleted when the course lookup failed")
		}
	})
}
