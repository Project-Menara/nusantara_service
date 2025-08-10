package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VoucherService struct {
	repo       repositories.VoucherRepository
	rdb        *redis.Client
	cloudinary cloudinary.CloudinaryService
}

func NewVoucherUsecase(repo repositories.VoucherRepository, rdb *redis.Client, cloudinary *cloudinary.CloudinaryService) services.VoucherService {
	return &VoucherService{repo: repo, rdb: rdb, cloudinary: *cloudinary}
}

// CreateVoucher implements services.VoucherService.
func (v *VoucherService) CreateVoucher(ctx context.Context, userId string, req dto.CreateVoucherRequest) (*entities.VoucherEntity, error) {
	existing, err := v.repo.FindByCode(ctx, req.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, response.NewCustomError(response.ErrNotFound, err.Error(), 404)
	}

	if existing != nil {
		return nil, response.NewCustomError(response.ErrExists, "code already exists", 409)
	}

	// Validasi field wajib
	if strings.TrimSpace(req.Code) == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "code is required", 400)
	}
	if req.MinimumSpend <= 0 {
		return nil, response.NewCustomError(response.ErrBadRequest, "minimum spend must be more than 0", 400)
	}
	if req.PointCost <= 0 {
		return nil, response.NewCustomError(response.ErrBadRequest, "point cost must be more than 0", 400)
	}
	if req.Quota <= 0 {
		return nil, response.NewCustomError(response.ErrBadRequest, "quota must be more than 0", 400)
	}
	if req.StartDate.IsZero() || req.EndDate.IsZero() {
		return nil, response.NewCustomError(response.ErrBadRequest, "start date and end date are required", 400)
	}
	if req.EndDate.Before(req.StartDate) {
		return nil, response.NewCustomError(response.ErrBadRequest, "end date cannot be before start date", 400)
	}
	if req.DiscountType != "amount" && req.DiscountType != "percent" {
		return nil, response.NewCustomError(response.ErrBadRequest, "discount type must be 'amount' or 'percent'", 400)
	}

	// Validasi sesuai tipe diskon & set field yang tidak dipakai jadi 0
	if req.DiscountType == "amount" {
		if req.DiscountAmount <= 0 {
			return nil, response.NewCustomError(response.ErrBadRequest, "discount amount must be more than 0", 400)
		}
		req.DiscountPercent = 0
	}

	// Strict rule: Tidak boleh mengisi keduanya sekaligus
	if req.DiscountAmount > 0 && req.DiscountPercent > 0 {
		return nil, response.NewCustomError(response.ErrBadRequest, "cannot set both discount amount and discount percent", 400)
	}

	if req.DiscountType == "percent" {
		if req.DiscountPercent <= 0 || req.DiscountPercent > 100 {
			return nil, response.NewCustomError(response.ErrBadRequest, "discount percent must be between 1 and 100", 400)
		}
		req.DiscountAmount = 0
	}

	if req.Description == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "description is required", 400)
	}

	if req.Status < 0 || req.Status > 1 {
		return nil, response.NewCustomError(response.ErrBadRequest, "Status must be 0 & 1", 400)
	}

	user, err := v.repo.FindByUserIDSuperAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	voucher := &entities.VoucherEntity{
		ID:              uuid.New(),
		Code:            strings.TrimSpace(req.Code),
		DiscountAmount:  req.DiscountAmount,
		DiscountPercent: req.DiscountPercent,
		MinimumSpend:    req.MinimumSpend,
		PointCost:       req.PointCost,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Quota:           req.Quota,
		Description:     req.Description,
		DiscountType:    req.DiscountType,
		Status:          req.Status,
		CreatedBy:       uuid.MustParse(userId),
		User:            *user,
	}
	_, err = v.repo.Create(ctx, voucher)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, err.Error(), 500)
	}

	v.InvalidateVoucherCache(ctx)

	return voucher, nil
}

