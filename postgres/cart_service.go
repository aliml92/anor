package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/cart"
	"github.com/aliml92/anor/postgres/repository/product"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"strconv"
)

var _ anor.CartService = (*CartService)(nil)

type CartService struct {
	cartRepository    cart.Repository
	productRepository product.Repository
}

func NewCartService(cr cart.Repository, pr product.Repository) *CartService {
	return &CartService{
		cartRepository:    cr,
		productRepository: pr,
	}
}

func (cs *CartService) GetGuestCartItems(ctx context.Context, cartID int64) ([]*anor.CartItem, error) {
	cartItems, err := cs.cartRepository.GetCartItemsByCartID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	items := make([]*anor.CartItem, len(cartItems))
	variantIds := make([]int64, len(cartItems))
	for index, item := range cartItems {
		it := &anor.CartItem{
			ID:           item.ID,
			CartID:       item.CartID,
			VariantID:    item.VariantID,
			Qty:          item.Qty,
			Price:        item.Price,
			CurrencyCode: item.CurrencyCode,
			Thumbnail:    item.Thumbnail,
			ProductName:  item.ProductName,
			ProductPath:  item.ProductPath,
			CreatedAt:    item.CreatedAt.Time,
			UpdatedAt:    item.UpdatedAt.Time,
		}

		attrs := make(map[string]string)
		if err := json.Unmarshal(item.VariantAttributes, &attrs); err != nil {
			return nil, err
		}
		it.VariantAttributes = attrs
		items[index] = it

		variantIds[index] = item.VariantID
	}

	if len(variantIds) > 0 {
		ph, err := cs.productRepository.GetProductVariantQtyInBulk(ctx, variantIds)
		if err != nil {
			return nil, err
		}

		for index := range items {
			items[index].AvailableQty = ph[index].Qty
		}
	}

	return items, nil
}

func (cs *CartService) CreateCart(ctx context.Context, userID int64) (anor.Cart, error) {
	var (
		c   *cart.Cart
		err error
	)
	if userID > 0 {
		c, err = cs.cartRepository.CreateCart(ctx, &userID)
		if err != nil {
			return anor.Cart{}, err
		}
	} else {
		c, err = cs.cartRepository.CreateGuestCart(ctx)
		if err != nil {
			return anor.Cart{}, err
		}
	}

	ac := anor.Cart{
		ID:        c.ID,
		Status:    anor.CartStatus(c.Status),
		CreatedAt: c.CreatedAt.Time,
		DeletedAt: c.DeletedAt.Time,
	}

	if c.UserID != nil {
		ac.UserID = *c.UserID
	}

	return ac, nil
}

func (cs *CartService) GetCart(ctx context.Context, userID int64, includeCartItems bool) (anor.Cart, error) {
	c, err := cs.cartRepository.GetCartByUserID(ctx, &userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.Cart{}, anor.ErrNotFound
		}
		return anor.Cart{}, err
	}

	ac := anor.Cart{
		ID:        c.ID,
		UserID:    *c.UserID,
		Status:    anor.CartStatus(c.Status),
		CreatedAt: c.CreatedAt.Time,
		DeletedAt: c.DeletedAt.Time,
	}

	if c.PiClientSecret != nil {
		ac.PIClientSecret = *c.PiClientSecret
	}

	if includeCartItems {
		cartItems, err := cs.cartRepository.GetCartItemsByCartID(ctx, c.ID)
		if err != nil {
			return anor.Cart{}, err
		}

		items := make([]*anor.CartItem, len(cartItems))
		variantIds := make([]int64, len(cartItems))
		for index, item := range cartItems {
			it := &anor.CartItem{
				ID:           item.ID,
				CartID:       item.CartID,
				VariantID:    item.VariantID,
				Qty:          item.Qty,
				Price:        item.Price,
				CurrencyCode: item.CurrencyCode,
				Thumbnail:    item.Thumbnail,
				ProductName:  item.ProductName,
				CreatedAt:    item.CreatedAt.Time,
				UpdatedAt:    item.UpdatedAt.Time,
			}

			attrs := make(map[string]string)
			if err := json.Unmarshal(item.VariantAttributes, &attrs); err != nil {
				return anor.Cart{}, err
			}
			it.VariantAttributes = attrs
			items[index] = it

			variantIds[index] = item.VariantID
		}

		if len(variantIds) > 0 {
			ph, err := cs.productRepository.GetProductVariantQtyInBulk(ctx, variantIds)
			if err != nil {
				return anor.Cart{}, err
			}

			for index := range items {
				items[index].AvailableQty = ph[index].Qty
			}
		}

		ac.CartItems = items
	}

	return ac, nil
}

