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

type mockProductRepository struct {
	mock.Mock
}

func (m *mockProductRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *mockProductRepository) FetchListProducts(ctx context.Context) ([]*domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *mockProductRepository) FetchProductById(ctx context.Context, id *bson.ObjectID) (*domain.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *mockProductRepository) UpdateProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) error {
	args := m.Called(ctx, id, product)
	return args.Error(0)
}

func (m *mockProductRepository) DeleteProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) error {
	args := m.Called(ctx, id, product)
	return args.Error(0)
}

var _ outbound.ProductRepository = (*mockProductRepository)(nil)

func TestProduct_CreateProduct_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	product := &domain.Product{
		Sku:        "SKU001",
		Name:       "Test Product",
		Descrption: "A test product",
		Price:      99.99,
		StockQty:   100,
	}

	mockRepo.On("CreateProduct", mock.Anything, product).Return(nil)

	err := usecase.CreateProduct(context.Background(), product)

	assert.NoError(t, err)
	assert.NotNil(t, product.Id)
	assert.False(t, product.CreatedAt.IsZero())
	assert.False(t, product.UpdatedAt.IsZero())
	mockRepo.AssertExpectations(t)
}

func TestProduct_CreateProduct_RepoError(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	product := &domain.Product{
		Sku:      "SKU001",
		Name:     "Test Product",
		Price:    99.99,
		StockQty: 100,
	}

	mockRepo.On("CreateProduct", mock.Anything, product).Return(errors.New("db error"))

	err := usecase.CreateProduct(context.Background(), product)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProduct_FetchListProducts_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	expectedProducts := []*domain.Product{
		{Sku: "SKU001", Name: "Product A", Price: 10.0, StockQty: 50},
		{Sku: "SKU002", Name: "Product B", Price: 20.0, StockQty: 30},
	}

	mockRepo.On("FetchListProducts", mock.Anything).Return(expectedProducts, nil)

	products, err := usecase.FetchListProducts(context.Background())

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestProduct_FetchListProducts_RepoError(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	mockRepo.On("FetchListProducts", mock.Anything).Return([]*domain.Product{}, errors.New("db error"))

	products, err := usecase.FetchListProducts(context.Background())

	assert.Error(t, err)
	assert.Empty(t, products)
	mockRepo.AssertExpectations(t)
}

func TestProduct_FetchProductById_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	expectedProduct := &domain.Product{
		Id:   &id,
		Sku:  "SKU001",
		Name: "Test Product",
	}

	mockRepo.On("FetchProductById", mock.Anything, &id).Return(expectedProduct, nil)

	product, err := usecase.FetchProductById(context.Background(), &id)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestProduct_FetchProductById_NotFound(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	id := bson.NewObjectID()

	mockRepo.On("FetchProductById", mock.Anything, &id).Return((*domain.Product)(nil), errors.New("not found"))

	product, err := usecase.FetchProductById(context.Background(), &id)

	assert.Error(t, err)
	assert.Nil(t, product)
	mockRepo.AssertExpectations(t)
}

func TestProduct_UpdateProductById_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	updatedProduct := &domain.Product{
		Sku:      "SKU001",
		Name:     "Updated Product",
		Price:    149.99,
		StockQty: 80,
	}

	mockRepo.On("UpdateProductById", mock.Anything, &id, mock.AnythingOfType("*domain.Product")).Return(nil)
	mockRepo.On("FetchProductById", mock.Anything, &id).Return(&domain.Product{
		Id:       &id,
		Sku:      "SKU001",
		Name:     "Updated Product",
		Price:    149.99,
		StockQty: 80,
	}, nil)

	result, err := usecase.UpdateProductById(context.Background(), &id, updatedProduct)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id, *result.Id)
	mockRepo.AssertExpectations(t)
}

func TestProduct_UpdateProductById_RepoError(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	updatedProduct := &domain.Product{
		Sku:   "SKU001",
		Name:  "Updated Product",
		Price: 149.99,
	}

	mockRepo.On("UpdateProductById", mock.Anything, &id, mock.AnythingOfType("*domain.Product")).Return(errors.New("update failed"))

	result, err := usecase.UpdateProductById(context.Background(), &id, updatedProduct)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProduct_DeleteProductById_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	product := &domain.Product{
		Sku:  "SKU001",
		Name: "Test Product",
	}

	mockRepo.On("DeleteProductById", mock.Anything, &id, mock.AnythingOfType("*domain.Product")).Return(nil)

	err := usecase.DeleteProductById(context.Background(), &id, product)

	assert.NoError(t, err)
	assert.Equal(t, domain.PRODUCT_STATUS_INACTIVE, product.Status)
	assert.NotNil(t, product.DeletedAt)
	mockRepo.AssertExpectations(t)
}

func TestProduct_DeleteProductById_RepoError(t *testing.T) {
	mockRepo := new(mockProductRepository)
	usecase := services.NewProductUsecaseImpl(mockRepo)

	id := bson.NewObjectID()
	product := &domain.Product{
		Sku:  "SKU001",
		Name: "Test Product",
	}

	mockRepo.On("DeleteProductById", mock.Anything, &id, mock.AnythingOfType("*domain.Product")).Return(errors.New("delete failed"))

	err := usecase.DeleteProductById(context.Background(), &id, product)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
