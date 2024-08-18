package session

import (
	"context"
)

const UserKey = "userKey"

const (
	IsAuthKey                    = "is_auth"
	UserIDKey                    = "user_id"
	UserFirstnameKey             = "firstname"
	CartIDKey                    = "cart_id"
	ShippingAddressIDKey         = "shipping_address_id"
	BillingAddressIDKey          = "billing_address_id"
	PaymentMethodKey             = "payment_method"
	StripeConfirmationTokenIDKey = "stripe_confirmation_token_id"
)

type User struct {
	ID                        int64
	IsAuth                    bool
	Firstname                 string
	CartID                    int64
	ShippingAddressID         int64
	BillingAddressID          int64
	PaymentMethod             string
	StripeConfirmationTokenID string
}

// NewContextWithUser NewContext returns a new Context that carries value u.
func NewContextWithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, UserKey, u)
}

// UserFromContext FromContext returns the User value stored in ctx, if any.
func UserFromContext(ctx context.Context) *User {
	u, ok := ctx.Value(UserKey).(*User)
	if !ok {
		return nil
	}
	return u
}
