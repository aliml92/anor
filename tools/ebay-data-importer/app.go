package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/aliml92/go-typesense/typesense"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	errors2 "github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/shopspring/decimal"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/pkg/utils"
	"github.com/aliml92/anor/postgres/store/category"
	"github.com/aliml92/anor/postgres/store/product"
	"github.com/aliml92/anor/postgres/store/sellerstore"
	"github.com/aliml92/anor/postgres/store/user"
	ts2 "github.com/aliml92/anor/typesense"
)

const (
	// dataset files' extension
	dotJsonl = ".jsonl"
)

var (
	// flags
	source   string
	database string
	ts       string
	force    bool
)

type ProductJSON struct {
	Categories []string            `json:"categories"`
	Name       string              `json:"name"`
	Price      string              `json:"price"`
	Available  string              `json:"available"`
	Sold       string              `json:"sold"`
	ImageUrls  []string            `json:"image_links"`
	Specs      map[string]string   `json:"specs"`
	Attributes map[string][]string `json:"attributes"`
}

func init() {
	flag.StringVar(&source, "source", "", "path to the dataset folder")
	flag.StringVar(&database, "database", "", "PostgreSQL database connection string")
	flag.StringVar(&ts, "typesense", "", "Typesense server url")
	flag.BoolVar(&force, "force", false, "import all json files forcefully")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
}

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if source == "" {
		return fmt.Errorf("source path is required")
	}

	if database == "" {
		return fmt.Errorf("PostgreSQL database connection string is required")
	}

	if ts == "" {
		return fmt.Errorf("Typesense server url string is required")
	}

	db, err := pgxpool.New(ctx, database)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}
	defer db.Close()

	client, _ := typesense.NewClient(nil, ts)
	client = client.WithAPIKey("xyz")

	// Check if there is a file to keep track of imported files
	trackFile := "./tools/ebay-data-importer/imported_files.txt"
	importedFiles := make(map[string]bool)
	if _, err := os.Stat(trackFile); !os.IsNotExist(err) {
		// Read the existing imported files
		file, err := os.Open(trackFile)
		if err != nil {
			return fmt.Errorf("failed to open imported files file: %w", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			importedFiles[scanner.Text()] = true
		}
	}

	err = walkDataset(ctx, source, importedFiles, db, client)
	if err != nil {
		return err
	}

	// Write the updated list of imported files to the track file
	file, err := os.Create(trackFile)
	if err != nil {
		return fmt.Errorf("failed to create track file: %w", err)
	}
	defer file.Close()

	for path := range importedFiles {
		fmt.Fprintf(file, "%s\n", path)
	}

	return nil
}

