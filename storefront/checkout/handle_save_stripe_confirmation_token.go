package checkout

import (
	"encoding/json"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/session"
	"github.com/invopop/validation"
	"net/http"
)

type ConfirmationTokenRequest struct {
	ConfirmationTokenID string `json:"confirmation_token_id"`
}

func (req *ConfirmationTokenRequest) Bind(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	return err
}

func (req *ConfirmationTokenRequest) Validate() error {
	err := validation.Errors{
		"confirmationToken": validation.Validate(req.ConfirmationTokenID, validation.Required),
	}.Filter()

	return err
}

func (h *Handler) SaveStripeConfirmationToken(w http.ResponseWriter, r *http.Request) {
	var req ConfirmationTokenRequest
	err := anor.BindValid(r, &req)
	if err != nil {
		h.jsonClientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	h.session.Put(ctx, session.StripeConfirmationTokenIDKey, req.ConfirmationTokenID)

	w.WriteHeader(http.StatusOK)
}
