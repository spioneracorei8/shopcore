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

func TestProduct_CreateProduct(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewProductRepoImpl(db)

	t.Run("success creates product", func(t *testing.T) {
		product := &domain.Product{
			Sku:        "SKU001",
			Name:       "Test Product",
			Descrption: "A test product",
			Price:      99.99,
			StockQty:   100,
		}
		product.GenObjectID()

		err := repo.CreateProduct(context.Background(), product)
		require.NoError(t, err)
		require.NotNil(t, product.Id)

		fetched, err := repo.FetchProductById(context.Background(), product.Id)
		require.NoError(t, err)
		assert.Equal(t, product.Sku, fetched.Sku)
		assert.Equal(t, product.Name, fetched.Name)
		assert.Equal(t, product.Price, fetched.Price)
	})
}

func TestProduct_FetchListProducts(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewProductRepoImpl(db)

	t.Run("returns all non-deleted products", func(t *testing.T) {
		for _, p := range []*domain.Product{
			{Sku: "SKU001", Name: "Product A", Price: 10.0, StockQty: 50},
			{Sku: "SKU002", Name: "Product B", Price: 20.0, StockQty: 30},
		} {
			p.GenObjectID()
			require.NoError(t, repo.CreateProduct(context.Background(), p))
		}

		products, err := repo.FetchListProducts(context.Background())
		require.NoError(t, err)
		assert.Len(t, products, 2)
	})

	t.Run("excludes deleted products", func(t *testing.T) {
		active := &domain.Product{Sku: "SKU003", Name: "Active", Price: 15.0, StockQty: 10}
		active.GenObjectID()
		require.NoError(t, repo.CreateProduct(context.Background(), active))

		deleted := &domain.Product{Sku: "SKU004", Name: "Deleted", Price: 25.0, StockQty: 5}
		deleted.GenObjectID()
		deleted.SetDeletedAt()
		require.NoError(t, repo.CreateProduct(context.Background(), deleted))

		products, err := repo.FetchListProducts(context.Background())
		require.NoError(t, err)

		for _, p := range products {
			assert.Nil(t, p.DeletedAt)
		}
	})
}

func TestProduct_FetchProductById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewProductRepoImpl(db)

	t.Run("returns product by id", func(t *testing.T) {
		product := &domain.Product{Sku: "SKU005", Name: "Find Me", Price: 50.0, StockQty: 20}
		product.GenObjectID()
		require.NoError(t, repo.CreateProduct(context.Background(), product))

		found, err := repo.FetchProductById(context.Background(), product.Id)
		require.NoError(t, err)
		assert.Equal(t, product.Id, found.Id)
		assert.Equal(t, product.Sku, found.Sku)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		fakeID := bson.NewObjectID()
		_, err := repo.FetchProductById(context.Background(), &fakeID)
		assert.Error(t, err)
	})
}

func TestProduct_UpdateProductById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewProductRepoImpl(db)

	t.Run("updates product fields", func(t *testing.T) {
		product := &domain.Product{Sku: "SKU006", Name: "Old Name", Price: 10.0, StockQty: 5}
		product.GenObjectID()
		require.NoError(t, repo.CreateProduct(context.Background(), product))

		product.Name = "New Name"
		product.Price = 199.99
		product.SetUpdatedAt()

		err := repo.UpdateProductById(context.Background(), product.Id, product)
		require.NoError(t, err)

		updated, err := repo.FetchProductById(context.Background(), product.Id)
		require.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, 199.99, updated.Price)
	})
}

func TestProduct_DeleteProductById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewProductRepoImpl(db)

	t.Run("soft deletes product", func(t *testing.T) {
		product := &domain.Product{Sku: "SKU007", Name: "Delete Me", Price: 30.0, StockQty: 10}
		product.GenObjectID()
		require.NoError(t, repo.CreateProduct(context.Background(), product))

		product.SetDeletedAt()
		product.Status = domain.PRODUCT_STATUS_INACTIVE

		err := repo.DeleteProductById(context.Background(), product.Id, product)
		require.NoError(t, err)

		fetched, err := repo.FetchProductById(context.Background(), product.Id)
		require.NoError(t, err)
		assert.Equal(t, domain.PRODUCT_STATUS_INACTIVE, fetched.Status)
		assert.NotNil(t, fetched.DeletedAt)
	})
}
