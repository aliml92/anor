package session

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	defaultCtxTimeout = 2 * time.Second
	defaultKeyPrefix  = "keys:app:session:user:authenticated:"
)

type RedisStore struct {
	client     *redis.Client
	prefix     string
	ctxTimeout time.Duration
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client, prefix: defaultKeyPrefix, ctxTimeout: defaultCtxTimeout}
}

func (s *RedisStore) WithPrefix(prefix string) *RedisStore {
	s.prefix = prefix
	return s
}

func (s *RedisStore) Find(token string) (b []byte, exists bool, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	cmd := s.client.Get(ctx, s.prefix+token)
	err = cmd.Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, err
	}

	b, err = cmd.Bytes()
	if err != nil {
		return nil, false, err
	}

	return b, true, nil
}

func (s *RedisStore) Commit(token string, b []byte, expiry time.Time) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	cmd := s.client.Set(ctx, s.prefix+token, b, expiry.Sub(time.Now()))
	return cmd.Err()
}

func (s *RedisStore) Delete(token string) (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	err = s.client.Del(ctx, s.prefix+token).Err()
	return err
}

func (s *RedisStore) All() (map[string][]byte, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*s.ctxTimeout)
	defer cancel()

	keys, err := s.client.Keys(ctx, s.prefix+"*").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	sessions := make(map[string][]byte, len(keys))
	for _, key := range keys {
		token := key[len(s.prefix)+1:]
		data, exists, err := s.Find(token)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, nil
			}
			return nil, err
		}
		if exists {
			sessions[key] = data
		}
	}

	return sessions, nil
}

func (s *RedisStore) FindCtx(ctx context.Context, token string) (b []byte, exists bool, err error) {
	cmd := s.client.Get(ctx, s.prefix+token)
	err = cmd.Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, err
	}

	b, err = cmd.Bytes()
	if err != nil {
		return nil, false, err
	}

	return b, true, nil
}

func (s *RedisStore) CommitCtx(ctx context.Context, token string, b []byte, expiry time.Time) error {
	cmd := s.client.Set(ctx, s.prefix+token, b, expiry.Sub(time.Now()))
	return cmd.Err()
}

func (s *RedisStore) DeleteCtx(ctx context.Context, token string) (err error) {
	err = s.client.Del(ctx, s.prefix+token).Err()
	return err
}

func (s *RedisStore) AllCtx(ctx context.Context) (map[string][]byte, error) {
	keys, err := s.client.Keys(ctx, s.prefix+"*").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	sessions := make(map[string][]byte, len(keys))
	for _, key := range keys {
		token := key[len(s.prefix)+1:]
		data, exists, err := s.FindCtx(ctx, token)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, nil
			}
			return nil, err
		}
		if exists {
			sessions[key] = data
		}
	}

	return sessions, nil
}
