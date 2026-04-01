# Testing Strategy

ShopCore separates tests into two layers aligned with hexagonal architecture:

| Layer | Type | Tool | Purpose |
|-------|------|------|---------|
| **Usecase/Service** | Unit Test | `stretchr/testify/mock` | Verify business logic in isolation |
| **Repository/Adapter** | Integration Test | `testcontainers-go` | Verify real MongoDB queries against real database |

## Why Not Mock MongoDB in Repository Tests?

Mocking MongoDB queries (`FindOne`, `InsertOne`, `FindOneAndUpdate`) gives false confidence:

```
Mock says: "update succeeded"
Real MongoDB says: "schema validation failed" / "duplicate key" / "write concern timeout"
```

Integration tests with testcontainers-go run against **real MongoDB**, catching issues that mocks never will:
- Actual BSON serialization behavior
- Real query filter logic (`$set`, `$ne`, `$exists`)
- MongoDB driver version compatibility
- Write concern and read preference behavior
- Index constraints and unique key violations

---

## Usecase/Service Layer — Unit Tests

### Location

```
internal/core/services/
├── customer_service_test.go
├── product_service_test.go
├── order_service_test.go
└── run_number_service_test.go
```

### Pattern: Mock Repository via testify/mock

Each test file defines a mock implementation of the **outbound port interface**:

```go
package services_test

import (
    "context"
    "shopcore/internal/core/domain"
    "shopcore/internal/core/ports/outbound"
    "github.com/stretchr/testify/mock"
    "go.mongodb.org/mongo-driver/v2/bson"
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

// Compile-time interface verification
var _ outbound.CustomerRepository = (*mockCustomerRepository)(nil)
```

### Pattern: Test Case Structure

```go
func TestCustomerUsecase_CreateCustomer_Success(t *testing.T) {
    mockRepo := new(mockCustomerRepository)
    usecase := services.NewCustomerUsecaseImpl(mockRepo)

    customer := &domain.Customer{
        Email:     "test@example.com",
        FirstName: "John",
        LastName:  "Doe",
        Phone:     "1234567890",
    }

    // Arrange: set mock expectations
    mockRepo.On("CreateCustomer", mock.Anything, customer).Return(nil)

    // Act: execute usecase
    err := usecase.CreateCustomer(context.Background(), customer)

    // Assert: verify behavior
    assert.NoError(t, err)
    assert.NotNil(t, customer.Id)
    assert.Equal(t, domain.CUSTOMER_STATUS_ACTIVE, customer.Status)
    assert.False(t, customer.CreatedAt.IsZero())
    assert.False(t, customer.UpdatedAt.IsZero())
    mockRepo.AssertExpectations(t)
}
```

### What Unit Tests Verify

| Test | Verifies |
|------|----------|
| `CreateCustomer_Success` | ID generation, status set to ACTIVE, timestamps set, repo called |
| `CreateCustomer_RepoError` | Error propagation from repository |
| `FetchListCustomers_Success` | Data passthrough from repository |
| `UpdateCustomerById_Success` | Timestamp updated, repo called, result returned |
| `DeleteCustomerById_Success` | Status set to INACTIVE, soft delete timestamp, repo called |

### Running Unit Tests

```bash
go test ./internal/core/services/... -v

# Output:
# === RUN   TestCustomerUsecase_CreateCustomer_Success
# --- PASS: TestCustomerUsecase_CreateCustomer_Success (0.00s)
# === RUN   TestProductUsecase_CreateProduct_Success
# --- PASS: TestProductUsecase_CreateProduct_Success (0.00s)
# ...
# ok  	shopcore/internal/core/services	0.003s
```

No Docker required. No MongoDB connection. Fast feedback.

---

## Repository Layer — Integration Tests

### What is testcontainers-go?

