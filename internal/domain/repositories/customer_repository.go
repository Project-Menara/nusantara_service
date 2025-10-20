package repositories

import (
	"context"
	"nusantara_service/internal/domain/entities"
	cartresponse "nusantara_service/internal/dto/responses/cart_response"
	shopresponse "nusantara_service/internal/dto/responses/shop_response"
	"time"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	FindByPhoneCustomer(ctx context.Context, phone string) (*entities.UserEntity, error)
	FindByUsername(ctx context.Context, username string) (*entities.UserEntity, error)
	FindByEmail(ctx context.Context, email string) (*entities.UserEntity, error)
	FindByPhone(ctx context.Context, phone string) (*entities.UserEntity, error)
	FindRoleByName(ctx context.Context, role string) (*entities.RoleEntity, error)
	CreateCustomer(ctx context.Context, user *entities.UserEntity) (*entities.UserEntity, error)
	UpdateStatusCustomer(ctx context.Context, userID string, status int) error
	UpdatePinCustomer(ctx context.Context, userID string, pin string) error
	FindByUseIDCustomer(ctx context.Context, userID string) (*entities.UserEntity, error)
	UpdateCustomer(ctx context.Context, userId string, data *entities.UserEntity) (*entities.UserEntity, error)

	//Change Phone
	ChangePhoneCustomer(ctx context.Context, userId string, data *entities.UserEntity) (*entities.UserEntity, error)

	CreateCustomerPoint(ctx context.Context, userPoint *entities.UserPointEntity) error
	AddVoucher(ctx context.Context, userVoucher *entities.UserVoucherEntity) (*entities.UserVoucherEntity, error)
	AddDetailVoucher(ctx context.Context, detailVoucher *entities.UserVoucherDetailEntity) (*entities.UserVoucherDetailEntity, error)
	GetCustomerPoint(ctx context.Context, customerID uuid.UUID) (*entities.UserPointEntity, error)
	UpdateCustomerPoint(ctx context.Context, customerID uuid.UUID, points int) error
	CreatePointHistory(ctx context.Context, history *entities.UserPointHistoriesEntity) error
	FindUserPoint(ctx context.Context, customerID uuid.UUID) (*entities.UserPointEntity, int, *time.Time, error)
	MarkPointsAsExpired(ctx context.Context, userID uuid.UUID, expiredDate time.Time) error
	DecreaseTotalPoints(ctx context.Context, customerID uuid.UUID, amount int) error
	FindUserPointHistory(ctx context.Context, customerID uuid.UUID) ([]*entities.UserPointHistoriesEntity, error)
	FindUserVoucherClaimed(ctx context.Context, userId uuid.UUID) ([]*entities.UserVoucherEntity, error)

	FindShopById(ctx context.Context, offset int, limit int, search string, typeID uuid.UUID, shopID uuid.UUID) (*shopresponse.ShopCustomerResponse, int, error)
	GetMyCart(ctx context.Context, customerID uuid.UUID) (*cartresponse.CartResponse, error)
}
