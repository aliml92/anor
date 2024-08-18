package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/pkg/utils"
	"github.com/aliml92/anor/postgres/repository/category"
	"github.com/aliml92/anor/postgres/repository/product"
	"github.com/aliml92/anor/postgres/repository/store"
	"github.com/aliml92/anor/postgres/repository/user"
	ts "github.com/aliml92/anor/typesense"
	"github.com/aliml92/go-typesense/typesense"
	"github.com/brianvoe/gofakeit"
	"github.com/gosimple/slug"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"math/rand"
	"time"
)

const (
	fakeUserPassword = "Password1@"
)

var discounts = []float32{
	0, 0, 0.02, 0, 0, 0, 0.03, 0, 0.04, 0, 0, 0.05, 0, 0, 0.07, 0, 0, 0, 0.08, 0, 0, 0, 0.09, 0,
	0.10, 0, 0.11, 0, 0.12, 0, 0.13, 0, 0.14, 0, 0.15, 0, 0.20, 0, 0, 0.25, 0, 0, 0, 0, 0, 0, 0,
	0.30, 0, 0, 0.35, 0, 0.40, 0, 0.45, 0, 0.50, 0.55, 0, 0.60, 0, 0.65, 0, 0.70, 0, 0, 0, 0, 0,
	0.15, 0, 0.12, 0, 0, 0, 0.05, 0.07, 0.20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var shortInfo = []string{
	"Trendy design for versatile wear",
	"Quality materials for durability",
	"Available in various sizes and colors",
}

type config struct {
	datasetDir     string
	database       string
	imported       string
	typesense      string
	ignoreImported bool
}

type DataImporter struct {
	cfg                  config
	database             *pgxpool.Pool
	typesense            *typesense.Client
	datasetList          []string
	datasetListSucceeded []string

	// temporary
	userRepository     user.Repository
	productRepository  product.Repository
	categoryRepository category.Repository
	storeRepository    store.Repository
	searcher           *ts.Searcher
}

func (d *DataImporter) safeImportData(ctx context.Context) (err error) {
	// Defer a function to catch any panic
	defer func() {
		if r := recover(); r != nil {
			// Convert the panic to an error
			err = fmt.Errorf("panic in importData: %v", r)

			// Write importedFilesList to file
			fileErr := updateImported(d.cfg.imported, d.datasetListSucceeded)
			if fileErr != nil {
				// If there's an error updating the file, include it in the returned error
				err = fmt.Errorf("%v; additionally, failed to update file: %v", err, fileErr)
			}
		}
	}()

	err = d.importData(ctx)
	// Write importedFilesList to file
	fileErr := updateImported(d.cfg.imported, d.datasetListSucceeded)
	if fileErr != nil {
		// If there's an error updating the file, include it in the returned error
		err = fmt.Errorf("%v; additionally, failed to update file: %v", err, fileErr)
	}

	return err
}

func (d *DataImporter) importData(ctx context.Context) error {
	for _, dataset := range d.datasetList {
		// parse products data from dataset file
		products, err := parse(dataset)
		if err != nil {
			return fmt.Errorf("error parsing dataset %s: %v", dataset, err)
		}

		products = clean(products)
		if len(products) == 0 {
			fmt.Printf("No products found in %s\n", dataset)
		}

		// create random seller users
		userIDs, err := d.createSellerUsers(ctx, 10)
		if err != nil {
			return fmt.Errorf("error creating seller users: %v", err)
		}

		// create stores for newly created sellers
		sellerIDs, err := d.createStoresAndIndex(ctx, userIDs)
		if err != nil {
			return fmt.Errorf("error creating stores: %v", err)
		}

		// each .jsonl file has the same top category which is a root category
		// therefore, we create root category using the first product data
		rootCategoryID, err := d.storeRootCategoryAndIndex(ctx, products[0].Categories[0])
		if err != nil {
			return fmt.Errorf("error creating root category: %v", err)
		}

		for _, productItem := range products {
			time.Sleep(15 * time.Millisecond)

			productItem.RootCategoryID = rootCategoryID
			if err := d.storeProductDataAndIndex(ctx, productItem, sellerIDs); err != nil {
				return fmt.Errorf("error importing product: %v", err)
			}
		}

		d.datasetListSucceeded = append(d.datasetListSucceeded, dataset)
	}

	return nil
}

func (d *DataImporter) createSellerUsers(ctx context.Context, n int) ([]int64, error) {
	userIDs := make([]int64, n)
	for i := 0; i < n; i++ {
		// save a default user and get its id
		hashedPassword, _ := utils.HashPassword(fakeUserPassword)

		email := gofakeit.Email()
		fname := gofakeit.Name()
		status := user.UserStatusActive

		userID, err := d.userRepository.CreateSeller(ctx, email, hashedPassword, fname, status)
		if err != nil {
			return nil, err
		}

		userIDs[i] = userID
	}

	return userIDs, nil
}

func (d *DataImporter) createStoresAndIndex(ctx context.Context, userIDs []int64) ([]int32, error) {
	storeIDs := make([]int32, len(userIDs))
	for index, userID := range userIDs {
		// save a default repository and get its id
		name := gofakeit.Company()
		handle := slug.Make(name)
		description := gofakeit.Sentence(20)

		storeID, err := d.storeRepository.CreateStore(ctx, handle, userID, name, description)
		if err != nil {
			return nil, err
		}

		err = d.searcher.IndexStore(ctx, anor.Store{
			ID:     storeID,
			Name:   name,
			Handle: handle,
		})
		if err != nil {
			return nil, err
		}

		storeIDs[index] = storeID
	}

	return storeIDs, nil
}

func (d *DataImporter) storeRootCategoryAndIndex(ctx context.Context, rootCategory string) (int32, error) {
	// check if root category is created in the db
	c, err := d.categoryRepository.GetCategoryByName(ctx, rootCategory)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		// save the top category in db and get retrieve its id
		c, err := d.categoryRepository.CreateTopCategory(ctx,
			rootCategory,
			utils.CreateHandle(rootCategory),
		)
		if err != nil {
			return 0, err
		}

		// index root category
		err = d.searcher.IndexCategory(ctx, anor.Category{
			ID:             c.ID,
			Category:       c.Category,
			Handle:         c.Handle,
			ParentID:       0,
			IsLeafCategory: false,
		})

		if err != nil {
			return 0, err
		}

		return c.ID, nil
	}
	return c.ID, nil
}

