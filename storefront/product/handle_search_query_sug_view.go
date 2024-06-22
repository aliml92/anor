package product

import (
	"fmt"
	"github.com/aliml92/anor/html/dtos/partials"
	"html/template"
	"net/http"
	"net/url"
)

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

	sqsl := partials.SearchQuerySuggestionsList{}

	sugs := make([]template.HTML, len(ss.ProductNameSuggestions))
	for index, sug := range ss.ProductNameSuggestions {
		s := template.HTML(sug)
		sugs[index] = s
	}

	sqsl.ProductNameSuggestions = sugs

	h.view.RenderComponent(w, "partials/header/search-query-suggestions-list.gohtml", sqsl)
}
