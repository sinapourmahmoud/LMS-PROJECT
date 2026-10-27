package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/internal/enrollment/domain/entity"
	"github.com/lms-project/backend/internal/enrollment/domain/repository"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type enrollmentRepository struct {
	db *sqlx.DB
}

func NewEnrollmentRepository(db *sqlx.DB) repository.EnrollmentRepository {

	return &enrollmentRepository{
		db: db,
	}

}

func (r *enrollmentRepository) Create(ctx context.Context, enrollment *entity.Enrollment) *errorhandler.Response {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	query := `INSERT INTO enrollments (id,student_id,course_id,enrolled_at) VALUES ($1,$2,$3,$4)`

	_, err := r.db.ExecContext(ctx, query, enrollment.GetID(), enrollment.GetStudentID(), enrollment.GetCourseID(), enrollment.GetEnrolledAt())

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &errorhandler.Response{
				Code:       errorhandler.AlreadyEnrolledErrorCode,
				Message:    errorhandler.AlreadyEnrolledErrorMessage,
				StatusCode: 409,
			}
		}
		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    errorhandler.DatabaseErrorMessage,
			StatusCode: 500,
		}
	}
	return nil

}
func (r *enrollmentRepository) FindByStudentAndCourse(ctx context.Context, studentID uuid.UUID, courseID uuid.UUID) (*entity.Enrollment, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()
	query := `SELECT  id,student_id,course_id,enrolled_at FROM enrollments WHERE student_id = $1 AND course_id=$2`
	var (
		id         uuid.UUID
		sID, cID   uuid.UUID
		enrolledAt time.Time
	)

	err := r.db.QueryRowContext(ctx, query, studentID, courseID).
		Scan(&id, &sID, &cID, &enrolledAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return entity.NewEnrollment(id, sID, cID, enrolledAt), nil

}

func (r *enrollmentRepository) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Enrollment, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	query := `SELECT id,student_id,course_id,enrolled_at FROM enrollments WHERE course_id=$1`

	rows, err := r.db.QueryContext(ctx, query, courseID)

	if err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	var enrollments []*entity.Enrollment
	defer rows.Close()

	for rows.Next() {

		var (
			id         uuid.UUID
			sID, cID   uuid.UUID
			enrolledAt time.Time
		)
		if err := rows.Scan(&id, &sID, &cID, &enrolledAt); err != nil {
			return nil, &errorhandler.Response{
				Code:       errorhandler.DatabaseErrorCode,
				Message:    err.Error(),
				StatusCode: 500,
			}
		}

		enrollments = append(enrollments, entity.NewEnrollment(id, sID, cID, enrolledAt))

	}

	if err := rows.Err(); err != nil {
		return nil, &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, Message: errorhandler.DatabaseErrorMessage, StatusCode: 500}
	}

	return enrollments, nil

}

func (r *enrollmentRepository) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.Enrollment, *errorhandler.Response) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `SELECT id,student_id,course_id,enrolled_at FROM enrollments WHERE student_id=$1`

	rows, err := r.db.QueryContext(ctx, query, studentID)

	if err != nil {
		return nil, &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	var enrollments []*entity.Enrollment
	defer rows.Close()

	for rows.Next() {

		var (
			id         uuid.UUID
			sID, cID   uuid.UUID
			enrolledAt time.Time
		)
		if err := rows.Scan(&id, &sID, &cID, &enrolledAt); err != nil {
			return nil, &errorhandler.Response{
				Code:       errorhandler.DatabaseErrorCode,
				Message:    err.Error(),
				StatusCode: 500,
			}
		}

		enrollments = append(enrollments, entity.NewEnrollment(id, sID, cID, enrolledAt))

	}

	if err := rows.Err(); err != nil {
		return nil, &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, Message: errorhandler.DatabaseErrorMessage, StatusCode: 500}
	}
	return enrollments, nil

}

func (r *enrollmentRepository) Delete(ctx context.Context, studentID uuid.UUID, courseID uuid.UUID) *errorhandler.Response {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `DELETE FROM enrollments WHERE student_id=$1 AND course_id=$2`

	_, err := r.db.DB.ExecContext(ctx, query, studentID, courseID)

	if err != nil {
		return &errorhandler.Response{
			Code:       errorhandler.DatabaseErrorCode,
			Message:    errorhandler.DatabaseErrorMessage,
			StatusCode: 500,
		}
	}

	return nil

}
