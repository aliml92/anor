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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/oops"
	"strconv"
	"time"
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

func (s *CartService) Create(ctx context.Context, params anor.CartCreateParams) (anor.Cart, error) {
	if params.UserID == 0 && params.ExpiresAt.IsZero() {
		return anor.Cart{}, oops.Errorf("UserID and ExpiresAt cannot be nil at the same time")
	}

	var (
		c   *cart.Cart
		err error
	)
	if params.UserID != 0 {
		c, err = s.cartRepository.CreateCart(ctx, anor.Int64(params.UserID))
		if err != nil {
			return anor.Cart{}, oops.Wrap(err)
		}
	}

	if !params.ExpiresAt.IsZero() {
		if time.Now().After(params.ExpiresAt) {
			return anor.Cart{}, oops.Errorf("invalid expiration time")
		}
		c, err = s.cartRepository.CreateGuestCart(ctx, convertToPgTimestamptz(params.ExpiresAt))
		if err != nil {
			return anor.Cart{}, oops.Wrap(err)
		}
	}

	ac := anor.Cart{
		ID:        c.ID,
		UserID:    anor.Int64Value(c.UserID),
		Status:    anor.CartStatus(c.Status),
		ExpiresAt: c.ExpiresAt.Time,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}

	return ac, nil
}

func (s *CartService) Get(ctx context.Context, id int64, includeItems bool) (anor.Cart, error) {
	c, err := s.cartRepository.GetCartByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.Cart{}, anor.ErrCartNotFound
		}
		return anor.Cart{}, oops.Wrap(err)
	}

	ac := anor.Cart{
		ID:        c.ID,
		UserID:    anor.Int64Value(c.UserID),
		Status:    anor.CartStatus(c.Status),
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}

	if includeItems {
		cartItems, err := s.cartRepository.GetCartItemsByCartID(ctx, c.ID)
		if err != nil {
			return anor.Cart{}, oops.Wrap(err)
		}

		items := make([]*anor.CartItem, len(cartItems))
		variantIds := make([]int64, len(cartItems))
		for index, item := range cartItems {
			it := &anor.CartItem{
				ID:          item.ID,
				CartID:      item.CartID,
				VariantID:   item.VariantID,
				Qty:         item.Qty,
				Price:       item.Price,
				Currency:    item.Currency,
				Thumbnail:   item.Thumbnail,
				ProductName: item.ProductName,
				ProductPath: item.ProductPath,
				CreatedAt:   item.CreatedAt.Time,
				UpdatedAt:   item.UpdatedAt.Time,
			}

			attrs := make(map[string]string)
			if err := json.Unmarshal(item.VariantAttributes, &attrs); err != nil {
				return anor.Cart{}, oops.Wrap(err)
			}
			it.VariantAttributes = attrs
			items[index] = it

			variantIds[index] = item.VariantID
		}

		if len(variantIds) > 0 {
			ph, err := s.productRepository.GetProductVariantQtyInBulk(ctx, variantIds)
			if err != nil {
				return anor.Cart{}, oops.Wrap(err)
			}

			for index := range items {
				items[index].AvailableQty = ph[index].Qty
			}
		}

		ac.CartItems = items
	}

	return ac, nil
}

