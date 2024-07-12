package postgres

import (
	"context"
	"errors"
	"github.com/aliml92/anor/cache"
	"github.com/samber/oops"

	"github.com/jackc/pgx/v5"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/category"
)

var _ anor.CategoryService = (*CategoryService)(nil)

type CategoryService struct {
	categoryRepository category.Repository
	categoryCache      cache.CategoryCache
}

func NewCategoryService(cr category.Repository, cache cache.CategoryCache) *CategoryService {
	return &CategoryService{
		categoryRepository: cr,
		categoryCache:      cache,
	}
}

func (s *CategoryService) GetCategory(ctx context.Context, categoryID int32) (anor.Category, error) {
	var c anor.Category
	c, err := s.categoryCache.GetCachedCategory(ctx, categoryID)
	if err != nil {
		if errors.Is(err, cache.ErrKeyNotExist) {
			rc, err := s.categoryRepository.GetCategory(ctx, categoryID) // retrieved category
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return c, anor.ErrNotFound
				}
				return c, oops.Errorf("failed to fetch category: %v", err)
			}

			c.ID = rc.ID
			c.Category = rc.Category
			c.Handle = rc.Handle
			c.IsLeafCategory = rc.IsLeafCategory

			if rc.ParentID != nil {
				c.ParentID = *rc.ParentID
			}

			if err := s.categoryCache.CacheCategory(ctx, categoryID, c, 0); err != nil {
				return c, oops.Errorf("failed to cache category: %v", err)
			}

			return c, nil
		}

		return c, err
	}

	return c, nil
}

func (s *CategoryService) GetAncestorCategories(ctx context.Context, categoryID int32) ([]anor.Category, error) {
	// get product's category along with ancestor categories
	ac, err := s.categoryRepository.GetCategoryWithAncestors(ctx, categoryID)
	if err != nil {
		return nil, oops.Errorf("failed to fetch category with ancestors: %v", err)
	}

	var categories []anor.Category
	for _, c := range ac {
		// exclude category from its ancestors
		if c.ID != categoryID {
			cat := anor.Category{
				ID:       c.ID,
				Category: c.Category,
				Handle:   c.Handle,
			}
			if c.ParentID != nil {
				cat.ParentID = *c.ParentID
			}
			categories = append(categories, cat)
		}
	}

	return categories, nil
}

func (s *CategoryService) GetSiblingCategories(ctx context.Context, categoryID int32) ([]anor.Category, error) {
	sc, err := s.categoryRepository.GetCategoryWithSiblings(ctx, categoryID)
	if err != nil {
		return nil, oops.Errorf("failed to fetch category with siblings: %v", err)
	}

	var categories []anor.Category
	for _, c := range sc {
		// exclude category from its ancestors
		if c.ID != categoryID {
			cat := anor.Category{
				ID:       c.ID,
				Category: c.Category,
				Handle:   c.Handle,
				ParentID: *c.ParentID,
			}
			categories = append(categories, cat)
		}
	}

	return categories, nil
}

func (s *CategoryService) GetChildCategories(ctx context.Context, categoryID int32) ([]anor.Category, error) {
	ac, err := s.categoryRepository.GetChildCategories(ctx, &categoryID)
	if err != nil {
		return nil, oops.Errorf("failed to fetch child categories: %v", err)
	}

	categories := make([]anor.Category, len(ac))
	for index, c := range ac {
		cat := anor.Category{
			ID:       c.ID,
			Category: c.Category,
			Handle:   c.Handle,
			ParentID: *c.ParentID,
		}
		categories[index] = cat
	}

	return categories, nil
}

func (s *CategoryService) GetRootCategories(ctx context.Context) ([]anor.Category, error) {
	rc, err := s.categoryRepository.GetRootCategories(ctx)
	if err != nil {
		return nil, oops.Errorf("failed to fetch root categories: %v", err)
	}

	categories := make([]anor.Category, len(rc))
	for index, c := range rc {
		cat := anor.Category{
			ID:       c.ID,
			Category: c.Category,
			Handle:   c.Handle,
		}
		categories[index] = cat
	}

	return categories, nil
}

func (s *CategoryService) GetCategoryHierarchy(ctx context.Context, c anor.Category) (anor.CategoryHierarchy, error) {
	var h anor.CategoryHierarchy
	h, err := s.categoryCache.GetCachedCategoryHierarchy(ctx, c.ID)
	if err != nil {
		if errors.Is(err, cache.ErrKeyNotExist) {
			if c.IsLeaf() {
				siblingCategories, err := s.GetSiblingCategories(ctx, c.ID)
				if err != nil {
					return h, oops.Errorf("failed to fetch sibling categories: %v", err)
				}
				h.SiblingCategories = siblingCategories
			} else {
				childCategories, err := s.GetChildCategories(ctx, c.ID)
				if err != nil {
					return h, oops.Errorf("failed to fetch child categories: %v", err)
				}
				h.ChildCategories = childCategories
			}

			if !c.IsRoot() {
				ancestorCategories, err := s.GetAncestorCategories(ctx, c.ID)
				if err != nil {
					return h, oops.Errorf("failed to fetch ancestor categories: %v", err)
				}
				h.AncestorCategories = ancestorCategories
			}

			if err := s.categoryCache.CacheCategoryHierarchy(ctx, c.ID, h, 0); err != nil {
				return h, oops.Errorf("failed to cache category hierarchy: %v", err)
			}

			return h, nil
		}

		return h, oops.Errorf("failed to get cached category hierarchy: %v", err)
	}

	return h, nil
}

func (s *CategoryService) GetCategoryWithAncestors(ctx context.Context, categoryID int32) (anor.CategoryWithAncestors, error) {
	var ca anor.CategoryWithAncestors
	ca, err := s.categoryCache.GetCachedCategoryWithAncestors(ctx, categoryID)
	if err != nil {
		if errors.Is(err, cache.ErrKeyNotExist) {
			cat, err := s.GetCategory(ctx, categoryID)
			if err != nil {
				return ca, oops.Errorf("failed to fetch category: %v", err)
			}
			ca.Category = cat

			if !cat.IsRoot() {
				ancestorCategories, err := s.GetAncestorCategories(ctx, categoryID)
				if err != nil {
					return ca, oops.Errorf("failed to fetch ancestor categories: %v", err)
				}
				ca.AncestorCategories = ancestorCategories
			}

			if err := s.categoryCache.CacheCategoryWithAncestors(ctx, categoryID, ca, 0); err != nil {
				return ca, oops.Errorf("failed to cache category with ancestors: %v", err)
			}

			return ca, nil
		}

		return ca, oops.Errorf("failed to get cached category with ancestors: %v", err)
	}

	return ca, nil
}
