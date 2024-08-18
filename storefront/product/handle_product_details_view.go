package product

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	notfound "github.com/aliml92/anor/html/templates/pages/not_found"
	productdetails "github.com/aliml92/anor/html/templates/pages/product_details"
	"github.com/aliml92/anor/html/templates/pages/product_details/components"
	"github.com/aliml92/anor/html/templates/shared"
	"github.com/aliml92/anor/relation"
	"net/http"
	"strconv"
	"strings"
)

// ProductDetailsView handles the request for the product details page.
func (h *Handler) ProductDetailsView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handle := r.PathValue("handle")
	productID, err := extractProductID(handle)
	if err != nil {
		h.renderNotFound(w, r, "Product not found")
		return
	}

	relationSet := relation.New(relation.Store, relation.Pricing, relation.ProductAttribute, relation.ProductVariant)
	p, err := h.productService.Get(ctx, productID, relationSet)
	if err != nil {
		if errors.Is(err, anor.ErrNotFound) {
			h.renderNotFound(w, r, "Product not found")
			return
		}
		h.serverInternalError(w, err)
		return
	}

	c, err := h.categoryService.GetCategory(ctx, p.CategoryID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	ac, err := h.categoryService.GetAncestorCategories(ctx, p.CategoryID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	pdc := productdetails.Content{
		CategoryBreadcrumb: shared.CategoryBreadcrumb{
			Category:           c,
			AncestorCategories: ac,
		},
		ProductMain: components.ProductMain{
			Product: p,
		},
		ProductSpecs: components.ProductSpecs{
			Product: p,
		},
		ProductVariantMatrix: constructProductVariantMatrix(p.ProductVariants, p.Attributes),
	}

	h.Render(w, r, "pages/product_details/content.gohtml", pdc)
}

// renderNotFound renders `not found` page with given message
func (h *Handler) renderNotFound(w http.ResponseWriter, r *http.Request, message string) {
	c := notfound.Content{Message: message}
	h.Render(w, r, "pages/not_found/content.gohtml", c)
}

// extractProductID extracts product id from handle
func extractProductID(handle string) (int64, error) {
	lastIndex := strings.LastIndex(handle, "-")
	if lastIndex == -1 {
		// assume handle is number string
		id, err := strconv.Atoi(handle)
		if err != nil {
			return 0, fmt.Errorf("failed to convert ID to integer: %s", handle)
		}
		return int64(id), nil
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

	return int64(id), nil
}

// constructProductVariantMatrix creates a matrix of product variants based on attributes.
// It handles cases for 0 to 3 attributes, returning different structures accordingly.
func constructProductVariantMatrix(variants []anor.ProductVariant, attributes []anor.ProductAttribute) interface{} {
	switch len(attributes) {
	case 0:
		return []anor.ProductVariant{variants[0]}
	case 1:
		a := attributes[0]

		m := make(map[string]int)
		for index, v := range a.Values {
			m[v] = index
		}

		matrix := make([]anor.ProductVariant, len(m))
		for _, sku := range variants {
			val := sku.Attributes[a.Attribute]
			index := m[val]
			matrix[index] = sku
		}

		return matrix
	case 2:
		a1 := attributes[0]
		a2 := attributes[1]

		m1 := make(map[string]int)
		for index, v := range a1.Values {
			m1[v] = index
		}

		m2 := make(map[string]int)
		for index, v := range a2.Values {
			m2[v] = index
		}

		matrix := create2DProductVariantMatrix(len(m1), len(m2))
		for _, variant := range variants {
			val1 := variant.Attributes[a1.Attribute]
			val2 := variant.Attributes[a2.Attribute]
			index1 := m1[val1]
			index2 := m2[val2]
			matrix[index1][index2] = variant
		}
		return matrix
	case 3:
		a1 := attributes[0]
		a2 := attributes[1]
		a3 := attributes[2]

		m1 := make(map[string]int)
		for index, v := range a1.Values {
			m1[v] = index
		}

		m2 := make(map[string]int)
		for index, v := range a2.Values {
			m2[v] = index
		}

		m3 := make(map[string]int)
		for index, v := range a3.Values {
			m3[v] = index
		}

		matrix := create3DProductVariantMatrix(len(m1), len(m2), len(m3))
		for _, variant := range variants {
			val1 := variant.Attributes[a1.Attribute]
			val2 := variant.Attributes[a2.Attribute]
			val3 := variant.Attributes[a3.Attribute]
			index1 := m1[val1]
			index2 := m2[val2]
			index3 := m3[val3]
			matrix[index1][index2][index3] = variant
		}
		return matrix
	default:
		// case 4 or more not handled for now
		return nil
	}
}

// create2DProductVariantMatrix initializes a 2D matrix of ProductVariants.
func create2DProductVariantMatrix(n, m int) [][]anor.ProductVariant {
	var matrix [][]anor.ProductVariant
	for i := 0; i < n; i++ {
		row := make([]anor.ProductVariant, m)
		matrix = append(matrix, row)
	}

	return matrix
}

// create3DProductVariantMatrix initializes a 3D matrix of ProductVariants.
func create3DProductVariantMatrix(n, m, l int) [][][]anor.ProductVariant {
	var matrix [][][]anor.ProductVariant
	for i := 0; i < n; i++ {
		var innerMatrix [][]anor.ProductVariant
		for j := 0; j < m; j++ {
			row := make([]anor.ProductVariant, l)
			innerMatrix = append(innerMatrix, row)
		}
		matrix = append(matrix, innerMatrix)
	}
	return matrix
}