// GetAllVoucher implements services.VoucherService.
func (v *VoucherService) GetAllVoucher(ctx context.Context, page int, limit int, search string) ([]*entities.VoucherEntity, int, error) {
	cacheKey := fmt.Sprintf("vouchers:search:%s:page:%d:limit:%d", search, page, limit)
	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil {
		var result struct {
			Data  []*entities.VoucherEntity `json:"data"`
			Total int                       `json:"total"`
		}

		_ = json.Unmarshal([]byte(cached), &result)
		return result.Data, result.Total, nil
	}

	offset := (page - 1) * limit

	vouchers, err := v.repo.FindAll(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "Voucher not found", 404)
	}

	total, err := v.repo.CountAll(ctx, search)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "failed to count voucher", 400)
	}

	dataToCache, _ := json.Marshal(map[string]interface{}{
		"data":  vouchers,
		"total": total,
	})

	_ = configs.SetRedis(ctx, cacheKey, dataToCache, time.Minute*30)
	return vouchers, total, nil
}

// GetByIdVoucher implements services.VoucherService.
func (v *VoucherService) GetByIdVoucher(ctx context.Context, id uuid.UUID) (*entities.VoucherEntity, error) {
	cacheKey := fmt.Sprintf("voucher:%s", id.String())

	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil && cached != "" {
		var voucher entities.VoucherEntity
		if err := json.Unmarshal([]byte(cached), &voucher); err != nil {
			return &voucher, nil
		}
	}

	voucher, err := v.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "Voucher not found", 404)
		}
		return nil, err
	}

	dataCache, _ := json.Marshal(voucher)
	_ = configs.SetRedis(ctx, cacheKey, dataCache, time.Minute*30)

	return voucher, nil
}

// UpdateVoucher implements services.VoucherService.
func (v *VoucherService) UpdateVoucher(ctx context.Context, userId string, id uuid.UUID, req dto.UpdateVoucherRequest) (*entities.VoucherEntity, error) {
	voucher, err := v.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "voucher not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, err.Error(), 500)
	}

	// Validasi field wajib
	if strings.TrimSpace(req.Code) == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "code is required", 400)
	}
	if req.MinimumSpend <= 0 {
		return nil, response.NewCustomError(response.ErrBadRequest, "minimum spend must be more than 0", 400)
	}
	if req.PointCost <= 0 {
		return nil, response.NewCustomError(response.ErrBadRequest, "point cost must be more than 0", 400)
	}
	if req.Quota <= 0 {
		return nil, response.NewCustomError(response.ErrBadRequest, "quota must be more than 0", 400)
	}
	if req.StartDate.IsZero() || req.EndDate.IsZero() {
		return nil, response.NewCustomError(response.ErrBadRequest, "start date and end date are required", 400)
	}
	if req.EndDate.Before(req.StartDate) {
		return nil, response.NewCustomError(response.ErrBadRequest, "end date cannot be before start date", 400)
	}
	if req.DiscountType != "amount" && req.DiscountType != "percent" {
		return nil, response.NewCustomError(response.ErrBadRequest, "discount type must be 'amount' or 'percent'", 400)
	}
	if req.DiscountAmount > 0 && req.DiscountPercent > 0 {
		return nil, response.NewCustomError(response.ErrBadRequest, "cannot set both discount amount and discount percent", 400)
	}

	// Validasi sesuai tipe diskon
	if req.DiscountType == "amount" {
		if req.DiscountAmount <= 0 {
			return nil, response.NewCustomError(response.ErrBadRequest, "discount amount must be more than 0", 400)
		}
		req.DiscountPercent = 0
	}
	if req.DiscountType == "percent" {
		if req.DiscountPercent <= 0 || req.DiscountPercent > 100 {
			return nil, response.NewCustomError(response.ErrBadRequest, "discount percent must be between 1 and 100", 400)
		}
		req.DiscountAmount = 0
	}

	user, err := v.repo.FindByUserIDSuperAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}

	// Update data
	voucher.Code = strings.TrimSpace(req.Code)
	voucher.DiscountAmount = req.DiscountAmount
	voucher.DiscountPercent = req.DiscountPercent
	voucher.MinimumSpend = req.MinimumSpend
	voucher.PointCost = req.PointCost
	voucher.StartDate = req.StartDate
	voucher.EndDate = req.EndDate
	voucher.Quota = req.Quota
	voucher.Description = req.Description
	voucher.DiscountType = req.DiscountType
	voucher.CreatedBy = uuid.MustParse(user.ID)

	updatedVoucher, err := v.repo.Update(ctx, id, voucher)
	if err != nil {
		return nil, response.NewCustomError(response.ErrInternal, "Failed to update voucher", 500)
	}

	v.InvalidateVoucherCache(ctx)

	return updatedVoucher, nil
}

