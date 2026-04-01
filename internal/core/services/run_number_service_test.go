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

type mockRunNumberRepository struct {
	mock.Mock
}

func (m *mockRunNumberRepository) CreateRunNumber(ctx context.Context, rn *domain.RunNumber) error {
	args := m.Called(ctx, rn)
	return args.Error(0)
}

func (m *mockRunNumberRepository) FetchRunNumber(ctx context.Context) (*domain.RunNumber, error) {
	args := m.Called(ctx)
	return args.Get(0).(*domain.RunNumber), args.Error(1)
}

func (m *mockRunNumberRepository) UpdateRunNumber(ctx context.Context, rn *domain.RunNumber) error {
	args := m.Called(ctx, rn)
	return args.Error(0)
}

var _ outbound.RunNumberRepository = (*mockRunNumberRepository)(nil)

func TestRunNumber_CreateRunNumber_Success(t *testing.T) {
	mockRepo := new(mockRunNumberRepository)
	usecase := services.NewRunNumberUsecaseImpl(mockRepo)

	rn := &domain.RunNumber{
		Prefix:  "ORD",
		Running: 1,
	}

	mockRepo.On("CreateRunNumber", mock.Anything, rn).Return(nil)

	err := usecase.CreateRunNumber(context.Background(), rn)

	assert.NoError(t, err)
	assert.NotNil(t, rn.Id)
	mockRepo.AssertExpectations(t)
}

func TestRunNumber_CreateRunNumber_RepoError(t *testing.T) {
	mockRepo := new(mockRunNumberRepository)
	usecase := services.NewRunNumberUsecaseImpl(mockRepo)

	rn := &domain.RunNumber{
		Prefix:  "ORD",
		Running: 1,
	}

	mockRepo.On("CreateRunNumber", mock.Anything, rn).Return(errors.New("db error"))

	err := usecase.CreateRunNumber(context.Background(), rn)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRunNumber_FetchRunNumber_Success(t *testing.T) {
	mockRepo := new(mockRunNumberRepository)
	usecase := services.NewRunNumberUsecaseImpl(mockRepo)

	expectedRn := &domain.RunNumber{
		Prefix:  "ORD",
		Running: 42,
	}

	mockRepo.On("FetchRunNumber", mock.Anything).Return(expectedRn, nil)

	rn, err := usecase.FetchRunNumber(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedRn, rn)
	mockRepo.AssertExpectations(t)
}

func TestRunNumber_FetchRunNumber_NotFound(t *testing.T) {
	mockRepo := new(mockRunNumberRepository)
	usecase := services.NewRunNumberUsecaseImpl(mockRepo)

	mockRepo.On("FetchRunNumber", mock.Anything).Return((*domain.RunNumber)(nil), errors.New("not found"))

	rn, err := usecase.FetchRunNumber(context.Background())

	assert.Error(t, err)
	assert.Nil(t, rn)
	mockRepo.AssertExpectations(t)
}

func TestRunNumber_UpdateRunNumber_Success(t *testing.T) {
	mockRepo := new(mockRunNumberRepository)
	usecase := services.NewRunNumberUsecaseImpl(mockRepo)

	existingRn := &domain.RunNumber{
		Id:      ptrObjectID(bson.NewObjectID()),
		Prefix:  "ORD",
		Running: 42,
	}

	mockRepo.On("FetchRunNumber", mock.Anything).Return(existingRn, nil)
	mockRepo.On("UpdateRunNumber", mock.Anything, mock.MatchedBy(func(rn *domain.RunNumber) bool {
		return rn.Running == 43
	})).Return(nil)

	err := usecase.UpdateRunNumber(context.Background())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRunNumber_UpdateRunNumber_FetchError(t *testing.T) {
	mockRepo := new(mockRunNumberRepository)
	usecase := services.NewRunNumberUsecaseImpl(mockRepo)

	mockRepo.On("FetchRunNumber", mock.Anything).Return((*domain.RunNumber)(nil), errors.New("fetch failed"))

	err := usecase.UpdateRunNumber(context.Background())

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRunNumber_UpdateRunNumber_UpdateError(t *testing.T) {
	mockRepo := new(mockRunNumberRepository)
	usecase := services.NewRunNumberUsecaseImpl(mockRepo)

	existingRn := &domain.RunNumber{
		Id:      ptrObjectID(bson.NewObjectID()),
		Prefix:  "ORD",
		Running: 42,
	}

	mockRepo.On("FetchRunNumber", mock.Anything).Return(existingRn, nil)
	mockRepo.On("UpdateRunNumber", mock.Anything, mock.AnythingOfType("*domain.RunNumber")).Return(errors.New("update failed"))

	err := usecase.UpdateRunNumber(context.Background())

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func ptrObjectID(id bson.ObjectID) *bson.ObjectID {
	return &id
}
