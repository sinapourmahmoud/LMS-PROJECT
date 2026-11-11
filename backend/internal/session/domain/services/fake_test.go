package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	courseEntity "github.com/lms-project/backend/internal/course/domain/entity"
	courseRepository "github.com/lms-project/backend/internal/course/domain/repository"
	enrollmentEntity "github.com/lms-project/backend/internal/enrollment/domain/entity"
	enrollmentRepository "github.com/lms-project/backend/internal/enrollment/domain/repository"
	"github.com/lms-project/backend/internal/session/domain/entity"
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

type fakeSessionRepo struct {
	session     *entity.Session
	getErr      *errorhandler.Response
	list        []*entity.Session
	total       int
	listErr     *errorhandler.Response
	listCalls   int
	gotLimit    int
	gotOffset   int
	createErr   *errorhandler.Response
	createCalls int
	created     *entity.Session
	updateErr   *errorhandler.Response
	updateCalls int
	updated     *entity.Session
	deleteErr   *errorhandler.Response
	deleteCalls int
}

func (f *fakeSessionRepo) Create(ctx context.Context, s *entity.Session) *errorhandler.Response {
	f.createCalls++
	f.created = s
	return f.createErr
}

func (f *fakeSessionRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Session, *errorhandler.Response) {
	return f.session, f.getErr
}

func (f *fakeSessionRepo) Update(ctx context.Context, s *entity.Session) *errorhandler.Response {
	f.updateCalls++
	f.updated = s
	return f.updateErr
}

func (f *fakeSessionRepo) Delete(ctx context.Context, id uuid.UUID) *errorhandler.Response {
	f.deleteCalls++
	return f.deleteErr
}

func (f *fakeSessionRepo) ListByCourse(ctx context.Context, courseID uuid.UUID, limit int, offset int) ([]*entity.Session, int, *errorhandler.Response) {
	f.listCalls++
	f.gotLimit = limit
	f.gotOffset = offset
	return f.list, f.total, f.listErr
}

// newTestCourse returns a published course owned by teacherID.
func newTestCourse(courseID, teacherID uuid.UUID) *courseEntity.Course {
	return courseEntity.NewCourse(courseID, "title", "desc", teacherID, courseEntity.StatusPublished, time.Now(), nil)
}

// newTestSession returns a session of courseID with the given status.
func newTestSession(courseID uuid.UUID, status entity.Status) *entity.Session {
	return entity.NewSession(uuid.New(), courseID, "Live Q&A", time.Now().Add(time.Hour), status, nil, nil, time.Now(), nil)
}
