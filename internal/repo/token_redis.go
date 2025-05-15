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

const (
	tokenPrefix        string = "token_"
	refreshTokenPrefix string = "refresh_"
)

func NewTokenRedis(rdb *redis.Client) Token {
	return &TokenRedis{
		rdb: rdb,
	}
}
func (a *TokenRedis) SetToken(ctx context.Context, userId uuid.UUID, secret string, ttl time.Duration) error {
	return a.rdb.Set(ctx, tokenPrefix+userId.String(), secret, ttl).Err()

}
func (a *TokenRedis) SetRefreshToken(ctx context.Context, userId uuid.UUID, secret string, ttl time.Duration) error {
	return a.rdb.Set(ctx, refreshTokenPrefix+userId.String(), secret, ttl).Err()
}
func (a *TokenRedis) GetTokenByUserId(ctx context.Context, userId uuid.UUID) (string, error) {
	result, err := a.rdb.Get(ctx, tokenPrefix+userId.String()).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
func (a *TokenRedis) GetRefreshTokenByUserId(ctx context.Context, userId uuid.UUID) (string, error) {
	result, err := a.rdb.Get(ctx, refreshTokenPrefix+userId.String()).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
func (a *TokenRedis) DeleteToken(ctx context.Context, userId uuid.UUID) error {
	_, err := a.rdb.Del(ctx, tokenPrefix+userId.String()).Result()
	if err != nil {
		return err
	}
	return nil
}
func (a *TokenRedis) DeleteRefreshToken(ctx context.Context, userId uuid.UUID) error {
	_, err := a.rdb.Del(ctx, refreshTokenPrefix+userId.String()).Result()
	if err != nil {
		return err
	}
	return nil
}
