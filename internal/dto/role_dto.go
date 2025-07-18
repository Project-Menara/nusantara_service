package dto

import "github.com/google/uuid"

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type UpdateRoleRequest struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
