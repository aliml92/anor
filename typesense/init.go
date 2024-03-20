package typesense

import (
	"context"

	"github.com/aliml92/go-typesense/typesense"
	"github.com/pkg/errors"

	"github.com/aliml92/anor"
)

func InitCollections(ctx context.Context, c *typesense.Client) error {
	if err := createProductsCollection(ctx, c); err != nil {
		return errors.Wrap(err, "failed to create `products` collection")
	}
	if err := createCategoriesCollection(ctx, c); err != nil {
		return errors.Wrap(err, "failed to create `categories` collection")
	}
	if err := createSellerStoresCollection(ctx, c); err != nil {
		return errors.Wrap(err, "failed to create `sellerstores` collection")
	}

	return nil
}

func createProductsCollection(ctx context.Context, c *typesense.Client) error {
	productsSchema := &typesense.CollectionSchema{
		Name:                anor.INDEXPRODUCTS,
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
				Name:  "price_discounted_amount",
				Type:  "float",
				Facet: typesense.Bool(true),
				Index: typesense.Bool(true),
			},
			{
				Name:     "slug",
				Type:     "string",
				Facet:    typesense.Bool(false),
				Index:    typesense.Bool(false),
				Optional: typesense.Bool(true),
			},
			{
				Name:     "thumb_img_urls",
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
	return err
}

func createCategoriesCollection(ctx context.Context, c *typesense.Client) error {
	categoriesSchema := &typesense.CollectionSchema{
		Name:               anor.INDEXCATEGORIES,
		EnableNestedFields: typesense.Bool(false),
		Fields: []*typesense.Field{
			{
				Name:  "category",
				Type:  "string",
				Facet: typesense.Bool(false),
				Index: typesense.Bool(true),
			},
			{
				Name:     "slug",
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
	return err
}

func createSellerStoresCollection(ctx context.Context, c *typesense.Client) error {
	sellerStoreSchema := &typesense.CollectionSchema{
		Name:               anor.INDEXSELLERSTORES,
		EnableNestedFields: typesense.Bool(false),
		Fields: []*typesense.Field{
			{
				Name:  "name",
				Type:  "string",
				Facet: typesense.Bool(false),
				Index: typesense.Bool(true),
			},
			{
				Name:     "public_id",
				Type:     "string",
				Facet:    typesense.Bool(false),
				Index:    typesense.Bool(false),
				Optional: typesense.Bool(true),
			},
		},
	}

	_, err := c.Collections.Create(ctx, sellerStoreSchema)
	return err
}
