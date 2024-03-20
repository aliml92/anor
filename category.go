package anor

import "context"

type Category struct {
	ID       int32
	Category string
	Slug     string
	ParentID int32

	IsLeafCategory bool
}

func (c Category) IsLeaf() bool {
	return c.IsLeafCategory
}

func (c Category) IsRoot() bool {
	if c.ParentID == 0 {
		return true
	}

	return false
}

type CategoryService interface {
	GetCategory(ctx context.Context, id int32) (Category, error)
	GetAncestorCategories(ctx context.Context, id int32) ([]Category, error)
	GetSiblingCategories(ctx context.Context, id int32) ([]Category, error)
	GetChildCategories(ctx context.Context, id int32) ([]Category, error)
}
