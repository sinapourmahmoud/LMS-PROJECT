package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	courseEntity "github.com/lms-project/backend/internal/course/domain/entity"
	courseRepository "github.com/lms-project/backend/internal/course/domain/repository"
	"github.com/lms-project/backend/internal/session/domain/entity"
	"github.com/lms-project/backend/internal/session/domain/repository"
	"github.com/lms-project/backend/internal/session/domain/values"
	userEntity "github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type SessionCommandServices interface {
	Create(ctx context.Context, creatorID uuid.UUID, courseID uuid.UUID, title string, startsAt time.Time) (*values.SessionResponse, *errorhandler.Response)
	Update(ctx context.Context, sessionID uuid.UUID, title string, startsAt time.Time, requesterID uuid.UUID) (*values.SessionResponse, *errorhandler.Response)
	Start(ctx context.Context, sessionID uuid.UUID, requesterID uuid.UUID) (*values.SessionResponse, *errorhandler.Response)
	End(ctx context.Context, sessionID uuid.UUID, requesterID uuid.UUID) (*values.SessionResponse, *errorhandler.Response)
	Delete(ctx context.Context, sessionID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) *errorhandler.Response
}

type sessionCommandServices struct {
	sessionRepository repository.SessionRepository
	courseRepository  courseRepository.CourseRepository
}

func NewSessionCommandServices(sessionRepository repository.SessionRepository, courseRepository courseRepository.CourseRepository) SessionCommandServices {
	return &sessionCommandServices{
		sessionRepository: sessionRepository,
		courseRepository:  courseRepository,
	}
}

func forbiddenError() *errorhandler.Response {
	return &errorhandler.Response{
		Code:       errorhandler.ForbiddenErrorCode,
		Message:    errorhandler.ForbiddenMessage,
		StatusCode: 403,
	}
}

func courseNotFoundError() *errorhandler.Response {
	return &errorhandler.Response{
		Code:       errorhandler.CourseNotFoundErrorCode,
		Message:    errorhandler.CourseNotFoundMessage,
		StatusCode: 404,
	}
}

func sessionNotFoundError() *errorhandler.Response {
	return &errorhandler.Response{
		Code:       errorhandler.SessionNotFoundErrorCode,
		Message:    errorhandler.SessionNotFoundMessage,
		StatusCode: 404,
	}
}

func invalidStateError() *errorhandler.Response {
	return &errorhandler.Response{
		Code:       errorhandler.SessionInvalidStateErrorCode,
		Message:    errorhandler.SessionInvalidStateMessage,
		StatusCode: 409,
	}
}

func startsAtInPastError() *errorhandler.Response {
	return &errorhandler.Response{
		Code:       errorhandler.ValidationErrorCode,
		Message:    "starts_at must be in the future",
		StatusCode: 400,
	}
}

// findSessionAndCourse loads a session and the course it belongs to.
func (r *sessionCommandServices) findSessionAndCourse(ctx context.Context, sessionID uuid.UUID) (*entity.Session, *courseEntity.Course, *errorhandler.Response) {

	session, err := r.sessionRepository.GetByID(ctx, sessionID)

	if err != nil {
		return nil, nil, err
	}

	if session == nil {
		return nil, nil, sessionNotFoundError()
	}

	course, err := r.courseRepository.GetByID(ctx, session.GetCourseID())

	if err != nil {
		return nil, nil, err
	}

	if course == nil {
		return nil, nil, courseNotFoundError()
	}

	return session, course, nil
}

func (r *sessionCommandServices) Create(ctx context.Context, creatorID uuid.UUID, courseID uuid.UUID, title string, startsAt time.Time) (*values.SessionResponse, *errorhandler.Response) {

	course, err := r.courseRepository.GetByID(ctx, courseID)

	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, courseNotFoundError()
	}

	if course.GetTeacherID() != creatorID {
		return nil, forbiddenError()
	}

	if !startsAt.After(time.Now()) {
		return nil, startsAtInPastError()
	}

	session := entity.NewSession(uuid.New(), courseID, title, startsAt, entity.StatusScheduled, nil, nil, time.Now(), nil)

	err = r.sessionRepository.Create(ctx, session)

	if err != nil {
		return nil, err
	}

	return values.NewSessionResponse(session), nil
}

func (r *sessionCommandServices) Update(ctx context.Context, sessionID uuid.UUID, title string, startsAt time.Time, requesterID uuid.UUID) (*values.SessionResponse, *errorhandler.Response) {

	session, course, err := r.findSessionAndCourse(ctx, sessionID)

	if err != nil {
		return nil, err
	}

	if course.GetTeacherID() != requesterID {
		return nil, forbiddenError()
	}

	if !startsAt.After(time.Now()) {
		return nil, startsAtInPastError()
	}

	if updateErr := session.Update(title, startsAt); updateErr != nil {
		return nil, mapStateError(updateErr)
	}

	err = r.sessionRepository.Update(ctx, session)

	if err != nil {
		return nil, err
	}

	return values.NewSessionResponse(session), nil
}

func (r *sessionCommandServices) Start(ctx context.Context, sessionID uuid.UUID, requesterID uuid.UUID) (*values.SessionResponse, *errorhandler.Response) {

	session, course, err := r.findSessionAndCourse(ctx, sessionID)

	if err != nil {
		return nil, err
	}

	if course.GetTeacherID() != requesterID {
		return nil, forbiddenError()
	}

	if startErr := session.Start(); startErr != nil {
		return nil, mapStateError(startErr)
	}

	err = r.sessionRepository.Update(ctx, session)

	if err != nil {
		return nil, err
	}

	return values.NewSessionResponse(session), nil
}

func (r *sessionCommandServices) End(ctx context.Context, sessionID uuid.UUID, requesterID uuid.UUID) (*values.SessionResponse, *errorhandler.Response) {

	session, course, err := r.findSessionAndCourse(ctx, sessionID)

	if err != nil {
		return nil, err
	}

	if course.GetTeacherID() != requesterID {
		return nil, forbiddenError()
	}

	if endErr := session.End(); endErr != nil {
		return nil, mapStateError(endErr)
	}

	err = r.sessionRepository.Update(ctx, session)

	if err != nil {
		return nil, err
	}

	return values.NewSessionResponse(session), nil
}

func (r *sessionCommandServices) Delete(ctx context.Context, sessionID uuid.UUID, requesterID uuid.UUID, requesterRole userEntity.Role) *errorhandler.Response {

	session, course, err := r.findSessionAndCourse(ctx, sessionID)

	if err != nil {
		return err
	}

	if requesterRole != userEntity.RoleAdmin && course.GetTeacherID() != requesterID {
		return forbiddenError()
	}

	// a class that is running right now must be ended first
	if session.GetStatus() == entity.StatusLive {
		return invalidStateError()
	}

	return r.sessionRepository.Delete(ctx, sessionID)
}

// mapStateError turns the entity's state error into the API error.
func mapStateError(err error) *errorhandler.Response {
	if errors.Is(err, entity.ErrInvalidState) {
		return invalidStateError()
	}

	return &errorhandler.Response{
		Code:       errorhandler.InternalSystemErrorCode,
		Message:    errorhandler.InternalSystemErrorMessage,
		StatusCode: 500,
	}
}
