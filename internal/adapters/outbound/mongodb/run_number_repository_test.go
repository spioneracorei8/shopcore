package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"shopcore/internal/adapters/outbound/mongodb"
	"shopcore/internal/core/domain"
)

func TestRunNumber_CreateRunNumber(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewRunNumberRepoImpl(db)

	t.Run("success creates run number", func(t *testing.T) {
		rn := &domain.RunNumber{
			Prefix:  "ORD",
			Running: 1,
		}
		rn.GenObjectID()

		err := repo.CreateRunNumber(context.Background(), rn)
		require.NoError(t, err)
		require.NotNil(t, rn.Id)

		fetched, err := repo.FetchRunNumber(context.Background())
		require.NoError(t, err)
		assert.Equal(t, rn.Prefix, fetched.Prefix)
		assert.Equal(t, rn.Running, fetched.Running)
	})
}

func TestRunNumber_FetchRunNumber(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewRunNumberRepoImpl(db)

	t.Run("returns run number", func(t *testing.T) {
		rn := &domain.RunNumber{Prefix: "INV", Running: 42}
		rn.GenObjectID()
		require.NoError(t, repo.CreateRunNumber(context.Background(), rn))

		fetched, err := repo.FetchRunNumber(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "INV", fetched.Prefix)
		assert.Equal(t, 42, fetched.Running)
	})

	t.Run("returns first document when multiple exist", func(t *testing.T) {
		rn2 := &domain.RunNumber{Prefix: "ORD", Running: 999}
		rn2.GenObjectID()
		require.NoError(t, repo.CreateRunNumber(context.Background(), rn2))

		fetched, err := repo.FetchRunNumber(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "INV", fetched.Prefix)
		assert.Equal(t, 42, fetched.Running)
	})
}

func TestRunNumber_UpdateRunNumber(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewRunNumberRepoImpl(db)

	t.Run("increments running number", func(t *testing.T) {
		rn := &domain.RunNumber{Prefix: "ORD", Running: 100}
		rn.GenObjectID()
		require.NoError(t, repo.CreateRunNumber(context.Background(), rn))

		rn.Running++
		err := repo.UpdateRunNumber(context.Background(), rn)
		require.NoError(t, err)

		fetched, err := repo.FetchRunNumber(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 101, fetched.Running)
	})
}
