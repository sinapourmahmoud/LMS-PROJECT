package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	enrollmentEntity "github.com/lms-project/backend/internal/enrollment/domain/entity"
	"github.com/lms-project/backend/internal/lesson/domain/entity"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

func TestLessonListByCourse(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)
	lessons := []*entity.Lesson{newTestLesson(courseID, 1), newTestLesson(courseID, 2)}
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
			lessonRepo := &fakeLessonRepo{list: lessons}
			svc := NewLessonQueryServices(lessonRepo, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{existing: tt.enrollment})

			res, resErr := svc.ListByCourse(context.Background(), courseID, tt.requesterID, tt.role)

			if tt.wantStatus == 0 {
				if resErr != nil {
					t.Fatalf("expected success, got %+v", resErr)
				}
				if len(res) != 2 {
					t.Errorf("got %d lessons, want 2", len(res))
				}
				return
			}

			if resErr == nil || resErr.StatusCode != tt.wantStatus || resErr.Code != errorhandler.ForbiddenErrorCode {
				t.Fatalf("want 403 forbidden, got %+v", resErr)
			}
			if lessonRepo.listCalls != 0 {
				t.Error("the lessons must not be queried when access is denied")
			}
		})
	}

	t.Run("course does not exist", func(t *testing.T) {
		svc := NewLessonQueryServices(&fakeLessonRepo{}, &fakeCourseRepo{course: nil}, &fakeEnrollmentRepo{})

		_, resErr := svc.ListByCourse(context.Background(), courseID, teacherID, userEntity.RoleTeacher)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.CourseNotFoundErrorCode {
			t.Fatalf("want 404 course not found, got %+v", resErr)
		}
	})

	t.Run("a course without lessons gives an empty slice, not nil", func(t *testing.T) {
		svc := NewLessonQueryServices(&fakeLessonRepo{}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		res, resErr := svc.ListByCourse(context.Background(), courseID, teacherID, userEntity.RoleTeacher)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res == nil || len(res) != 0 {
			t.Errorf("want empty non-nil slice (JSON [] not null), got %#v", res)
		}
	})

	t.Run("enrollment lookup error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewLessonQueryServices(&fakeLessonRepo{list: lessons}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{err: dbErr})

		_, resErr := svc.ListByCourse(context.Background(), courseID, studentID, userEntity.RoleStudent)

		if resErr != dbErr {
			t.Errorf("want the database error, got %+v", resErr)
		}
	})

	t.Run("lessons repository error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewLessonQueryServices(&fakeLessonRepo{listErr: dbErr}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		_, resErr := svc.ListByCourse(context.Background(), courseID, teacherID, userEntity.RoleTeacher)

		if resErr != dbErr {
			t.Errorf("want the database error, got %+v", resErr)
		}
	})
}

func TestLessonGetByID(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	course := newTestCourse(courseID, teacherID)
	lesson := newTestLesson(courseID, 1)
	studentID := uuid.New()
	enrollment := enrollmentEntity.NewEnrollment(uuid.New(), studentID, courseID, time.Now())

	t.Run("enrolled student can read the lesson", func(t *testing.T) {
		svc := NewLessonQueryServices(&fakeLessonRepo{lesson: lesson}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{existing: enrollment})

		res, resErr := svc.GetByID(context.Background(), lesson.GetID(), studentID, userEntity.RoleStudent)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
		if res.ID != lesson.GetID() || res.Content != "content" {
			t.Errorf("unexpected response: %+v", res)
		}
	})

	t.Run("owner teacher can read the lesson", func(t *testing.T) {
		svc := NewLessonQueryServices(&fakeLessonRepo{lesson: lesson}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		_, resErr := svc.GetByID(context.Background(), lesson.GetID(), teacherID, userEntity.RoleTeacher)

		if resErr != nil {
			t.Fatalf("unexpected error: %+v", resErr)
		}
	})

	t.Run("student who is not enrolled is forbidden", func(t *testing.T) {
		svc := NewLessonQueryServices(&fakeLessonRepo{lesson: lesson}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		res, resErr := svc.GetByID(context.Background(), lesson.GetID(), uuid.New(), userEntity.RoleStudent)

		if res != nil || resErr == nil || resErr.StatusCode != 403 || resErr.Code != errorhandler.ForbiddenErrorCode {
			t.Fatalf("want 403 forbidden, got res=%v err=%+v", res, resErr)
		}
	})

	t.Run("lesson does not exist", func(t *testing.T) {
		svc := NewLessonQueryServices(&fakeLessonRepo{lesson: nil}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		_, resErr := svc.GetByID(context.Background(), uuid.New(), teacherID, userEntity.RoleTeacher)

		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.LessonNotFoundErrorCode {
			t.Fatalf("want 404 lesson not found, got %+v", resErr)
		}
	})

	t.Run("lesson repository error is passed through", func(t *testing.T) {
		dbErr := &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, StatusCode: 500}
		svc := NewLessonQueryServices(&fakeLessonRepo{getErr: dbErr}, &fakeCourseRepo{course: course}, &fakeEnrollmentRepo{})

		_, resErr := svc.GetByID(context.Background(), lesson.GetID(), teacherID, userEntity.RoleTeacher)

		if resErr != dbErr {
			t.Errorf("want the database error, got %+v", resErr)
		}
	})
}
