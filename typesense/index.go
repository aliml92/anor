package typesense

import (
	"context"
	"fmt"
	"github.com/aliml92/anor/search"
	"strconv"
	"strings"

	"github.com/aliml92/anor"
)

func (ts TsSearcher) IndexProduct(ctx context.Context, p anor.Product) error {
	doc := map[string]interface{}{
		"id":               strconv.Itoa(int(p.ID)),
		"name":             p.Name,
		"category_id":      p.CategoryID,
		"brand":            p.Brand,
		"handle":           p.Handle,
		"base_price":       p.Pricing.BasePrice,
		"discount":         p.Pricing.Discount,
		"discounted_price": p.Pricing.DiscountedPrice,
		"image_urls":       p.ImageUrls,
		"rating":           p.Rating.Rating,
		"num_reviews":      p.Rating.RatingCount,
		"created_at":       p.CreatedAt.Unix(),
		"updated_at":       p.UpdatedAt.Unix(),
	}

	for _, a := range p.Attributes {
		attr := fmt.Sprintf("%s_attribute", strings.ToLower(a.Attribute))
		doc[attr] = a.Values
	}

	_, err := ts.tsClient.Documents.Create(ctx, search.INDEXPRODUCTS, doc)
	return err
}

func (ts TsSearcher) IndexCategory(ctx context.Context, c anor.Category) error {
	doc := map[string]interface{}{
		"id":        strconv.Itoa(int(c.ID)),
		"category":  c.Category,
		"handle":    c.Handle,
		"parent_id": c.ParentID,
		"is_leaf":   c.IsLeaf(),
	}

	_, err := ts.tsClient.Documents.Create(ctx, search.INDEXCATEGORIES, doc)
	return err
}

func (ts TsSearcher) IndexStore(ctx context.Context, ss anor.Store) error {
	doc := map[string]interface{}{
		"id":     strconv.Itoa(int(ss.ID)),
		"name":   ss.Name,
		"handle": ss.Handle,
	}

	_, err := ts.tsClient.Documents.Create(ctx, search.INDEXSTORES, doc)
	return err
}
