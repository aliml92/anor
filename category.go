package anor

import "context"

type Category struct {
	ID       int32
	Category string
	Handle   string
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

type CategoryHierarchy struct {
	AncestorCategories []Category `json:"ac,omitempty"`
	SiblingCategories  []Category `json:"sc,omitempty"`
	ChildCategories    []Category `json:"cc,omitempty"`
}

type CategoryWithAncestors struct {
	Category           Category   `json:"c"`
	AncestorCategories []Category `json:"ac"`
}

type CategoryService interface {
	GetCategory(ctx context.Context, id int32) (Category, error)
	GetAncestorCategories(ctx context.Context, id int32) ([]Category, error)
	GetSiblingCategories(ctx context.Context, id int32) ([]Category, error)
	GetChildCategories(ctx context.Context, id int32) ([]Category, error)
	GetRootCategories(ctx context.Context) ([]Category, error)
	GetCategoryHierarchy(ctx context.Context, category Category) (CategoryHierarchy, error)
	GetCategoryWithAncestors(ctx context.Context, categoryID int32) (CategoryWithAncestors, error)
}
