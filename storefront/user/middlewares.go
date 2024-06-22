package user

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aliml92/anor"
)

func (h *Handler) AuthInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *anor.User
		id := h.session.Auth.GetInt64(r.Context(), "authenticatedUserID")
		if id != 0 {
			// retrieve user by id from database
			ru, err := h.userSvc.GetUser(r.Context(), id) // retrieved user
			if errors.Is(err, anor.ErrNotFound) {
				h.logger.LogAttrs(
					r.Context(),
					slog.LevelWarn,
					"ERR_INCONSISTENT_SESSION_DATA",
					slog.String("user_id", fmt.Sprint(id)),
					slog.String("user_data_in_session", "found"),
					slog.String("user_data_in_db", "not found"),
				)
			} else if err != nil {
				h.serverInternalError(w, err)
				return
			} else {
				user = &ru
			}
		}

		ctx := anor.NewContextWithUser(r.Context(), user)
		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
