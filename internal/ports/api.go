package ports

import (
	"context"

	"github.com/aholake/order-service/internal/application/core/domain"
)

type APIPort interface {
	PlaceOrder(context.Context, domain.Order) (domain.Order, error)
}
