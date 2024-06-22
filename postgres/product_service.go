package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"

	"github.com/jackc/pgx/v5"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/category"
	"github.com/aliml92/anor/postgres/repository/product"
)

// Ensure service implements interface.
var _ anor.ProductService = (*ProductService)(nil)

type ProductService struct {
	productRepository  product.Repository
	categoryRepository category.Repository
}

func NewProductService(ps product.Repository, cs category.Repository) *ProductService {
	return &ProductService{
		productRepository:  ps,
		categoryRepository: cs,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*anor.Product, error) {
	// get product, its store, and pricing data
	p, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, anor.ErrNotFound
		}
		return nil, err
	}

	// get product attributes along with their values
	pas, err := s.productRepository.GetProductAttributesByProductID(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	productAttributes := make([]anor.ProductAttribute, len(pas))
	for index, pa := range pas {
		a := anor.ProductAttribute{
			Attribute: pa.AttributeName,
			Values:    pa.AttributeValues,
		}
		productAttributes[index] = a
	}

	// get variants, there will be at least one
	pvs, err := s.productRepository.GetProductVariantsByProductID(ctx, id)
	if err != nil {
		return nil, err
	}

	var qty int32
	variants := make([]anor.ProductVariant, len(pvs))
	for index, pv := range pvs {
		variant := anor.ProductVariant{
			ID:               pv.VariantID,
			SKU:              pv.Sku,
			Qty:              pv.Qty,
			IsCustomPriced:   pv.IsCustomPriced,
			ImageIdentifiers: pv.ImageIdentifiers,
		}

		if pv.Attributes != nil {
			var attrMap map[string]string
			if err := json.Unmarshal(pv.Attributes, &attrMap); err != nil {
				return nil, err
			}
			variant.Attributes = attrMap
		}
		variants[index] = variant

		qty = qty + pv.Qty
	}

	pm := anor.BaseProduct{
		ID:               p.Product.ID,
		StoreID:          p.Product.StoreID,
		CategoryID:       p.Product.CategoryID,
		Name:             p.Product.Name,
		Brand:            *p.Product.Brand,
		Handle:           p.Product.Handle,
		ShortInformation: p.Product.ShortInformation,
		ImageUrls:        p.Product.ImageUrls,
		Specifications:   p.Product.Specifications,
	}

	productDetails := anor.Product{
		BaseProduct: pm,
		Store: anor.Store{
			Handle:      p.Store.Handle,
			Name:        p.Store.Name,
			Description: p.Store.Description,
			CreatedAt:   p.Store.CreatedAt.Time,
		},
		Pricing:         anor.ProductPricing(p.ProductPricing),
		ProductVariants: variants,
		Attributes:      productAttributes,
		LeftCount:       int(qty),
	}

	return &productDetails, nil
}

func (s *ProductService) GetProductsByLeafCategoryID(
	ctx context.Context,
	categoryID int32,
	params anor.GetProductsByCategoryParams,
) ([]anor.Product, int64, error) {
	var (
		products []anor.Product
		count    int64
	)

	rp, err := s.productRepository.GetProductsByLeafCategoryID(ctx, categoryID, params)
	if err != nil {
		return products, count, err
	}

	products = make([]anor.Product, len(rp))
	for i, v := range rp {
		p := anor.Product{
			BaseProduct: anor.BaseProduct{
				ID:        v.ID,
				Name:      v.Name,
				Handle:    v.Handle,
				ImageUrls: v.ImageUrls,
			},
			Pricing: anor.ProductPricing{
				BasePrice:       v.BasePrice,
				Discount:        v.Discount,
				DiscountedPrice: v.DiscountedPrice,
				CurrencyCode:    v.CurrencyCode,
			},
		}
		products[i] = p
	}

	if len(rp) > 0 {
		count = rp[0].Count
	}

	return products, count, nil
}

func (s *ProductService) GetProductsByNonLeafCategoryID(
	ctx context.Context,
	categoryID int32,
	params anor.GetProductsByCategoryParams,
) ([]anor.Product, int64, error) {
	var (
		products []anor.Product
		count    int64
	)
	leafCategoryIDs, err := s.categoryRepository.GetLeafCategoryIDs(ctx, categoryID)
	if err != nil {
		return products, count, err
	}

	// get products by leaf category ids
	rp, err := s.productRepository.GetProductsByLeafCategoryIDs(ctx, leafCategoryIDs[0], params)
	if err != nil {
		return products, count, err
	}

	products = make([]anor.Product, len(rp))
	for i, v := range rp {
		p := anor.Product{
			BaseProduct: anor.BaseProduct{
				ID:        v.ID,
				Name:      v.Name,
				Handle:    v.Handle,
				ImageUrls: v.ImageUrls,
			},
			Pricing: anor.ProductPricing{
				BasePrice:       v.BasePrice,
				Discount:        v.Discount,
				DiscountedPrice: v.DiscountedPrice,
				CurrencyCode:    v.CurrencyCode,
			},
		}
		products[i] = p
	}

	if len(rp) > 0 {
		count = rp[0].Count
	}

	return products, count, nil
}

