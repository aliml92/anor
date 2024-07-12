package main

import (
	"github.com/agnivade/levenshtein"
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode"
)

var uniqueBrands = make(map[string]string)

func clean(products []ProductJSON) []ProductJSON {
	var cleanedProducts []ProductJSON
	for _, p := range products {
		name, skip := cleanProductName(p.Name)
		if skip {
			continue
		}

		cleanedCategories, skip := cleanCategories(p.Categories)
		if skip {
			continue
		}

		cleanedPrice := cleanPriceStr(p.Price)

		cleanedImageURLs := cleanProductImageUrls(p.ImageUrls)

		cleanedSpecs := cleanProductSpecifications(p.Specs)

		brand := getBrandFromSpecs(cleanedSpecs)

		cleanedAttributes := cleanProductAttributes(p.Attributes)

		cp := ProductJSON{
			Categories: cleanedCategories,
			Name:       name,
			Price:      cleanedPrice,
			Available:  p.Available,
			Sold:       p.Sold,
			ImageUrls:  cleanedImageURLs,
			Specs:      cleanedSpecs,
			Brand:      brand,
			Attributes: cleanedAttributes,
		}

		cleanedProducts = append(cleanedProducts, cp)
	}
	return cleanedProducts
}

func cleanProductName(name string) (string, bool) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", true
	}

	return name, false
}

func cleanCategories(categories []string) ([]string, bool) {
	if len(categories) <= 1 || len(categories) > 4 {
		return categories, true
	}

	var cleanedCategories []string
	for _, c := range categories {
		if strings.TrimSpace(c) == "" {
			return categories, true
		}
		if strings.Contains(c, "Underwear") {
			return categories, true
		}
		if strings.Contains(c, "See more") {
			c = strings.ReplaceAll(c, "See more", "")
		}
		cleanedCategories = append(cleanedCategories, c)
	}

	return cleanedCategories, false
}

func cleanPriceStr(p string) string {
	if strings.TrimSpace(p) == "" {
		return "14.99"
	}

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

	// set hardcoded values to priceStr if it is not convertible to shopspring.Decimal
	_, err := decimal.NewFromString(priceStr)
	if err != nil {
		priceStr = "14.99"
	}

	return priceStr

}

func cleanProductImageUrls(urls []string) []string {
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
	return images
}

func cleanProductSpecifications(specs map[string]string) map[string]string {
	for k, v := range specs {
		if strings.Contains(v, "...") {
			v := strings.Split(v, "...")[0]
			if v == "" {
				delete(specs, k)
			} else {
				specs[k] = v
			}
		}
	}

	return specs
}

func cleanProductAttributes(attrs map[string][]string) map[string][]string {
	cleanedAttrMap := make(map[string][]string)
	for k, v := range attrs {

		k = strings.TrimSpace(k)
		if k == "" || len(v) == 0 {
			continue
		}

		// filter out "choose" attributes
		k := strings.ToLower(k)
		if k == "choose" {
			continue
		}

		k = strings.ToLower(k)
		if strings.Contains(k, "color") || strings.Contains(k, "colour") {
			k = "Color"
		}
		if strings.Contains(k, "size") {
			if !strings.Contains(k, "eu") && !strings.Contains(k, "us") {
				k = "Size"
			}
		}

		var vals []string
		for _, attVal := range v {
			if strings.Contains(attVal, "(Out Of Stock)") {
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

func getBrandFromSpecs(specs map[string]string) string {
	var b string
	if len(specs) > 0 {
		b = specs["Brand"]
	}
	b = strings.TrimSpace(b)
	if b != "" {
		b = strings.ToLower(b)
		value := cases.Title(language.English).String(b)
		bs := strings.Fields(b)
		var builder strings.Builder
		for _, s := range bs {
			builder.WriteString(s)
		}
		key := builder.String()
		if val, ok := uniqueBrands[key]; ok {
			return val
		} else {
			// use Levenshtein distance
			for k, v := range uniqueBrands {
				distance := levenshtein.ComputeDistance(k, key)
				if distance <= 1 {
					return v
				}
			}

			uniqueBrands[key] = value
			return value
		}
	}

	return "Unbranded"
}
