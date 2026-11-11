package dto

import "time"

type SessionRequest struct {
	Title    string    `json:"title" validate:"required"`
	StartsAt time.Time `json:"starts_at" validate:"required"`
}
