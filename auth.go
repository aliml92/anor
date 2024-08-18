package anor

import (
	"context"
)

type AuthService interface {
	Signup(ctx context.Context, name, email, password string) (User, error)
	SignupConfirm(ctx context.Context, otp, email string) error
	Signin(ctx context.Context, email, password string) (User, error)
	SigninWithGoogle(ctx context.Context, email, name string) (User, error)
	ResendOTP(ctx context.Context, email string) error
	SendResetPasswordLink(ctx context.Context, email string) error
	VerifyResetPasswordToken(ctx context.Context, token string) (bool, error)
	ResetPassword(ctx context.Context, token string, password string) error
	GetUser(ctx context.Context, id int64) (*User, error)
}
