package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliml92/anor/cache/auth"
	"github.com/redis/go-redis/v9"
	"time"
)

var _ auth.SignupConfirmationOTPCache = (*SignupConfirmationOTPCacher)(nil)

const basePrefixSignupOTP = "keys:app:signup_otp"

type SignupConfirmationOTPCacher struct {
	client *redis.Client
}

func NewSignupConfirmationOTPCacher(client *redis.Client) *SignupConfirmationOTPCacher {
	return &SignupConfirmationOTPCacher{client: client}
}

func (c *SignupConfirmationOTPCacher) CacheOTP(ctx context.Context, otp string, email string, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s", basePrefixSignupOTP, otp)
	err := c.client.Set(ctx, key, email, ttl).Err()
	return err
}

func (c *SignupConfirmationOTPCacher) GetEmailByOTP(ctx context.Context, otp string) (string, error) {
	key := fmt.Sprintf("%s:%s", basePrefixSignupOTP, otp)
	cmd := c.client.Get(ctx, key)
	err := cmd.Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}

	return cmd.Val(), nil
}