func (d *DataImporter) storeCategoryTreeAndIndex(ctx context.Context, rootCategoryID int32, tree []string) (anor.Category, error) {
	var ac anor.Category
	id := rootCategoryID

	// loop over categories
	for index, c := range tree {
		// save every child category along with its parent id if not exists
		cat, err := d.categoryRepository.CreateChildCategoryIfNotExists(ctx, c, utils.CreateHandle(c), &id)
		if err != nil {
			return ac, err
		}

		id = cat.ID

		ac = anor.Category{
			ID:       cat.ID,
			Category: cat.Category,
			Handle:   cat.Handle,
			ParentID: *cat.ParentID,
		}

		if index == len(tree)-1 {
			ac.IsLeafCategory = true
		}

		err = d.searcher.IndexCategory(ctx, ac)
		if err != nil {
			var terr *typesense.ApiError
			if errors.As(err, &terr) && terr.StatusCode == 409 {
				continue
			}
			return ac, err
		}
	}

	return ac, nil
}

func (d *DataImporter) storeProductDataAndIndex(ctx context.Context, p ProductJSON, sellerIDs []int32) error {
	leafCategory, err := d.storeCategoryTreeAndIndex(ctx, p.RootCategoryID, p.Categories[1:])
	if err != nil {
		return fmt.Errorf("failed to create descendant categories: %w", err)
	}

	imageURLs := make(map[int]string, len(p.ImageUrls))
	for i, url := range p.ImageUrls {
		imageURLs[i] = url
	}

	specs := make(map[string]string, len(p.Specs))
	for k, v := range p.Specs {
		specs[k] = v
	}

	// get random seller repository id
	sellerStoreID := sellerIDs[utils.GenRandomNum(0, len(sellerIDs)-1)]

	// create product
	handle := utils.CreateHandle(p.Name)
	savedProduct, err := d.productRepository.CreateProduct(
		ctx,
		sellerStoreID,
		leafCategory.ID,
		p.Name,
		&p.Brand,
		handle,
		shortInfo,
		imageURLs,
		specs,
		product.ProductStatusPublished,
	)
	if err != nil {
		return err
	}

	// create product pricing
	dp, _ := decimal.NewFromString(p.Price)
	pp, err := d.storeProductPricing(ctx, savedProduct.ID, dp)
	if err != nil {
		return err
	}

	// no create skus, productItem attributes and sku productItem attr values
	// cleanedAttMap might have 0 or more attrs
	// for ease of use I will handle for cases; 0, 1 , 2 and 3
	attrLen := len(p.Attributes)
	skuGen := utils.NewSKUGenerator(leafCategory.Category, p.Name)
	if attrLen == 0 {
		_, err := d.storeProductVariantRecursive(ctx, savedProduct.ID, skuGen.GetBaseSKU(), 0)
		if err != nil {
			return err
		}

	} else {
		attrIDMap := make(map[string]int64) // { "Color": 266 }
		for k := range p.Attributes {
			id, err := d.productRepository.CreateProductAttribute(ctx, savedProduct.ID, k)
			if err != nil {
				return err
			}
			attrIDMap[k] = id
		}

		var slices [][]map[string]string
		for k, v := range p.Attributes {
			var slice []map[string]string
			for _, attr := range v {
				i := map[string]string{
					k: attr,
				}
				slice = append(slice, i)
			}
			slices = append(slices, slice)
		}

		crPr := cartesianProduct(slices...)
		for _, cp := range crPr {
			var attrVals []string
			for _, attr := range cp {
				attrVals = append(attrVals, attr)
			}
			sku := skuGen.GenerateSKU(attrVals)
			randNum := utils.GenRandomNum(0, 40) - 10
			if randNum < 0 {
				randNum = 0
			}

			skuID, err := d.productRepository.CreateProductVariant(ctx,
				savedProduct.ID,
				sku,
				int32(randNum),
				false,
				[]int16{},
			)
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					if pgErr.Code == pgerrcode.UniqueViolation {
						skuID, err = d.productRepository.CreateProductVariant(ctx,
							savedProduct.ID,
							sku,
							int32(randNum),
							false,
							[]int16{},
						)
						if err != nil {
							return err
						}
					}
				}
			}

			// loop through map
			for k, v := range cp {
				err = d.productRepository.CreateProductVariantAttributeValues(ctx, skuID, attrIDMap[k], v)
				if err != nil {
					return err
				}
			}
			fmt.Printf("product variant created: %v\n", skuID)
		}
	}

	err = d.searcher.IndexProduct(ctx, anor.Product{
		BaseProduct: anor.BaseProduct{
			ID:         savedProduct.ID,
			Name:       p.Name,
			CategoryID: leafCategory.ID,
			Brand:      p.Brand,
			Handle:     handle,
			ImageUrls:  imageURLs,
			CreatedAt:  savedProduct.CreatedAt.Time,
			UpdatedAt:  savedProduct.UpdatedAt.Time,
		},
		Pricing: anor.ProductPricing{
			BasePrice:       pp.BasePrice,
			Discount:        pp.Discount,
			DiscountedPrice: pp.DiscountedPrice,
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("product indexed: %v\n", savedProduct.ID)
	return nil
}

func (d *DataImporter) storeProductPricing(ctx context.Context, productID int64, price decimal.Decimal) (anor.ProductPricing, error) {
	var (
		discount        decimal.Decimal
		discountedPrice decimal.Decimal
		isOnSale        bool
	)
	rndDiscount := generateRandomDiscount()
	if rndDiscount != 0 {
		discount = decimal.NewFromFloat32(rndDiscount)
		discountedAmount := price.Mul(discount).Round(2)
		discountedPrice = price.Sub(discountedAmount)
		isOnSale = true
	} else {
		discountedPrice = price
	}

	err := d.productRepository.CreateProductPricing(
		ctx,
		productID,
		price,
		"USD",
		discount,
		discountedPrice,
		isOnSale,
	)
	if err != nil {
		return anor.ProductPricing{}, err
	}

	return anor.ProductPricing{
		BasePrice:       price,
		Discount:        discount,
		DiscountedPrice: discountedPrice,
	}, nil
}

func (d *DataImporter) storeProductVariantRecursive(ctx context.Context, productID int64, skuValue string, retryCount int) (int64, error) {
	variantID, err := d.productRepository.CreateProductVariant(ctx,
		productID,
		skuValue,
		int32(utils.GenRandomNum(0, 50)),
		false,
		[]int16{},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			if retryCount >= 10 { // Set a maximum retry limit to avoid infinite recursion
				return 0, fmt.Errorf("max retries reached for creating product variant: %w", err)
			}
			// Recursive call with incremented retry count
			return d.storeProductVariantRecursive(ctx, productID, skuValue, retryCount+1)
		}
		return 0, err // Return any other error
	}

	return variantID, nil
}

func cartesianProduct(input ...[]map[string]string) []map[string]string {
	var result []map[string]string

	// Recursive function to generate Cartesian product
	var generate func(int, map[string]string)
	generate = func(index int, current map[string]string) {
		if index == len(input) {
			// Add the current combination to the result
			result = append(result, copyMap(current))
			return
		}

		// Iterate over the values for the current map
		for _, value := range input[index] {
			// Relation the current map in the combination
			for k, v := range value {
				current[k] = v
			}
			// Recursively generate combinations for the next map
			generate(index+1, current)
			// Backtrack and remove the values for the next iteration
			for k := range value {
				delete(current, k)
			}
		}
	}

	// Start the generation process
	generate(0, make(map[string]string))

	return result
}

// Helper function to copy a map
func copyMap(original map[string]string) map[string]string {
	c := make(map[string]string)
	for key, value := range original {
		c[key] = value
	}
	return c
}

func generateRandomDiscount() float32 {
	idx := rand.Intn(len(discounts))
	return discounts[idx]
}
