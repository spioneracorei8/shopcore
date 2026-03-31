package services

import (
	"context"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/inbound"
	"shopcore/internal/core/ports/outbound"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type productUsecase struct {
	productRepo outbound.ProductRepository
}

func NewProductUsecaseImpl(productRepo outbound.ProductRepository) inbound.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	product.GenObjectID()
	product.SetCreatedAt()
	product.SetUpdatedAt()
	return u.productRepo.CreateProduct(ctx, product)
}

func (u *productUsecase) FetchListProducts(ctx context.Context) ([]*domain.Product, error) {
	return u.productRepo.FetchListProducts(ctx)
}

func (u *productUsecase) FetchProductById(ctx context.Context, id *bson.ObjectID) (*domain.Product, error) {
	return u.productRepo.FetchProductById(ctx, id)
}

func (u *productUsecase) UpdateProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) (*domain.Product, error) {
	product.Id = id
	product.SetUpdatedAt()
	if err := u.productRepo.UpdateProductById(ctx, id, product); err != nil {
		return nil, err
	}
	return u.productRepo.FetchProductById(ctx, id)
}

func (u *productUsecase) DeleteProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) error {
	product.Id = id
	product.Status = domain.PRODUCT_STATUS_INACTIVE
	product.SetUpdatedAt()
	product.SetDeletedAt()
	return u.productRepo.DeleteProductById(ctx, id, product)
}
