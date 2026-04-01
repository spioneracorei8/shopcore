package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/inbound"
	"shopcore/internal/core/ports/outbound"
	"shopcore/internal/core/services"
)

type mockOrderRepository struct {
	mock.Mock
}

func (m *mockOrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *mockOrderRepository) FetchListOrders(ctx context.Context) ([]*domain.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *mockOrderRepository) FetchOrderById(ctx context.Context, id *bson.ObjectID) (*domain.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *mockOrderRepository) UpdateOrderById(ctx context.Context, id *bson.ObjectID, order *domain.Order) error {
	args := m.Called(ctx, id, order)
	return args.Error(0)
}

func (m *mockOrderRepository) DeleteOrderById(ctx context.Context, id *bson.ObjectID, order *domain.Order) error {
	args := m.Called(ctx, id, order)
	return args.Error(0)
}

var _ outbound.OrderRepository = (*mockOrderRepository)(nil)

type mockRunNumberUsecase struct {
	mock.Mock
}

func (m *mockRunNumberUsecase) CreateRunNumber(ctx context.Context, rn *domain.RunNumber) error {
	args := m.Called(ctx, rn)
	return args.Error(0)
}

func (m *mockRunNumberUsecase) FetchRunNumber(ctx context.Context) (*domain.RunNumber, error) {
	args := m.Called(ctx)
	return args.Get(0).(*domain.RunNumber), args.Error(1)
}

func (m *mockRunNumberUsecase) UpdateRunNumber(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

var _ inbound.RunNumberUsecase = (*mockRunNumberUsecase)(nil)

type mockProductUsecase struct {
	mock.Mock
}

func (m *mockProductUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *mockProductUsecase) FetchListProducts(ctx context.Context) ([]*domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *mockProductUsecase) FetchProductById(ctx context.Context, id *bson.ObjectID) (*domain.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *mockProductUsecase) UpdateProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, id, product)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *mockProductUsecase) DeleteProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) error {
	args := m.Called(ctx, id, product)
	return args.Error(0)
}

var _ inbound.ProductUsecase = (*mockProductUsecase)(nil)

func TestOrder_CreateOrder_Success(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	productID := bson.NewObjectID()
	order := &domain.Order{
		CustomerId:     ptrObjectID(bson.NewObjectID()),
		Status:         domain.ORDER_STATUS_PENDING,
		Subtotal:       100.0,
		DiscountAmount: 10.0,
		ShippingFee:    5.0,
		TotalAmount:    95.0,
		OrderItems: []*domain.OrderItems{
			{
				ProductId:           &productID,
				SkuSnapshot:         "SKU001",
				ProductNameSnapshot: "Product A",
				UnitPrice:           50.0,
				Qty:                 2,
				LineTotal:           100.0,
			},
		},
	}

	rn := &domain.RunNumber{
		Id:      ptrObjectID(bson.NewObjectID()),
		Prefix:  "ORD",
		Running: 100,
	}

	product := &domain.Product{
		Id:       &productID,
		Sku:      "SKU001",
		Name:     "Product A",
		StockQty: 50,
	}

	mockRnUs.On("FetchRunNumber", mock.Anything).Return(rn, nil)
	mockRnUs.On("UpdateRunNumber", mock.Anything).Return(nil)
	mockProductUs.On("FetchProductById", mock.Anything, &productID).Return(product, nil)
	mockProductUs.On("UpdateProductById", mock.Anything, &productID, mock.AnythingOfType("*domain.Product")).Return(product, nil)
	mockOrderRepo.On("CreateOrder", mock.Anything, order).Return(nil)

	err := usecase.CreateOrder(context.Background(), order)

	assert.NoError(t, err)
	assert.NotNil(t, order.Id)
	assert.NotEmpty(t, order.OrderNo)
	assert.False(t, order.CreatedAt.IsZero())
	assert.False(t, order.UpdatedAt.IsZero())
	assert.Equal(t, 48, product.StockQty)
	mockOrderRepo.AssertExpectations(t)
	mockRnUs.AssertExpectations(t)
	mockProductUs.AssertExpectations(t)
}

func TestOrder_CreateOrder_FetchRunNumberError(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	order := &domain.Order{
		Status:      domain.ORDER_STATUS_PENDING,
		Subtotal:    100.0,
		TotalAmount: 100.0,
	}

	mockRnUs.On("FetchRunNumber", mock.Anything).Return((*domain.RunNumber)(nil), errors.New("fetch failed"))

	err := usecase.CreateOrder(context.Background(), order)

	assert.Error(t, err)
	mockRnUs.AssertExpectations(t)
}

func TestOrder_CreateOrder_UpdateRunNumberError(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	order := &domain.Order{
		Status:      domain.ORDER_STATUS_PENDING,
		Subtotal:    100.0,
		TotalAmount: 100.0,
	}

	rn := &domain.RunNumber{
		Prefix:  "ORD",
		Running: 100,
	}

	mockRnUs.On("FetchRunNumber", mock.Anything).Return(rn, nil)
	mockRnUs.On("UpdateRunNumber", mock.Anything).Return(errors.New("update failed"))

	err := usecase.CreateOrder(context.Background(), order)

	assert.Error(t, err)
	mockRnUs.AssertExpectations(t)
}

