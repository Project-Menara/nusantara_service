package shopresponse

import (
	"fmt"
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

type ShopNearbyResponse struct {
	ShopResponse
	Distance string `json:"distance"`
}

func ToShopNearbyResponse(shop entities.ShopEntity, distance float64) ShopNearbyResponse {
	return ShopNearbyResponse{
		ShopResponse: ToShopResponse(shop),
		Distance:     fmt.Sprintf("%.2f Km", distance),
	}
}

type ShopNameCashierResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ShopCashierResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Cover       string    `json:"cover"`
	Description string    `json:"description"`
	FullAddress string    `json:"full_address"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lang"`
	Status      int       `json:"status"`
	ShopImages  []string  `json:"shop_images"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"update_at"`
}

func ToShopCashierResponse(shop entities.ShopEntity) ShopCashierResponse {
	shopImages := []string{}
	for _, shopImage := range shop.ShopImages {
		shopImages = append(shopImages, shopImage.Image.ImagePath)
	}
	return ShopCashierResponse{
		ID:          shop.ID,
		Name:        shop.Name,
		Cover:       shop.Cover,
		Description: shop.Description,
		FullAddress: shop.FullAddress,
		Lat:         shop.Lat,
		Lng:         shop.Lng,
		Status:      shop.Status,
		ShopImages:  shopImages,
		CreatedAt:   shop.CreatedAt,
		UpdatedAt:   shop.UpdatedAt,
	}
}

type ShopProductCashierResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	Code          string    `json:"code"`
	Price         int       `json:"price"`
	Unit          string    `json:"unit"`
	Stock         int       `json:"stock"`
	Description   string    `json:"description"`
	Status        int       `json:"status"`
	TypeProduct   string    `json:"type_product"`
	ProductImages []string  `json:"product_images"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ToShopProductCashierResponse(shopProducts *entities.ShopProductEntity) ShopProductCashierResponse {
	productImages := []string{}
	for _, productImage := range shopProducts.Product.ProductImages {
		productImages = append(productImages, productImage.Image.ImagePath)
	}

	return ShopProductCashierResponse{
		ID:            shopProducts.ID,
		Name:          shopProducts.Product.Name,
		Image:         shopProducts.Product.Image.ImagePath,
		Code:          shopProducts.Product.Code,
		Price:         int(shopProducts.Price),
		Unit:          shopProducts.Product.Unit,
		Stock:         shopProducts.Stock,
		Description:   shopProducts.Product.Description,
		Status:        shopProducts.Status,
		TypeProduct:   shopProducts.Product.TypeProduct.Name,
		ProductImages: productImages,
		CreatedAt:     shopProducts.CreatedAt,
		UpdatedAt:     shopProducts.UpdatedAt,
	}
}

type ShopCustomerResponse struct {
	ID           uuid.UUID                         `json:"id"`
	Name         string                            `json:"name"`
	Cover        string                            `json:"cover"`
	Description  string                            `json:"description"`
	FullAddress  string                            `json:"full_address"`
	Lat          float64                           `json:"lat"`
	Lng          float64                           `json:"lang"`
	Status       int                               `json:"status"`
	CreatedAt    time.Time                         `json:"created_at"`
	UpdatedAt    time.Time                         `json:"update_at"`
	DeletedAt    gorm.DeletedAt                    `json:"deleted_at"`
	ShopImages   []string                          `json:"shop_images"`
	ShopProducts []productresponse.ProductResponse `json:"shop_product"`
}

func ToShopCustomerResponse(shop entities.ShopEntity) ShopCustomerResponse {
	shopImages := []string{}
	for _, shopImage := range shop.ShopImages {
		shopImages = append(shopImages, shopImage.Image.ImagePath)
	}

	shopProducts := []productresponse.ProductResponse{}
	for _, shopProduct := range shop.ShopProducts {
		productRes := productresponse.ToProductResponse(&shopProduct.Product)
		shopProducts = append(shopProducts, productRes)
	}
	return ShopCustomerResponse{
		ID:           shop.ID,
		Name:         shop.Name,
		Cover:        shop.Cover,
		Description:  shop.Description,
		FullAddress:  shop.FullAddress,
		Lat:          shop.Lat,
		Lng:          shop.Lng,
		Status:       shop.Status,
		CreatedAt:    shop.CreatedAt,
		UpdatedAt:    shop.UpdatedAt,
		DeletedAt:    shop.DeletedAt,
		ShopImages:   shopImages,
		ShopProducts: shopProducts,
	}
}
