package cache

import (
	"context"
	"github.com/aliml92/anor"
	"time"
)

type CategoryCache interface {
	CacheCategory(ctx context.Context, categoryID int32, category anor.Category, ttl time.Duration) error
	CacheCategoryHierarchy(ctx context.Context, categoryID int32, hierarchy anor.CategoryHierarchy, ttl time.Duration) error
	CacheCategoryWithAncestors(ctx context.Context, categoryID int32, cwa anor.CategoryWithAncestors, ttl time.Duration) error
	GetCachedCategory(ctx context.Context, categoryID int32) (anor.Category, error)
	GetCachedCategoryHierarchy(ctx context.Context, categoryID int32) (anor.CategoryHierarchy, error)
	GetCachedCategoryWithAncestors(ctx context.Context, categoryID int32) (anor.CategoryWithAncestors, error)
}
