package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/category"
	"github.com/aliml92/anor/postgres/repository/featured_selection"
	"github.com/aliml92/anor/postgres/repository/product"
	"github.com/samber/oops"
	"strconv"
	"strings"
)

var _ anor.FeaturedSelectionService = (*FeaturedSelectionService)(nil)

type FeaturedSelectionService struct {
	featuredSelectionRepository featured_selection.Repository
	productRepository           product.Repository
	categoryRepository          category.Repository
}

func NewFeaturedSelectionService(fsr featured_selection.Repository, pr product.Repository, cr category.Repository) *FeaturedSelectionService {
	return &FeaturedSelectionService{
		featuredSelectionRepository: fsr,
		productRepository:           pr,
		categoryRepository:          cr,
	}
}

func (s *FeaturedSelectionService) ListAllActive(ctx context.Context) ([]anor.FeaturedSelection, error) {
	fs, err := s.featuredSelectionRepository.ListAllActive(ctx)
	if err != nil {
		return nil, oops.Wrap(err)
	}

	res := make([]anor.FeaturedSelection, len(fs))
	for i, f := range fs {
		res[i] = anor.FeaturedSelection{
			ID:           f.ID,
			ResourcePath: f.ResourcePath,
			ImageUrl:     f.ImageUrl,
			DisplayOrder: int(anor.Int32Value(f.DisplayOrder)),
			StartDate:    f.StartDate.Time,
			EndDate:      f.EndDate.Time,
		}

		bi, err := convertByteToMap(f.BannerInfo)
		if err != nil {
			return nil, oops.Wrap(err)
		}
		qp, err := convertByteToMap(f.QueryParams)
		if err != nil {
			return nil, oops.Wrap(err)
		}

		res[i].BannerInfo = bi
		res[i].QueryParams = qp

		var products []anor.Product

		if strings.Contains(f.ResourcePath, "/products/") {
			continue
		}

		cID, err := extractCategoryID(f.ResourcePath)
		if err != nil {
			return nil, oops.Wrap(err)
		}

		c, err := s.categoryRepository.GetCategory(ctx, cID)
		if err != nil {
			return nil, oops.Wrap(err)
		}
		if c.IsLeafCategory {
			rp, err := s.productRepository.GetProductsByLeafCategoryID(ctx, c.ID, anor.ListByCategoryParams{
				Paging: anor.Paging{
					Page:     1,
					PageSize: 10,
				},
			})
			if err != nil {
				return nil, err
			}
			products = make([]anor.Product, len(rp))
			for j, v := range rp {
				p := anor.Product{
					ID:        v.ID,
					Name:      v.Name,
					Handle:    v.Handle,
					ImageUrls: v.ImageUrls,
					Pricing: anor.ProductPricing{
						BasePrice:       v.BasePrice,
						Discount:        v.Discount,
						DiscountedPrice: v.DiscountedPrice,
						CurrencyCode:    v.CurrencyCode,
					},
				}
				products[j] = p
			}

		} else {
			leafCategoryIDs, err := s.categoryRepository.GetLeafCategoryIDs(ctx, c.ID)
			if err != nil {
				return nil, oops.Wrap(err)
			}

			// get products by leaf category ids
			rp, err := s.productRepository.GetProductsByLeafCategoryIDs(ctx, leafCategoryIDs[0], anor.ListByCategoryParams{
				Paging: anor.Paging{
					Page:     1,
					PageSize: 10,
				},
			})
			if err != nil {
				return nil, oops.Wrap(err)
			}

			products = make([]anor.Product, len(rp))
			for j, v := range rp {
				p := anor.Product{
					ID:        v.ID,
					Name:      v.Name,
					Handle:    v.Handle,
					ImageUrls: v.ImageUrls,
					Pricing: anor.ProductPricing{
						BasePrice:       v.BasePrice,
						Discount:        v.Discount,
						DiscountedPrice: v.DiscountedPrice,
						CurrencyCode:    v.CurrencyCode,
					},
				}
				products[j] = p
			}
		}

		res[i].Products = products
	}

	return res, nil
}

func convertByteToMap(b []byte) (map[string]string, error) {
	if b == nil {
		return nil, nil
	}
	m := make(map[string]string)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func extractCategoryID(handle string) (int32, error) {
	lastIndex := strings.LastIndex(handle, "-")
	if lastIndex == -1 {
		// assume handle is number string
		id, err := strconv.Atoi(handle)
		if err != nil {
			return 0, fmt.Errorf("failed to convert ID to integer: %s", handle)
		}
		return int32(id), nil
	}

	if lastIndex == len(handle)-1 {
		return 0, fmt.Errorf("invalid handle format: %s", handle)
	}

	// Extract the ID part after the last "-"
	idStr := handle[lastIndex+1:]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert ID to integer: %s", idStr)
	}

	return int32(id), nil
}
