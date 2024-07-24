package anor

import (
	"context"
	"time"
)

type FeaturedPromotion struct {
	ID           int64
	Title        string
	ImageUrl     string
	Type         string
	TargetID     int32
	FilterParams map[string]string
	StartDate    time.Time
	EndDate      time.Time
	DisplayOrder int32

	Category   Category
	Collection Collection
}

type FeaturedPromotionService interface {
	GetFeaturedPromotions(ctx context.Context, limit int) ([]FeaturedPromotion, error)
}
