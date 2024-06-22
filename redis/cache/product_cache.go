package cache

import (
	"context"
	"github.com/aliml92/anor/cache"
	"github.com/redis/go-redis/v9"
	"time"
)

var _ cache.ProductCache = (*ProductCacher)(nil)

const (
	basePrefixProductBrand = "keys:app:product:brands"
)

type ProductCacher struct {
	client *redis.Client
}

func NewProductCacher(client *redis.Client) *ProductCacher {
	return &ProductCacher{client: client}
}

func (c *ProductCacher) CacheProductBrandsByCategoryID(ctx context.Context, categoryID int32, brands []string, ttl time.Duration) error {

	return nil
}

func (c *ProductCacher) GetCachedProductBrandsByCategoryID(ctx context.Context, categoryID int32) ([]string, error) {

	return nil, nil
}
