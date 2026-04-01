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

func TestCustomer_CreateCustomer(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewCustomerRepoImpl(db)

	t.Run("success creates customer", func(t *testing.T) {
		customer := &domain.Customer{
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
		}
		customer.GenObjectID()

		err := repo.CreateCustomer(context.Background(), customer)
		require.NoError(t, err)
		require.NotNil(t, customer.Id)

		fetched, err := repo.FetchCustomerById(context.Background(), customer.Id)
		require.NoError(t, err)
		assert.Equal(t, customer.Email, fetched.Email)
		assert.Equal(t, customer.FirstName, fetched.FirstName)
	})
}

func TestCustomer_FetchListCustomers(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewCustomerRepoImpl(db)

	t.Run("returns all non-deleted customers", func(t *testing.T) {
		for _, c := range []*domain.Customer{
			{Email: "a@test.com", FirstName: "A", LastName: "B", Phone: "111"},
			{Email: "c@test.com", FirstName: "C", LastName: "D", Phone: "222"},
		} {
			c.GenObjectID()
			require.NoError(t, repo.CreateCustomer(context.Background(), c))
		}

		customers, err := repo.FetchListCustomers(context.Background())
		require.NoError(t, err)
		assert.Len(t, customers, 2)
	})

	t.Run("excludes deleted customers", func(t *testing.T) {
		active := &domain.Customer{Email: "active@test.com", FirstName: "Active", LastName: "User", Phone: "333"}
		active.GenObjectID()
		require.NoError(t, repo.CreateCustomer(context.Background(), active))

		deleted := &domain.Customer{Email: "deleted@test.com", FirstName: "Deleted", LastName: "User", Phone: "444"}
		deleted.GenObjectID()
		deleted.SetDeletedAt()
		require.NoError(t, repo.CreateCustomer(context.Background(), deleted))

		customers, err := repo.FetchListCustomers(context.Background())
		require.NoError(t, err)

		for _, c := range customers {
			assert.Nil(t, c.DeletedAt)
		}
	})
}

func TestCustomer_FetchCustomerById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewCustomerRepoImpl(db)

	t.Run("returns customer by id", func(t *testing.T) {
		customer := &domain.Customer{Email: "find@test.com", FirstName: "Find", LastName: "Me", Phone: "555"}
		customer.GenObjectID()
		require.NoError(t, repo.CreateCustomer(context.Background(), customer))

		found, err := repo.FetchCustomerById(context.Background(), customer.Id)
		require.NoError(t, err)
		assert.Equal(t, customer.Id, found.Id)
		assert.Equal(t, customer.Email, found.Email)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		fakeID := bson.NewObjectID()
		_, err := repo.FetchCustomerById(context.Background(), &fakeID)
		assert.Error(t, err)
	})
}

func TestCustomer_UpdateCustomerById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewCustomerRepoImpl(db)

	t.Run("updates customer fields", func(t *testing.T) {
		customer := &domain.Customer{Email: "update@test.com", FirstName: "Old", LastName: "Name", Phone: "666"}
		customer.GenObjectID()
		require.NoError(t, repo.CreateCustomer(context.Background(), customer))

		customer.FirstName = "New"
		customer.SetUpdatedAt()

		err := repo.UpdateCustomerById(context.Background(), customer.Id, customer)
		require.NoError(t, err)

		updated, err := repo.FetchCustomerById(context.Background(), customer.Id)
		require.NoError(t, err)
		assert.Equal(t, "New", updated.FirstName)
	})
}

func TestCustomer_DeleteCustomerById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewCustomerRepoImpl(db)

	t.Run("soft deletes customer", func(t *testing.T) {
		customer := &domain.Customer{Email: "delete@test.com", FirstName: "Delete", LastName: "Me", Phone: "777"}
		customer.GenObjectID()
		require.NoError(t, repo.CreateCustomer(context.Background(), customer))

		customer.SetDeletedAt()
		customer.Status = domain.CUSTOMER_STATUS_INACTIVE

		err := repo.DeleteCustomerById(context.Background(), customer.Id, customer)
		require.NoError(t, err)

		fetched, err := repo.FetchCustomerById(context.Background(), customer.Id)
		require.NoError(t, err)
		assert.Equal(t, domain.CUSTOMER_STATUS_INACTIVE, fetched.Status)
		assert.NotNil(t, fetched.DeletedAt)
	})
}
