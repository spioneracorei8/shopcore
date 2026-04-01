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

	mongoContainer, err := mongodb.Run(ctx, "mongo:7")
	require.NoError(t, err)

	uri, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	require.NoError(t, err)

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	dbName := t.Name()
	db := client.Database(dbName)

	cleanup := func() {
		client.Disconnect(ctx)
		mongoContainer.Terminate(ctx)
	}

	return db, cleanup
}
