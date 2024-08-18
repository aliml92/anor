package postgres

import (
	"context"
	"encoding/json"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/order"
	"github.com/aliml92/anor/relation"
	"github.com/samber/lo"
	"github.com/samber/oops"
	"github.com/shopspring/decimal"
	"math/rand"
	"time"
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

func (s *OrderService) Create(ctx context.Context, params anor.OrderCreateParams) (int64, error) {
	co, err := s.orderRepository.Create(
		ctx,
		params.Cart.UserID,
		params.Cart.ID,
		order.PaymentMethod(params.PaymentMethod),
		order.PaymentStatus(params.PaymentStatus),
		order.OrderStatusProcessing,
		params.ShippingAddressID,
		params.IsPickup,
		params.Amount,
		params.Currency,
	)
	if err != nil {
		return 0, oops.Wrap(err)
	}

	// WORKAROUND: sqlc copyfrom not respecting query_parameter_limit #3388
	arg := make([]order.CreateOrderItemsParams, len(params.Cart.CartItems))
	for i, v := range params.Cart.CartItems {
		vab, err := json.Marshal(v.VariantAttributes)
		if err != nil {
			return co.ID, oops.Wrap(err)
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
	_, err = s.orderRepository.CreateOrderItems(ctx, arg)
	if err != nil {
		return co.ID, oops.Wrap(err)
	}

	return co.ID, nil
}

func calculateTotalAmount(cartItems []*anor.CartItem) decimal.Decimal {
	var total decimal.Decimal
	for _, cartItem := range cartItems {
		total = total.Add(cartItem.Price)
	}
	return total
}

func (s *OrderService) Get(ctx context.Context, orderID int64, withItems bool) (anor.Order, error) {
	o, err := s.orderRepository.Get(ctx, orderID)
	if err != nil {
		return anor.Order{}, oops.Wrap(err)
	}

	res := convertOrder(o)

	if withItems {
		panic("not implemented")
	}

	return res, nil
}

func (s *OrderService) List(ctx context.Context, params anor.OrderListParams) ([]anor.Order, error) {
	page := lo.Ternary(params.Page < 1, 1, params.Page)
	offset := (page - 1) * params.PageSize

	var (
		orders      []*order.Order
		addresses   []order.Address
		payments    []order.StripeCardPayment
		catchAllErr error

		withPayment         = params.WithRelations.Has(relation.StripeCardPayment)
		withShippingAddress = params.WithRelations.Has(relation.ShippingAddress)
		withOrderItems      = params.WithRelations.Has(relation.OrderItems)

		//TODO: add other includes if necessary
	)

	if withPayment && withShippingAddress {
		var oap []*order.ListWithPaymentAndAddressByUserIDRow
		oap, catchAllErr = s.orderRepository.ListWithPaymentAndAddressByUserID(ctx, params.UserID, int32(params.PageSize), int32(offset))
		if catchAllErr == nil {
			for _, v := range oap {
				orders = append(orders, &v.Order)
				addresses = append(addresses, v.Address)
				payments = append(payments, v.StripeCardPayment)
			}
		}
	} else {
		orders, catchAllErr = s.orderRepository.ListByUserID(ctx, params.UserID, int32(params.PageSize), int32(offset))
	}

	// Handle any errors that occurred during the switch statement
	if catchAllErr != nil {
		return nil, oops.Wrap(catchAllErr)
	}

	if len(orders) == 0 {
		return []anor.Order{}, nil
	}

	res := make([]anor.Order, len(orders))
	if withOrderItems {
		orderIDs := make([]int64, len(orders))
		for i, o := range orders {
			orderIDs[i] = o.ID
		}

		orderItems, err := s.orderRepository.GetItemByOrderIds(ctx, orderIDs)
		if err != nil {
			return nil, oops.Wrap(err)
		}

		itemsByOrderID := make(map[int64][]*order.OrderItem)
		for _, item := range orderItems {
			itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
		}

		for i, o := range orders {
			res[i] = convertOrder(o)

			items := itemsByOrderID[o.ID]
			res[i].OrderItems = make([]*anor.OrderItem, len(items))
			for j, item := range items {
				orderItem, err := convertOrderItem(item)
				if err != nil {
					return nil, oops.Wrap(err)
				}
				res[i].OrderItems[j] = orderItem
			}

			if withPayment && withShippingAddress {
				res[i].ShippingAddress = convertAddress(addresses[i])
				res[i].StripeCardPayment = convertCardPayment(payments[i])
			}

		}
	} else {
		for i, o := range orders {
			res[i] = convertOrder(o)
			if withPayment && withShippingAddress {
				res[i].ShippingAddress = convertAddress(addresses[i])
				res[i].StripeCardPayment = convertCardPayment(payments[i])
			}

		}
	}

	return res, nil
}

func (s *OrderService) ListActive(ctx context.Context, params anor.OrderListParams) ([]anor.Order, error) {
	page := lo.Ternary(params.Page < 1, 1, params.Page)
	offset := (page - 1) * params.PageSize

	var (
		orders      []*order.Order
		addresses   []order.Address
		payments    []order.StripeCardPayment
		catchAllErr error

		withPayment         = params.WithRelations.Has(relation.StripeCardPayment)
		withShippingAddress = params.WithRelations.Has(relation.ShippingAddress)
		withOrderItems      = params.WithRelations.Has(relation.OrderItems)

		//TODO: add other includes if necessary
	)

	if withPayment && withShippingAddress {
		var oap []*order.ListActiveWithPaymentAndAddressByUserIDRow
		oap, catchAllErr = s.orderRepository.ListActiveWithPaymentAndAddressByUserID(ctx, params.UserID, int32(params.PageSize), int32(offset))
		if catchAllErr == nil {
			for _, v := range oap {
				orders = append(orders, &v.Order)
				addresses = append(addresses, v.Address)
				payments = append(payments, v.StripeCardPayment)
			}
		}
	} else {
		orders, catchAllErr = s.orderRepository.ListActiveByUserID(ctx, params.UserID, int32(params.PageSize), int32(offset))
	}

	if catchAllErr != nil {
		return nil, oops.Wrap(catchAllErr)
	}

	if len(orders) == 0 {
		return []anor.Order{}, nil
	}

	res := make([]anor.Order, len(orders))
	if withOrderItems {
		orderIDs := make([]int64, len(orders))
		for i, o := range orders {
			orderIDs[i] = o.ID
		}

		orderItems, err := s.orderRepository.GetItemByOrderIds(ctx, orderIDs)
		if err != nil {
			return nil, oops.Wrap(err)
		}

		itemsByOrderID := make(map[int64][]*order.OrderItem)
		for _, item := range orderItems {
			itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
		}

		for i, o := range orders {
			res[i] = convertOrder(o)

			items := itemsByOrderID[o.ID]
			res[i].OrderItems = make([]*anor.OrderItem, len(items))
			for j, item := range items {
				orderItem, err := convertOrderItem(item)
				if err != nil {
					return nil, oops.Wrap(err)
				}
				res[i].OrderItems[j] = orderItem
			}

			if withPayment && withShippingAddress {
				res[i].ShippingAddress = convertAddress(addresses[i])
				res[i].StripeCardPayment = convertCardPayment(payments[i])
			}

		}
	} else {
		for i, o := range orders {
			res[i] = convertOrder(o)
			if withPayment && withShippingAddress {
				res[i].ShippingAddress = convertAddress(addresses[i])
				res[i].StripeCardPayment = convertCardPayment(payments[i])
			}

		}
	}

	return res, nil
}

func (s *OrderService) ListUnpaid(ctx context.Context, params anor.OrderListParams) ([]anor.Order, error) {
	page := lo.Ternary(params.Page < 1, 1, params.Page)
	offset := (page - 1) * params.PageSize

	var (
		orders      []*order.Order
		addresses   []order.Address
		payments    []order.StripeCardPayment
		catchAllErr error

		withPayment         = params.WithRelations.Has(relation.StripeCardPayment)
		withShippingAddress = params.WithRelations.Has(relation.ShippingAddress)
		withOrderItems      = params.WithRelations.Has(relation.OrderItems)

		//TODO: add other includes if necessary
	)

	if withPayment && withShippingAddress {
		var oap []*order.ListUnpaidWithPaymentAndAddressByUserIDRow
		oap, catchAllErr = s.orderRepository.ListUnpaidWithPaymentAndAddressByUserID(ctx, params.UserID, int32(params.PageSize), int32(offset))
		if catchAllErr == nil {
			for _, v := range oap {
				orders = append(orders, &v.Order)
				addresses = append(addresses, v.Address)
				payments = append(payments, v.StripeCardPayment)
			}
		}
	} else {
		orders, catchAllErr = s.orderRepository.ListUnpaidPaymentByUserID(ctx, params.UserID, int32(params.PageSize), int32(offset))
	}

	// Handle any errors that occurred during the switch statement
	if catchAllErr != nil {
		return nil, oops.Wrap(catchAllErr)
	}

	if len(orders) == 0 {
		return []anor.Order{}, nil
	}

	res := make([]anor.Order, len(orders))
	if withOrderItems {
		orderIDs := make([]int64, len(orders))
		for i, o := range orders {
			orderIDs[i] = o.ID
		}

		orderItems, err := s.orderRepository.GetItemByOrderIds(ctx, orderIDs)
		if err != nil {
			return nil, oops.Wrap(err)
		}

		itemsByOrderID := make(map[int64][]*order.OrderItem)
		for _, item := range orderItems {
			itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
		}

		for i, o := range orders {
			res[i] = convertOrder(o)

			items := itemsByOrderID[o.ID]
			res[i].OrderItems = make([]*anor.OrderItem, len(items))
			for j, item := range items {
				orderItem, err := convertOrderItem(item)
				if err != nil {
					return nil, oops.Wrap(err)
				}
				res[i].OrderItems[j] = orderItem
			}

			if withPayment && withShippingAddress {
				res[i].ShippingAddress = convertAddress(addresses[i])
				res[i].StripeCardPayment = convertCardPayment(payments[i])
			}

		}
	} else {
		for i, o := range orders {
			res[i] = convertOrder(o)
			if withPayment && withShippingAddress {
				res[i].ShippingAddress = convertAddress(addresses[i])
				res[i].StripeCardPayment = convertCardPayment(payments[i])
			}

		}
	}

	return res, nil
}

func convertOrder(o *order.Order) anor.Order {
	return anor.Order{
		ID:                o.ID,
		CartID:            o.CartID,
		UserID:            o.UserID,
		Status:            anor.OrderStatus(o.Status),
		PaymentMethod:     anor.PaymentMethod(o.PaymentMethod),
		PaymentStatus:     anor.PaymentStatus(o.PaymentStatus),
		ShippingAddressID: o.ShippingAddressID,
		IsPickup:          o.IsPickup,
		Amount:            o.Amount,
		Currency:          o.Currency,
		DeliveryDate:      generateEstimatedDeliveryDate(), //TODO: remove fake data
		CreatedAt:         o.CreatedAt.Time,
		UpdatedAt:         o.UpdatedAt.Time,
	}
}

func convertOrderItem(i *order.OrderItem) (*anor.OrderItem, error) {
	var variantAttributes map[string]string
	if err := json.Unmarshal(i.VariantAttributes, &variantAttributes); err != nil {
		return nil, err
	}
	return &anor.OrderItem{
		ID:                i.ID,
		OrderID:           i.OrderID,
		VariantID:         i.VariantID,
		Qty:               i.Qty,
		Price:             i.Price,
		Thumbnail:         i.Thumbnail,
		ProductName:       i.ProductName,
		VariantAttributes: variantAttributes,
		CreatedAt:         i.CreatedAt.Time,
		UpdatedAt:         i.UpdatedAt.Time,
	}, nil
}

func convertAddress(a order.Address) anor.Address {
	return anor.Address{
		ID:            a.ID,
		UserID:        a.UserID,
		DefaultFor:    anor.AddressDefaultType(a.DefaultFor.AddressDefaultType),
		Name:          a.Name,
		AddressLine1:  a.AddressLine1,
		AddressLine2:  anor.StringValue(a.AddressLine2),
		City:          a.City,
		StateProvince: anor.StringValue(a.StateProvince),
		PostalCode:    anor.StringValue(a.PostalCode),
		Country:       anor.StringValue(a.Country),
		Phone:         anor.StringValue(a.Phone),
		CreatedAt:     a.CreatedAt.Time,
		UpdatedAt:     a.UpdatedAt.Time,
	}
}

func convertCardPayment(p order.StripeCardPayment) anor.StripeCardPayment {
	return anor.StripeCardPayment{
		ID:               p.ID,
		OrderID:          p.OrderID,
		UserID:           anor.Int64Value(p.UserID),
		BillingAddressID: p.BillingAddressID,
		PaymentIntentID:  p.PaymentIntentID,
		PaymentMethodID:  anor.StringValue(p.PaymentMethodID),
		Amount:           p.Amount,
		Currency:         p.Currency,
		Status:           p.Status,
		ClientSecret:     anor.StringValue(p.ClientSecret),
		LastError:        anor.StringValue(p.LastError),
		CardLast4:        p.CardLast4,
		CardBrand:        p.CardBrand,
		CreatedAt:        p.CreatedAt.Time,
		UpdatedAt:        p.UpdatedAt.Time,
	}
}

func generateEstimatedDeliveryDate() time.Time {
	now := time.Now()

	minDays := 3
	maxDays := 7
	deliveryDays := rand.Intn(maxDays-minDays+1) + minDays

	return now.AddDate(0, 0, deliveryDays)
}
