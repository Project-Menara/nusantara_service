package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	dto "nusantara_service/internal/dto/request"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"nusantara_service/internal/workers/payload"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CashierService struct {
	cashierRepo repositories.CashierRepository
	roleRepo    repositories.RoleRepository
	rdb         *redis.Client
	cloudinary  cloudinary.CloudinaryService
}

func NewCashierUsecase(cashierRepo repositories.CashierRepository, roleRepo repositories.RoleRepository, rdb *redis.Client, cloudinary *cloudinary.CloudinaryService) services.CashierService {
	return &CashierService{cashierRepo: cashierRepo, roleRepo: roleRepo, rdb: rdb, cloudinary: *cloudinary}
}

func (c *CashierService) CreateCashier(ctx context.Context, req dto.CreateCashierRequest) error {
	role, err := c.roleRepo.FindByName(ctx, "admin")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	name := strings.TrimSpace(req.Name)
	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	existing, err := c.cashierRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existing != nil {
		return response.NewCustomError(response.ErrExists, "email already exists", 409)
	}

	existingByUsername, err := c.cashierRepo.FindByUsername(ctx, username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingByUsername != nil {
		return response.NewCustomError(response.ErrExists, "username already exists", 409)
	}
	if name == "" {
		return response.NewCustomError(response.ErrBadRequest, "name is required", 400)
	}
	if username == "" {
		return response.NewCustomError(response.ErrBadRequest, "username is required", 400)
	}

	if email == "" {
		return response.NewCustomError(response.ErrBadRequest, "email is required", 400)
	}
	if !strings.HasSuffix(strings.ToLower(email), "@gmail.com") {
		return response.NewCustomError(response.ErrBadRequest, "only Gmail addresses are allowed", 400)
	}
	if password == "" {
		return response.NewCustomError(response.ErrBadRequest, "password is required", 400)
	}

	if req.Photo == nil {
		return response.NewCustomError(response.ErrBadRequest, "image is required", 400)
	}

	if req.Status < 0 || req.Status > 1 {
		return response.NewCustomError(response.ErrBadRequest, "status must be 0 & 1", 400)
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to hash password", 500)
	}

	newCashier := &entities.UserEntity{
		ID:       uuid.NewString(),
		Name:     name,
		Username: username,
		Email:    email,
		Status:   req.Status,
		Password: hashed,
		RoleID:   role.ID,
	}

	err = c.cashierRepo.Create(ctx, newCashier)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to create cashier", 500)
	}

	file, err := req.Photo.Open()
	if err != nil {
		fmt.Println("Error opening image file:", err)
	} else {
		defer file.Close()
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading image file:", err)
		} else {
			payload := payload.ImageSendTask{
				UserID:    uuid.MustParse(newCashier.ID),
				FileBytes: fileBytes,
				Folder:    "nusantara_service/cashier",
				Filename:  fmt.Sprintf("cashier_%s", newCashier.ID),
			}
			err = rabbitmq.PublishToQueue("", rabbitmq.SendImageQueueName, payload)
			if err != nil {
				fmt.Println("Error publishing to RabbitMQ:", err)
			}
		}
	}

	c.InvalidateCashierCache(ctx)

	return nil
}

func (c *CashierService) GetCashierAll(ctx context.Context, page int, limit int, search string) ([]*entities.UserEntity, int, error) {
	cacheKey := fmt.Sprintf("cashiers:search:%s:page:%d:limit:%d", search, page, limit)
	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil {
		var result struct {
			Data  []*entities.UserEntity `json:"data"`
			Total int                    `json:"total"`
		}

		_ = json.Unmarshal([]byte(cached), &result)
		return result.Data, result.Total, nil
	}

	offset := (page - 1) * limit

	cashiers, total, err := c.cashierRepo.FindByAll(ctx, offset, limit, search)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "Cashier Not Found", 404)
	}

	dataCache, _ := json.Marshal(map[string]interface{}{
		"data":  cashiers,
		"total": total,
	})

	_ = configs.SetRedis(ctx, cacheKey, dataCache, time.Minute*30)
	return cashiers, total, nil
}

func (c *CashierService) GetCashierById(ctx context.Context, id uuid.UUID) (*entities.UserEntity, error) {
	cacheKey := fmt.Sprintf("cashier:%s", id)
	cached, err := configs.GetRedis(ctx, cacheKey)
	if err == nil && cached != "" {
		var cashier entities.UserEntity
		if err := json.Unmarshal([]byte(cached), &cashier); err == nil {
			return &cashier, nil
		}
	}

	cashier, err := c.cashierRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "cashier not found", 404)
		}
		return nil, err
	}

	dataCache, _ := json.Marshal(cashier)
	_ = configs.SetRedis(ctx, cacheKey, dataCache, time.Minute*30)

	return cashier, nil
}

func (c *CashierService) UpdateCashier(ctx context.Context, id uuid.UUID, req dto.UpdateCashierRequest) error {
	cashier, err := c.cashierRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "cashier not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, "failed to find cashier", 500)
	}

	if req.Name != nil {
		cashier.Name = *req.Name
	}
	if req.Username != nil {
		cashier.Username = *req.Username
	}
	if req.Email != nil {
		cashier.Email = *req.Email
	}
	if req.Status != nil {
		cashier.Status = *req.Status
	}

	if req.Photo != nil {
		file, err := req.Photo.Open()
		if err != nil {
			fmt.Println("Error opening image file for update:", err)
		} else {
			defer file.Close()
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				fmt.Println("Error reading image file for update:", err)
			} else {
				taskPayload := payload.ImageSendTask{
					UserID:    uuid.MustParse(cashier.ID),
					FileBytes: fileBytes,
					Folder:    "nusantara_service/cashier",
					Filename:  fmt.Sprintf("cashier_%s", id),
				}
				err = rabbitmq.PublishToQueue("", rabbitmq.SendImageQueueName, taskPayload)
				if err != nil {
					fmt.Println("Error publishing update to RabbitMQ:", err)
				}
			}
		}
	}

	err = c.cashierRepo.Update(ctx, id, cashier)
	if err != nil {
		return response.NewCustomError(response.ErrInternal, "failed to update cashier data", 500)
	}

	c.InvalidateCashierCache(ctx)

	return nil
}

func (c *CashierService) Delete(ctx context.Context, id uuid.UUID) error {
	cashier, err := c.cashierRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "cashier not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, "failed to find cashier", 500)
	}

	if cashier.Photo != nil {
		publicID := utils.ExtractPublicIDFromCloudinaryURL(*cashier.Photo)
		if publicID != "" {
			if err := c.cloudinary.DestroyImage(ctx, publicID); err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to delete image", 500)
			}
		}
	}

	c.InvalidateCashierCache(ctx)

	return c.cashierRepo.Delete(ctx, id)
}

func (c *CashierService) InvalidateCashierCache(ctx context.Context) {
	iter := c.rdb.Scan(ctx, 0, "cashiers:*", 0).Iterator()
	for iter.Next(ctx) {
		c.rdb.Del(ctx, iter.Val())
	}

	iterID := c.rdb.Scan(ctx, 0, "cashier:*", 0).Iterator()
	for iterID.Next(ctx) {
		c.rdb.Del(ctx, iterID.Val())
	}
}
