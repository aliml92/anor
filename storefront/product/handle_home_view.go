package product

import (
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/pages/home"
	"github.com/aliml92/anor/html/templates/pages/home/components"
	"net/http"
)

func (h *Handler) HomeView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	featuredSelections, err := h.featuredSelectionService.ListAllActive(ctx)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	// TODO: this service does not bring actual popular products
	popular, err := h.productService.ListPopularProducts(ctx, anor.PopularProductListParams{})

	// Prepare home page content
	hc := home.Content{
		Featured: components.Featured{
			Selections: featuredSelections,
		},
		Popular: components.Collection{
			Title:        "Popular Products",
			ResourcePath: "/categories/popular-products-101",
			Products:     popular,
		},
	}

	h.Render(w, r, "pages/home/content.gohtml", hc)
}
