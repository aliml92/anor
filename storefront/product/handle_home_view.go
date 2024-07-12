package product

import (
	homepage "github.com/aliml92/anor/html/dtos/pages/home"
	"github.com/aliml92/anor/html/dtos/pages/home/components"
	"net/http"
)

// HomeView handles the request for the home page.
func (h *Handler) HomeView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	newArrivals, err := h.productSvc.GetNewArrivals(ctx, 10)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	// Prepare home page content
	hc := homepage.Content{
		Featured: components.Featured{
			Products: nil, // TODO: get featured products for carousel
		},
		NewArrivals: components.Collection{
			Products: newArrivals,
		},
		Popular: components.Collection{
			Products: nil, // TODO: get popular products
		},
	}

	h.Render(w, r, "pages/home", hc)
}
