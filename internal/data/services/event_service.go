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
	UpdateEvent(ctx context.Context, superAdminId uuid.UUID, id uuid.UUID, req dto.UpdateEventRequest) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	UpdateStatusEvent(ctx context.Context, id uuid.UUID, req dto.UpdateStatusEventRequest) error

	GetAllEventPublic(ctx context.Context) ([]*entities.EventEntity, error)
}