func (s *CartService) GetByUserID(ctx context.Context, userID int64, includeItems bool) (anor.Cart, error) {
	c, err := s.cartRepository.GetCartByUserID(ctx, &userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.Cart{}, anor.ErrCartNotFound
		}
		return anor.Cart{}, oops.Wrap(err)
	}

	ac := anor.Cart{
		ID:        c.ID,
		UserID:    anor.Int64Value(c.UserID),
		Status:    anor.CartStatus(c.Status),
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}

	if includeItems {
		cartItems, err := s.cartRepository.GetCartItemsByCartID(ctx, c.ID)
		if err != nil {
			return anor.Cart{}, oops.Wrap(err)
		}

		items := make([]*anor.CartItem, len(cartItems))
		variantIds := make([]int64, len(cartItems))
		for index, item := range cartItems {
			it := &anor.CartItem{
				ID:          item.ID,
				CartID:      item.CartID,
				VariantID:   item.VariantID,
				Qty:         item.Qty,
				Price:       item.Price,
				Currency:    item.Currency,
				Thumbnail:   item.Thumbnail,
				ProductName: item.ProductName,
				ProductPath: item.ProductPath,
				CreatedAt:   item.CreatedAt.Time,
				UpdatedAt:   item.UpdatedAt.Time,
			}

			attrs := make(map[string]string)
			if err := json.Unmarshal(item.VariantAttributes, &attrs); err != nil {
				return anor.Cart{}, oops.Wrap(err)
			}
			it.VariantAttributes = attrs
			items[index] = it

			variantIds[index] = item.VariantID
		}

		if len(variantIds) > 0 {
			ph, err := s.productRepository.GetProductVariantQtyInBulk(ctx, variantIds)
			if err != nil {
				return anor.Cart{}, oops.Wrap(err)
			}

			for index := range items {
				items[index].AvailableQty = ph[index].Qty
			}
		}

		ac.CartItems = items
	}

	return ac, nil
}

func (s *CartService) Update(ctx context.Context, id int64, params anor.CartUpdateParams) error {
	rowsAffected, err := s.cartRepository.UpdateCartStatus(ctx, id, cart.CartStatus(params.Status))
	if err != nil {
		return oops.Wrap(err)
	}

	if rowsAffected == 0 {
		return anor.ErrCartNotFound
	}

	return nil
}

func (s *CartService) Merge(ctx context.Context, params anor.CartMergeParams) (anor.Cart, error) {
	userID := anor.Int64(params.UserID)
	c, err := s.cartRepository.GetCartByUserID(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		c, err = s.cartRepository.CreateCart(ctx, userID)
		if err != nil {
			return anor.Cart{}, oops.Wrap(err)
		}
	}
	if err != nil {
		return anor.Cart{}, oops.Wrap(err)
	}
	err = s.cartRepository.MergeGuestCartWithUserCart(ctx, params.GuestCartID, c.ID)
	if err != nil {
		return anor.Cart{}, oops.Wrap(err)
	}

	_, err = s.cartRepository.UpdateCartStatus(ctx, params.GuestCartID, cart.CartStatusMerged)
	if err != nil {
		return anor.Cart{}, oops.Wrap(err)
	}

	return anor.Cart{
		ID:        c.ID,
		UserID:    anor.Int64Value(c.UserID),
		Status:    anor.CartStatus(c.Status),
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}, nil
}

func (s *CartService) ListItems(ctx context.Context, params anor.CartItemListParams) ([]anor.CartItem, error) {
	cartItems, err := s.cartRepository.GetCartItemsByCartID(ctx, params.CartID)
	if err != nil {
		return nil, oops.Wrap(err)
	}

	items := make([]anor.CartItem, len(cartItems))
	for index, item := range cartItems {
		it := anor.CartItem{
			ID:          item.ID,
			CartID:      item.CartID,
			VariantID:   item.VariantID,
			Qty:         item.Qty,
			Price:       item.Price,
			Currency:    item.Currency,
			Thumbnail:   item.Thumbnail,
			ProductName: item.ProductName,
			ProductPath: item.ProductPath,
			CreatedAt:   item.CreatedAt.Time,
			UpdatedAt:   item.UpdatedAt.Time,
		}

		attrs := make(map[string]string)
		if err := json.Unmarshal(item.VariantAttributes, &attrs); err != nil {
			return nil, oops.Wrap(err)
		}
		it.VariantAttributes = attrs
		items[index] = it
	}

	return items, nil
}

