package product

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	notfound "github.com/aliml92/anor/html/dtos/pages/not-found"
	productdetails "github.com/aliml92/anor/html/dtos/pages/product-details"
	"github.com/aliml92/anor/html/dtos/pages/product-details/components"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/aliml92/anor/html/dtos/shared"
	"net/http"
)

// ProductDetailsView handles the HTTP request for displaying the product details page.
func (h *Handler) ProductDetailsView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := anor.UserFromContext(ctx)

	handle := r.PathValue("handle")
	productID, err := extractProductID(handle)
	if err != nil {
		h.viewRenderNotFound(w, r, "Product not found")
		return
	}

	p, err := h.productSvc.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, anor.ErrNotFound) {
			h.viewRenderNotFound(w, r, "Product not found")
			return
		}
		h.serverInternalError(w, err)
		return
	}

	c, err := h.categorySvc.GetCategory(ctx, p.CategoryID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	ac, err := h.categorySvc.GetAncestorCategories(ctx, p.CategoryID)
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

	if isHXRequest(r) {
		h.view.Render(w, "pages/product-details/content.gohtml", pdc)
		return
	}

	headerContent, err := h.getHeaderContent(ctx, u)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	pdb := productdetails.Base{
		Header:  headerContent,
		Content: pdc,
	}

	h.view.Render(w, "pages/product-details/base.gohtml", pdb)
}

func (h *Handler) viewRenderNotFound(w http.ResponseWriter, r *http.Request, message string) {
	ctx := r.Context()
	u := anor.UserFromContext(ctx)
	nc := notfound.Content{Message: message}
	if isHXRequest(r) {
		h.view.Render(w, "pages/not-found/content.gohtml", nc)
		return
	}

	headerContent, err := h.getHeaderContent(ctx, u)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	nb := notfound.Base{
		Header:  headerContent,
		Content: nc,
	}

	h.view.Render(w, "pages/not-found/base.gohtml", nb)
}

func (h *Handler) getHeaderContent(ctx context.Context, u *anor.User) (partials.Header, error) {
	headerContent := partials.Header{User: u}
	if u != nil {
		ac, err := h.userSvc.GetUserActivityCounts(ctx, u.ID)
		if err != nil {
			return headerContent, err
		}

		headerContent.ActiveOrdersCount = ac.ActiveOrdersCount
		headerContent.WishlistItemsCount = ac.WishlistItemsCount
		headerContent.CartItemsCount = ac.CartItemsCount

	} else {
		cartId := h.session.Guest.GetInt64(ctx, "guest_cart_id")
		if cartId != 0 {

			guestCartItemCount, err := h.cartSvc.CountCartItems(ctx, cartId)
			if err != nil {
				return headerContent, err
			}

			headerContent.CartItemsCount = int(guestCartItemCount)
		}
	}

	return headerContent, nil
}

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

		fmt.Printf("len of m1: %d\n", len(m1))
		fmt.Printf("len of m2: %d\n", len(m2))

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

func create2DProductVariantMatrix(n, m int) [][]anor.ProductVariant {
	var matrix [][]anor.ProductVariant
	for i := 0; i < n; i++ {
		row := make([]anor.ProductVariant, m)
		matrix = append(matrix, row)
	}

	return matrix
}

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
