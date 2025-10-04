package services

import (
	"context"
	dto "nusantara_service/internal/dto/request"

	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(ctx context.Context, superAdminId uuid.UUID, req dto.CreateEventRequest) error
}
