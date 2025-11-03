package cashierresponse

import (
	"nusantara_service/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type CashierResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Photo     string    `json:"photo"`
	Status    int       `json:"status"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func ToCashierResponse(user entities.UserEntity) CashierResponse {
	var photo string
	if user.Photo != nil {
		photo = *user.Photo
	}
	return CashierResponse{
		ID:        uuid.MustParse(user.ID),
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Photo:     photo,
		Status:    user.Status,
		Role:      user.Role.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt.Time,
	}
}
