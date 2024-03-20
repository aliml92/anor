package typesense

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aliml92/go-typesense/typesense"

	"github.com/aliml92/anor"
)

type Searcher struct {
	client *typesense.Client
}

func NewSearcher(c *typesense.Client) *Searcher {
	return &Searcher{client: c}
}

func (s Searcher) IndexProduct(ctx context.Context, p anor.Product) error {
	doc := map[string]interface{}{
		"id":                      strconv.Itoa(int(p.ID)),
		"name":                    p.Name,
		"category_id":             p.CategoryID,
		"brand":                   p.Brand,
		"slug":                    p.Slug,
		"base_price":              p.Pricing.BasePrice,
		"price_discounted_amount": p.Pricing.DiscountedAmount,
		"thumb_img_urls":          p.ThumbImgUrls,
		"rating":                  p.Rating.Rating,
		"num_reviews":             p.Rating.RatingCount,
		"created_at":              p.CreatedAt.Unix(),
		"updated_at":              p.UpdatedAt.Unix(),
	}

	for attr, values := range p.Attributes {
		attr = fmt.Sprintf("%s_attribute", strings.ToLower(attr))
		doc[attr] = values
	}

	_, err := s.client.Documents.Create(ctx, anor.INDEXPRODUCTS, doc)
	return err
}

func (s Searcher) IndexCategory(ctx context.Context, c anor.Category) error {
	doc := map[string]interface{}{
		"id":        strconv.Itoa(int(c.ID)),
		"category":  c.Category,
		"slug":      c.Slug,
		"parent_id": c.ParentID,
		"is_leaf":   c.IsLeaf(),
	}

	_, err := s.client.Documents.Create(ctx, anor.INDEXCATEGORIES, doc)
	return err
}

func (s Searcher) IndexSellerStore(ctx context.Context, ss anor.SellerStore) error {
	doc := map[string]interface{}{
		"id":        strconv.Itoa(int(ss.ID)),
		"name":      ss.Name,
		"public_id": ss.PublicID,
	}

	_, err := s.client.Documents.Create(ctx, anor.INDEXSELLERSTORES, doc)
	return err
}
