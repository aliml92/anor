package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliml92/anor/cache/auth"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var _ auth.ResetPasswordTokenCache = (*ResetPasswordTokenCacher)(nil)

const basePrefixResetToken = "keys:app:reset_token"

type ResetPasswordTokenCacher struct {
	client *redis.Client
}

func NewResetPasswordTokenCacher(client *redis.Client) *ResetPasswordTokenCacher {
	return &ResetPasswordTokenCacher{client: client}
}

func (c *ResetPasswordTokenCacher) CacheToken(ctx context.Context, token string, userID int64, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s", basePrefixResetToken, token)
	err := c.client.Set(ctx, key, userID, ttl).Err()
	return err
}

func (c *ResetPasswordTokenCacher) TokenExists(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("%s:%s", basePrefixResetToken, token)
	cmd := c.client.Get(ctx, key)
	err := cmd.Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *ResetPasswordTokenCacher) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	key := fmt.Sprintf("%s:%s", basePrefixResetToken, token)
	cmd := c.client.Get(ctx, key)
	err := cmd.Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}

	return strconv.ParseInt(cmd.Val(), 10, 64)
}
