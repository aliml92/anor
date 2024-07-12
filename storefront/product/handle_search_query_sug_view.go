package product

import (
	"fmt"
	"github.com/aliml92/anor/html/dtos/partials"
	"html/template"
	"net/http"
	"net/url"
)

// SearchQuerySuggestionsView handles search query suggestions requests.
// It parses the query parameter, fetches product, category and store suggestions,
// and renders them as a component.
func (h *Handler) SearchQuerySuggestionsView(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
	}

	q := values.Get("q")

	ctx := r.Context()
	ss, err := h.searcher.SearchQuerySuggestions(ctx, q)
	if err != nil {
		fmt.Println(err)
	}

	pns := make([]template.HTML, len(ss.ProductNameSuggestions))
	for idx, sug := range ss.ProductNameSuggestions {
		pns[idx] = template.HTML(sug)
	}

	sqsl := partials.SearchQuerySuggestionsList{
		ProductNameSuggestions: pns,
	}

	h.view.RenderComponent(w, "partials/header/search_query_suggestions_list.gohtml", sqsl)
}
