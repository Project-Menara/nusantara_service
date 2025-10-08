package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"

	"github.com/google/uuid"
)

type EventRepository interface {
	Create(ctx context.Context, data *entities.EventEntity) error
	FindUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error)
	FindByName(ctx context.Context, name string) (*entities.EventEntity, error)
	AssignEventProduct(ctx context.Context, eventID uuid.UUID, items []entities.EventProductEntity) error
	AssignEventBundleBuy(ctx context.Context, eventID uuid.UUID, items []entities.EventBundleBuyEntity) error
	AssignEventBundleReward(ctx context.Context, eventID uuid.UUID, items []entities.EventBundleRewardEntity) error
	CreateImage(ctx context.Context, img *entities.ImageEntity) error
	UpdateEventCover(ctx context.Context, eventID uuid.UUID, coverURL string) error
	FindById(ctx context.Context, id uuid.UUID) (*entities.EventEntity, error)
	FindAll(ctx context.Context, offset, limit int, search string) ([]*entities.EventEntity, int, error)
	Update(ctx context.Context, id uuid.UUID, data *entities.EventEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status int) error
	DeleteEventProductsByEventID(ctx context.Context, eventID uuid.UUID) error
	DeleteEventBundleBuysByEventID(ctx context.Context, eventID uuid.UUID) error
	DeleteEventBundleRewardsByEventID(ctx context.Context, eventID uuid.UUID) error
}
