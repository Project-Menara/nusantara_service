package response

import (
	"nusantara_service/internal/domain/entities"
	"time"
)

type CheckPhoneResult struct {
	Action string               `json:"action"` // register, verify_otp, login
	User   *entities.UserEntity `json:"user,omitempty"`
	Ttl    *time.Duration       `json:"ttl,omitempty"`
}
