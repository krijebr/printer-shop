package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenRedis struct {
	rdb *redis.Client
}

func NewAuthRedis(rdb *redis.Client) Token {
	return &TokenRedis{
		rdb: rdb,
	}
}
func (a *TokenRedis) SetToken(ctx context.Context, userId uuid.UUID, secret string, ttl time.Duration) error {
	err := a.rdb.Set(ctx, "token_"+userId.String(), secret, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}
func (a *TokenRedis) SetRefreshToken(ctx context.Context, userId uuid.UUID, secret string, ttl time.Duration) error {
	err := a.rdb.Set(ctx, "refresh_"+userId.String(), secret, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}
func (a *TokenRedis) GetTokenByUserId(ctx context.Context, userId uuid.UUID) (string, error) {
	result, err := a.rdb.Get(ctx, "token"+userId.String()).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
func (a *TokenRedis) GetRefreshTokenByUserId(ctx context.Context, userId uuid.UUID) (string, error) {
	result, err := a.rdb.Get(ctx, "refresh"+userId.String()).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
func (a *TokenRedis) DeleteToken(ctx context.Context, userId uuid.UUID) error {
	_, err := a.rdb.Del(ctx, "token_"+userId.String()).Result()
	if err != nil {
		return err
	}
	return nil
}
func (a *TokenRedis) DeleteRefreshToken(ctx context.Context, userId uuid.UUID) error {
	_, err := a.rdb.Del(ctx, "refresh_"+userId.String()).Result()
	if err != nil {
		return err
	}
	return nil
}
