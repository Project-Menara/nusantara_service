package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepositoryImpl struct {
	db *gorm.DB
}

func NewEventRepositoryImpl(db *gorm.DB) repositories.EventRepository {
	return &EventRepositoryImpl{db: db}
}

func (e *EventRepositoryImpl) preloadRelations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("EventProducts.Product").
		Preload("EventProducts.Product.Image").
		Preload("EventProducts.Product.ProductImages.Image").
		Preload("EventProducts.Product.TypeProduct").
		Preload("EventProducts.Product.User").
		Preload("EventProducts.Product.User.Role").
		Preload("EventBundleBuys.Product").
		Preload("EventBundleBuys.Product.Image").
		Preload("EventBundleBuys.Product.ProductImages.Image").
		Preload("EventBundleBuys.Product.TypeProduct").
		Preload("EventBundleBuys.Product.User").
		Preload("EventBundleBuys.Product.User.Role").
		Preload("EventBundleRewards.Product").
		Preload("EventBundleRewards.Product.Image").
		Preload("EventBundleRewards.Product.ProductImages.Image").
		Preload("EventBundleRewards.Product.TypeProduct").
		Preload("EventBundleRewards.Product.User").
		Preload("EventBundleRewards.Product.User.Role").
		Preload("User").
		Preload("User.Role")
}

// Create implements repositories.EventRepository.
func (e *EventRepositoryImpl) Create(ctx context.Context, data *entities.EventEntity) error {
	return e.preloadRelations(e.db.WithContext(ctx)).Create(data).Error
}

// FindByName implements repositories.EventRepository.
func (e *EventRepositoryImpl) FindByName(ctx context.Context, name string) (*entities.EventEntity, error) {
	var event entities.EventEntity
	if err := e.preloadRelations(e.db.WithContext(ctx)).First(&event, "name = ? AND deleted_at IS NULL", name).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// FindUserIDSuperAdmin implements repositories.EventRepository.
func (e *EventRepositoryImpl) FindUserIDSuperAdmin(ctx context.Context, userID string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := e.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", userID, "superadmin").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// AssignEventBundleBuy implements repositories.EventRepository.
func (e *EventRepositoryImpl) AssignEventBundleBuy(ctx context.Context, eventID uuid.UUID, items []entities.EventBundleBuyEntity) error {
	if len(items) == 0 {
		return nil
	}

	for i := range items {
		items[i].ID = uuid.New()
		items[i].EventID = eventID
	}

	return e.db.WithContext(ctx).Create(&items).Error
}

// AssignEventBundleReward implements repositories.EventRepository.
func (e *EventRepositoryImpl) AssignEventBundleReward(ctx context.Context, eventID uuid.UUID, items []entities.EventBundleRewardEntity) error {
	if len(items) == 0 {
		return nil
	}

	for i := range items {
		items[i].ID = uuid.New()
		items[i].EventID = eventID
	}

	return e.db.WithContext(ctx).Create(&items).Error
}

// AssignEventProduct implements repositories.EventRepository.
func (e *EventRepositoryImpl) AssignEventProduct(ctx context.Context, eventID uuid.UUID, items []entities.EventProductEntity) error {
	if len(items) == 0 {
		return nil
	}

	for i := range items {
		items[i].ID = uuid.New()
		items[i].EventID = eventID
	}

	return e.db.WithContext(ctx).Create(&items).Error
}

// CreateImage implements repositories.EventRepository.
func (e *EventRepositoryImpl) CreateImage(ctx context.Context, img *entities.ImageEntity) error {
	return e.db.WithContext(ctx).Create(img).Error
}

// UpdateEventCover implements repositories.EventRepository.
func (e *EventRepositoryImpl) UpdateEventCover(ctx context.Context, eventID uuid.UUID, coverURL string) error {
	return e.db.WithContext(ctx).Model(&entities.EventEntity{}).
		Where("id = ?", eventID).
		Update("cover", coverURL).Error
}

// FindById implements repositories.EventRepository.
func (e *EventRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*entities.EventEntity, error) {
	var event entities.EventEntity
	if err := e.preloadRelations(e.db.WithContext(ctx)).First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}
