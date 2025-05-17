package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"unsafe"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
func (a *auth) generateRandomKey() string {
	const (
		n             = 10                                                     // key length
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // dictionary
		letterIdxBits = 6                                                      // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits                                     // # of letter indices fitting in 63 bits
	)

	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
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
	if !a.userUseCase.ValidatePassword(password, user.PasswordHash) {
		return "", "", ErrWrongPassword
	}
	secret := a.generateRandomKey()
	expTime := time.Now().Add(a.tokenTTL)
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": user.Id.String(),
		"exp": expTime.Unix(),
	})
	err = a.tokenRepo.SetToken(ctx, user.Id, secret, a.tokenTTL)
	if err != nil {
		return "", "", err
	}
	token, err := tokenObj.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	secret = a.generateRandomKey()
	expTime = time.Now().Add(a.refreshTokenTTL)
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": user.Id.String(),
		"exp": expTime.Unix(),
	})
	err = a.tokenRepo.SetRefreshToken(ctx, user.Id, secret, a.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := refreshTokenObj.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
func (a *auth) ValidateToken(ctx context.Context, token string) (*entity.User, error) {

	jwToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		userId, err := token.Claims.GetIssuer()
		if err != nil {
			return nil, ErrInvalidToken
		}
		uid, err := uuid.Parse(userId)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		secret, err := a.tokenRepo.GetTokenByUserId(ctx, uid)
		if err != nil {
			if errors.Is(err, repo.ErrTokenNotFound) {
				return nil, ErrInvalidToken
			}
			return nil, err
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	userId, err := jwToken.Claims.GetIssuer()
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	user, err := a.userRepo.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (a *auth) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	return "", "", ErrNotImplemented
}
