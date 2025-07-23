package services

import (
	"context"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/dto"
)

type CustomerService interface {
	CheckPhone(ctx context.Context, req dto.CheckPhoneRequest) (*entities.UserEntity, error)
}
