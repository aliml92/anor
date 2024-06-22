package anor

import "errors"

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
