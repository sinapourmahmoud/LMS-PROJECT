package services
 
import (
	"context"
	"testing"
 
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/course/domain/entity"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)
 
func TestGetCourseByID(t *testing.T) {
	course := newCourse(uuid.New(), entity.StatusPublished)
 
	t.Run("found", func(t *testing.T) {
		svc := NewCourseQueryServices(&fakeCourseRepo{course: course}, &fakeUserRepo{})
 
		res, resErr := svc.GetByID(context.Background(), course.GetID())
 
		if resErr != nil || res.ID != course.GetID() {
			t.Fatalf("got %+v, %+v; want the course", res, resErr)
		}
	})
 
	t.Run("not found", func(t *testing.T) {
		svc := NewCourseQueryServices(&fakeCourseRepo{}, &fakeUserRepo{})
 
		_, resErr := svc.GetByID(context.Background(), uuid.New())
 
		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.CourseNotFoundErrorCode {
			t.Errorf("want 404 course not found, got %+v", resErr)
		}
	})
 
	t.Run("repository error is passed through", func(t *testing.T) {
		svc := NewCourseQueryServices(&fakeCourseRepo{err: dbErr}, &fakeUserRepo{})
 
		_, resErr := svc.GetByID(context.Background(), uuid.New())
 
		if resErr == nil || resErr.StatusCode != 500 {
			t.Errorf("want the 500 error, got %+v", resErr)
		}
	})
}
 
// The three list methods behave the same way, so one table covers all of them.
func TestCourseLists(t *testing.T) {
	teacher := newUser(userEntity.RoleTeacher)
 
	lists := map[string]func(svc CourseQueryDomainServices) (int, *errorhandler.Response, bool){
		"ListPublished": func(svc CourseQueryDomainServices) (int, *errorhandler.Response, bool) {
			res, err := svc.ListPublished(context.Background())
			return len(res), err, res == nil
		},
		"ListByTeacher": func(svc CourseQueryDomainServices) (int, *errorhandler.Response, bool) {
			res, err := svc.ListByTeacher(context.Background(), teacher.GetID())
			return len(res), err, res == nil
		},
		"ListPending": func(svc CourseQueryDomainServices) (int, *errorhandler.Response, bool) {
			res, err := svc.ListPending(context.Background())
			return len(res), err, res == nil
		},
	}
 
	for name, call := range lists {
		t.Run(name+" returns the courses", func(t *testing.T) {
			repo := &fakeCourseRepo{list: []*entity.Course{newCourse(teacher.GetID(), entity.StatusPublished), newCourse(teacher.GetID(), entity.StatusPublished)}}
			svc := NewCourseQueryServices(repo, &fakeUserRepo{user: teacher})
 
			n, resErr, _ := call(svc)
 
			if resErr != nil || n != 2 {
				t.Errorf("got %d courses, %+v; want 2", n, resErr)
			}
		})
 
		t.Run(name+" gives an empty slice, not nil", func(t *testing.T) {
			svc := NewCourseQueryServices(&fakeCourseRepo{}, &fakeUserRepo{user: teacher})
 
			n, resErr, isNil := call(svc)
 
			if resErr != nil || n != 0 || isNil {
				t.Errorf("want empty non-nil slice (JSON [] not null), got len=%d nil=%v err=%+v", n, isNil, resErr)
			}
		})
 
		t.Run(name+" passes repository errors through", func(t *testing.T) {
			svc := NewCourseQueryServices(&fakeCourseRepo{err: dbErr}, &fakeUserRepo{user: teacher})
 
			_, resErr, _ := call(svc)
 
			if resErr == nil || resErr.StatusCode != 500 {
				t.Errorf("want the 500 error, got %+v", resErr)
			}
		})
	}
 
	t.Run("ListByTeacher: teacher does not exist", func(t *testing.T) {
		svc := NewCourseQueryServices(&fakeCourseRepo{}, &fakeUserRepo{user: nil})
 
		_, resErr := svc.ListByTeacher(context.Background(), uuid.New())
 
		if resErr == nil || resErr.StatusCode != 404 || resErr.Code != errorhandler.UserNotFoundErrorCode {
			t.Errorf("want 404 user not found, got %+v", resErr)
		}
	})
}
 