// DeleteVoucher implements services.VoucherService.
func (v *VoucherService) DeleteVoucher(ctx context.Context, id uuid.UUID) error {
	_, err := v.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "voucher not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, err.Error(), 500)
	}

	v.InvalidateVoucherCache(ctx)
	return v.repo.Delete(ctx, id)
}

// UpdateStatusVoucher implements services.VoucherService.
func (v *VoucherService) UpdateStatusVoucher(ctx context.Context, id uuid.UUID, req dto.UpdateStatusVoucherRequest) error {
	voucher, err := v.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "voucher not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, err.Error(), 500)
	}
	if req.Status < 0 || req.Status > 1 {
		return response.NewCustomError(response.ErrBadRequest, "Status must be 0 or 1", 400)
	}

	voucher.Status = req.Status

	err = v.repo.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "Failed to update status voucher", 500)
	}

	v.InvalidateVoucherCache(ctx)

	return nil
}

func (v *VoucherService) InvalidateVoucherCache(ctx context.Context) {
	iter := v.rdb.Scan(ctx, 0, "vouchers:*", 0).Iterator()
	for iter.Next(ctx) {
		v.rdb.Del(ctx, iter.Val())
	}

	iterID := v.rdb.Scan(ctx, 0, "voucher:*", 0).Iterator()
	for iterID.Next(ctx) {
		v.rdb.Del(ctx, iterID.Val())
	}
}

// GetAllVoucherCustomer implements services.VoucherService.
func (v *VoucherService) GetAllVoucherCustomer(ctx context.Context, page, limit int) ([]*entities.VoucherEntity, int, error) {
	cacheKey := fmt.Sprintf("vouchers:page:%d:limit:%d", page, limit)
	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil {
		var result struct {
			Data  []*entities.VoucherEntity `json:"data"`
			Total int                       `json:"total"`
		}

		_ = json.Unmarshal([]byte(cached), &result)
		return result.Data, result.Total, nil
	}

	offset := (page - 1) * limit

	vouchers, err := v.repo.GetAllVoucherCustomer(ctx, limit, offset)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "Voucher not found", 404)
	}

	total, err := v.repo.CountAll(ctx, "")
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "failed to count voucher", 400)
	}

	dataToCache, _ := json.Marshal(map[string]interface{}{
		"data":  vouchers,
		"total": total,
	})

	_ = configs.SetRedis(ctx, cacheKey, dataToCache, time.Minute*30)
	return vouchers, total, nil
}

// GetByIdVoucherCustomer implements services.VoucherService.
func (v *VoucherService) GetByIdVoucherCustomer(ctx context.Context, id uuid.UUID) (*entities.VoucherEntity, error) {
	cacheKey := fmt.Sprintf("voucher:%s", id.String())

	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil && cached != "" {
		var voucher entities.VoucherEntity
		if err := json.Unmarshal([]byte(cached), &voucher); err != nil {
			return &voucher, nil
		}
	}

	voucher, err := v.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "Voucher not found", 404)
		}
		return nil, err
	}

	dataCache, _ := json.Marshal(voucher)
	_ = configs.SetRedis(ctx, cacheKey, dataCache, time.Minute*30)

	return voucher, nil
}
