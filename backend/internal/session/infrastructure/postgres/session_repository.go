package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/internal/session/domain/entity"
	"github.com/lms-project/backend/internal/session/domain/repository"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) repository.SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func databaseError() *errorhandler.Response {
	return &errorhandler.Response{
		Code:       errorhandler.DatabaseErrorCode,
		Message:    errorhandler.DatabaseErrorMessage,
		StatusCode: 500,
	}
}

func sessionNotFoundError() *errorhandler.Response {
	return &errorhandler.Response{
		Code:       errorhandler.SessionNotFoundErrorCode,
		Message:    errorhandler.SessionNotFoundMessage,
		StatusCode: 404,
	}
}

func (r *sessionRepository) Create(ctx context.Context, session *entity.Session) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `INSERT INTO sessions(id,course_id,title,starts_at,status,created_at) VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := r.db.ExecContext(ctx, query,
		session.GetID(),
		session.GetCourseID(),
		session.GetTitle(),
		session.GetStartsAt(),
		session.GetStatus(),
		session.GetCreatedAt(),
	)

	if err != nil {
		return databaseError()
	}

	return nil
}

func (r *sessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*entity.Session, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `SELECT id,course_id,title,starts_at,status,started_at,ended_at,created_at,updated_at FROM sessions WHERE id=$1`

	var (
		id        uuid.UUID
		courseID  uuid.UUID
		title     string
		startsAt  time.Time
		status    entity.Status
		startedAt sql.NullTime
		endedAt   sql.NullTime
		createdAt time.Time
		updatedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(&id, &courseID, &title, &startsAt, &status, &startedAt, &endedAt, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, databaseError()
	}

	return entity.NewSession(id, courseID, title, startsAt, status, nullTimeToPtr(startedAt), nullTimeToPtr(endedAt), createdAt, nullTimeToPtr(updatedAt)), nil
}

func (r *sessionRepository) Update(ctx context.Context, session *entity.Session) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `UPDATE sessions
	SET title = $1, starts_at = $2, status = $3, started_at = $4, ended_at = $5, updated_at = $6
	WHERE id = $7`

	result, err := r.db.ExecContext(ctx, query,
		session.GetTitle(),
		session.GetStartsAt(),
		session.GetStatus(),
		session.GetStartedAt(),
		session.GetEndedAt(),
		session.GetUpdatedAt(),
		session.GetID(),
	)

	if err != nil {
		return databaseError()
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return databaseError()
	}

	if rows == 0 {
		return sessionNotFoundError()
	}

	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `DELETE FROM sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, sessionID)

	if err != nil {
		return databaseError()
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return databaseError()
	}

	if rows == 0 {
		return sessionNotFoundError()
	}

	return nil
}

func (r *sessionRepository) ListByCourse(ctx context.Context, courseID uuid.UUID, limit int, offset int) ([]*entity.Session, int, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE course_id=$1`, courseID).Scan(&total); err != nil {
		return nil, 0, databaseError()
	}

	query := `SELECT id,course_id,title,starts_at,status,started_at,ended_at,created_at,updated_at
	FROM sessions WHERE course_id=$1
	ORDER BY starts_at, id
	LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, courseID, limit, offset)

	if err != nil {
		return nil, 0, databaseError()
	}

	defer rows.Close()

	sessions := make([]*entity.Session, 0, limit)

	for rows.Next() {
		var (
			id        uuid.UUID
			cID       uuid.UUID
			title     string
			startsAt  time.Time
			status    entity.Status
			startedAt sql.NullTime
			endedAt   sql.NullTime
			createdAt time.Time
			updatedAt sql.NullTime
		)

		if err := rows.Scan(&id, &cID, &title, &startsAt, &status, &startedAt, &endedAt, &createdAt, &updatedAt); err != nil {
			return nil, 0, databaseError()
		}

		sessions = append(sessions, entity.NewSession(id, cID, title, startsAt, status, nullTimeToPtr(startedAt), nullTimeToPtr(endedAt), createdAt, nullTimeToPtr(updatedAt)))
	}

	if err := rows.Err(); err != nil {
		return nil, 0, databaseError()
	}

	return sessions, total, nil
}

func nullTimeToPtr(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
