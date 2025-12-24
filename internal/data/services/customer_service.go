package services

import (
	"context"
	"mime/multipart"
	"nusantara_service/internal/domain/entities"
	dto "nusantara_service/internal/dto/request"
	cartresponse "nusantara_service/internal/dto/responses/cart_response"
	favoriteresponse "nusantara_service/internal/dto/responses/favorite_response"
	shopresponse "nusantara_service/internal/dto/responses/shop_response"
	"nusantara_service/internal/response"
	"time"

	"github.com/google/uuid"
)

type CustomerService interface {
	CheckPhone(ctx context.Context, req dto.CheckPhoneRequest) (*response.CheckPhoneResult, error)
	RegisterCustomer(ctx context.Context, req dto.RegisterCustomerRequest) (*entities.UserEntity, time.Duration, error)
	ResendCodeOTPVerify(ctx context.Context, req dto.ResendOTPRequest) error
	VerifyCodeOTP(ctx context.Context, req dto.VerifyOTPRequest) error
	NewPinCustomer(ctx context.Context, req dto.NewPinRequest) error
	ConfirmPinCustomer(ctx context.Context, req dto.ConfirmPinRequest) (*entities.UserEntity, string, error)
	GetProfileCustomer(ctx context.Context, userId string, token string) (*entities.UserEntity, error)
	CheckTokenBlacklisted(ctx context.Context, token string) (bool, error)
	LogoutCustomer(ctx context.Context, userId, token string) error
	LoginCustomer(ctx context.Context, req dto.LoginCustomerRequest) (string, error)
	UpdateProfileCustomer(ctx context.Context, userId string, req dto.UpdateCustomerRequest, photoFileHeader *multipart.FileHeader) (*entities.UserEntity, error)
	FogotPIN(ctx context.Context, req dto.ForgotPINRequest) (string, error)
	ValidateTokenForgotPIN(ctx context.Context, token string) (string, error)
	CreateNewPIN(ctx context.Context, token string, req dto.CreateNewPinRequest) error
	CreateNewConfirmPIN(ctx context.Context, tokenPIN string, req dto.CreateConfirmPinRequest) (*entities.UserEntity, string, error)

	//Change Phone
	VerifyPINCustomer(ctx context.Context, userId string, req dto.VerifyPINCustomerRequest) error
	NewPINCustomer(ctx context.Context, userId string, req dto.NewPINCustomer) error
	ConfirmationPINCustomerUpdate(ctx context.Context, userId string, req dto.ConfirmNewPINCustomer) (*entities.UserEntity, error)
	NewPhoneCustomer(ctx context.Context, userId string, req dto.NewPhoneCustomerRequest) error
	VerifyCodeOTPCustomerUpdate(ctx context.Context, userId string, req dto.VerifyOTPCustomerUpdateRequest) (*entities.UserEntity, error)

	ClaimVoucherCustomer(ctx context.Context, customerID uuid.UUID, voucherID uuid.UUID) (*entities.UserVoucherEntity, error)
	GetCustomerPoint(ctx context.Context, customerID uuid.UUID) (*entities.UserPointEntity, int, *time.Time, error)
	GetCustomerPointHistory(ctx context.Context, customerID uuid.UUID) ([]*entities.UserPointHistoriesEntity, error)
	GetCustomerVouchersClaimed(ctx context.Context, customerID uuid.UUID) ([]*entities.UserVoucherEntity, error)

	GetDetailShopById(ctx context.Context, page, limit int, search string, typeID uuid.UUID, shopID uuid.UUID) (*shopresponse.ShopCustomerResponse, int, error)

	GetMyCart(ctx context.Context, customerID uuid.UUID) (*cartresponse.CartResponse, error)
	AddProductToMyCart(ctx context.Context, customerID uuid.UUID, shopId uuid.UUID, req dto.AddCartItemRequest) error
	DeleteCartItem(ctx context.Context, customerID uuid.UUID, productID uuid.UUID) error

	GetMyFavorite(ctx context.Context, customerID uuid.UUID) (*favoriteresponse.FavoriteResponse, error)
	AddProductToFavorite(ctx context.Context, customerID uuid.UUID, req dto.AddFavoriteItemRequest) error
	DeleteFavoriteItem(ctx context.Context, customerID uuid.UUID, productID uuid.UUID) error
}
