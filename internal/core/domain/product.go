package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const PRODUCT_COLLECTION = "product"

type ProductStatus string

const (
	PRODUCT_STATUS_ACTIVE         ProductStatus = "ACTIVE"         // สินค้าพร้อมขาย
	PRODUCT_STATUS_INACTIVE       ProductStatus = "INACTIVE"       // สินค้าปิดการขาย
	PRODUCT_STATUS_OUT_OF_STOCK   ProductStatus = "OUT_OF_STOCK"   // สินค้าหมดสต็อก
	PRODUCT_STATUS_DISCONTINUED   ProductStatus = "DISCONTINUED"   // สินค้ายกเลิกการผลิต
	PRODUCT_STATUS_PENDING_REVIEW ProductStatus = "PENDING_REVIEW" // สินค้าอยู่ระหว่างตรวจสอบ
	PRODUCT_STATUS_COMING_SOON    ProductStatus = "COMING_SOON"    // สินค้าจะเปิดตัวเร็วๆนี้
	PRODUCT_STATUS_PRE_ORDER      ProductStatus = "PRE_ORDER"      // สินค้าพรีออเดอร์
	PRODUCT_STATUS_LOW_STOCK      ProductStatus = "LOW_STOCK"      // สินค้าเหลือน้อย
)

type Product struct {
	Id         *bson.ObjectID `bson:"_id" json:"_id"`
	Sku        string         `bson:"sku,omitempty" json:"sku" validate:"required"`
	Name       string         `bson:"name,omitempty" json:"name" validate:"required"`
	Descrption string         `bson:"descrption,omitempty" json:"descrption" validate:"required"`
	Price      float64        `bson:"price,omitempty" json:"price" validate:"required"`
	StockQty   int            `bson:"stockQty,omitempty" json:"stockQty" validate:"required"`
	Status     ProductStatus  `bson:"status,omitempty" json:"status"`
	CreatedAt  time.Time      `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt  time.Time      `bson:"updatedAt,omitempty" json:"updatedAt"`
	DeletedAt  *time.Time     `bson:"deletedAt" json:"deletedAt"`
}

func (p *Product) GenObjectID() {
	id := bson.NewObjectID()
	p.Id = &id
}

func (p *Product) SetCreatedAt() {
	now := time.Now()
	p.CreatedAt = now
}

func (p *Product) SetUpdatedAt() {
	now := time.Now()
	p.UpdatedAt = now
}

func (p *Product) SetDeletedAt() {
	now := time.Now()
	p.DeletedAt = &now
}
