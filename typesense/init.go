package typesense

import (
	"context"
	"github.com/aliml92/anor/search"

	"github.com/aliml92/go-typesense/typesense"
	"github.com/pkg/errors"
)

func setupQuerySuggestions(ctx context.Context, c *typesense.Client) error {
	// initialize `product_queries` schema definition
	pqSchema := &typesense.CollectionSchema{
		Name: search.INDEXPRODUCTQUERIES,
		Fields: []*typesense.Field{
			{
				Name: "q",
				Type: "string",
			},
			{
				Name: "count",
				Type: "int32",
			},
		},
	}

	// create `product_queries` collection
	_, err := c.Collections.Create(ctx, pqSchema)
	if err != nil {
		var terr *typesense.ApiError
		if errors.As(err, &terr) {
			switch terr.StatusCode {
			case 409:
				goto createAnalyticsRule
			}
		}
		return err
	}

createAnalyticsRule:
	// rule name
	pqRuleName := "product_queries_aggregation"
	pqRuleSchema := &typesense.AnalyticsRuleUpsertSchema{
		Type: "popular_queries",
		Params: struct {
			Source struct {
				Collections []string `json:"collections"`
			} `json:"source"`
			Destination struct {
				Collection string `json:"collection"`
			} `json:"destination"`
			Limit int `json:"limit"`
		}{
			Source: struct {
				Collections []string `json:"collections"`
			}{
				Collections: []string{"products"},
			},
			Destination: struct {
				Collection string `json:"collection"`
			}{
				Collection: "product_queries",
			},
			Limit: 1000,
		},
	}
	// create `popular_queries` analytics rule
	_, err = c.AnalyticsRules.Upsert(ctx, pqRuleName, pqRuleSchema)
	if err != nil {
		var terr *typesense.ApiError
		if errors.As(err, &terr) {
			switch terr.StatusCode {
			case 409:
				return nil
			}
		}
		return err
	}
	return nil
}

func InitCollections(ctx context.Context, c *typesense.Client) error {
	if err := createProductsCollection(ctx, c); err != nil {
		return errors.Wrap(err, "failed to create `products` collection")
	}
	if err := createCategoriesCollection(ctx, c); err != nil {
		return errors.Wrap(err, "failed to create `categories` collection")
	}
	if err := createStoresCollection(ctx, c); err != nil {
		return errors.Wrap(err, "failed to create `stores` collection")
	}

	if err := setupQuerySuggestions(ctx, c); err != nil {
		return errors.Wrap(err, "failed to setup `product_queries` collection")
	}

	return nil
}

func createProductsCollection(ctx context.Context, c *typesense.Client) error {
	productsSchema := &typesense.CollectionSchema{
		Name:                search.INDEXPRODUCTS,
		EnableNestedFields:  typesense.Bool(true),
		DefaultSortingField: typesense.String("updated_at"),
		Fields: []*typesense.Field{
			{
				Name:  "name",
				Type:  "string",
				Facet: typesense.Bool(false),
				Index: typesense.Bool(true),
			},
			{
				Name:  "category_id",
				Type:  "int32",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:  "brand",
				Type:  "string",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:  "base_price",
				Type:  "float",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:  "discount",
				Type:  "float",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:  "discounted_price",
				Type:  "float",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:     "handle",
				Type:     "string",
				Facet:    typesense.Bool(false),
				Index:    typesense.Bool(false),
				Optional: typesense.Bool(true),
			},
			{
				Name:     "image_urls",
				Type:     "object",
				Facet:    typesense.Bool(false),
				Index:    typesense.Bool(false),
				Optional: typesense.Bool(true),
			},
			{
				Name:  "rating",
				Type:  "float",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:  "num_reviews",
				Type:  "int32",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:  "created_at",
				Type:  "int64",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:  "updated_at",
				Type:  "int64",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:     ".*_attribute",
				Type:     "auto",
				Facet:    typesense.Bool(true),
				Index:    typesense.Bool(true),
				Optional: typesense.Bool(true),
			},
		},
	}

	_, err := c.Collections.Create(ctx, productsSchema)
	if err != nil {
		var terr *typesense.ApiError
		if errors.As(err, &terr) {
			switch terr.StatusCode {
			case 409:
				return nil
			}
		}
		return err
	}
	return nil
}

func createCategoriesCollection(ctx context.Context, c *typesense.Client) error {
	categoriesSchema := &typesense.CollectionSchema{
		Name:               search.INDEXCATEGORIES,
		EnableNestedFields: typesense.Bool(false),
		Fields: []*typesense.Field{
			{
				Name:  "category",
				Type:  "string",
				Facet: typesense.Bool(false),
				Index: typesense.Bool(true),
			},
			{
				Name:     "handle",
				Type:     "string",
				Facet:    typesense.Bool(false),
				Index:    typesense.Bool(false),
				Optional: typesense.Bool(true),
			},
			{
				Name:  "parent_id",
				Type:  "int32",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:     "is_leaf",
				Type:     "bool",
				Facet:    typesense.Bool(false),
				Index:    typesense.Bool(false),
				Optional: typesense.Bool(true),
			},
		},
	}

	_, err := c.Collections.Create(ctx, categoriesSchema)
	if err != nil {
		var terr *typesense.ApiError
		if errors.As(err, &terr) {
			switch terr.StatusCode {
			case 409:
				return nil
			}
		}
		return err
	}
	return nil
}

func createStoresCollection(ctx context.Context, c *typesense.Client) error {
	sellerStoreSchema := &typesense.CollectionSchema{
		Name:               search.INDEXSTORES,
		EnableNestedFields: typesense.Bool(false),
		Fields: []*typesense.Field{
			{
				Name:  "name",
				Type:  "string",
				Facet: typesense.Bool(false),
				Index: typesense.Bool(true),
			},
			{
				Name:     "handle",
				Type:     "string",
				Facet:    typesense.Bool(false),
				Index:    typesense.Bool(false),
				Optional: typesense.Bool(true),
			},
		},
	}

	_, err := c.Collections.Create(ctx, sellerStoreSchema)
	if err != nil {
		var terr *typesense.ApiError
		if errors.As(err, &terr) {
			switch terr.StatusCode {
			case 409:
				return nil
			}
		}
		return err
	}
	return nil
}
