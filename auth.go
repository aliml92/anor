package anor

import (
	"context"
)

type AuthService interface {
	Signup(ctx context.Context, name, email, password string) error
	SignupConfirm(ctx context.Context, otp, email string) error
	Signin(ctx context.Context, email, password string) (int64, error)
}