func TestOrder_CreateOrder_FetchProductError(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	productID := bson.NewObjectID()
	order := &domain.Order{
		Status:      domain.ORDER_STATUS_PENDING,
		Subtotal:    100.0,
		TotalAmount: 100.0,
		OrderItems: []*domain.OrderItems{
			{
				ProductId: &productID,
				Qty:       1,
			},
		},
	}

	rn := &domain.RunNumber{
		Prefix:  "ORD",
		Running: 100,
	}

	mockRnUs.On("FetchRunNumber", mock.Anything).Return(rn, nil)
	mockRnUs.On("UpdateRunNumber", mock.Anything).Return(nil)
	mockProductUs.On("FetchProductById", mock.Anything, &productID).Return((*domain.Product)(nil), errors.New("product not found"))

	err := usecase.CreateOrder(context.Background(), order)

	assert.Error(t, err)
	mockProductUs.AssertExpectations(t)
}

func TestOrder_FetchListOrders_Success(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	expectedOrders := []*domain.Order{
		{OrderNo: "ORD-001", Status: domain.ORDER_STATUS_PENDING},
		{OrderNo: "ORD-002", Status: domain.ORDER_STATUS_PAID},
	}

	mockOrderRepo.On("FetchListOrders", mock.Anything).Return(expectedOrders, nil)

	orders, err := usecase.FetchListOrders(context.Background())

	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	assert.Equal(t, expectedOrders, orders)
	mockOrderRepo.AssertExpectations(t)
}

func TestOrder_FetchListOrders_RepoError(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	mockOrderRepo.On("FetchListOrders", mock.Anything).Return([]*domain.Order{}, errors.New("db error"))

	orders, err := usecase.FetchListOrders(context.Background())

	assert.Error(t, err)
	assert.Empty(t, orders)
	mockOrderRepo.AssertExpectations(t)
}

func TestOrder_FetchOrderById_Success(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	id := bson.NewObjectID()
	expectedOrder := &domain.Order{
		Id:       &id,
		OrderNo:  "ORD-001",
		Status:   domain.ORDER_STATUS_PENDING,
		Subtotal: 100.0,
	}

	mockOrderRepo.On("FetchOrderById", mock.Anything, &id).Return(expectedOrder, nil)

	order, err := usecase.FetchOrderById(context.Background(), &id)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	mockOrderRepo.AssertExpectations(t)
}

func TestOrder_FetchOrderById_NotFound(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	id := bson.NewObjectID()

	mockOrderRepo.On("FetchOrderById", mock.Anything, &id).Return((*domain.Order)(nil), errors.New("not found"))

	order, err := usecase.FetchOrderById(context.Background(), &id)

	assert.Error(t, err)
	assert.Nil(t, order)
	mockOrderRepo.AssertExpectations(t)
}

func TestOrder_UpdateOrderById_Success(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	id := bson.NewObjectID()
	updatedOrder := &domain.Order{
		Status:   domain.ORDER_STATUS_PAID,
		Subtotal: 150.0,
	}

	mockOrderRepo.On("UpdateOrderById", mock.Anything, &id, mock.AnythingOfType("*domain.Order")).Return(nil)
	mockOrderRepo.On("FetchOrderById", mock.Anything, &id).Return(&domain.Order{
		Id:       &id,
		OrderNo:  "ORD-001",
		Status:   domain.ORDER_STATUS_PAID,
		Subtotal: 150.0,
	}, nil)

	result, err := usecase.UpdateOrderById(context.Background(), &id, updatedOrder)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id, *result.Id)
	assert.Equal(t, domain.ORDER_STATUS_PAID, result.Status)
	mockOrderRepo.AssertExpectations(t)
}

func TestOrder_UpdateOrderById_RepoError(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	id := bson.NewObjectID()
	updatedOrder := &domain.Order{
		Status: domain.ORDER_STATUS_PAID,
	}

	mockOrderRepo.On("UpdateOrderById", mock.Anything, &id, mock.AnythingOfType("*domain.Order")).Return(errors.New("update failed"))

	result, err := usecase.UpdateOrderById(context.Background(), &id, updatedOrder)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockOrderRepo.AssertExpectations(t)
}

func TestOrder_DeleteOrderById_Success(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	id := bson.NewObjectID()
	order := &domain.Order{
		OrderNo:  "ORD-001",
		Status:   domain.ORDER_STATUS_PENDING,
		Subtotal: 100.0,
	}

	mockOrderRepo.On("DeleteOrderById", mock.Anything, &id, mock.AnythingOfType("*domain.Order")).Return(nil)

	err := usecase.DeleteOrderById(context.Background(), &id, order)

	assert.NoError(t, err)
	assert.Equal(t, domain.ORDER_STATUS_CANCELLED, order.Status)
	assert.NotNil(t, order.DeletedAt)
	mockOrderRepo.AssertExpectations(t)
}

func TestOrder_DeleteOrderById_RepoError(t *testing.T) {
	mockOrderRepo := new(mockOrderRepository)
	mockRnUs := new(mockRunNumberUsecase)
	mockProductUs := new(mockProductUsecase)

	usecase := services.NewOrderUsecaseImpl(mockRnUs, mockProductUs, mockOrderRepo)

	id := bson.NewObjectID()
	order := &domain.Order{
		OrderNo: "ORD-001",
	}

	mockOrderRepo.On("DeleteOrderById", mock.Anything, &id, mock.AnythingOfType("*domain.Order")).Return(errors.New("delete failed"))

	err := usecase.DeleteOrderById(context.Background(), &id, order)

	assert.Error(t, err)
	mockOrderRepo.AssertExpectations(t)
}
