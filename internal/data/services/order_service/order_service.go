package orderservice

import (
	"context"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID uuid.UUID) error
}
