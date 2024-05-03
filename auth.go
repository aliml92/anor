package anor

import (
	"context"
)

type AuthService interface {
	Signup(ctx context.Context, name, email, password string) error
	SignupConfirm(ctx context.Context, otp, email string) error
	Signin(ctx context.Context, email, password string) (int64, error)
	ResendOTP(ctx context.Context, email string) error
	SendResetPasswordLink(ctx context.Context, email string) error
	VerifyResetPasswordToken(ctx context.Context, token string) (bool, error)
	ResetPassword(ctx context.Context, token string, password string) error
}
