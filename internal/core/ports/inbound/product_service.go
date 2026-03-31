package inbound

import (
	"context"
	"shopcore/internal/core/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	FetchListProducts(ctx context.Context) ([]*domain.Product, error)
	FetchProductById(ctx context.Context, id *bson.ObjectID) (*domain.Product, error)
	UpdateProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) (*domain.Product, error)
	DeleteProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) error
}
