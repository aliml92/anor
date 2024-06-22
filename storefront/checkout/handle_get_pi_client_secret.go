package checkout

import (
	"encoding/json"
	"github.com/aliml92/anor"
	"net/http"
)

type GetPIClientSecretResponseBody struct {
	ClientSecret string `json:"clientSecret"`
}

func (h *Handler) GetPIClientSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := anor.UserFromContext(ctx)
	res := &GetPIClientSecretResponseBody{}

	if u != nil {
		c, err := h.cartSvc.GetCart(ctx, u.ID, false)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

		res.ClientSecret = c.PIClientSecret
	} else {
		// TODO implement guest order later
	}

	b, err := json.Marshal(res)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
