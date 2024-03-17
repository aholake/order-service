package ports

import (
	"context"

	"github.com/aholake/order-service/internal/application/core/domain"
)

type PaymentClientPort interface {
	Charge(ctx context.Context, order *domain.Order) error
}
