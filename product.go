package anor

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type ProductStatus string

const (
	ProductStatusDraft           ProductStatus = "Draft"
	ProductStatusPendingApproval ProductStatus = "PendingApproval"
	ProductStatusPublished       ProductStatus = "Published"
)

type SortParam string

const (
	SortParamPopular        SortParam = "popular"
	SortParamPriceHighToLow SortParam = "price_high_to_low"
	SortParamPriceLowToHigh SortParam = "price_low_to_high"
	SortParamHighestRated   SortParam = "highest_rated"
	SortParamNewArrivals    SortParam = "new_arrivals"
	SortParamBestSellers    SortParam = "best_sellers"
)

type FilterParam struct {
	Rating    int
	PriceFrom decimal.Decimal
	PriceTo   decimal.Decimal
	Brands    []string
	Colors    []string
	Sizes     []string
}

func (p FilterParam) IsZero() bool {
	return p.Rating == 0 &&
		p.PriceFrom.IsZero() &&
		p.PriceTo.IsZero() &&
		len(p.Brands) == 0
}

type Paging struct {
	Page     int
	PageSize int
}

type BaseProduct struct {
	ID               int64
	StoreID          int32
	CategoryID       int32
	Name             string
	Brand            string
	Handle           string
	ShortInformation []string
	ImageUrls        map[int]string
	Specifications   map[string]string // product item specifications, e.g. Ebay
	Status           ProductStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Product struct {
	BaseProduct
	Store           Store
	Category        Category
	ProductVariants []ProductVariant
	Pricing         ProductPricing
	Rating          Rating
	Reviews         []ProductRating // Product Reviews
	SoldCount       int
	LeftCount       int

	Attributes []ProductAttribute
}

type ProductAttribute struct {
	Attribute string
	Values    []string
}

type ProductVariant struct {
	ID               int64                 `json:"id"`
	SKU              string                `json:"sku"`
	Qty              int32                 `json:"quantity"`
	IsCustomPriced   bool                  `json:"-"` // TODO: implement
	ImageIdentifiers []int16               `json:"-"` // TODO: implement
	Attributes       map[string]string     `json:"attributes"`
	Pricing          ProductVariantPricing `json:"-"` // TODO: implement
	CreatedAt        time.Time             `json:"-"`
	UpdatedAt        time.Time             `json:"-"`
}

type ProductPricing struct {
	ProductID       int64
	BasePrice       decimal.Decimal
	CurrencyCode    string
	Discount        decimal.Decimal
	DiscountedPrice decimal.Decimal
	IsOnSale        bool
}

type ProductVariantPricing struct {
	VariantID       int64
	BasePrice       decimal.Decimal
	CurrencyCode    string
	Discount        decimal.Decimal
	DiscountedPrice decimal.Decimal
	IsOnSale        bool
}

type ProductRating struct {
	ID        int64
	ProductID int64
	UserID    int64
	Rating    int16
	Review    string
	ImageUrls []string
	CreatedAt time.Time
}

type Rating struct {
	Rating      float32
	RatingCount int
}

type GetProductsByCategoryParams struct {
	Filter FilterParam
	Sort   SortParam
	Paging Paging
}

type ProductService interface {
	GetProduct(ctx context.Context, id int64) (*Product, error)
	GetProductsByLeafCategoryID(ctx context.Context, categoryID int32, params GetProductsByCategoryParams) ([]Product, int64, error)
	GetProductsByNonLeafCategoryID(ctx context.Context, categoryID int32, params GetProductsByCategoryParams) ([]Product, int64, error)
	GetProductsByCategory(ctx context.Context, category Category, params GetProductsByCategoryParams) ([]Product, int64, error)
	GetProductBrandsByCategory(ctx context.Context, category Category) ([]string, error)
	GetMinMaxPricesByCategory(ctx context.Context, category Category) ([2]decimal.Decimal, error)

	GetNewArrivals(ctx context.Context, limit int) ([]Product, error)
}
