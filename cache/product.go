package cache

import (
	"context"
	"time"
)

type ProductCache interface {
	CacheProductBrandsByCategoryID(ctx context.Context, categoryID int32, brands []string, ttl time.Duration) error
	GetCachedProductBrandsByCategoryID(ctx context.Context, categoryID int32) ([]string, error)
}
