package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // registers "pgx" driver with database/sql
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/internal/user/domain/repository"
	"github.com/lms-project/backend/pkg/errorhandler"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {

	return &userRepository{
		db: db,
	}

}

func (r *userRepository) Create(ctx context.Context, user *entity.User) *errorhandler.Response {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO users (id,username,email,hashed_password,role,created_at,updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.GetID(),
		user.GetUsername(),
		user.GetEmail(),
		user.GetHashedPassword(),
		user.GetRole(),
		user.GetCreatedAt(),
		user.GetUpdatedAt(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return &errorhandler.Response{Code: errorhandler.EmailAlreadyExistsErrorCode, Message: errorhandler.EmailAlreadyExistsMessage, StatusCode: 409}
			case "users_username_key":
				return &errorhandler.Response{Code: errorhandler.UsernameAlreadyExistsErrorCode, Message: errorhandler.UsernameAlreadyExistsMessage, StatusCode: 409}
			}
		}
		// anything that wasn't one of the specific cases above is a real failure
		return &errorhandler.Response{Code: errorhandler.DatabaseErrorCode, Message: errorhandler.DatabaseErrorMessage, StatusCode: 500}
	}

	return nil

}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, *errorhandler.Response) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT id, username, email, hashed_password, role, created_at, updated_at FROM users WHERE email = $1
	`

	var (
		id             uuid.UUID
		username       string
		userEmail      string
		hashedPassword string
		role           entity.Role
		createdAt      time.Time
		updatedAt      sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&id, &username, &userEmail, &hashedPassword, &role, &createdAt, &updatedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, &errorhandler.Response{
			Code:    errorhandler.DatabaseErrorCode,
			Message: err.Error(),
			StatusCode: 500,
		}

	}

	var updatedAtPtr *time.Time

	if updatedAt.Valid {
		updatedAtPtr = &updatedAt.Time
	}

	return entity.NewUser(id, username, userEmail, hashedPassword, role, createdAt, updatedAtPtr), nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, *errorhandler.Response) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT id, username, email, hashed_password, role, created_at, updated_at FROM users WHERE username = $1
	`
	var (
		id             uuid.UUID
		userUsername   string
		email          string
		hashedPassword string
		role           entity.Role
		createdAt      time.Time
		updatedAt      sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, username).Scan(&id, &userUsername, &email, &hashedPassword, &role, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, &errorhandler.Response{
			Code:    errorhandler.DatabaseErrorCode,
			Message: err.Error(),
			StatusCode: 500,
		}
	}

	var updatedAtPtr *time.Time

	if updatedAt.Valid {
		updatedAtPtr = &updatedAt.Time
	}

	return entity.NewUser(id, userUsername, email, hashedPassword, role, createdAt, updatedAtPtr), nil
}
