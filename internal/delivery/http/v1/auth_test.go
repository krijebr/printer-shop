package v1

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/usecase"
	mock_usecase "github.com/krijebr/printer-shop/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var someErr = errors.New("some error")

func TestAuthHandlers_register(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockAuth, ctx context.Context, user entity.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            entity.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"first_name":"Ivan","last_name":"Ivanov","email":"ivan@gmail.com","password":"12345678910"}`,
			inputUser: entity.User{
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
			},
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, user entity.User) {
				s.EXPECT().Register(ctx, user).Return(&entity.User{
					Id:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					FirstName:    "Ivan",
					LastName:     "Ivanov",
					Email:        "ivan@gmail.com",
					PasswordHash: "12345678910",
					Status:       entity.UserStatusActive,
					Role:         entity.UserRoleCustomer,
					CreatedAt:    time.Date(2025, 6, 25, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":"00000000-0000-0000-0000-000000000001","first_name":"Ivan","last_name":"Ivanov","email":"ivan@gmail.com","status":"active","role":"customer","created_at":"2025-06-25T00:00:00Z"}`,
		},
		{
			name:      "User with email already exists",
			inputBody: `{"first_name":"Ivan","last_name":"Ivanov","email":"ivan@gmail.com","password":"12345678910"}`,
			inputUser: entity.User{
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
			},
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, user entity.User) {
				s.EXPECT().Register(ctx, user).Return(nil, usecase.ErrEmailAlreadyExists)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":9,"message":"user with this email already exists"}`,
		},
		{
			name:      "some internal error",
			inputBody: `{"first_name":"Ivan","last_name":"Ivanov","email":"ivan@gmail.com","password":"12345678910"}`,
			inputUser: entity.User{
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
			},
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, user entity.User) {
				s.EXPECT().Register(ctx, user).Return(nil, someErr)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":4,"message":"internal error"}`,
		},
		{
			name:      "empty input",
			inputBody: `{"first_name":"","last_name":"","email":"","password":""}`,
			inputUser: entity.User{
				FirstName:    "",
				LastName:     "",
				Email:        "",
				PasswordHash: "",
			},
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, user entity.User) {
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":8,"message":"validation error"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_usecase.NewMockAuth(c)
			testCase.mockBehavior(auth, context.Background(), testCase.inputUser)

			a := NewAuthHandlers(auth)

			r := echo.New()
			r.POST("/auth", a.register())

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.TrimRight(w.Body.String(), "\n"))
		})
	}
}

func TestAuthHandlers_login(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockAuth, ctx context.Context, email string, password string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputEmail           string
		inputPassword        string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "OK",
			inputBody:     `{"email":"ivan@gmail.com","password":"12345678910"}`,
			inputEmail:    "ivan@gmail.com",
			inputPassword: "12345678910",
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, email string, password string) {
				s.EXPECT().Login(ctx, email, password).
					Return(
						"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTExOTMyMzYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9._odGZZsBspT5OopI3ei8kZGBoIca-opL-N19zcf-X5U",
						"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTExOTY1MzYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.o_HRawWX6dWsTGab3CSktgpoU4jR1zdE1DdangzbiSE",
						nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTExOTMyMzYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9._odGZZsBspT5OopI3ei8kZGBoIca-opL-N19zcf-X5U","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTExOTY1MzYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.o_HRawWX6dWsTGab3CSktgpoU4jR1zdE1DdangzbiSE"}`,
		},
		{
			name:          "wrong email",
			inputBody:     `{"email":"ivan1@gmail.com","password":"12345678910"}`,
			inputEmail:    "ivan1@gmail.com",
			inputPassword: "12345678910",
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, email string, password string) {
				s.EXPECT().Login(ctx, email, password).
					Return("", "", usecase.ErrUserNotFound)
			},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"error":11,"message":"wrong email or password"}`,
		},
		{
			name:          "wrong password",
			inputBody:     `{"email":"ivan@gmail.com","password":"1234567891"}`,
			inputEmail:    "ivan@gmail.com",
			inputPassword: "1234567891",
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, email string, password string) {
				s.EXPECT().Login(ctx, email, password).
					Return("", "", usecase.ErrWrongPassword)
			},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"error":11,"message":"wrong email or password"}`,
		},
		{
			name:          "empty input",
			inputBody:     `{"email":"","password":""}`,
			inputEmail:    "",
			inputPassword: "",
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, email string, password string) {
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":8,"message":"validation error"}`,
		},
		{
			name:          "some error",
			inputBody:     `{"email":"ivan@gmail.com","password":"12345678910"}`,
			inputEmail:    "ivan@gmail.com",
			inputPassword: "12345678910",
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, email string, password string) {
				s.EXPECT().Login(ctx, email, password).
					Return("", "", someErr)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":4,"message":"internal error"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_usecase.NewMockAuth(c)
			testCase.mockBehavior(auth, context.Background(), testCase.inputEmail, testCase.inputPassword)

			a := NewAuthHandlers(auth)

			r := echo.New()
			r.POST("/auth", a.login())

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.TrimRight(w.Body.String(), "\n"))
		})
	}
}

