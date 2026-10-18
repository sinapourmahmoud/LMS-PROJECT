package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/internal/course/domain/entity"
	"github.com/lms-project/backend/internal/course/domain/repository"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type courseRepository struct {
	db *sqlx.DB
}

func NewCourseRepository(db *sqlx.DB) repository.CourseRepository {
	return &courseRepository{
		db: db,
	}
}

func (r *courseRepository) Create(ctx context.Context, course *entity.Course) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	query := `INSERT INTO courses (id,title,description,teacher_id,status,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_, err := r.db.ExecContext(ctx, query, course.GetID(),
		course.GetTitle(),
		course.GetDescription(),
		course.GetTeacherID(),
		course.GetStatus(),
		course.GetCreatedAt(),
		course.GetUpdatedAt(),
	)

	if err != nil {
		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    errorhandler.DatabaseErrorMessage,
			StatusCode: 500,
		}
	}
	return nil

}
func (r *courseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	query := `SELECT id,title,description,teacher_id,status,created_at,updated_at FROM courses WHERE id=$1`

	var (
		courseID    uuid.UUID
		title       string
		status      entity.Status
		teacherID   uuid.UUID
		description string
		createdAt   time.Time
		updatedAt   sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(&courseID, &title, &description, &teacherID, &status, &createdAt, &updatedAt)

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

	return entity.NewCourse(courseID, title, description, teacherID, status, createdAt, updatedAtPtr), nil

}

func (r *courseRepository) ListPublished(ctx context.Context) ([]*entity.Course, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	query := `SELECT id,title,description,teacher_id,status,created_at,updated_at FROM courses WHERE status=$1`

	rows, err := r.db.QueryContext(ctx, query, entity.StatusPublished)

	if err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	defer rows.Close()

	var courses []*entity.Course

	for rows.Next() {

		var (
			id          uuid.UUID
			title       string
			description string
			teacherID   uuid.UUID
			status      entity.Status
			createdAt   time.Time
			updatedAt   sql.NullTime
		)
		if err := rows.Scan(&id, &title, &description, &teacherID, &status, &createdAt, &updatedAt); err != nil {
			return nil, &errorhandler.Response{
				Code:       errorhandler.DatabaseErrorCode,
				Message:    err.Error(),
				StatusCode: 500,
			}
		}
		var updatedAtPtr *time.Time
		if updatedAt.Valid {
			updatedAtPtr = &updatedAt.Time
		}

		courses = append(courses, entity.NewCourse(id, title, description, teacherID, status, createdAt, updatedAtPtr))

	}

	if err := rows.Err(); err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return courses, nil
}

func (r *courseRepository) ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]*entity.Course, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	query := `SELECT id,title,description,teacher_id,status,created_at,updated_at FROM courses WHERE teacher_id=$1`

	rows, err := r.db.QueryContext(ctx, query, teacherID)

	if err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	defer rows.Close()

	var courses []*entity.Course

	for rows.Next() {

		var (
			id           uuid.UUID
			title        string
			description  string
			teacherIDVar uuid.UUID
			status       entity.Status
			createdAt    time.Time
			updatedAt    sql.NullTime
		)
		if err := rows.Scan(&id, &title, &description, &teacherIDVar, &status, &createdAt, &updatedAt); err != nil {
			return nil, &errorhandler.Response{
				Code:       errorhandler.DatabaseErrorCode,
				Message:    err.Error(),
				StatusCode: 500,
			}
		}

		var updatedAtPtr *time.Time
		if updatedAt.Valid {
			updatedAtPtr = &updatedAt.Time
		}

		courses = append(courses, entity.NewCourse(id, title, description, teacherIDVar, status, createdAt, updatedAtPtr))

	}

	if err := rows.Err(); err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return courses, nil
}

func (r *courseRepository) ListPending(ctx context.Context) ([]*entity.Course, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	query := `SELECT id,title,description,teacher_id,status,created_at,updated_at FROM courses WHERE status=$1`

	rows, err := r.db.QueryContext(ctx, query, entity.StatusPending)

	if err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	defer rows.Close()

	var courses []*entity.Course

	for rows.Next() {

		var (
			id          uuid.UUID
			title       string
			description string
			teacherID   uuid.UUID
			status      entity.Status
			createdAt   time.Time
			updatedAt   sql.NullTime
		)
		if err := rows.Scan(&id, &title, &description, &teacherID, &status, &createdAt, &updatedAt); err != nil {
			return nil, &errorhandler.Response{
				Code:       errorhandler.DatabaseErrorCode,
				Message:    err.Error(),
				StatusCode: 500,
			}
		}
		var updatedAtPtr *time.Time
		if updatedAt.Valid {
			updatedAtPtr = &updatedAt.Time
		}

		courses = append(courses, entity.NewCourse(id, title, description, teacherID, status, createdAt, updatedAtPtr))

	}

	if err := rows.Err(); err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return courses, nil
}

func (r *courseRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.Status) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	query := `UPDATE courses SET status=$1 WHERE id=$2`

	_, err := r.db.ExecContext(ctx, query,status, id)

	if err != nil {
		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return nil
}
func (r *courseRepository) Delete(ctx context.Context, id uuid.UUID) *errorhandler.Response {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM courses WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,id)

	if err != nil {
		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	if rowsAffected == 0 {
		return &errorhandler.Response{Code: errorhandler.CourseNotFoundErrorCode, Message: errorhandler.CourseNotFoundMessage, StatusCode: 404}
	}

	return nil
}