[testcontainers-go](https://github.com/testcontainers/testcontainers-go) is a Go library for running throwaway Docker containers during automated tests. It provides a programmatic API to start, configure, and clean up containerized services like MongoDB, PostgreSQL, Redis, etc.

Unlike mocking a database driver, testcontainers gives you:

| Mock Driver | testcontainers |
|-------------|----------------|
| Simulated query behavior | Real MongoDB engine |
| Fake return values | Actual BSON serialization |
| Never catches driver bugs | Catches driver/API mismatches |
| No network I/O | Real TCP connections |
| Assumes query is correct | Validates query against real DB |

### Location

```
internal/adapters/outbound/mongodb/
├── customer_repository_test.go
├── product_repository_test.go
├── order_repository_test.go
├── run_number_repository_test.go
└── repository_test_helper_test.go
```

### Pattern: testcontainers-go Setup

Shared helper spins up a real MongoDB container per test function:

```go
package mongodb_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go/modules/mongodb"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
    t.Helper()
    ctx := context.Background()

    // Start MongoDB 7 container
    mongoContainer, err := mongodb.Run(ctx, "mongo:7")
    require.NoError(t, err)

    // Get connection string
    uri, err := mongoContainer.ConnectionString(ctx)
    require.NoError(t, err)

    // Connect real MongoDB client
    client, err := mongo.Connect(options.Client().ApplyURI(uri))
    require.NoError(t, err)

    err = client.Ping(ctx, nil)
    require.NoError(t, err)

    // Isolated database per test function
    dbName := t.Name()
    db := client.Database(dbName)

    cleanup := func() {
        client.Disconnect(ctx)
        mongoContainer.Terminate(ctx)
    }

    return db, cleanup
}
```

### How Isolation Works

Each test function gets:
- Its own MongoDB container
- A database named after the test function

```
TestCustomerRepository_CreateCustomer/
  → database: "TestCustomerRepository_CreateCustomer"

TestCustomerRepository_FetchListCustomers/
  → database: "TestCustomerRepository_FetchListCustomers"
```

No data leaks between tests. No shared state. No cleanup scripts.

### Pattern: Test Case Structure

```go
func TestCustomerRepository_CreateCustomer(t *testing.T) {
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

        // Act: real MongoDB InsertOne
        err := repo.CreateCustomer(context.Background(), customer)
        require.NoError(t, err)
        require.NotNil(t, customer.Id)

        // Assert: real MongoDB FindOne
        fetched, err := repo.FetchCustomerById(context.Background(), customer.Id)
        require.NoError(t, err)
        assert.Equal(t, customer.Email, fetched.Email)
        assert.Equal(t, customer.FirstName, fetched.FirstName)
    })
}
```

### What Integration Tests Verify

| Test | Real MongoDB Operation |
|------|----------------------|
| `CreateCustomer` | `Collection.InsertOne` with BSON serialization |
| `FetchListCustomers` | `Collection.Find` with filter `{deletedAt: nil}` |
| `FetchCustomerById` | `Collection.FindOne` with `{_id: id}` filter |
| `UpdateCustomerById` | `Collection.FindOneAndUpdate` with `$set` operator |
| `DeleteCustomerById` | `Collection.UpdateOne` with `$set` for soft delete |
| `FetchListCustomers/excludes_deleted` | Verifies `deletedAt` filter actually works |

### Running Integration Tests

```bash
go test ./internal/adapters/outbound/mongodb/... -v

# Output:
# === RUN   TestCustomerRepository_CreateCustomer
# 🐳 Creating container for image mongo:7
# ✅ Container started
# === RUN   TestCustomerRepository_CreateCustomer/success_creates_customer
# --- PASS: TestCustomerRepository_CreateCustomer/success_creates_customer (0.01s)
# ...
# ok  	shopcore/internal/adapters/outbound/mongodb	11.143s
```

Requires Docker. Slower but verifies production behavior.

---

## Test Coverage Summary

```
Total: 61 tests

Usecase Layer (Unit Tests):     39 tests
├── CustomerUsecase:            10 tests
├── ProductUsecase:             10 tests
├── OrderUsecase:               12 tests
└── RunNumberUsecase:            7 tests

Repository Layer (Integration): 22 tests
├── CustomerRepository:          6 tests
├── ProductRepository:           7 tests
├── OrderRepository:             6 tests
└── RunNumberRepository:         3 tests
```

## Running All Tests

```bash
# Everything (requires Docker)
go test ./...

# Only unit tests (no Docker)
go test ./internal/core/services/... -v

# Only integration tests (requires Docker)
go test ./internal/adapters/outbound/mongodb/... -v
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/stretchr/testify` | Assertions (`assert`) and mocking (`mock`) |
| `github.com/testcontainers/testcontainers-go/modules/mongodb` | MongoDB container lifecycle |
| `github.com/testcontainers/testcontainers-go` | Core testcontainers runtime |
