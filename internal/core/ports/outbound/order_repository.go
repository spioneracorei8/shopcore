package outbound

import (
	"context"
	"shopcore/internal/core/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	FetchListOrders(ctx context.Context) ([]*domain.Order, error)
	FetchOrderById(ctx context.Context, id *bson.ObjectID) (*domain.Order, error)
	UpdateOrderById(ctx context.Context, id *bson.ObjectID, order *domain.Order) error
	DeleteOrderById(ctx context.Context, id *bson.ObjectID, order *domain.Order) error
}