func TestAuthHandlers_refreshTokens(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockAuth, ctx context.Context, refreshToken string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputRefreshToken    string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:              "OK",
			inputBody:         `{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTEyOTMwMjYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.nPVpyjwPpK2pmiStpjccmKdVdR2DRxqYsbW1gFiZzCc"}`,
			inputRefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTEyOTMwMjYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.nPVpyjwPpK2pmiStpjccmKdVdR2DRxqYsbW1gFiZzCc",
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, refreshToken string) {
				s.EXPECT().RefreshToken(ctx, refreshToken).
					Return(
						"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTExOTMyMzYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9._odGZZsBspT5OopI3ei8kZGBoIca-opL-N19zcf-X5U",
						"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTExOTY1MzYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.o_HRawWX6dWsTGab3CSktgpoU4jR1zdE1DdangzbiSE",
						nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTExOTMyMzYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9._odGZZsBspT5OopI3ei8kZGBoIca-opL-N19zcf-X5U","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTExOTY1MzYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.o_HRawWX6dWsTGab3CSktgpoU4jR1zdE1DdangzbiSE"}`,
		},
		{
			name:              "invalid refresh token",
			inputBody:         `{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTEyOTMwMjYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.nPVpyjwPpK2pmiStpjccmKdVdR2DRxqYsbW1gFiZzCq"}`,
			inputRefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTEyOTMwMjYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.nPVpyjwPpK2pmiStpjccmKdVdR2DRxqYsbW1gFiZzCq",
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, refreshToken string) {
				s.EXPECT().RefreshToken(ctx, refreshToken).
					Return("", "", usecase.ErrInvalidToken)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":2,"message":"invalid refresh token"}`,
		},
		{
			name:              "empty input",
			inputBody:         `{""}`,
			inputRefreshToken: "",
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, refreshToken string) {
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":7,"message":"invalid request"}`,
		},
		{
			name:              "OK",
			inputBody:         `{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTEyOTMwMjYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.nPVpyjwPpK2pmiStpjccmKdVdR2DRxqYsbW1gFiZzCc"}`,
			inputRefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTEyOTMwMjYsImlzcyI6IjYxMmRkOTcxLWQ3NDMtNDU1Mi05MDZlLWM3ZjJlZGI5YmM3YyJ9.nPVpyjwPpK2pmiStpjccmKdVdR2DRxqYsbW1gFiZzCc",
			mockBehavior: func(s *mock_usecase.MockAuth, ctx context.Context, refreshToken string) {
				s.EXPECT().RefreshToken(ctx, refreshToken).
					Return("", "", someErr)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":4,"message":"internal error"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_usecase.NewMockAuth(c)
			testCase.mockBehavior(auth, context.Background(), testCase.inputRefreshToken)

			a := NewAuthHandlers(auth)

			r := echo.New()
			r.POST("/auth", a.refreshTokens())

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.TrimRight(w.Body.String(), "\n"))
		})
	}
}
