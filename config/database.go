package config

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func ConnectDatabase() *mongo.Database {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(MONGO_DB_URI))
	if err != nil {
		log.Fatal().Err(err).Msg("Error while connecting to database...")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal().Err(err).Msg("Error while ping database...")
	}

	return client.Database("shopcore")
}
