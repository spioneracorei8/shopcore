package mongodb

import (
	"context"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/outbound"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type runNumberRepository struct {
	client *mongo.Database
}

func NewRunNumberRepoImpl(client *mongo.Database) outbound.RunNumberRepository {
	return &productRepository{
		client: client,
	}
}

func (r *productRepository) CreateRunNumber(ctx context.Context, rn *domain.RunNumber) error {
	_, err := r.client.Collection(domain.RUN_NUMBER_COLLECTION).InsertOne(ctx, rn)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) FetchRunNumber(ctx context.Context) (*domain.RunNumber, error) {
	var rn domain.RunNumber
	filter := bson.M{}
	if err := r.client.Collection(domain.RUN_NUMBER_COLLECTION).FindOne(ctx, filter).Decode(&rn); err != nil {
		return nil, err
	}
	return &rn, nil
}

func (r *productRepository) UpdateRunNumber(ctx context.Context, rn *domain.RunNumber) error {
	filter := bson.M{"_id": rn.Id}
	update := bson.M{"$set": rn}
	_, err := r.client.Collection(domain.RUN_NUMBER_COLLECTION).UpdateOne(ctx, filter, update)

	if err != nil {
		return nil
	}

	return nil
}
