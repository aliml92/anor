package redis

import (
	"context"
	"github.com/aliml92/anor/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Host + ":" + cfg.Port,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MinIdleConns: cfg.MinIdleConns,
		MaxIdleConns: cfg.MaxIdleConns,
	})

	statusCmd := rdb.Ping(ctx)
	if statusCmd.Err() != nil {
		return nil, statusCmd.Err()
	}

	return rdb, nil
}
