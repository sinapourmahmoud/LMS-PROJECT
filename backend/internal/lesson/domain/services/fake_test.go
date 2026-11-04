package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	courseEntity "github.com/lms-project/backend/internal/course/domain/entity"
	courseRepository "github.com/lms-project/backend/internal/course/domain/repository"
	enrollmentEntity "github.com/lms-project/backend/internal/enrollment/domain/entity"
	enrollmentRepository "github.com/lms-project/backend/internal/enrollment/domain/repository"
	"github.com/lms-project/backend/internal/lesson/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

// The fakes embed the interface so we only write the methods the tests use.
// Calling any other method would panic (nil interface), which is fine for a test.

type fakeCourseRepo struct {
	courseRepository.CourseRepository
	course *courseEntity.Course
	err    *errorhandler.Response
}

func (f *fakeCourseRepo) GetByID(ctx context.Context, id uuid.UUID) (*courseEntity.Course, *errorhandler.Response) {
	return f.course, f.err
}

type fakeEnrollmentRepo struct {
	enrollmentRepository.EnrollmentRepository
	existing *enrollmentEntity.Enrollment
	err      *errorhandler.Response
}

func (f *fakeEnrollmentRepo) FindByStudentAndCourse(ctx context.Context, studentID, courseID uuid.UUID) (*enrollmentEntity.Enrollment, *errorhandler.Response) {
	return f.existing, f.err
}

type fakeLessonRepo struct {
	lesson      *entity.Lesson
	getErr      *errorhandler.Response
	list        []*entity.Lesson
	listErr     *errorhandler.Response
	listCalls   int
	createErr   *errorhandler.Response
	createCalls int
	created     *entity.Lesson
	updateErr   *errorhandler.Response
	updateCalls int
	updated     *entity.Lesson
	deleteErr   *errorhandler.Response
	deleteCalls int
}

func (f *fakeLessonRepo) Create(ctx context.Context, l *entity.Lesson) *errorhandler.Response {
	f.createCalls++
	f.created = l
	return f.createErr
}

func (f *fakeLessonRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, *errorhandler.Response) {
	return f.lesson, f.getErr
}

func (f *fakeLessonRepo) Update(ctx context.Context, l *entity.Lesson) *errorhandler.Response {
	f.updateCalls++
	f.updated = l
	return f.updateErr
}

func (f *fakeLessonRepo) Delete(ctx context.Context, id uuid.UUID) *errorhandler.Response {
	f.deleteCalls++
	return f.deleteErr
}

func (f *fakeLessonRepo) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, *errorhandler.Response) {
	f.listCalls++
	return f.list, f.listErr
}

// newTestCourse returns a published course owned by teacherID.
func newTestCourse(courseID, teacherID uuid.UUID) *courseEntity.Course {
	return courseEntity.NewCourse(courseID, "title", "desc", teacherID, courseEntity.StatusPublished, time.Now(), nil)
}

// newTestLesson returns a lesson that belongs to courseID.
func newTestLesson(courseID uuid.UUID, position int) *entity.Lesson {
	return entity.NewLesson(uuid.New(), courseID, "title", "content", position, time.Now(), nil)
}
