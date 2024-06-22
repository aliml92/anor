package cache

import (
	"context"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/cache"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

var _ cache.CategoryCache = (*CategoryCacher)(nil)

const (
	basePrefixCategory              = "keys:app:category:"
	basePrefixCategoryHierarchy     = "keys:app:category:hierarchy"
	basePrefixCategoryWithAncestors = "keys:app:category:with:ancestors"
)

type CategoryCacher struct {
	client *redis.Client
}

func NewCategoryCacher(client *redis.Client) *CategoryCacher {
	return &CategoryCacher{client: client}
}

func (c *CategoryCacher) CacheCategory(ctx context.Context, categoryID int32, category anor.Category, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%d", basePrefixCategory, categoryID)
	b, err := msgpack.Marshal(&category)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, key, b, ttl).Err()
	return err
}

func (c *CategoryCacher) CacheCategoryHierarchy(
	ctx context.Context,
	categoryID int32,
	hierarchy anor.CategoryHierarchy,
	ttl time.Duration,
) error {
	key := fmt.Sprintf("%s:%d", basePrefixCategoryHierarchy, categoryID)
	b, err := msgpack.Marshal(&hierarchy)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, key, b, ttl).Err()
	return err
}

func (c *CategoryCacher) CacheCategoryWithAncestors(
	ctx context.Context,
	categoryID int32,
	cwa anor.CategoryWithAncestors,
	ttl time.Duration,
) error {
	key := fmt.Sprintf("%s:%d", basePrefixCategoryWithAncestors, categoryID)
	b, err := msgpack.Marshal(&cwa)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, key, b, ttl).Err()
	return err
}

func (c *CategoryCacher) GetCachedCategory(ctx context.Context, categoryID int32) (anor.Category, error) {
	var cat anor.Category
	key := fmt.Sprintf("%s:%d", basePrefixCategory, categoryID)
	cmd := c.client.Get(ctx, key)
	err := cmd.Err()
	if err != nil {
		if err == redis.Nil {
			return cat, cache.ErrKeyNotExist
		}
		return cat, err
	}

	b, err := cmd.Bytes()
	if err != nil {
		return cat, err
	}

	err = msgpack.Unmarshal(b, &cat)
	if err != nil {
		return cat, err
	}

	return cat, nil
}

func (c *CategoryCacher) GetCachedCategoryHierarchy(ctx context.Context, categoryID int32) (anor.CategoryHierarchy, error) {
	var h anor.CategoryHierarchy
	key := fmt.Sprintf("%s:%d", basePrefixCategoryHierarchy, categoryID)
	cmd := c.client.Get(ctx, key)
	err := cmd.Err()
	if err != nil {
		if err == redis.Nil {
			return h, cache.ErrKeyNotExist
		}
		return h, err
	}

	b, err := cmd.Bytes()
	if err != nil {
		return h, err
	}

	err = msgpack.Unmarshal(b, &h)
	if err != nil {
		return h, err
	}

	return h, nil
}

func (c *CategoryCacher) GetCachedCategoryWithAncestors(ctx context.Context, categoryID int32) (anor.CategoryWithAncestors, error) {
	var ca anor.CategoryWithAncestors
	key := fmt.Sprintf("%s:%d", basePrefixCategoryWithAncestors, categoryID)
	cmd := c.client.Get(ctx, key)
	err := cmd.Err()
	if err != nil {
		if err == redis.Nil {
			return ca, cache.ErrKeyNotExist
		}
		return ca, err
	}

	b, err := cmd.Bytes()
	if err != nil {
		return ca, err
	}

	err = msgpack.Unmarshal(b, &ca)
	if err != nil {
		return ca, err
	}

	return ca, nil
}
