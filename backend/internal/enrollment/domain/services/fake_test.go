package services

import (
	"context"

	"github.com/google/uuid"
	courseEntity "github.com/lms-project/backend/internal/course/domain/entity"
	courseRepository "github.com/lms-project/backend/internal/course/domain/repository"
	"github.com/lms-project/backend/internal/enrollment/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

// fakeCourseRepo embeds the interface so we only write the method the tests use.
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
	existing    *entity.Enrollment
	findErr     *errorhandler.Response
	createErr   *errorhandler.Response
	createCalls int
	list        []*entity.Enrollment
	listErr     *errorhandler.Response
	listCalls   int
	deleteErr   *errorhandler.Response
	deleteCalls int
}

func (f *fakeEnrollmentRepo) Create(ctx context.Context, e *entity.Enrollment) *errorhandler.Response {
	f.createCalls++
	return f.createErr
}

func (f *fakeEnrollmentRepo) FindByStudentAndCourse(ctx context.Context, studentID, courseID uuid.UUID) (*entity.Enrollment, *errorhandler.Response) {
	return f.existing, f.findErr
}

func (f *fakeEnrollmentRepo) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Enrollment, *errorhandler.Response) {
	f.listCalls++
	return f.list, f.listErr
}

func (f *fakeEnrollmentRepo) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.Enrollment, *errorhandler.Response) {
	f.listCalls++
	return f.list, f.listErr
}

func (f *fakeEnrollmentRepo) Delete(ctx context.Context, studentID, courseID uuid.UUID) *errorhandler.Response {
	f.deleteCalls++
	return f.deleteErr
}
