package postgres

import (
	"context"
	"encoding/json"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/order"
	"github.com/samber/oops"
	"github.com/shopspring/decimal"
)

var _ anor.OrderService = (*OrderService)(nil)

type OrderService struct {
	orderRepository order.Repository
}

func NewOrderService(or order.Repository) *OrderService {
	return &OrderService{
		orderRepository: or,
	}
}

func (os *OrderService) ConvertCartToOrder(ctx context.Context, c anor.Cart, piID string) (anor.Order, error) {

	// TODO: implement shipping/billing address
	shippingAddr := map[string]string{"Address": "123 Maple Street\nApt 4B\nSpringfield, IL 62704\nUSA"}
	shippingAddrByte, err := json.Marshal(shippingAddr)
	if err != nil {
		return anor.Order{}, oops.Wrap(err)
	}

	billingAddr := map[string]string{"Address": "456 Oak Avenue\nSuite 300\nMetropolis, NY 10001\nUSA"}
	billingAddrByte, err := json.Marshal(billingAddr)
	if err != nil {
		return anor.Order{}, oops.Wrap(err)
	}

	co, err := os.orderRepository.CreateOrder(
		ctx,
		c.ID,
		&c.UserID,
		calculateOrderTotalAmount(c.CartItems),
		piID,
		shippingAddrByte,
		billingAddrByte,
	)
	if err != nil {
		return anor.Order{}, oops.Errorf("failed to create order: %v", err)
	}

	// WORKAROUND: sqlc copyfrom not respecting query_parameter_limit #3388
	arg := make([]order.CreateOrderItemsParams, len(c.CartItems))
	for i, v := range c.CartItems {
		vab, err := json.Marshal(v.VariantAttributes)
		if err != nil {
			return anor.Order{}, oops.Wrap(err)
		}

		arg[i] = order.CreateOrderItemsParams{
			OrderID:           co.ID,
			VariantID:         v.VariantID,
			Qty:               v.Qty,
			Price:             v.Price,
			Thumbnail:         v.Thumbnail,
			ProductName:       v.ProductName,
			VariantAttributes: vab,
		}
	}
	_, err = os.orderRepository.CreateOrderItems(ctx, arg)
	if err != nil {
		return anor.Order{}, oops.Errorf("failed to create order items: %v", err)
	}

	return anor.Order{
		ID:              co.ID,
		CartID:          co.ID,
		UserID:          co.ID,
		Status:          anor.OrderStatus(co.Status),
		TotalAmount:     co.TotalAmount,
		PaymentIntentID: co.PaymentIntentID,
		ShippingAddress: shippingAddr,
		BillingAddress:  billingAddr,
		CreatedAt:       co.CreatedAt.Time,
		UpdatedAt:       co.UpdatedAt.Time,
	}, nil
}

func calculateOrderTotalAmount(cartItems []*anor.CartItem) decimal.Decimal {
	var totalPrice decimal.Decimal
	for _, cartItem := range cartItems {
		totalPrice = totalPrice.Add(cartItem.Price)
	}

	return totalPrice
}
