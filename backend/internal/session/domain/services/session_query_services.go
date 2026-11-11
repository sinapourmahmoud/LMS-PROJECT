package services

import (
	"context"

	"github.com/google/uuid"
	courseEntity "github.com/lms-project/backend/internal/course/domain/entity"
	courseRepository "github.com/lms-project/backend/internal/course/domain/repository"
	enrollmentRepository "github.com/lms-project/backend/internal/enrollment/domain/repository"
	"github.com/lms-project/backend/internal/session/domain/repository"
	"github.com/lms-project/backend/internal/session/domain/values"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
	"github.com/lms-project/backend/pkg/pagination"
)

type SessionQueryServices interface {
	ListByCourse(ctx context.Context, courseID uuid.UUID, params pagination.Params, requesterID uuid.UUID, requesterRole userEntity.Role) (*pagination.Page[*values.SessionResponse], *errorhandler.Response)
	GetByID(ctx context.Context, sessionID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) (*values.SessionResponse, *errorhandler.Response)
}

type sessionQueryServices struct {
	sessionRepository    repository.SessionRepository
	courseRepository     courseRepository.CourseRepository
	enrollmentRepository enrollmentRepository.EnrollmentRepository
}

func NewSessionQueryServices(sessionRepository repository.SessionRepository, courseRepository courseRepository.CourseRepository, enrollmentRepository enrollmentRepository.EnrollmentRepository) SessionQueryServices {
	return &sessionQueryServices{
		sessionRepository:    sessionRepository,
		courseRepository:     courseRepository,
		enrollmentRepository: enrollmentRepository,
	}
}

// canView says whether the requester may see the sessions of a course:
// an admin, the teacher of the course, or a student enrolled in it.
func (r *sessionQueryServices) canView(ctx context.Context, course *courseEntity.Course, requesterID uuid.UUID, requesterRole userEntity.Role) (bool, *errorhandler.Response) {

	if requesterRole == userEntity.RoleAdmin || requesterID == course.GetTeacherID() {
		return true, nil
	}

	enrollment, err := r.enrollmentRepository.FindByStudentAndCourse(ctx, requesterID, course.GetID())

	if err != nil {
		return false, err
	}

	return enrollment != nil, nil
}

func (r *sessionQueryServices) ListByCourse(ctx context.Context, courseID uuid.UUID, params pagination.Params, requesterID uuid.UUID, requesterRole userEntity.Role) (*pagination.Page[*values.SessionResponse], *errorhandler.Response) {

	course, err := r.courseRepository.GetByID(ctx, courseID)

	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, courseNotFoundError()
	}

	allowed, err := r.canView(ctx, course, requesterID, requesterRole)

	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, forbiddenError()
	}

	sessionsEntity, total, err := r.sessionRepository.ListByCourse(ctx, courseID, params.Limit, params.Offset())

	if err != nil {
		return nil, err
	}

	sessions := make([]*values.SessionResponse, 0, len(sessionsEntity))

	for _, session := range sessionsEntity {
		sessions = append(sessions, values.NewSessionResponse(session))
	}

	return pagination.NewPage(sessions, params, total), nil
}

func (r *sessionQueryServices) GetByID(ctx context.Context, sessionID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) (*values.SessionResponse, *errorhandler.Response) {

	session, err := r.sessionRepository.GetByID(ctx, sessionID)

	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, sessionNotFoundError()
	}

	course, err := r.courseRepository.GetByID(ctx, session.GetCourseID())

	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, sessionNotFoundError()
	}

	allowed, err := r.canView(ctx, course, requesterID, requesterRole)

	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, forbiddenError()
	}

	return values.NewSessionResponse(session), nil
}
