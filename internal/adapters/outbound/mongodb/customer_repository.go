package mongodb

import (
	"context"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/outbound"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type customerRepository struct {
	client *mongo.Database
}

func NewCustomerRepoImpl(client *mongo.Database) outbound.CustomerRepository {
	return &customerRepository{
		client: client,
	}
}

func (r *customerRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	_, err := r.client.Collection(domain.CUSTOMER_COLLECTION).InsertOne(ctx, customer)
	if err != nil {
		return err
	}

	return nil
}

func (r *customerRepository) FetchListCustomers(ctx context.Context) ([]*domain.Customer, error) {
	var customers = make([]*domain.Customer, 0)
	// filter := bson.M{
	// 	"deletedAt": bson.M{"$ne": nil},
	// }

	filter := bson.M{
		"deletedAt": nil,
	}
	cursor, err := r.client.Collection(domain.CUSTOMER_COLLECTION).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *customerRepository) FetchCustomerById(ctx context.Context, id *bson.ObjectID) (*domain.Customer, error) {
	var customer domain.Customer
	filter := bson.M{
		"_id": id,
	}
	if err := r.client.Collection(domain.CUSTOMER_COLLECTION).FindOne(ctx, filter).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) UpdateCustomerById(ctx context.Context, id *bson.ObjectID, customer *domain.Customer) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": customer}

	if err := r.client.Collection(domain.CUSTOMER_COLLECTION).FindOneAndUpdate(ctx, filter, update); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *customerRepository) DeleteCustomerById(ctx context.Context, id *bson.ObjectID, customer *domain.Customer) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": customer}

	if err := r.client.Collection(domain.CUSTOMER_COLLECTION).FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		return err
	}

	return nil
}
