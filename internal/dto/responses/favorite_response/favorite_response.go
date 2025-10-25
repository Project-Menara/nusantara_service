package favoriteresponse

import (
	"nusantara_service/internal/domain/entities"
	productresponse "nusantara_service/internal/dto/responses/product_response"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FavoriteResponse struct {
	ID            uuid.UUID              `json:"id"`
	User          string                 `json:"user"`
	FavoriteItems []FavoriteItemResponse `json:"favorite_item"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	DeletedAt     gorm.DeletedAt         `json:"deleted_at"`
}

type FavoriteItemResponse struct {
	ID       uuid.UUID                       `json:"id"`
	Selected bool                            `json:"selected"`
	Product  productresponse.ProductResponse `gorm:"foreignKey:ProductID;OnDelete:CASCADE"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func ToFavoriteItemResponse(favoriteItem entities.FavoriteItemEntity) FavoriteItemResponse {
	productItem := productresponse.ToProductResponse(&favoriteItem.Product)
	return FavoriteItemResponse{
		ID:        favoriteItem.ID,
		Selected:  favoriteItem.Selected,
		Product:   productItem,
		CreatedAt: favoriteItem.CreatedAt,
		UpdatedAt: favoriteItem.UpdatedAt,
		DeletedAt: favoriteItem.DeletedAt,
	}
}

func ToFavoriteResponse(favorite entities.FavoriteEntity) FavoriteResponse {
	favoriteItems := []FavoriteItemResponse{}
	for _, favoriteItem := range favorite.FavoriteItems {
		favoriteItemRes := ToFavoriteItemResponse(favoriteItem)
		favoriteItems = append(favoriteItems, favoriteItemRes)
	}

	return FavoriteResponse{
		ID:            favorite.ID,
		User:          favorite.User.Name,
		FavoriteItems: favoriteItems,
		CreatedAt:     favorite.CreatedAt,
		UpdatedAt:     favorite.UpdatedAt,
		DeletedAt:     favorite.DeletedAt,
	}
}
