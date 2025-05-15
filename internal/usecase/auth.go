package usecase

import (
	"context"
	"time"

	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type auth struct {
	UserRepo        repo.User
	TokenRepo       repo.Token
	tokenTTL        time.Duration
	refreshTokenTTL time.Duration
	hashSalt        string
}

func NewAuth(u repo.User, t repo.Token, tokenTTL time.Duration, refreshTokenTTL time.Duration, hashSalt string) Auth {
	return &auth{
		UserRepo:        u,
		TokenRepo:       t,
		tokenTTL:        tokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		hashSalt:        hashSalt,
	}
}
func (a *auth) Login(ctx context.Context, email, password string) (string, string, error) {
	return "", "", ErrNotImplemented
}
func (a *auth) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (a *auth) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	return "", "", ErrNotImplemented
}
