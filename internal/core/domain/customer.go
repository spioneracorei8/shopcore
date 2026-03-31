package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const CUSTOMER_COLLECTION = "customer"

type CustomerStatus string

const (
	CUSTOMER_STATUS_ACTIVE   CustomerStatus = "ACTIVE"
	CUSTOMER_STATUS_INACTIVE CustomerStatus = "INACTIVE"
)

type Customer struct {
	Id        *bson.ObjectID `bson:"_id" json:"_id"`
	Email     string         `bson:"email,omitempty" json:"email" validate:"required,email"`
	FirstName string         `bson:"firstName,omitempty" json:"firstName"  validate:"required,max=255"`
	LastName  string         `bson:"lastName,omitempty" json:"lastName"  validate:"required,max=255"`
	Phone     string         `bson:"phone,omitempty" json:"phone" validate:"required,max=10"`
	Status    CustomerStatus `bson:"status,omitempty" json:"status"`
	CreatedAt time.Time      `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time      `bson:"updatedAt,omitempty" json:"updatedAt"`
	DeletedAt *time.Time     `bson:"deletedAt" json:"deletedAt"`
}

func (c *Customer) GenObjectID() {
	id := bson.NewObjectID()
	c.Id = &id
}

func (c *Customer) SetCreatedAt() {
	now := time.Now()
	c.CreatedAt = now
}

func (c *Customer) SetUpdatedAt() {
	now := time.Now()
	c.UpdatedAt = now
}

func (c *Customer) SetDeletedAt() {
	now := time.Now()
	c.DeletedAt = &now
}
