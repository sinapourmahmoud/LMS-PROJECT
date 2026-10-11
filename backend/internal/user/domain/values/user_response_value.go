package values

import (
	"github.com/google/uuid"
	"github.com/lms-project/backend/internal/user/domain/entity"
)

type UserResponse struct {
	ID       uuid.UUID   `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     entity.Role `json:"role"`
}

func NewUserResponse(user *entity.User) *UserResponse {

	return &UserResponse{

		ID:       user.GetID(),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Role:     user.GetRole(),
	}

}
