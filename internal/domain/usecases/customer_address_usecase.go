package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"nusantara_service/configs"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	dto "nusantara_service/internal/dto/request"
	"nusantara_service/internal/response"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CustomerAddressService struct {
	custAddrRepo repositories.CustomerAddressRepository
	rdb          *redis.Client
}

func NewCustomerAddressUsecase(custAddrRepo repositories.CustomerAddressRepository, rdb *redis.Client) services.CustomerAddressService {
	return &CustomerAddressService{custAddrRepo: custAddrRepo, rdb: rdb}
}

// CreateAddres implements services.CustomerAddressService.
func (c *CustomerAddressService) CreateAddres(ctx context.Context, custId uuid.UUID, req dto.CreateAddressRequest) error {
	if strings.TrimSpace(req.Label) == "" {
		return response.NewCustomError(response.ErrBadRequest, "label is required", 400)
	}
	if strings.TrimSpace(req.AddressText) == "" {
		return response.NewCustomError(response.ErrBadRequest, "Address text address is required", 400)
	}
	if req.Lat == nil {
		return response.NewCustomError(response.ErrBadRequest, "lat is required", 400)
	}
	if req.Lng == nil {
		return response.NewCustomError(response.ErrBadRequest, "lng is required", 400)
	}

	customer, err := c.custAddrRepo.FindCustomerID(ctx, custId.String())
	if err != nil {
		return response.NewCustomError(response.ErrNotFound, "user not permission", 404)
	}

	address := &entities.CustomerAddressEntity{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(customer.ID),
		Label:       req.Label,
		AddressText: req.AddressText,
		Lat:         *req.Lat,
		Lng:         *req.Lng,
		CreatedAt:   time.Now(),
	}

	count, _ := c.custAddrRepo.FindByUser(ctx, uuid.MustParse(customer.ID))
	if len(count) == 0 {
		address.IsDefault = true
	}

	if err := c.custAddrRepo.Create(ctx, address); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to create address", 500)
	}

	c.invalidateCache(ctx, custId)
	return nil
}

// GetAllAddress implements services.CustomerAddressService.
func (c *CustomerAddressService) GetAllAddress(ctx context.Context, custID uuid.UUID) ([]*entities.CustomerAddressEntity, error) {
	customer, err := c.custAddrRepo.FindCustomerID(ctx, custID.String())
	if err != nil {
		return nil, response.NewCustomError(response.ErrNotFound, "user not permission", 404)
	}
	cacheKey := fmt.Sprintf("addresses:%s", customer.ID)

	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var items []*entities.CustomerAddressEntity
		if json.Unmarshal([]byte(cached), &items) == nil {
			return items, nil
		}
	}

	addresses, err := c.custAddrRepo.FindByUser(ctx, uuid.MustParse(customer.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*entities.CustomerAddressEntity{}, nil
		}
		return nil, err
	}

	buf, _ := json.Marshal(addresses)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return addresses, nil
}

// GetByIdAddress implements services.CustomerAddressService.
func (c *CustomerAddressService) GetByIdAddress(ctx context.Context, id uuid.UUID, custID uuid.UUID) (*entities.CustomerAddressEntity, error) {
	customer, err := c.custAddrRepo.FindCustomerID(ctx, custID.String())
	if err != nil {
		return nil, response.NewCustomError(response.ErrNotFound, "user not permission", 404)
	}
	address, err := c.custAddrRepo.FindByID(ctx, id, uuid.MustParse(customer.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "address not found", 404)
		}
		return nil, err
	}
	return address, nil
}

