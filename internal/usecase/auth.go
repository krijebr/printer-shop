package usecase
import (
	"context"
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
)
type auth struct{	
}
type NewAuth() Auth {
	return &auth{}
}
func (a *auth) Login(ctx context.Context, email, password string) (string, string, error){
	return "","", err
}
func (a *auth)	ValidateToken(ctx context.Context, token string) (*entity.User, error){
	return nil, err
}
func (a *auth)	RefreshToken(ctx context.Context, refreshToken string) (string, string, error){
	return "","",err
}
