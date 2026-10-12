package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
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
		return &errorhandler.Response{
			Code:    errorhandler.DatabaseErrorCode,
			Message: err.Error(),
		}
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
		}
	}

	var updatedAtPtr *time.Time

	if updatedAt.Valid {
		updatedAtPtr = &updatedAt.Time
	}

	return entity.NewUser(id, userUsername, email, hashedPassword, role, createdAt, updatedAtPtr), nil
}
