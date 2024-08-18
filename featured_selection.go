package anor

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

// FeaturedSelection represents a highlighted item or group of items.
// It can be used in various contexts, such as:
// - Banner carousel items on the home page
// - Highlighted categories of products
// - Curated search results with specific filters and sorting
// - Single product promotions
// - External links to campaigns or special offers
//
// FeaturedSelections are versatile and can represent either a group of products,
// a single product, or an external link, depending on the context and configuration.
type FeaturedSelection struct {
	ID           int64
	ResourcePath string
	BannerInfo   map[string]string
	ImageUrl     string
	QueryParams  map[string]string
	StartDate    time.Time
	EndDate      time.Time
	DisplayOrder int
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz

	Products []Product
}

type FeaturedSelectionService interface {
	ListAllActive(ctx context.Context) ([]FeaturedSelection, error)
}
