package inbound

import (
	"context"
	"shopcore/internal/core/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CustomerUsecase interface {
	CreateCustomer(ctx context.Context, customer *domain.Customer) error
	FetchListCustomers(ctx context.Context) ([]*domain.Customer, error)
	FetchCustomerById(ctx context.Context, id *bson.ObjectID) (*domain.Customer, error)
	UpdateCustomerById(ctx context.Context, id *bson.ObjectID, customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomerById(ctx context.Context, id *bson.ObjectID, customer *domain.Customer) error
}
