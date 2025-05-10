package usecase

import (
	"context"

	"github.com/krijebr/printer-shop/internal/entity"
)

type auth struct {
}

func NewAuth() Auth {
	return &auth{}
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
