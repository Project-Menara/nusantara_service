package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	dto "nusantara_service/internal/dto/request"

	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(ctx context.Context, superAdminId uuid.UUID, req dto.CreateEventRequest) error
	GetAllEvents(ctx context.Context, page, limit int, search string) ([]*entities.EventEntity, int, error)
	GetEventById(ctx context.Context, id uuid.UUID) (*entities.EventEntity, error)
}