func (cs *CartService) AddCartItem(ctx context.Context, cartID int64, p anor.AddCartItemParam) (anor.CartItem, error) {
	var ci anor.CartItem
	cd, err := cs.productRepository.CollectCartData(ctx, p.VariantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.CartItem{}, anor.ErrNotFound
		}
	}

	// TODO: research on returning error when pv.Qty is 0, currently 'optimistic update' applied
	// EXPLANATION: customer clicks 'add to cart' button on product details page, depending on how long
	// it took from the time this page opened until 'add to cart' button clicked there might be a certain
	// change another customer already bought this product variant. We don't check remaining stock when
	// we add this cart item to cart, Let's be optimistic.
	if cd.IsCustomPriced {
		pvPrice, err := cs.productRepository.GetProductVariantPricingByVariantID(ctx, p.VariantID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ci, anor.ErrProductVariantPricingNotFound
			}
			return ci, err
		}
		ci.Price = pvPrice.DiscountedPrice
		ci.CurrencyCode = pvPrice.CurrencyCode

	} else {
		pPrice, err := cs.productRepository.GetProductPricingByProductID(ctx, cd.ProductID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ci, anor.ErrProductPricingNotFound
			}
			return ci, err
		}
		ci.Price = pPrice.DiscountedPrice
		ci.CurrencyCode = pPrice.CurrencyCode
	}

	// get thumbnail
	imgIds := cd.ImageIdentifiers
	imgId := 0
	if len(imgIds) != 0 {
		imgId = int(imgIds[0])
	}

	ci.Thumbnail = cd.ProductImageUrls[imgId]
	ci.CartID = cartID
	ci.VariantID = p.VariantID
	ci.Qty = int32(p.Qty)
	ci.ProductName = cd.ProductName
	ci.ProductPath = cd.ProductHandle + "-" + strconv.FormatInt(cd.ProductID, 10)

	var attr map[string]string
	if err := json.Unmarshal(cd.VariantAttributes, &attr); err != nil {
		return anor.CartItem{}, err
	}
	ci.VariantAttributes = attr

	newCartItem, err := cs.cartRepository.AddCartItem(ctx,
		ci.CartID,
		ci.VariantID,
		ci.Qty,
		ci.Price,
		ci.CurrencyCode,
		ci.Thumbnail,
		ci.ProductName,
		ci.ProductPath,
		cd.VariantAttributes,
	)

	if err != nil {
		if cartItemAlreadyExists(err) {
			_, err := cs.cartRepository.IncrementCartItemQty(ctx, ci.CartID, ci.VariantID, ci.Qty)
			if err != nil {
				return ci, err
			}
		} else {
			return ci, err
		}
	}

	ci.ID = newCartItem.ID
	ci.CreatedAt = newCartItem.CreatedAt.Time
	ci.UpdatedAt = newCartItem.UpdatedAt.Time

	return ci, nil
}

func (cs *CartService) UpdateCart(ctx context.Context, c anor.Cart) error {
	if c.ID == 0 || c.UserID == 0 {
		return anor.ErrInvalidCart
	}
	if c.PIClientSecret != "" {
		if err := cs.cartRepository.UpdateCartClientSecret(ctx, c.ID, &c.UserID, &c.PIClientSecret); err != nil {
			return err
		}
	}

	return nil
}

func (cs *CartService) UpdateCartItem(ctx context.Context, cartItemID int64, p anor.UpdateCartItemParam) error {
	err := cs.cartRepository.UpdateCartItemQty(ctx, cartItemID, int32(p.Qty))
	return err
}

func (cs *CartService) DeleteCartItem(ctx context.Context, cartItemID int64) error {
	err := cs.cartRepository.DeleteCartItem(ctx, cartItemID)
	return err
}

func (cs *CartService) CountCartItems(ctx context.Context, cartID int64) (int64, error) {
	count, err := cs.cartRepository.CountCartItemsByCartID(ctx, cartID)
	return count, err
}

func (cs *CartService) IsCartItemOwner(ctx context.Context, userID int64, cartItemId int64) (bool, error) {
	ok, err := cs.cartRepository.CartItemExistsByUserID(ctx, cartItemId, userID)
	return ok, err
}

func (cs *CartService) IsGuestCartItemOwner(ctx context.Context, cartID int64, cartItemId int64) (bool, error) {
	ok, err := cs.cartRepository.CartItemExistsByCartID(ctx, cartID, cartItemId)
	return ok, err
}

func cartItemAlreadyExists(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return true
		}
	}

	return false
}
