package services

import (
	"context"
	"fmt"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/inbound"
	"shopcore/internal/core/ports/outbound"
	helper "shopcore/pkg/helpers"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type orderUsecase struct {
	rnUs      inbound.RunNumberUsecase
	productUs inbound.ProductUsecase
	orderRepo outbound.OrderRepository
}

func NewOrderUsecaseImpl(rnUs inbound.RunNumberUsecase, productUs inbound.ProductUsecase, orderRepo outbound.OrderRepository) inbound.OrderUsecase {
	return &orderUsecase{
		rnUs:      rnUs,
		productUs: productUs,
		orderRepo: orderRepo,
	}
}

func (u *orderUsecase) formatOrderNumber(rn *domain.RunNumber) string {
	formatted := fmt.Sprintf("%s-%s-%04d", rn.Prefix, helper.GetOrderDate(), rn.Running)
	return formatted
}

func (u *orderUsecase) CreateOrder(ctx context.Context, order *domain.Order) error {
	rn, err := u.rnUs.FetchRunNumber(ctx)
	if err != nil {
		return err
	}

	if err := u.rnUs.UpdateRunNumber(ctx); err != nil {
		return err
	}

	order.GenObjectID()
	order.OrderNo = u.formatOrderNumber(rn)
	order.SetCreatedAt()
	order.SetUpdatedAt()

	for _, item := range order.OrderItems {
		product, err := u.productUs.FetchProductById(ctx, item.ProductId)
		if err != nil {
			return err
		}
		product.StockQty -= item.Qty
		u.productUs.UpdateProductById(ctx, item.ProductId, product)
	}

	return u.orderRepo.CreateOrder(ctx, order)
}

func (u *orderUsecase) FetchListOrders(ctx context.Context) ([]*domain.Order, error) {
	return u.orderRepo.FetchListOrders(ctx)
}

func (u *orderUsecase) FetchOrderById(ctx context.Context, id *bson.ObjectID) (*domain.Order, error) {
	return u.orderRepo.FetchOrderById(ctx, id)
}

func (u *orderUsecase) UpdateOrderById(ctx context.Context, id *bson.ObjectID, order *domain.Order) (*domain.Order, error) {
	order.Id = id
	order.SetUpdatedAt()

	for _, item := range order.OrderItems {
		product, err := u.productUs.FetchProductById(ctx, item.ProductId)
		if err != nil {
			return nil, err
		}
		product.StockQty -= item.Qty
		u.productUs.UpdateProductById(ctx, item.ProductId, product)
	}

	if err := u.orderRepo.UpdateOrderById(ctx, id, order); err != nil {
		return nil, err
	}
	return u.orderRepo.FetchOrderById(ctx, id)
}

func (u *orderUsecase) DeleteOrderById(ctx context.Context, id *bson.ObjectID, order *domain.Order) error {
	order.Id = id
	order.Status = domain.ORDER_STATUS_CANCELLED
	order.SetUpdatedAt()
	order.SetDeletedAt()
	return u.orderRepo.DeleteOrderById(ctx, id, order)
}
