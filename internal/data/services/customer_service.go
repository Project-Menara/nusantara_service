package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
)

type CustomerService interface {
	CheckPhone(ctx context.Context, req dto.CheckPhoneRequest) (*response.CheckPhoneResult, error)
	RegisterCustomer(ctx context.Context, req dto.RegisterCustomerRequest) (*entities.UserEntity, error)
	ResendCodeOTPVerify(ctx context.Context, req dto.ResendOTPRequest) error
	VerifyCodeOTP(ctx context.Context, req dto.VerifyOTPRequest) error
	NewPinCustomer(ctx context.Context, req dto.NewPinRequest) error
	ConfirmPinCustomer(ctx context.Context, req dto.ConfirmPinRequest) (*entities.UserEntity, string, error)
	GetProfileCustomer(ctx context.Context, userId string, token string) (*entities.UserEntity, error)
	CheckTokenBlacklisted(ctx context.Context, token string) (bool, error)
	LogoutCustomer(ctx context.Context, userId, token string) error
	LoginCustomer(ctx context.Context, req dto.LoginCustomerRequest) (string, error)
}
