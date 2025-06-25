package v1

import (
	"bytes"
	"context"
	"fmt"
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
				s.EXPECT().Register(ctx, user).Return(nil, fmt.Errorf("some error"))
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
