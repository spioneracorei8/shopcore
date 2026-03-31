package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const ORDER_COLLECTION string = "order"

type OrderStatus string

const (
	ORDER_STATUS_PENDING   OrderStatus = "PENDING"
	ORDER_STATUS_PAID      OrderStatus = "PAID"
	ORDER_STATUS_CANCELLED OrderStatus = "CANCELLED"
	ORDER_STATUS_SHIPPED   OrderStatus = "SHIPPED"
	ORDER_STATUS_COMPLETED OrderStatus = "COMPLETED"
)

type Order struct {
	Id             *bson.ObjectID `bson:"_id" json:"_id"`
	CustomerId     *bson.ObjectID `bson:"customerId" json:"customerId"`
	OrderNo        string         `bson:"orderNo,omitempty" json:"orderNo"`
	Status         OrderStatus    `bson:"status,omitempty" json:"status" validate:"required"`
	OrderItems     []*OrderItems  `bson:"items,omitempty" json:"items"`
	Subtotal       float64        `bson:"subtotal,omitempty" json:"subtotal" validate:"required"`
	DiscountAmount float64        `bson:"discountAmount,omitempty" json:"discountAmount" validate:"required"`
	ShippingFee    float64        `bson:"shippingFee,omitempty" json:"shippingFee" validate:"required"`
	TotalAmount    float64        `bson:"totalAmount,omitempty" json:"totalAmount" validate:"required"`
	CreatedAt      time.Time      `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt      time.Time      `bson:"updatedAt,omitempty" json:"updatedAt"`
	DeletedAt      *time.Time     `bson:"deletedAt" json:"deletedAt"`
}

type OrderItems struct {
	ProductId           *bson.ObjectID `bson:"productId" json:"productId"`
	SkuSnapshot         string         `bson:"skuSnapshot,omitempty" json:"skuSnapshot" validate:"required"`
	ProductNameSnapshot string         `bson:"productNameSnapshot,omitempty" json:"productNameSnapshot" validate:"required"`
	UnitPrice           float64        `bson:"unitPrice,omitempty" json:"unitPrice" validate:"required"`
	Qty                 int            `bson:"qty,omitempty" json:"qty" validate:"required"`
	LineTotal           float64        `bson:"lineTotal,omitempty" json:"lineTotal" validate:"required"`
}

func (o *Order) GenObjectID() {
	id := bson.NewObjectID()
	o.Id = &id
}

func (o *Order) SetCreatedAt() {
	now := time.Now()
	o.CreatedAt = now
}

func (o *Order) SetUpdatedAt() {
	now := time.Now()
	o.UpdatedAt = now
}

func (o *Order) SetDeletedAt() {
	now := time.Now()
	o.DeletedAt = &now
}
