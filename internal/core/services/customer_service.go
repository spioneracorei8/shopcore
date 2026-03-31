package services

import (
	"context"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/inbound"
	"shopcore/internal/core/ports/outbound"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type customerUsecase struct {
	customerRepo outbound.CustomerRepository
}

func NewCustomerUsecaseImpl(customerRepo outbound.CustomerRepository) inbound.CustomerUsecase {
	return &customerUsecase{
		customerRepo: customerRepo,
	}
}

func (u *customerUsecase) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	customer.GenObjectID()
	customer.Status = domain.CUSTOMER_STATUS_ACTIVE
	customer.SetCreatedAt()
	customer.SetUpdatedAt()
	return u.customerRepo.CreateCustomer(ctx, customer)
}

func (u *customerUsecase) FetchListCustomers(ctx context.Context) ([]*domain.Customer, error) {
	return u.customerRepo.FetchListCustomers(ctx)
}

func (u *customerUsecase) FetchCustomerById(ctx context.Context, id *bson.ObjectID) (*domain.Customer, error) {
	return u.customerRepo.FetchCustomerById(ctx, id)
}

func (u *customerUsecase) UpdateCustomerById(ctx context.Context, id *bson.ObjectID, customer *domain.Customer) (*domain.Customer, error) {
	customer.Id = id
	customer.SetUpdatedAt()
	if err := u.customerRepo.UpdateCustomerById(ctx, id, customer); err != nil {
		return nil, err
	}

	return u.customerRepo.FetchCustomerById(ctx, id)
}

func (u *customerUsecase) DeleteCustomerById(ctx context.Context, id *bson.ObjectID, customer *domain.Customer) error {
	customer.Id = id
	customer.Status = domain.CUSTOMER_STATUS_INACTIVE
	customer.SetUpdatedAt()
	customer.SetDeletedAt()
	return u.customerRepo.DeleteCustomerById(ctx, id, customer)
}
