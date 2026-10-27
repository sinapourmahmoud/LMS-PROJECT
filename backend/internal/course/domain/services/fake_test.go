package services
 
import (
	"context"
 
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/course/domain/entity"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	userRepository "github.com/lms-project/backend/internal/user/domain/repository"
	"github.com/lms-project/backend/pkg/errorhandler"
)
 
// fakeCourseRepo implements the whole CourseRepository interface.
// Each field decides what a method returns; the *Calls fields record what happened.
type fakeCourseRepo struct {
	course *entity.Course
	list   []*entity.Course
	err    *errorhandler.Response // returned by every method
 
	createCalls   int
	deleteCalls   int
	updateCalls   int
	updatedStatus entity.Status
}
 
func (f *fakeCourseRepo) Create(ctx context.Context, course *entity.Course) *errorhandler.Response {
	f.createCalls++
	return f.err
}
 
func (f *fakeCourseRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, *errorhandler.Response) {
	return f.course, f.err
}
 
func (f *fakeCourseRepo) ListPublished(ctx context.Context) ([]*entity.Course, *errorhandler.Response) {
	return f.list, f.err
}
 
func (f *fakeCourseRepo) ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]*entity.Course, *errorhandler.Response) {
	return f.list, f.err
}
 
func (f *fakeCourseRepo) ListPending(ctx context.Context) ([]*entity.Course, *errorhandler.Response) {
	return f.list, f.err
}
 
func (f *fakeCourseRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.Status) *errorhandler.Response {
	f.updateCalls++
	f.updatedStatus = status
	return f.err
}
 
func (f *fakeCourseRepo) Delete(ctx context.Context, id uuid.UUID) *errorhandler.Response {
	f.deleteCalls++
	return f.err
}
 
// fakeUserRepo embeds the interface: only GetByID is needed by the course services.
type fakeUserRepo struct {
	userRepository.UserRepository
	user *userEntity.User
	err  *errorhandler.Response
}
 
func (f *fakeUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*userEntity.User, *errorhandler.Response) {
	return f.user, f.err
}
 