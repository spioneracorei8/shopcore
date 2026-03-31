package mongodb

import (
	"context"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/outbound"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type productRepository struct {
	client *mongo.Database
}

func NewProductRepoImpl(client *mongo.Database) outbound.ProductRepository {
	return &productRepository{
		client: client,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	_, err := r.client.Collection(domain.PRODUCT_COLLECTION).InsertOne(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) FetchListProducts(ctx context.Context) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	filter := bson.M{
		"deletedAt": nil,
	}
	cursor, err := r.client.Collection(domain.PRODUCT_COLLECTION).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) FetchProductById(ctx context.Context, id *bson.ObjectID) (*domain.Product, error) {
	var product domain.Product
	filter := bson.M{"_id": id}
	if err := r.client.Collection(domain.PRODUCT_COLLECTION).FindOne(ctx, filter).Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": product}
	if err := r.client.Collection(domain.PRODUCT_COLLECTION).FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		return nil
	}

	return nil
}

func (r *productRepository) DeleteProductById(ctx context.Context, id *bson.ObjectID, product *domain.Product) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": product}
	_, err := r.client.Collection(domain.PRODUCT_COLLECTION).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
