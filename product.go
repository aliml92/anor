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
	ProductStatusActive          ProductStatus = "Active"
)

type SortOption string

const (
	SortOptionBestMatch      SortOption = "BestMatch"
	SortOptionPriceHighToLow SortOption = "PriceHighToLow"
	SortOptionPriceLowToHigh SortOption = "PriceLowToHigh"
	SortOptionHighestRated   SortOption = "HighestRated"
	SortOptionNewArrivals    SortOption = "NewArrivals"
	SortOptionBestSellers    SortOption = "BestSellers"
)

type FilterOption struct {
	Rating     int
	PriceRange [2]float32
	Brands     []string
	Colors     []string
	Sizes      []string
}

type BaseProduct struct {
	ID           int64
	StoreID      int32
	CategoryID   int32
	Name         string
	Brand        string
	Slug         string
	ShortInfo    []string
	ThumbImg     string // thumbnail image link
	ThumbImgUrls map[string]string
	ImageLinks   map[int]string
	Specs        map[string]string // product item specifications, e.g. Ebay
	Status       ProductStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Product struct {
	BaseProduct
	Store     SellerStore
	Category  Category
	SKUs      []SKU
	Pricing   ProductPricing
	Rating    Rating          // Rating is calculated using Reviews (in the future, using also some statistical method)
	Reviews   []ProductRating // Product Reviews
	SoldCount int
	LeftCount int

	Attributes map[string][]string // attributes with their values, e.g. color, size
}

type SKU struct {
	ID              int64
	SKU             string
	QuantityInStock int32
	HasSepPricing   bool
	ImageRefs       []int16
	Attributes      map[string]string
	Pricing         SKUPricing
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ProductPricing struct {
	ProductID        int64
	BasePrice        decimal.Decimal
	CurrencyCode     string
	DiscountLevel    decimal.Decimal
	DiscountedAmount decimal.Decimal
	IsOnSale         bool
}

type SKUPricing struct {
	SkuID            int64
	BasePrice        decimal.Decimal
	CurrencyCode     string
	DiscountLevel    decimal.Decimal
	DiscountedAmount decimal.Decimal
	IsOnSale         bool
}

type ProductRating struct {
	ID         int64
	ProductID  int64
	UserID     int64
	Rating     int16
	Review     string
	ImageLinks []string
	CreatedAt  time.Time
}

type Rating struct {
	Rating      float32
	RatingCount int
}

type GetProductsByCategoryParams struct {
	FilterOption FilterOption
	SortOption   SortOption
	Offset       int
	Limit        int
}

type ProductService interface {
	GetProduct(ctx context.Context, id int64) (*Product, error)
	GetProductsByLeafCategoryID(ctx context.Context, categoryID int32, params GetProductsByCategoryParams) ([]Product, int64, error)
	GetProductsByNonLeafCategoryID(ctx context.Context, categoryID int32, params GetProductsByCategoryParams) ([]Product, int64, error)
}
