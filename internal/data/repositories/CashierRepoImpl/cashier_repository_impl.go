package cashierrepoimpl

import (
	"context"
	"errors"
	"nusantara_service/internal/domain/entities"
	cashierrepo "nusantara_service/internal/domain/repositories/CashierRepo"
	shopresponse "nusantara_service/internal/dto/responses/shop_response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CashierRepositoryImpl struct {
	db *gorm.DB
}

func NewCashierRepositoryImpl(db *gorm.DB) cashierrepo.ICashierRepository {
	return &CashierRepositoryImpl{db: db}
}

// GetShop implements cashierrepo.ICashierRepository.
func (c *CashierRepositoryImpl) GetShop(ctx context.Context, cashierID uuid.UUID) ([]shopresponse.ShopNameCashierResponse, error) {
	var shopNames []shopresponse.ShopNameCashierResponse

	err := c.db.WithContext(ctx).
		Table("shops").
		Select("shops.id", "shops.name").
		Joins("JOIN shop_cashiers ON shop_cashiers.shop_id = shops.id").
		Where("shops.deleted_at IS NULL").
		Where("shop_cashiers.cashier_id = ?", cashierID).
		Where("shop_cashiers.deleted_at IS NULL").
		Order("shop_cashiers.created_at ASC").
		Scan(&shopNames).Error

	if err != nil {
		return nil, err
	}

	return shopNames, nil
}

// FindCashier implements cashierrepo.ICashierRepository.
func (c *CashierRepositoryImpl) FindCashier(ctx context.Context, userID string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := c.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", userID, "admin").
		Preload("Role").
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetDetailShop implements cashierrepo.ICashierRepository.
func (c *CashierRepositoryImpl) GetDetailShop(ctx context.Context, cashierID uuid.UUID, shopID uuid.UUID) (*shopresponse.ShopCashierResponse, error) {
	var shop shopresponse.ShopCashierResponse

	// --- 1. Ambil Detail Toko ---
	err := c.db.WithContext(ctx).
		Table("shops").
		Joins("JOIN shop_cashiers ON shop_cashiers.shop_id = shops.id").
		Where("shops.id = ?", shopID).
		Where("shop_cashiers.cashier_id = ?", cashierID).
		Where("shop_cashiers.deleted_at IS NULL").
		Select("shops.*").
		First(&shop).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shop not found or not assigned to cashier")
		}
		return nil, err
	}

	var imagePaths []string

	err = c.db.WithContext(ctx).
		Table("shop_images").
		Joins("JOIN images ON images.id = shop_images.image_id").
		Where("shop_images.shop_id = ?", shopID).
		Where("shop_images.deleted_at IS NULL").
		Pluck("images.image_path", &imagePaths).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	shop.ShopImages = imagePaths

	return &shop, nil
}

// GetDetailProduct implements cashierrepo.ICashierRepository.
func (c *CashierRepositoryImpl) GetDetailProduct(ctx context.Context, cashierID uuid.UUID, shopID uuid.UUID) ([]shopresponse.ShopProductCashierResponse, error) {
	var shopProducts []entities.ShopProductEntity

	err := c.db.WithContext(ctx).
		Joins("JOIN shops ON shops.id = shop_products.shop_id").
		Joins("JOIN shop_cashiers ON shop_cashiers.shop_id = shops.id").
		Where("shop_products.shop_id = ? AND shop_cashiers.cashier_id = ? AND shop_products.deleted_at IS NULL AND shop_cashiers.deleted_at IS NULL", shopID, cashierID).
		Preload("Product.Image").
		Preload("Product.ProductImages.Image").
		Preload("Product.TypeProduct").
		Find(&shopProducts).Error

	if err != nil {
		return nil, err
	}

	var responses []shopresponse.ShopProductCashierResponse
	for _, sp := range shopProducts {
		product := sp.Product

		// Ambil semua image tambahan dari ProductImages
		productImages := []string{}
		for _, productImage := range product.ProductImages {
			productImages = append(productImages, productImage.Image.ImagePath)
		}

		res := shopresponse.ShopProductCashierResponse{
			ID:            product.ID,
			Name:          product.Name,
			Image:         product.Image.ImagePath,
			Code:          product.Code,
			Price:         int(sp.Price),
			Unit:          product.Unit,
			Stock:         sp.Stock,
			Description:   product.Description,
			Status:        sp.Status,
			TypeProduct:   product.TypeProduct.Name,
			ProductImages: productImages,
		}

		responses = append(responses, res)
	}

	return responses, nil
}