func (s *CartService) AddItem(ctx context.Context, params anor.CartItemAddParams) (anor.CartItem, error) {
	var ci anor.CartItem
	cd, err := s.productRepository.CollectCartData(ctx, params.VariantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.CartItem{}, anor.ErrProductVariantNotFound
		}
	}

	// TODO: research on returning error when pv.Qty is 0, currently 'optimistic update' applied
	// EXPLANATION: customer clicks 'add to cart' button on product details page, depending on how long
	// it took from the time this page opened until 'add to cart' button clicked there might be a certain
	// change another customer already bought this product variant. We don't check remaining stock when
	// we add this cart item to cart, Let's be optimistic.
	if cd.IsCustomPriced {
		pvPrice, err := s.productRepository.GetProductVariantPricingByVariantID(ctx, params.VariantID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ci, anor.ErrProductVariantPricingNotFound
			}
			return ci, oops.Wrap(err)
		}
		ci.Price = pvPrice.DiscountedPrice
		ci.Currency = pvPrice.CurrencyCode

	} else {
		pPrice, err := s.productRepository.GetProductPricingByProductID(ctx, cd.ProductID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ci, anor.ErrProductPricingNotFound
			}
			return ci, oops.Wrap(err)
		}
		ci.Price = pPrice.DiscountedPrice
		ci.Currency = pPrice.CurrencyCode
	}

	// get thumbnail
	imgIds := cd.ImageIdentifiers
	imgId := 0
	if len(imgIds) != 0 {
		imgId = int(imgIds[0])
	}

	ci.Thumbnail = cd.ProductImageUrls[imgId]
	ci.CartID = params.CartID
	ci.VariantID = params.VariantID
	ci.Qty = int32(params.Qty)
	ci.ProductName = cd.ProductName
	ci.ProductPath = cd.ProductHandle + "-" + strconv.FormatInt(cd.ProductID, 10)

	var attr map[string]string
	if err := json.Unmarshal(cd.VariantAttributes, &attr); err != nil {
		return anor.CartItem{}, oops.Wrap(err)
	}
	ci.VariantAttributes = attr

	newCartItem, err := s.cartRepository.AddCartItem(ctx,
		ci.CartID,
		ci.VariantID,
		ci.Qty,
		ci.Price,
		ci.Currency,
		ci.Thumbnail,
		ci.ProductName,
		ci.ProductPath,
		cd.VariantAttributes,
	)

	if err != nil {
		if cartItemAlreadyExists(err) {
			_, err := s.cartRepository.IncrementCartItemQty(ctx, ci.CartID, ci.VariantID, ci.Qty)
			if err != nil {
				return ci, oops.Wrap(err)
			}
		} else {
			return ci, oops.Wrap(err)
		}
	}

	ci.ID = newCartItem.ID
	ci.CreatedAt = newCartItem.CreatedAt.Time
	ci.UpdatedAt = newCartItem.UpdatedAt.Time

	return ci, nil
}

func (s *CartService) UpdateItem(ctx context.Context, itemID int64, params anor.CartItemUpdateParams) error {
	rowsAffected, err := s.cartRepository.UpdateCartItemQty(ctx, itemID, int32(params.Qty))
	if err != nil {
		return oops.Wrap(err)
	}

	if rowsAffected == 0 {
		return anor.ErrCartItemNotFound
	}

	return nil
}

func (s *CartService) DeleteItem(ctx context.Context, itemID int64) error {
	rowsAffected, err := s.cartRepository.DeleteCartItem(ctx, itemID)
	if err != nil {
		return oops.Wrap(err)
	}

	if rowsAffected == 0 {
		return anor.ErrCartItemNotFound
	}

	return nil
}

func (s *CartService) CountItems(ctx context.Context, cartID int64) (int64, error) {
	count, err := s.cartRepository.CountCartItemsByCartID(ctx, cartID)
	if err != nil {
		return 0, oops.Wrap(err)
	}
	return count, nil
}

func (s *CartService) IsMyItem(ctx context.Context, params anor.IsMyCartItemParams) (bool, error) {
	ok, err := s.cartRepository.CartItemExistsByCartID(ctx, params.ItemID, params.CartID)
	if err != nil {
		return false, oops.Wrap(err)
	}
	return ok, nil
}

func convertToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
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
