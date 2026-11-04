package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/internal/lesson/domain/entity"
	"github.com/lms-project/backend/internal/lesson/domain/repository"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type lessonRepository struct {
	db *sqlx.DB
}

func NewLessonRepository(db *sqlx.DB) repository.LessonRepository {
	return &lessonRepository{
		db: db,
	}
}

func (r *lessonRepository) Create(ctx context.Context, lessonEntity entity.Lesson) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `INSERT INTO lessons(id,course_id,title,content,position,created_at) VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := r.db.ExecContext(ctx, query,
		lessonEntity.GetID(),
		lessonEntity.GetCourseID(),
		lessonEntity.GetTitle(),
		lessonEntity.GetContent(),
		lessonEntity.GetPosition(),
		lessonEntity.GetCreatedAt(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &errorhandler.Response{
				Code:       errorhandler.LessonPositionTakenErrorCode,
				Message:    errorhandler.LessonPositionTakenMessage,
				StatusCode: 409,
			}
		}
		return &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, Message: errorhandler.DatabaseErrorMessage, StatusCode: 500}
	}

	return nil

}
func (r *lessonRepository) GetByID(ctx context.Context, lessonID uuid.UUID) (*entity.Lesson, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()
	query := `SELECT id,course_id,title,content,position,created_at,updated_at FROM lessons WHERE id=$1`

	var (
		id        uuid.UUID
		courseID  uuid.UUID
		title     string
		content   string
		position  int
		createdAt time.Time
		updatedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(&id, &courseID, &title, &content, &position, &createdAt, &updatedAt)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    errorhandler.DatabaseErrorMessage,
			StatusCode: 500,
		}

	}

	var updatedAtPtr *time.Time
	if updatedAt.Valid {
		updatedAtPtr = &updatedAt.Time
	}

	return entity.NewLesson(id, courseID, title, content, position, createdAt, updatedAtPtr), nil
}

func (r *lessonRepository) Update(ctx context.Context, lesson *entity.Lesson) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `UPDATE lessons
	SET title = $1, content = $2, position = $3, updated_at = $4
	WHERE id = $5`

	result, err := r.db.ExecContext(ctx, query,
		lesson.GetTitle(),
		lesson.GetContent(),
		lesson.GetPosition(),
		lesson.GetUpdatedAt(),
		lesson.GetID(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &errorhandler.Response{
				Code:       errorhandler.LessonPositionTakenErrorCode,
				Message:    errorhandler.LessonPositionTakenMessage,
				StatusCode: 409,
			}
		}
		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    errorhandler.DatabaseErrorMessage,
			StatusCode: 500,
		}
	}

	rows, err := result.RowsAffected()
	if err == nil && rows == 0 {
		return &errorhandler.Response{
			Code:       errorhandler.LessonNotFoundErrorCode,
			Message:    errorhandler.LessonNotFoundMessage,
			StatusCode: 404,
		}
	}

	return nil
}

func (r *lessonRepository) Delete(ctx context.Context, id uuid.UUID) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `DELETE FROM lessons WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    errorhandler.DatabaseErrorMessage,
			StatusCode: 500,
		}
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    errorhandler.DatabaseErrorMessage,
			StatusCode: 500,
		}
	}

	if rows == 0 {
		return &errorhandler.Response{
			Code:       errorhandler.LessonNotFoundErrorCode,
			Message:    errorhandler.LessonNotFoundMessage,
			StatusCode: 404,
		}
	}

	return nil
}

func (r *lessonRepository) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Lesson, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	query := `SELECT (id,course_id,title,content,position,created_at,updated_at) FROM lessons WHERE course_id=$1 ORDER BY position`

	rows, err := r.db.QueryContext(ctx, query, courseID)

	if err != nil {

		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    errorhandler.DatabaseErrorMessage,
			StatusCode: 500,
		}
	}

	defer rows.Close()

	lessons := make([]*entity.Lesson, 0)

	for rows.Next() {
		var (
			id        uuid.UUID
			cID       uuid.UUID
			title     string
			content   string
			position  int
			createdAt time.Time
			updatedAt sql.NullTime
		)

		if err := rows.Scan(&id, &cID, &title, &content, &position, &createdAt, &updatedAt); err != nil {
			return nil, &errorhandler.Response{
				Code:       errorhandler.DatabaseErrorCode,
				Message:    errorhandler.DatabaseErrorMessage,
				StatusCode: 500,
			}
		}

		var updatedAtPtr *time.Time
		if updatedAt.Valid {
			updatedAtPtr = &updatedAt.Time
		}

		lessons = append(lessons, entity.NewLesson(id, cID, title, content, position, createdAt, updatedAtPtr))
	}

	if err := rows.Err(); err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    errorhandler.DatabaseErrorMessage,
			StatusCode: 500,
		}
	}

	return lessons, nil

}
