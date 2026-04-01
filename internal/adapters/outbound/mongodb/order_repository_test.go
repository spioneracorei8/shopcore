package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"shopcore/internal/adapters/outbound/mongodb"
	"shopcore/internal/core/domain"
)

func TestOrder_CreateOrder(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewOrderRepoImpl(db)

	t.Run("success creates order", func(t *testing.T) {
		customerID := bson.NewObjectID()
		order := &domain.Order{
			CustomerId:     &customerID,
			OrderNo:        "ORD-001",
			Status:         domain.ORDER_STATUS_PENDING,
			Subtotal:       100.0,
			DiscountAmount: 10.0,
			ShippingFee:    5.0,
			TotalAmount:    95.0,
			OrderItems: []*domain.OrderItems{
				{
					ProductId:           &bson.ObjectID{},
					SkuSnapshot:         "SKU001",
					ProductNameSnapshot: "Product A",
					UnitPrice:           50.0,
					Qty:                 2,
					LineTotal:           100.0,
				},
			},
		}
		order.GenObjectID()

		err := repo.CreateOrder(context.Background(), order)
		require.NoError(t, err)
		require.NotNil(t, order.Id)

		fetched, err := repo.FetchOrderById(context.Background(), order.Id)
		require.NoError(t, err)
		assert.Equal(t, order.OrderNo, fetched.OrderNo)
		assert.Equal(t, order.Status, fetched.Status)
		assert.Len(t, fetched.OrderItems, 1)
	})
}

func TestOrder_FetchListOrders(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewOrderRepoImpl(db)

	t.Run("returns all orders", func(t *testing.T) {
		for _, o := range []*domain.Order{
			{OrderNo: "ORD-001", Status: domain.ORDER_STATUS_PENDING, Subtotal: 100.0, DiscountAmount: 0, ShippingFee: 0, TotalAmount: 100.0},
			{OrderNo: "ORD-002", Status: domain.ORDER_STATUS_PAID, Subtotal: 200.0, DiscountAmount: 0, ShippingFee: 0, TotalAmount: 200.0},
		} {
			o.GenObjectID()
			require.NoError(t, repo.CreateOrder(context.Background(), o))
		}

		orders, err := repo.FetchListOrders(context.Background())
		require.NoError(t, err)
		assert.Len(t, orders, 2)
	})
}

func TestOrder_FetchOrderById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewOrderRepoImpl(db)

	t.Run("returns order by id", func(t *testing.T) {
		order := &domain.Order{
			OrderNo:        "ORD-003",
			Status:         domain.ORDER_STATUS_PENDING,
			Subtotal:       150.0,
			DiscountAmount: 0,
			ShippingFee:    0,
			TotalAmount:    150.0,
		}
		order.GenObjectID()
		require.NoError(t, repo.CreateOrder(context.Background(), order))

		found, err := repo.FetchOrderById(context.Background(), order.Id)
		require.NoError(t, err)
		assert.Equal(t, order.Id, found.Id)
		assert.Equal(t, order.OrderNo, found.OrderNo)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		fakeID := bson.NewObjectID()
		_, err := repo.FetchOrderById(context.Background(), &fakeID)
		assert.Error(t, err)
	})
}

func TestOrder_UpdateOrderById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewOrderRepoImpl(db)

	t.Run("updates order fields", func(t *testing.T) {
		order := &domain.Order{
			OrderNo:        "ORD-004",
			Status:         domain.ORDER_STATUS_PENDING,
			Subtotal:       100.0,
			DiscountAmount: 0,
			ShippingFee:    0,
			TotalAmount:    100.0,
		}
		order.GenObjectID()
		require.NoError(t, repo.CreateOrder(context.Background(), order))

		order.Status = domain.ORDER_STATUS_PAID
		order.SetUpdatedAt()

		err := repo.UpdateOrderById(context.Background(), order.Id, order)
		require.NoError(t, err)

		updated, err := repo.FetchOrderById(context.Background(), order.Id)
		require.NoError(t, err)
		assert.Equal(t, domain.ORDER_STATUS_PAID, updated.Status)
	})
}

func TestOrder_DeleteOrderById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewOrderRepoImpl(db)

	t.Run("soft deletes order", func(t *testing.T) {
		order := &domain.Order{
			OrderNo:        "ORD-005",
			Status:         domain.ORDER_STATUS_PENDING,
			Subtotal:       100.0,
			DiscountAmount: 0,
			ShippingFee:    0,
			TotalAmount:    100.0,
		}
		order.GenObjectID()
		require.NoError(t, repo.CreateOrder(context.Background(), order))

		order.SetDeletedAt()
		order.Status = domain.ORDER_STATUS_CANCELLED

		err := repo.DeleteOrderById(context.Background(), order.Id, order)
		require.NoError(t, err)

		fetched, err := repo.FetchOrderById(context.Background(), order.Id)
		require.NoError(t, err)
		assert.Equal(t, domain.ORDER_STATUS_CANCELLED, fetched.Status)
		assert.NotNil(t, fetched.DeletedAt)
	})
}
