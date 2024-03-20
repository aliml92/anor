package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/store/category"
)

var _ anor.CategoryService = (*CategoryService)(nil)

type CategoryService struct {
	categoryStore category.Store
}

func NewCategoryService(cs category.Store) *CategoryService {
	return &CategoryService{
		categoryStore: cs,
	}
}

func (s CategoryService) GetCategory(ctx context.Context, categoryID int32) (anor.Category, error) {
	var c anor.Category
	rc, err := s.categoryStore.GetCategory(ctx, categoryID) // retrieved category
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c, anor.ErrNotFound
		}
		return c, err
	}

	c.ID = rc.ID
	c.Category = rc.Category
	c.Slug = rc.Slug
	c.IsLeafCategory = rc.IsLeafCategory

	if rc.ParentID != nil {
		c.ParentID = *rc.ParentID
	}

	return c, nil
}

func (s CategoryService) GetAncestorCategories(ctx context.Context, categoryID int32) ([]anor.Category, error) {
	// get product's category along with ancestor categories
	ac, err := s.categoryStore.GetCategoryWithAncestors(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	var categories []anor.Category
	for _, c := range ac {
		// exclude category from its ancestors
		if c.ID != categoryID {
			cat := anor.Category{
				ID:       c.ID,
				Category: c.Category,
				Slug:     c.Slug,
			}
			if c.ParentID != nil {
				cat.ParentID = *c.ParentID
			}
			categories = append(categories, cat)
		}
	}

	return categories, nil
}

func (s CategoryService) GetSiblingCategories(ctx context.Context, categoryID int32) ([]anor.Category, error) {
	sc, err := s.categoryStore.GetCategoryWithSiblings(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	var categories []anor.Category
	for _, c := range sc {
		// exclude category from its ancestors
		if c.ID != categoryID {
			cat := anor.Category{
				ID:       c.ID,
				Category: c.Category,
				Slug:     c.Slug,
				ParentID: *c.ParentID,
			}
			categories = append(categories, cat)
		}
	}

	return categories, nil
}

func (s CategoryService) GetChildCategories(ctx context.Context, categoryID int32) ([]anor.Category, error) {
	ac, err := s.categoryStore.GetChildCategories(ctx, &categoryID)
	if err != nil {
		return nil, err
	}

	categories := make([]anor.Category, len(ac))
	for index, c := range ac {
		cat := anor.Category{
			ID:       c.ID,
			Category: c.Category,
			Slug:     c.Slug,
			ParentID: *c.ParentID,
		}
		categories[index] = cat
	}

	return categories, nil
}
