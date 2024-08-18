package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aliml92/anor/relation"
	"github.com/samber/oops"
	"github.com/shopspring/decimal"
	"os"
	"strconv"
	"strings"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/category"
	"github.com/aliml92/anor/postgres/repository/product"
	"github.com/jackc/pgx/v5"
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

func (s *ProductService) Get(ctx context.Context, id int64, with relation.Set) (*anor.Product, error) {
	var (
		withStore             = with.Has(relation.Store)
		withPricing           = with.Has(relation.Pricing)
		withProductAttributes = with.Has(relation.ProductAttribute)
		withProductVariants   = with.Has(relation.ProductVariant)
	)

	// first get product with its one-to-one relations
	var (
		p       product.Product
		store   product.Store
		pricing product.ProductPricing
	)
	if withStore && withPricing {
		// get product, its store, and pricing data
		psp, err := s.productRepository.GetProductWithStoreAndPricingByID(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, anor.ErrProductNotFound
			}
			return nil, err
		}
		p = psp.Product
		store = psp.Store
		pricing = psp.ProductPricing
	} else if withStore {
		// TODO: currently the app do not use this case
		panic("implement me")
	} else if withPricing {
		// TODO: currently the app do not use this case
		panic("implement me")
	} else {
		pr, err := s.productRepository.GetProductByID(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, anor.ErrProductNotFound
			}
			return nil, err
		}
		p = *pr
	}

	var attributes []anor.ProductAttribute
	if withProductAttributes {
		pas, err := s.productRepository.GetProductAttributesByProductID(ctx, id)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, oops.Wrap(err)
		}

		attributes = make([]anor.ProductAttribute, len(pas))
		for index, pa := range pas {
			a := anor.ProductAttribute{
				Attribute: pa.AttributeName,
				Values:    pa.AttributeValues,
			}
			attributes[index] = a
		}
	}

	var (
		variants []anor.ProductVariant
		qty      int
	)
	if withProductVariants {
		// get variants, there will be at least one
		pvs, err := s.productRepository.GetProductVariantsByProductID(ctx, id)
		if err != nil {
			return nil, oops.Wrap(err)
		}

		variants = make([]anor.ProductVariant, len(pvs))
		for index, pv := range pvs {
			variant := anor.ProductVariant{
				ID:               pv.VariantID,
				SKU:              pv.Sku,
				Qty:              pv.Qty,
				IsCustomPriced:   pv.IsCustomPriced, // TODO: if true get variant pricing
				ImageIdentifiers: pv.ImageIdentifiers,
				CreatedAt:        pv.CreatedAt.Time,
				UpdatedAt:        pv.UpdatedAt.Time,
			}

			if pv.Attributes != nil {
				var attrMap map[string]string
				if err := json.Unmarshal(pv.Attributes, &attrMap); err != nil {
					return nil, oops.Wrap(err)
				}
				variant.Attributes = attrMap
			}
			variants[index] = variant

			qty = qty + int(pv.Qty)
		}

	}

	res := anor.Product{
		ID:               p.ID,
		StoreID:          p.StoreID,
		CategoryID:       p.CategoryID,
		Name:             p.Name,
		Brand:            anor.StringValue(p.Brand),
		Handle:           p.Handle,
		ShortInformation: p.ShortInformation,
		ImageUrls:        p.ImageUrls,
		Specifications:   p.Specifications,
		Store: anor.Store{
			Handle:      store.Handle,
			Name:        store.Name,
			Description: store.Description,
			CreatedAt:   store.CreatedAt.Time,
		},
		Pricing:         anor.ProductPricing(pricing),
		ProductVariants: variants,
		Attributes:      attributes,
		LeftCount:       qty,
	}

	return &res, nil
}

