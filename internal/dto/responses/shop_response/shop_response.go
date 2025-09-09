package shopresponse

import (
	"nusantara_service/internal/domain/entities"
	cashierresponse "nusantara_service/internal/dto/responses/cashier_response"
	productresponse "nusantara_service/internal/dto/responses/product_response"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopResponse struct {
	ID           uuid.UUID                         `json:"id"`
	Name         string                            `json:"name"`
	Cover        string                            `json:"cover"`
	Description  string                            `json:"description"`
	FullAddress  string                            `json:"full_address"`
	Lat          float64                           `json:"lat"`
	Lng          float64                           `json:"lang"`
	Status       int                               `json:"status"`
	CreatedBy    string                            `json:"created_by"`
	CreatedAt    time.Time                         `json:"created_at"`
	UpdatedAt    time.Time                         `json:"update_at"`
	DeletedAt    gorm.DeletedAt                    `json:"deleted_at"`
	ShopImages   []string                          `json:"shop_images"`
	ShopProducts []productresponse.ProductResponse `json:"shop_product"`
	ShopCashier  []cashierresponse.CashierResponse `json:"shop_cashier"`
}

func ToShopResponse(shop entities.ShopEntity) ShopResponse {
	shopImages := []string{}
	for _, shopImage := range shop.ShopImages {
		shopImages = append(shopImages, shopImage.Image.ImagePath)
	}

	shopProducts := []productresponse.ProductResponse{}
	for _, shopProduct := range shop.ShopProducts {
		productRes := productresponse.ToProductResponse(&shopProduct.Product)
		shopProducts = append(shopProducts, productRes)
	}

	shopCashiers := []cashierresponse.CashierResponse{}
	for _, shopCashier := range shop.ShopCashiers {
		cashierRes := cashierresponse.ToCashierResponse(shopCashier.User)
		shopCashiers = append(shopCashiers, cashierRes)
	}

	return ShopResponse{
		ID:           shop.ID,
		Name:         shop.Name,
		Cover:        shop.Cover,
		Description:  shop.Description,
		FullAddress:  shop.FullAddress,
		Lat:          shop.Lat,
		Lng:          shop.Lng,
		Status:       shop.Status,
		CreatedBy:    shop.Creator.Name,
		CreatedAt:    shop.CreatedAt,
		UpdatedAt:    shop.UpdatedAt,
		DeletedAt:    shop.DeletedAt,
		ShopImages:   shopImages,
		ShopProducts: shopProducts,
		ShopCashier:  shopCashiers,
	}
}
