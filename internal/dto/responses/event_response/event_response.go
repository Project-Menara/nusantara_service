package eventresponse

import (
	"nusantara_service/internal/domain/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventResponse struct {
	ID                 uuid.UUID                   `json:"id"`
	Name               string                      `json:"name"`
	TypeEvent          entities.EventType          `json:"type_event"`
	StartDate          time.Time                   `json:"start_date"`
	EndDate            time.Time                   `json:"end_date"`
	Cover              string                      `json:"cover"`
	Status             int                         `json:"status"`
	EventProducts      []EventProductResponse      `json:"event_product"`
	EventBundleBuys    []EventBundleBuyResponse    `json:"event_bundle_buy"`
	EventBundleRewards []EventBundleRewardResponse `json:"event_bundle_reward"`
	CreatedBy          string                      `json:"created_by"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
	DeletedAt          gorm.DeletedAt              `json:"deleted_at"`
}

func ToEventResponse(event entities.EventEntity) EventResponse {
	eventProducts := []EventProductResponse{}
	for _, event_product := range event.EventProducts {
		event_product_res := ToEventProductResponse(event_product)
		eventProducts = append(eventProducts, event_product_res)
	}
	eventBundleBuys := []EventBundleBuyResponse{}
	for _, event_bundle_buy := range event.EventBundleBuys {
		event_bundle_buy_res := ToEvenBundleBuyResponse(&event_bundle_buy)
		eventBundleBuys = append(eventBundleBuys, event_bundle_buy_res)
	}
	eventBundleRewards := []EventBundleRewardResponse{}
	for _, event_bundle_reward := range event.EventBundleRewards {
		event_bundle_reward_res := ToEventBundleRewardResponse(&event_bundle_reward)
		eventBundleRewards = append(eventBundleRewards, event_bundle_reward_res)
	}
	switch event.TypeEvent {
	case entities.EventTypeDiskon:
		return EventResponse{
			ID:            event.ID,
			Name:          event.Name,
			TypeEvent:     event.TypeEvent,
			StartDate:     event.StartDate,
			EndDate:       event.EndDate,
			Cover:         event.Cover,
			Status:        event.Status,
			EventProducts: eventProducts,
			CreatedBy:     event.User.Name,
			CreatedAt:     event.CreatedAt,
			UpdatedAt:     event.UpdatedAt,
			DeletedAt:     event.DeletedAt,
		}
	case entities.EventTypeBundle:
		return EventResponse{
			ID:                 event.ID,
			Name:               event.Name,
			TypeEvent:          event.TypeEvent,
			StartDate:          event.StartDate,
			EndDate:            event.EndDate,
			Cover:              event.Cover,
			Status:             event.Status,
			EventBundleBuys:    eventBundleBuys,
			EventBundleRewards: eventBundleRewards,
			CreatedBy:          event.User.Name,
			CreatedAt:          event.CreatedAt,
			UpdatedAt:          event.UpdatedAt,
			DeletedAt:          event.DeletedAt,
		}
	default:
		return EventResponse{}
	}
}

type EventPublicResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Cover string    `json:"cover"`
}

func ToEventPublicResponse(event entities.EventEntity) EventPublicResponse {
	return EventPublicResponse{
		ID:    event.ID,
		Name:  event.Name,
		Cover: event.Cover,
	}
}
