package postgres

import (
	"context"
	"encoding/json"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/featured_promotion"
	"github.com/samber/oops"
)

var _ anor.FeaturedPromotionService = (*FeaturedPromotionService)(nil)

type FeaturedPromotionService struct {
	featuredPromotionRepository featured_promotion.Repository
}

func NewFeaturedPromotionService(fpr featured_promotion.Repository) *FeaturedPromotionService {
	return &FeaturedPromotionService{
		featuredPromotionRepository: fpr,
	}
}

func (s *FeaturedPromotionService) GetFeaturedPromotions(ctx context.Context, limit int) ([]anor.FeaturedPromotion, error) {
	fp, err := s.featuredPromotionRepository.GetFeaturedPromotions(ctx, int32(limit))
	if err != nil {
		return nil, oops.Errorf("failed to get featured promotions: %v", err)
	}

	f := make([]anor.FeaturedPromotion, len(fp))
	for i, fp := range fp {
		f[i] = anor.FeaturedPromotion{
			ID:        fp.ID,
			Title:     fp.Title,
			ImageUrl:  fp.ImageUrl,
			Type:      fp.Type,
			StartDate: fp.StartDate.Time,
			EndDate:   fp.EndDate.Time,
		}
		if fp.TargetID != nil {
			f[i].TargetID = *fp.TargetID
		}
		if fp.DisplayOrder != nil {
			f[i].DisplayOrder = *fp.DisplayOrder
		}

		err := json.Unmarshal(fp.FilterParams, &f[i].FilterParams)
		if err != nil {
			return nil, oops.Errorf("failed to unmarshal filter params: %v", err)
		}

		switch fp.Type {
		case "category":
			f[i].Category = anor.Category{
				Category: fp.TargetName,
				Handle:   fp.TargetHandle,
			}
		case "collection":
			f[i].Collection = anor.Collection{
				Name:   fp.TargetName,
				Handle: fp.TargetHandle,
			}
		}
	}

	return f, nil
}
