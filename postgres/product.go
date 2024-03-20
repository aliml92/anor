package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/store/category"
	"github.com/aliml92/anor/postgres/store/product"
)

// Ensure service implements interface.
var _ anor.ProductService = (*ProductService)(nil)

type ProductService struct {
	productStore  product.Store
	categoryStore category.Store
}

func NewProductService(ps product.Store, cs category.Store) *ProductService {
	return &ProductService{
		productStore:  ps,
		categoryStore: cs,
	}
}

func (s ProductService) GetProduct(ctx context.Context, id int64) (*anor.Product, error) {
	// get product, its store, and pricing data
	p, err := s.productStore.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, anor.ErrNotFound
		}
		return nil, err
	}

	// get skus, there will be at least one
	ss, err := s.productStore.GetSKUsByProductID(ctx, id)
	if err != nil {
		return nil, err
	}

	skus := make([]anor.SKU, len(ss))
	productAttributes := make(map[string][]string)
	if ss[0].Attributes != nil {
		for index, sk := range ss {
			sku := anor.SKU{
				ID:              sk.SkuID,
				SKU:             sk.Sku,
				QuantityInStock: sk.QuantityInStock,
				HasSepPricing:   sk.HasSepPricing,
				ImageRefs:       sk.ImageRefs,
			}

			var attrMap map[string]string
			if err := json.Unmarshal(sk.Attributes, &attrMap); err != nil {
				return nil, err
			}

			sku.Attributes = attrMap
			skus[index] = sku

			pas, err := s.productStore.GetProductAttributesByProductID(ctx, id)
			if err != nil {
				return nil, err
			}
			for _, pa := range pas {
				productAttributes[pa.AttributeName] = pa.AttributeValues
			}
		}
	} else {
		for index, s := range ss {
			sku := anor.SKU{
				ID:              s.SkuID,
				SKU:             s.Sku,
				QuantityInStock: s.QuantityInStock,
				HasSepPricing:   s.HasSepPricing,
				ImageRefs:       s.ImageRefs,
			}

			skus[index] = sku
		}
	}

	pm := anor.BaseProduct{
		ID:         p.Product.ID,
		StoreID:    p.Product.StoreID,
		CategoryID: p.Product.CategoryID,
		Name:       p.Product.Name,
		Brand:      *p.Product.Brand,
		Slug:       p.Product.Slug,
		ShortInfo:  p.Product.ShortInfo,
		ImageLinks: p.Product.ImageUrls,
		Specs:      p.Product.Specs,
	}

	productDetails := anor.Product{
		BaseProduct: pm,
		Store: anor.SellerStore{
			PublicID:    p.SellerStore.PublicID,
			Name:        p.SellerStore.Name,
			Description: p.SellerStore.Description,
			CreatedAt:   p.SellerStore.CreatedAt.Time,
		},
		Pricing:    anor.ProductPricing(p.ProductPricing),
		SKUs:       skus,
		Attributes: productAttributes,
	}

	return &productDetails, nil
}

func (s ProductService) GetProductsByLeafCategoryID(
	ctx context.Context,
	categoryID int32,
	params anor.GetProductsByCategoryParams,
) ([]anor.Product, int64, error) {
	var (
		products []anor.Product
		count    int64
	)

	rp, err := s.productStore.GetProductsByLeafCategoryID(ctx, categoryID, params)
	if err != nil {
		return products, count, err
	}

	products = make([]anor.Product, len(rp))
	for i, v := range rp {
		p := anor.Product{
			BaseProduct: anor.BaseProduct{
				ID:       v.ID,
				Name:     v.Name,
				Slug:     v.Slug,
				ThumbImg: v.ThumbImg,
			},
			Pricing: anor.ProductPricing{
				BasePrice:        v.BasePrice,
				DiscountedAmount: v.DiscountedAmount,
				CurrencyCode:     v.CurrencyCode,
			},
		}
		products[i] = p
	}

	if len(rp) > 0 {
		count = rp[0].Count
	}

	return products, count, nil
}

func (s ProductService) GetProductsByNonLeafCategoryID(
	ctx context.Context,
	categoryID int32,
	params anor.GetProductsByCategoryParams,
) ([]anor.Product, int64, error) {
	var (
		products []anor.Product
		count    int64
	)
	leafCategoryIDs, err := s.categoryStore.GetLeafCategoryIDs(ctx, categoryID)
	if err != nil {
		return products, count, err
	}

	// get products by leaf category ids
	rp, err := s.productStore.GetProductsByLeafCategoryIDs(ctx, leafCategoryIDs[0], params)
	if err != nil {
		return products, count, err
	}

	products = make([]anor.Product, len(rp))
	for i, v := range rp {
		p := anor.Product{
			BaseProduct: anor.BaseProduct{
				ID:       v.ID,
				Name:     v.Name,
				Slug:     v.Slug,
				ThumbImg: v.ThumbImg,
			},
			Pricing: anor.ProductPricing{
				BasePrice:        v.BasePrice,
				DiscountedAmount: v.DiscountedAmount,
				CurrencyCode:     v.CurrencyCode,
			},
		}
		products[i] = p
	}

	if len(rp) > 0 {
		count = rp[0].Count
	}

	return products, count, nil
}
