package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/data/model"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	dto "nusantara_service/internal/dto/request"
	"nusantara_service/internal/response"
	"nusantara_service/internal/workers/payload"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EventService struct {
	eventRepo   repositories.EventRepository
	productRepo repositories.ProductRepository
	rdb         *redis.Client
	cld         cloudinary.CloudinaryService
}

func NewEventUsecase(repo repositories.EventRepository, productRepo repositories.ProductRepository, rdb *redis.Client, cld cloudinary.CloudinaryService) services.EventService {
	return &EventService{
		eventRepo:   repo,
		productRepo: productRepo,
		rdb:         rdb,
		cld:         cld,
	}
}

func fileEventToBytes(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return io.ReadAll(f)
}

func (e *EventService) invalidateEventCache(ctx context.Context) {
	iter := e.rdb.Scan(ctx, 0, "events:*", 0).Iterator()
	for iter.Next(ctx) {
		e.rdb.Del(ctx, iter.Val())
	}

	iterID := e.rdb.Scan(ctx, 0, "event:*", 0).Iterator()
	for iterID.Next(ctx) {
		e.rdb.Del(ctx, iterID.Val())
	}
}

// CreateEvent implements services.EventService.
func (e *EventService) CreateEvent(ctx context.Context, superAdminId uuid.UUID, req dto.CreateEventRequest) error {
	superAdmin, err := e.eventRepo.FindUserIDSuperAdmin(ctx, superAdminId.String())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if superAdmin == nil {
		return response.NewCustomError(response.ErrForbidden, "only superadmin can create event", 403)
	}

	if strings.TrimSpace(req.Name) == "" {
		return response.NewCustomError(response.ErrBadRequest, "event name is required", 400)
	}
	existingEvent, err := e.eventRepo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingEvent != nil {
		return response.NewCustomError(response.ErrExists, "event name already exists", 409)
	}
	if strings.TrimSpace(req.Type) == "" {
		return response.NewCustomError(response.ErrBadRequest, "event type is required", 400)
	}
	if req.StartDate.IsZero() {
		return response.NewCustomError(response.ErrBadRequest, "start at is required", 400)
	}
	if req.EndDate.IsZero() {
		return response.NewCustomError(response.ErrBadRequest, "end at is required", 400)
	}
	if req.EndDate.Before(req.StartDate) {
		return response.NewCustomError(response.ErrBadRequest, "end at must be after start at", 400)
	}
	if req.Status != 0 && req.Status != 1 {
		return response.NewCustomError(response.ErrBadRequest, "status must be 0 or 1", 400)
	}
	if req.Cover == nil {
		return response.NewCustomError(response.ErrBadRequest, "cover is required", 400)
	}
	switch req.Type {
	case string(model.EventTypeDiskon):
		if req.EventProducts == nil {
			return response.NewCustomError(response.ErrBadRequest, "event products is required for diskon event", 400)
		}
	case string(model.EventTypeBundle):
		if req.EventBundleBuys == nil || req.EventBundleRewards == nil {
			return response.NewCustomError(response.ErrBadRequest, "event bundle buys and rewards are required for bundle event", 400)
		}
	}
	event := &entities.EventEntity{
		ID:        uuid.New(),
		Name:      req.Name,
		TypeEvent: entities.EventType(strings.ToUpper(req.Type)),
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    req.Status,
		CreatedBy: uuid.MustParse(superAdmin.ID),
	}
	if err := e.eventRepo.Create(ctx, event); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to create event", 500)
	}
	if req.Cover != nil {
		if binner, err := fileEventToBytes(req.Cover); err == nil && len(binner) > 0 {
			task := payload.ImageEventUploadPayload{
				EventID:   event.ID,
				Type:      "cover",
				FileBytes: binner,
				Folder:    "nusantara_service/events",
				Filename:  fmt.Sprintf("event_%s_cover", event.ID.String()),
			}
			_ = rabbitmq.PublishToQueue("", rabbitmq.SendImageEventQueueName, task)
		}
	}
	var eventProducts []entities.EventProductEntity
	var eventBundleBuys []entities.EventBundleBuyEntity
	var eventBundleRewards []entities.EventBundleRewardEntity
	switch req.Type {
	case string(model.EventTypeDiskon):
		for _, ep := range req.EventProducts {
			product, err := e.productRepo.GetByID(ctx, ep.ProductID)
			if err != nil {
				return response.NewCustomError(response.ErrNotFound, "product id "+ep.ProductID.String()+" not found", 404)
			}

			discount_amount := (ep.DiscountPercent * product.Price) / 100

			eventProducts = append(eventProducts, entities.EventProductEntity{
				ID:              uuid.New(),
				EventID:         event.ID,
				ProductID:       ep.ProductID,
				DiscountPercent: ep.DiscountPercent,
				DiscountAmount:  float64(discount_amount),
			})
		}
	case string(model.EventTypeBundle):
		for _, ebb := range req.EventBundleBuys {
			_, err := e.productRepo.GetByID(ctx, ebb.ProductID)
			if err != nil {
				return response.NewCustomError(response.ErrNotFound, "product id "+ebb.ProductID.String()+" not found", 404)
			}
			eventBundleBuys = append(eventBundleBuys, entities.EventBundleBuyEntity{
				ID:        uuid.New(),
				EventID:   event.ID,
				ProductID: ebb.ProductID,
				Quantity:  ebb.Quantity,
			})
		}

		for _, ebw := range req.EventBundleRewards {
			_, err := e.productRepo.GetByID(ctx, ebw.ProductID)
			if err != nil {
				return response.NewCustomError(response.ErrNotFound, "product id "+ebw.ProductID.String()+" not found", 404)
			}
			eventBundleRewards = append(eventBundleRewards, entities.EventBundleRewardEntity{
				ID:        uuid.New(),
				EventID:   event.ID,
				ProductID: ebw.ProductID,
				Quantity:  ebw.Quantity,
			})
		}
	}

	if len(eventProducts) > 0 {
		if err := e.eventRepo.AssignEventProduct(ctx, event.ID, eventProducts); err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to create event products", 500)
		}
	}
	if len(eventBundleBuys) > 0 {
		if err := e.eventRepo.AssignEventBundleBuy(ctx, event.ID, eventBundleBuys); err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to create event bundle buys", 500)
		}
	}
	if len(eventBundleRewards) > 0 {
		if err := e.eventRepo.AssignEventBundleReward(ctx, event.ID, eventBundleRewards); err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to create event bundle rewards", 500)
		}
	}
	e.invalidateEventCache(ctx)
	return nil
}

// GetAllEvents implements services.EventService.
func (e *EventService) GetAllEvents(ctx context.Context, page int, limit int, search string) ([]*entities.EventEntity, int, error) {
	cacheKey := fmt.Sprintf("shops:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*entities.EventEntity `json:"data"`
			Total int                     `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit
	events, total, err := e.eventRepo.FindAll(ctx, offset, limit, search)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, response.NewCustomError(response.ErrNotFound, "events not found", 404)
		}
		return nil, 0, response.NewCustomError(response.ErrInternal, "failed to fetch events", 500)
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  events,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return events, total, nil
}

// GetEventById implements services.EventService.
func (e *EventService) GetEventById(ctx context.Context, id uuid.UUID) (*entities.EventEntity, error) {
	cacheKey := fmt.Sprintf("event:%s", id)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var event entities.EventEntity
		if json.Unmarshal([]byte(cached), &event) == nil {
			return &event, nil
		}
	}

	event, err := e.eventRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "event not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to fetch event", 500)
	}

	buf, _ := json.Marshal(event)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return event, nil
}
