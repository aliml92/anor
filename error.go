package anor

import (
	"errors"
	"log/slog"
	"net/http"
)

const (
	EINTERNALMSG = "Something went wrong. Please try again later."
)

var (
	ErrInvalidCart = errors.New("invalid cart")

	ErrNotFound = errors.New("not found")

	ErrNoPaymentIntent               = errors.New("no payment intent found")
	ErrUserExists                    = errors.New("user exists")
	ErrProductNotFound               = errors.New("product not found")
	ErrProductPricingNotFound        = errors.New("product pricing not found")
	ErrCartNotFound                  = errors.New("cart not found")
	ErrProductVariantNotFound        = errors.New("product variant not found")
	ErrProductVariantPricingNotFound = errors.New("product variant pricing not found")
)

func ClientError(logger *slog.Logger, w http.ResponseWriter, err error, statusCode int) {
	logger.Error(
		err.Error(),
		slog.Any("error", err),
	)
	http.Error(w, err.Error(), statusCode)
}

func ServerInternalError(logger *slog.Logger, w http.ResponseWriter, err error) {
	logger.Error(
		err.Error(),
		slog.Any("error", err),
	)
	http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
}

func LogClientError(logger *slog.Logger, err error) {
	logger.Error(
		err.Error(),
		slog.Any("error", err),
	)
}