// UpdateAddress implements services.CustomerAddressService.
func (c *CustomerAddressService) UpdateAddress(ctx context.Context, id uuid.UUID, custID uuid.UUID, req dto.UpdateAddressRequest) error {
	customer, err := c.custAddrRepo.FindCustomerID(ctx, custID.String())
	if err != nil {
		return response.NewCustomError(response.ErrNotFound, "user not permission", 404)
	}
	addr, err := c.custAddrRepo.FindByID(ctx, id, uuid.MustParse(customer.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "address not found", 404)
		}
		return err
	}

	if req.Label != "" {
		addr.Label = req.Label
	}
	if req.AddressText != "" {
		addr.AddressText = req.AddressText
	}
	if req.Lat != nil {
		addr.Lat = *req.Lat
	}
	if req.Lng != nil {
		addr.Lng = *req.Lng
	}

	if err := c.custAddrRepo.Update(ctx, addr.ID, uuid.MustParse(customer.ID), addr); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to update address", 500)
	}

	c.invalidateCache(ctx, uuid.MustParse(customer.ID))
	return nil
}

// Delete implements services.CustomerAddressService.
func (c *CustomerAddressService) Delete(ctx context.Context, id uuid.UUID, custID uuid.UUID) error {
	customer, err := c.custAddrRepo.FindCustomerID(ctx, custID.String())
	if err != nil {
		return response.NewCustomError(response.ErrNotFound, "user not permission", 404)
	}
	addr, err := c.custAddrRepo.FindByID(ctx, id, uuid.MustParse(customer.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "address not found", 404)
		}
		return err
	}

	if err := c.custAddrRepo.Delete(ctx, addr.ID, uuid.MustParse(customer.ID)); err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to delete address", 500)
	}

	if addr.IsDefault {
		_ = c.custAddrRepo.SetFirstAsDefault(ctx, custID)
	}

	c.invalidateCache(ctx, uuid.MustParse(customer.ID))
	return nil
}

// GetDefaultAddress implements services.CustomerAddressService.
func (c *CustomerAddressService) GetDefaultAddress(ctx context.Context, userID uuid.UUID) (*entities.CustomerAddressEntity, error) {
	customer, err := c.custAddrRepo.FindCustomerID(ctx, userID.String())
	if err != nil {
		return nil, response.NewCustomError(response.ErrNotFound, "user not permission", 404)
	}

	addr, err := c.custAddrRepo.FindDefaultByUser(ctx, uuid.MustParse(customer.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "no default address", 404)
		}
		return nil, err
	}

	return addr, nil
}

// SetDefaultAddress implements services.CustomerAddressService.
func (c *CustomerAddressService) SetDefaultAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error {
	customer, err := c.custAddrRepo.FindCustomerID(ctx, userID.String())
	if err != nil {
		return response.NewCustomError(response.ErrNotFound, "user not permission", 404)
	}

	addr, err := c.custAddrRepo.FindByID(ctx, addressID, uuid.MustParse(customer.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "address not found", 404)
		}
		return err
	}

	if err := c.custAddrRepo.ClearDefaultByUser(ctx, uuid.MustParse(customer.ID)); err != nil {
		return err
	}

	if err := c.custAddrRepo.SetDefaultByUser(ctx, uuid.MustParse(customer.ID), addr.ID); err != nil {
		return err
	}

	c.invalidateCache(ctx, uuid.MustParse(customer.ID))
	return nil
}

// GetNearbyShops implements services.CustomerAddressService.
func (c *CustomerAddressService) GetNearbyShops(ctx context.Context, lat float64, lng float64, maxDistance float64) ([]*entities.ShopEntity, map[uuid.UUID]float64, error) {
	shops, dist, err := c.custAddrRepo.FindNearby(ctx, lat, lng, maxDistance)
	if err != nil {
		return nil, nil, response.NewCustomError(response.ErrInternal, "failed to get nearby shops", 500)
	}
	return shops, dist, nil
}

func (c *CustomerAddressService) invalidateCache(ctx context.Context, custID uuid.UUID) {
	iter := c.rdb.Scan(ctx, 0, fmt.Sprintf("addresses:%s*", custID), 0).Iterator()
	for iter.Next(ctx) {
		c.rdb.Del(ctx, iter.Val())
	}
}
