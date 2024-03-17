package ports

import "context"

type ShippingClientPort interface {
	Ship(ctx context.Context, orderId int64, address string) error
}
