package usecase

import (
	"context"

	"time"

	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
)

type auth struct {
	userRepo        repo.User
	tokenRepo       repo.Token
	tokenTTL        time.Duration
	refreshTokenTTL time.Duration
	userUseCase     User
}

func NewAuth(u repo.User, t repo.Token, tokenTTL time.Duration, refreshTokenTTL time.Duration, user User) Auth {
	return &auth{
		userRepo:        u,
		tokenRepo:       t,
		tokenTTL:        tokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		userUseCase:     user,
	}
}
func (a *auth) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		switch {
		case err == repo.ErrUserNotFound:
			return "", "", ErrUserNotFound
		default:
			return "", "", err
		}
	}
	if a.userUseCase.ValidatePassword(password, user.PasswordHash) {
		return "", "", nil
	}
	return "", "", ErrWrongPassword
}
func (a *auth) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	return nil, ErrNotImplemented
}
func (a *auth) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	return "", "", ErrNotImplemented
}
