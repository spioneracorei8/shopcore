package mongodb

import (
	"context"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/outbound"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type orderRepository struct {
	client *mongo.Database
}

func NewOrderRepoImpl(client *mongo.Database) outbound.OrderRepository {
	return &orderRepository{
		client: client,
	}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	_, err := r.client.Collection(domain.ORDER_COLLECTION).InsertOne(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) FetchListOrders(ctx context.Context) ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0)
	filter := bson.M{}
	cursor, err := r.client.Collection(domain.ORDER_COLLECTION).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) FetchOrderById(ctx context.Context, id *bson.ObjectID) (*domain.Order, error) {
	var order domain.Order
	filter := bson.M{"_id": id}
	if err := r.client.Collection(domain.ORDER_COLLECTION).FindOne(ctx, filter).Decode(&order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrderById(ctx context.Context, id *bson.ObjectID, order *domain.Order) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": order}
	if err := r.client.Collection(domain.ORDER_COLLECTION).FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		return nil
	}

	return nil
}

func (r *orderRepository) DeleteOrderById(ctx context.Context, id *bson.ObjectID, order *domain.Order) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": order}
	_, err := r.client.Collection(domain.ORDER_COLLECTION).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
