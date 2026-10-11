package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
)

type User struct {
	ID uuid.UUID

	username       string
	email          string
	hashedPassword string
	role           Role
	createdAt      time.Time
	updatedAt      *time.Time
}

func (u *User)GetID() uuid.UUID {
	return u.ID
}
func (u *User)GetUsername() string {
	return u.username
}
func (u *User)GetEmail() string {
	return u.email
}
func (u *User)GetHashedPassword() string {
	return u.hashedPassword
}

func (u *User)GetRole() Role {
	return u.role
}
func (u *User)GetCreatedAt() time.Time {
	return u.createdAt
}
func (u *User)GetUpdatedAt() *time.Time {
	return u.updatedAt
}

func NewUser(id uuid.UUID, username string, email string, hashedPassword string, role Role, createdAt time.Time, updatedAt *time.Time) *User {

	return &User{

		ID:             id,
		username:       username,
		email:          email,
		hashedPassword: hashedPassword,
		role:           role,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}

}