func (s *ProductService) ListByCategory(ctx context.Context, category anor.Category, params anor.ListByCategoryParams) ([]anor.Product, int64, error) {

	var (
		products []anor.Product
		count    int64
		//withPricing = params.With.Has(relation.Pricing)
	)

	if category.IsLeaf() {
		rp, err := s.productRepository.GetProductsByLeafCategoryID(ctx, category.ID, params)
		if err != nil {
			return products, count, err
		}
		products = make([]anor.Product, len(rp))
		for i, v := range rp {
			p := anor.Product{
				ID:        v.ID,
				Name:      v.Name,
				Handle:    v.Handle,
				ImageUrls: v.ImageUrls,
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
			return products, count, oops.Errorf("failed to get leaf category ids: %v", err)
		}

		// get products by leaf category ids
		rp, err := s.productRepository.GetProductsByLeafCategoryIDs(ctx, leafCategoryIDs[0], params)
		if err != nil {
			return products, count, oops.Errorf("failed to get products by leaf category ids: %v", err)
		}

		products = make([]anor.Product, len(rp))
		for i, v := range rp {
			p := anor.Product{
				ID:        v.ID,
				Name:      v.Name,
				Handle:    v.Handle,
				ImageUrls: v.ImageUrls,
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

func (s *ProductService) ListAllBrandsByCategory(ctx context.Context, category anor.Category) ([]string, error) {
	var brands []string

	if category.IsLeaf() {
		b, err := s.productRepository.GetProductBrandsByCategoryID(ctx, category.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return brands, anor.ErrNotFound
			}

			return brands, oops.Errorf("failed to get product brands: %v", err)
		}

		brands = b[0]
	} else {
		leafCategoryIDs, err := s.categoryRepository.GetLeafCategoryIDs(ctx, category.ID)
		if err != nil {
			return brands, oops.Errorf("failed to get leaf category ids: %v", err)
		}

		b, err := s.productRepository.GetProductBrandsByCategoryIDs(ctx, leafCategoryIDs[0])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return brands, anor.ErrNotFound
			}

			return brands, oops.Errorf("failed to get product brands: %v", err)
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

			return minmax, oops.Errorf("failed to get min and max prices by category id: %v", err)
		}

		minmax[0] = p.MinPrice
		minmax[1] = p.MaxPrice
	} else {
		leafCategoryIDs, err := s.categoryRepository.GetLeafCategoryIDs(ctx, category.ID)
		if err != nil {
			return minmax, oops.Errorf("failed to get leaf category ids: %v", err)
		}

		p, err := s.productRepository.GetMinMaxPricesByCategoryIDs(ctx, leafCategoryIDs[0])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return minmax, anor.ErrNotFound
			}

			return minmax, oops.Errorf("failed to get min and max prices by category ids: %v", err)
		}

		minmax[0] = p.MinPrice
		minmax[1] = p.MaxPrice
	}

	return minmax, nil
}

func (s *ProductService) ListPopularProducts(ctx context.Context, params anor.PopularProductListParams) ([]anor.Product, error) {
	// TODO: get real popular products
	ids, err := getPopularProductIDs()
	if err != nil {
		return nil, oops.Wrap(err)
	}
	p, err := s.productRepository.ListPopularProducts(ctx, 10, 0, ids)
	if err != nil {
		return nil, oops.Wrap(err)
	}
	products := make([]anor.Product, len(p))
	for i, v := range p {
		products[i] = anor.Product{
			ID:         v.Product.ID,
			StoreID:    v.Product.StoreID,
			CategoryID: v.Product.CategoryID,
			Name:       v.Product.Name,
			Brand:      *v.Product.Brand,
			Handle:     v.Product.Handle,
			ImageUrls:  v.Product.ImageUrls,
			CreatedAt:  v.Product.CreatedAt.Time,
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

func getPopularProductIDs() ([]int64, error) {
	idsStr := os.Getenv("FAKE_POPULAR_PRODUCT_IDS")
	strIDs := strings.Split(idsStr, ",")

	ids := make([]int64, 0, len(strIDs))

	for _, strID := range strIDs {
		id, err := strconv.ParseInt(strings.TrimSpace(strID), 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