func walkDataset(ctx context.Context, source string, importedFiles map[string]bool, db *pgxpool.Pool, client *typesense.Client) error {
	// create store objects
	var (
		userStore     = user.NewStore(db)
		productStore  = product.NewStore(db)
		categoryStore = category.NewStore(db)
		sellerStore   = sellerstore.NewSellerStore(db)
		searcher      = ts2.NewSearcher(client)
	)

	// walk every dataset files
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk through source directory: %w", err)
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == dotJsonl {
			// Check if the file is already imported
			if !force && importedFiles[path] {
				fmt.Printf("Skipping already imported file: %s\n", path)
				return nil
			}

			fmt.Printf("Processing file: %s\n", path)

			products, err := parse(path)
			if err != nil {
				return err
			}

			// create 10 seller users
			userIDs, err := createSellerUsers(ctx, userStore, 10)
			if err != nil {
				return err
			}

			// create sellers' stores
			sellerIDs, err := createSellerStores(ctx, sellerStore, userIDs, searcher)
			if err != nil {
				return errors2.Wrap(err, "failed to create seller store")
			}

			// dataset .jsonl files has the same top category (the root category)
			// therefore, get the first product data a sample
			// create top category
			rootCategoryID, err := createRootCategory(ctx, categoryStore, products[0], searcher)
			if err != nil {
				return errors2.Wrap(err, "failed to create root category")
			}

			shortInfo := []string{"Trendy design for versatile wear", "Quality materials for durability", "Available in various sizes and colors"}
			for _, productItem := range products {
				time.Sleep(10 * time.Millisecond)

				categoryID, skip, err := createDescendentCategories(ctx, categoryStore, productItem, rootCategoryID, searcher)
				if skip {
					continue
				}
				if err != nil {
					return errors2.Wrap(err, "failed to create descendent categories")
				}

				price, err := cleanProductPrice(productItem.Price)
				if err != nil {
					return err
				}

				urls := cleanProductImageUrls(productItem.ImageUrls)
				imageUrls, _ := urls.(product.ImageUrls)

				// clean productItem item specs data
				var brand string
				var s map[string]string
				s = productItem.Specs
				if len(s) > 0 {
					// clean empty key-value pairs
					delete(s, "")

					// get brand name
					brand = s["Brand"]
				}

				// transfer specs to Specs map
				specs := make(product.Specs)
				for idx, spec := range s {
					specs[idx] = spec
				}

				// get random seller store id
				sellerStoreID := sellerIDs[utils.GenRandomNum(0, len(sellerIDs)-1)]

				// create product
				slug := utils.CreateSlug(productItem.Name)
				savedProduct, err := productStore.CreateProduct(
					ctx,
					sellerStoreID,
					categoryID,
					productItem.Name,
					&brand,
					slug,
					shortInfo,
					imageUrls,
					specs,
					product.ProductStatusActive,
				)
				if err != nil {
					return err
				}

				// create product pricing
				pp, err := createProductPricing(ctx, productStore, price, savedProduct.ID)
				if err != nil {
					return err
				}

				// clean product attributes
				cleanedAttrMap := cleanProductAttributes(productItem.Attributes)

				// no create skus, productItem attributes and sku productItem attr values
				// cleanedAttMap might have 0 or more attrs
				// for ease of use I will handle for cases; 0, 1 , 2 and 3
				attrLen := len(cleanedAttrMap)
				prefix := "SKU-"
				if attrLen == 0 {
					_, err := createSKU(ctx, productStore, savedProduct.ID)
					if err != nil {
						return err
					}

				} else {
					attrIDMap := make(map[string]int64) // { "Color": 266 }
					for k := range cleanedAttrMap {
						id, err := productStore.CreateProductAttribute(ctx, savedProduct.ID, k)
						if err != nil {
							return err
						}
						attrIDMap[k] = id
					}

					var slices [][]map[string]string
					for k, v := range cleanedAttrMap {
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
						sku := prefix + xid.New().String()
						skuID, err := productStore.CreateSKU(ctx,
							savedProduct.ID,
							sku,
							int32(utils.GenRandomNum(0, 50)),
							false,
							[]int16{},
						)
						if err != nil {
							return err
						}

						// loop through map
						for k, v := range cp {
							err = productStore.CreateSKUProductAttributeValues(ctx, skuID, attrIDMap[k], v)
							if err != nil {
								return err
							}
						}

					}
				}
				err = searcher.IndexProduct(ctx, anor.Product{
					BaseProduct: anor.BaseProduct{
						ID:           savedProduct.ID,
						Name:         productItem.Name,
						CategoryID:   categoryID,
						Brand:        brand,
						Slug:         slug,
						ThumbImgUrls: map[string]string{"main": imageUrls[0]},
						CreatedAt:    savedProduct.CreatedAt.Time,
						UpdatedAt:    savedProduct.UpdatedAt.Time,
					},
					Pricing: anor.ProductPricing{
						BasePrice:        pp.BasePrice,
						DiscountedAmount: pp.DiscountedAmount,
					},
				})
				if err != nil {
					return err
				}

			}

			// Update the list of imported files
			importedFiles[path] = true
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func parse(filepath string) ([]ProductJSON, error) {
	// open json data
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	var products []ProductJSON
	// add productItem data to products
	d := json.NewDecoder(f)
	for d.More() {
		var line ProductJSON
		if err := d.Decode(&line); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		products = append(products, line)
	}
	return products, nil
}

func createRootCategory(ctx context.Context, cs category.Store, p ProductJSON, searcher *ts2.Searcher) (int32, error) {
	rootCategory := p.Categories[0]

	// check if root category is created in the db
	c, err := cs.GetCategoryByName(ctx, rootCategory)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		// save the top category in db and get retrieve its id
		c, err := cs.CreateTopCategory(ctx,
			rootCategory,
			utils.CreateSlug(rootCategory),
		)
		if err != nil {
			return 0, err
		}

		// index root category
		err = searcher.IndexCategory(ctx, anor.Category{
			ID:             c.ID,
			Category:       c.Category,
			Slug:           c.Slug,
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

func createDescendentCategories(ctx context.Context, cs category.Store, p ProductJSON, rootCategoryID int32, searcher *ts2.Searcher) (int32, bool, error) {
	id := rootCategoryID

	// exclude the first category (which is the root)
	// since it is already created
	categories := p.Categories[1:]

	// some products belong to two categories at the same time
	// therefore the product item with such categories not processed
	if len(categories) >= 4 {
		return 0, true, nil
	}

	// loop over categories
	for index, c := range categories {

		// save every child category along with its parent id if not exists
		cat, err := cs.CreateChildCategoryIfNotExists(ctx, c, utils.CreateSlug(c), &id)
		if err != nil {
			return 0, false, err
		}

		id = cat.ID

		ac := anor.Category{
			ID:       cat.ID,
			Category: cat.Category,
			Slug:     cat.Slug,
			ParentID: *cat.ParentID,
		}

		if index == len(categories)-1 {
			ac.IsLeafCategory = true
		}

		err = searcher.IndexCategory(ctx, ac)
		if err != nil {
			var terr *typesense.ApiError
			if errors.As(err, &terr) && terr.StatusCode == 409 {
				continue
			}
			return 0, false, err
		}
	}

	return id, false, nil
}

func cleanProductPrice(p string) (decimal.Decimal, error) {
	// split price data, e.g. 'GBP 15.20' or 'US $15' or 'US $15 - $19' or 'US $18/ea'
	// splitting is necessary for the third types of price data
	// since we only keep the first price for price ranges
	priceStr := strings.Split(p, " ")[1]

	// init builder to construct cleaned data
	var builder strings.Builder

	// construct price data by filter out all chars except digit and .
	for _, r := range priceStr {
		if unicode.IsDigit(r) || r == '.' {
			builder.WriteRune(r)
		}
	}

	// get price data as string
	priceStr = builder.String()

	// create decimal.Decimal from string
	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return price, nil
}

func cleanProductImageUrls(urls []string) interface{} {
	// drop half of the image links, since they are duplicate
	// with different sizes
	// images empty assign a default image link (not found image link)
	var images []string
	half := len(urls) / 2
	if half != 0 {
		images = urls[:half]
	} else {
		images = []string{"https://ir.ebaystatic.com/pictures/aw/pics/stockimage1.jpg"}
	}

	// transfer images to ImageUrls map
	imageUrls := make(product.ImageUrls)
	for idx, img := range images {
		imageUrls[idx] = img
	}

	return imageUrls
}

func createProductPricing(ctx context.Context, ps product.Store, price decimal.Decimal, productID int64) (anor.ProductPricing, error) {
	var (
		discountLevel    decimal.Decimal
		discountedAmount decimal.Decimal
		isOnSale         bool
	)
	rndDiscount := generateRandomDiscount()
	if rndDiscount != 0 {
		discountLevel = decimal.NewFromFloat32(rndDiscount)
		discountedAmount = price.Mul(discountLevel).Round(2)
		isOnSale = true
	}

	err := ps.CreateProductPricing(
		ctx,
		productID,
		price,
		"USD",
		discountLevel,
		discountedAmount,
		isOnSale,
	)
	if err != nil {
		return anor.ProductPricing{}, err
	}

	return anor.ProductPricing{
		BasePrice:        price,
		DiscountedAmount: discountedAmount,
	}, nil
}

func cleanProductAttributes(attrs map[string][]string) map[string][]string {
	cleanedAttrMap := make(map[string][]string)
	for k, v := range attrs {

		k = strings.TrimSpace(k)
		if k == "" || len(v) == 0 {
			continue
		}

		// filter out "choose" attributes
		low := strings.ToLower(k)
		if low == "choose" {
			continue
		}

		vals := []string{}
		for _, attVal := range v {
			if strings.Index(attVal, "(Out Of Stock)") != -1 {
				attVal = strings.ReplaceAll(attVal, "(Out Of Stock)", "")
			}
			attVal = strings.TrimSpace(attVal)
			if attVal != "" {
				vals = append(vals, attVal)
			}
		}

		if len(vals) != 0 {
			cleanedAttrMap[k] = vals
		}
	}

	return cleanedAttrMap
}

func createSKU(ctx context.Context, ps product.Store, productID int64) (int64, error) {
	// create a random sku
	sku := "SKU-" + xid.New().String()
	skuID, err := ps.CreateSKU(ctx,
		productID,
		sku,
		int32(utils.GenRandomNum(0, 50)),
		false,
		[]int16{},
	)

	return skuID, err
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
			// Include the current map in the combination
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