func (s *ProductService) GetProductsByCategory(
	ctx context.Context,
	category anor.Category,
	params anor.GetProductsByCategoryParams,
) ([]anor.Product, int64, error) {

	var (
		products []anor.Product
		count    int64
	)

	if category.IsLeaf() {
		rp, err := s.productRepository.GetProductsByLeafCategoryID(ctx, category.ID, params)
		if err != nil {
			return products, count, err
		}
		products = make([]anor.Product, len(rp))
		for i, v := range rp {
			p := anor.Product{
				BaseProduct: anor.BaseProduct{
					ID:        v.ID,
					Name:      v.Name,
					Handle:    v.Handle,
					ImageUrls: v.ImageUrls,
				},
				Pricing: anor.ProductPricing{
					BasePrice:       v.BasePrice,
					Discount:        v.Discount,
					DiscountedPrice: v.DiscountedPrice,
					CurrencyCode:    v.CurrencyCode,
				},
			}
			products[i] = p
		}

		if len(rp) > 0 {
			count = rp[0].Count
		}
	} else {
		leafCategoryIDs, err := s.categoryRepository.GetLeafCategoryIDs(ctx, category.ID)
		if err != nil {
			return products, count, err
		}

		// get products by leaf category ids
		rp, err := s.productRepository.GetProductsByLeafCategoryIDs(ctx, leafCategoryIDs[0], params)
		if err != nil {
			return products, count, err
		}

		products = make([]anor.Product, len(rp))
		for i, v := range rp {
			p := anor.Product{
				BaseProduct: anor.BaseProduct{
					ID:        v.ID,
					Name:      v.Name,
					Handle:    v.Handle,
					ImageUrls: v.ImageUrls,
				},
				Pricing: anor.ProductPricing{
					BasePrice:       v.BasePrice,
					Discount:        v.Discount,
					DiscountedPrice: v.DiscountedPrice,
					CurrencyCode:    v.CurrencyCode,
				},
			}
			products[i] = p
		}

		if len(rp) > 0 {
			count = rp[0].Count
		}
	}

	return products, count, nil
}

func (s *ProductService) GetProductBrandsByCategory(ctx context.Context, category anor.Category) ([]string, error) {
	var brands []string

	if category.IsLeaf() {
		b, err := s.productRepository.GetProductBrandsByCategoryID(ctx, category.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return brands, anor.ErrNotFound
			}

			return brands, err
		}

		brands = b[0]
	} else {
		leafCategoryIDs, err := s.categoryRepository.GetLeafCategoryIDs(ctx, category.ID)
		if err != nil {
			return brands, err
		}

		b, err := s.productRepository.GetProductBrandsByCategoryIDs(ctx, leafCategoryIDs[0])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return brands, anor.ErrNotFound
			}

			return brands, err
		}

		brands = b[0]
	}

	return brands, nil
}

func (s *ProductService) GetMinMaxPricesByCategory(ctx context.Context, category anor.Category) ([2]decimal.Decimal, error) {
	var minmax [2]decimal.Decimal

	if category.IsLeaf() {
		p, err := s.productRepository.GetMinMaxPricesByCategoryID(ctx, category.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return minmax, anor.ErrNotFound
			}

			return minmax, err
		}

		minmax[0] = p.MinPrice
		minmax[1] = p.MaxPrice
	} else {
		leafCategoryIDs, err := s.categoryRepository.GetLeafCategoryIDs(ctx, category.ID)
		if err != nil {
			return minmax, err
		}

		p, err := s.productRepository.GetMinMaxPricesByCategoryIDs(ctx, leafCategoryIDs[0])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return minmax, anor.ErrNotFound
			}

			return minmax, err
		}

		minmax[0] = p.MinPrice
		minmax[1] = p.MaxPrice
	}

	return minmax, nil
}

func (s *ProductService) GetNewArrivals(ctx context.Context, limit int) ([]anor.Product, error) {
	p, err := s.productRepository.GetProductsByCreatedAtDesc(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	products := make([]anor.Product, len(p))
	for i, v := range p {
		products[i] = anor.Product{
			BaseProduct: anor.BaseProduct{
				ID:         v.Product.ID,
				StoreID:    v.Product.StoreID,
				CategoryID: v.Product.CategoryID,
				Name:       v.Product.Name,
				Brand:      *v.Product.Brand,
				Handle:     v.Product.Handle,
				ImageUrls:  v.Product.ImageUrls,
				CreatedAt:  v.Product.CreatedAt.Time,
			},
			Pricing: anor.ProductPricing{
				BasePrice:       v.ProductPricing.BasePrice,
				Discount:        v.ProductPricing.Discount,
				DiscountedPrice: v.ProductPricing.DiscountedPrice,
				CurrencyCode:    v.ProductPricing.CurrencyCode,
				IsOnSale:        v.ProductPricing.IsOnSale,
			},
		}
	}

	return products, nil
}
