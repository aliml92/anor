package cache

import (
	"context"
	"time"
)

type ResetPasswordTokenCacher interface {
	CacheToken(ctx context.Context, token string, userID int64, ttl time.Duration) error
	TokenExists(ctx context.Context, token string) (bool, error)
	GetUserIDByToken(ctx context.Context, token string) (int64, error)
}
