package values

import (
	"time"

	"github.com/google/uuid"
)

type RosterEntry struct {
	StudentID  uuid.UUID `json:"student_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	EnrolledAt time.Time `json:"enrolled_at"`
}
