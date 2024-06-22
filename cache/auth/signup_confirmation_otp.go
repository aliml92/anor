package auth

import (
	"context"
	"time"
)

type SignupConfirmationOTPCache interface {
	CacheOTP(ctx context.Context, otp string, email string, ttl time.Duration) error
	GetEmailByOTP(ctx context.Context, otp string) (string, error)
}
