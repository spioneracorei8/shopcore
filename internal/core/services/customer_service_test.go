package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/outbound"
	"shopcore/internal/core/services"
)

type mockCustomerRepository struct {
	mock.Mock
}

func (m *mockCustomerRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *mockCustomerRepository) FetchListCustomers(ctx context.Context) ([]*domain.Customer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Customer), args.Error(1)
}

func (m *mockCustomerRepository) FetchCustomerById(ctx context.Context, id *bson.ObjectID) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *mockCustomerRepository) UpdateCustomerById(ctx context.Context, id *bson.ObjectID, customer *domain.Customer) error {
	args := m.Called(ctx, id, customer)
	return args.Error(0)
}

func (m *mockCustomerRepository) DeleteCustomerById(ctx context.Context, id *bson.ObjectID, customer *domain.Customer) error {
	args := m.Called(ctx, id, customer)
	return args.Error(0)
}

var _ outbound.CustomerRepository = (*mockCustomerRepository)(nil)

func TestCustomer_CreateCustomer_Success(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	customer := &domain.Customer{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
	}

	mockRepo.On("CreateCustomer", mock.Anything, customer).Return(nil)

	err := usecase.CreateCustomer(context.Background(), customer)

	assert.NoError(t, err)
	assert.NotNil(t, customer.Id)
	assert.Equal(t, domain.CUSTOMER_STATUS_ACTIVE, customer.Status)
	assert.False(t, customer.CreatedAt.IsZero())
	assert.False(t, customer.UpdatedAt.IsZero())
	mockRepo.AssertExpectations(t)
}

func TestCustomer_CreateCustomer_RepoError(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	customer := &domain.Customer{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
	}

	mockRepo.On("CreateCustomer", mock.Anything, customer).Return(errors.New("db error"))

	err := usecase.CreateCustomer(context.Background(), customer)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomer_FetchListCustomers_Success(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	expectedCustomers := []*domain.Customer{
		{Email: "a@example.com", FirstName: "A", LastName: "B", Phone: "123"},
		{Email: "c@example.com", FirstName: "C", LastName: "D", Phone: "456"},
	}

	mockRepo.On("FetchListCustomers", mock.Anything).Return(expectedCustomers, nil)

	customers, err := usecase.FetchListCustomers(context.Background())

	assert.NoError(t, err)
	assert.Len(t, customers, 2)
	assert.Equal(t, expectedCustomers, customers)
	mockRepo.AssertExpectations(t)
}

func TestCustomer_FetchListCustomers_RepoError(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	mockRepo.On("FetchListCustomers", mock.Anything).Return([]*domain.Customer{}, errors.New("db error"))

	customers, err := usecase.FetchListCustomers(context.Background())

	assert.Error(t, err)
	assert.Empty(t, customers)
	mockRepo.AssertExpectations(t)
}

func TestCustomer_FetchCustomerById_Success(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	expectedCustomer := &domain.Customer{
		Id:        &id,
		Email:     "test@example.com",
		FirstName: "John",
	}

	mockRepo.On("FetchCustomerById", mock.Anything, &id).Return(expectedCustomer, nil)

	customer, err := usecase.FetchCustomerById(context.Background(), &id)

	assert.NoError(t, err)
	assert.Equal(t, expectedCustomer, customer)
	mockRepo.AssertExpectations(t)
}

func TestCustomer_FetchCustomerById_NotFound(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	id := bson.NewObjectID()

	mockRepo.On("FetchCustomerById", mock.Anything, &id).Return((*domain.Customer)(nil), errors.New("not found"))

	customer, err := usecase.FetchCustomerById(context.Background(), &id)

	assert.Error(t, err)
	assert.Nil(t, customer)
	mockRepo.AssertExpectations(t)
}

func TestCustomer_UpdateCustomerById_Success(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	updatedCustomer := &domain.Customer{
		Email:     "updated@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
		Phone:     "0987654321",
	}

	mockRepo.On("UpdateCustomerById", mock.Anything, &id, mock.AnythingOfType("*domain.Customer")).Return(nil)
	mockRepo.On("FetchCustomerById", mock.Anything, &id).Return(&domain.Customer{
		Id:        &id,
		Email:     "updated@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
		Phone:     "0987654321",
	}, nil)

	result, err := usecase.UpdateCustomerById(context.Background(), &id, updatedCustomer)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id, *result.Id)
	mockRepo.AssertExpectations(t)
}

func TestCustomer_UpdateCustomerById_RepoError(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	updatedCustomer := &domain.Customer{
		Email: "updated@example.com",
	}

	mockRepo.On("UpdateCustomerById", mock.Anything, &id, mock.AnythingOfType("*domain.Customer")).Return(errors.New("update failed"))

	result, err := usecase.UpdateCustomerById(context.Background(), &id, updatedCustomer)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCustomer_DeleteCustomerById_Success(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	customer := &domain.Customer{
		Email:     "test@example.com",
		FirstName: "John",
	}

	mockRepo.On("DeleteCustomerById", mock.Anything, &id, mock.AnythingOfType("*domain.Customer")).Return(nil)

	err := usecase.DeleteCustomerById(context.Background(), &id, customer)

	assert.NoError(t, err)
	assert.Equal(t, domain.CUSTOMER_STATUS_INACTIVE, customer.Status)
	assert.NotNil(t, customer.DeletedAt)
	mockRepo.AssertExpectations(t)
}

func TestCustomer_DeleteCustomerById_RepoError(t *testing.T) {
	mockRepo := new(mockCustomerRepository)
	usecase := services.NewCustomerUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	customer := &domain.Customer{
		Email: "test@example.com",
	}

	mockRepo.On("DeleteCustomerById", mock.Anything, &id, mock.AnythingOfType("*domain.Customer")).Return(errors.New("delete failed"))

	err := usecase.DeleteCustomerById(context.Background(), &id, customer)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
