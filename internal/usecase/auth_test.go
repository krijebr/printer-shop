package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
	mock_repo "github.com/krijebr/printer-shop/internal/repo/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Register(t *testing.T) {
	type mockBehavior func(s *mock_repo.MockUser, ctx context.Context, user entity.User)
	var createdUser entity.User
	testTable := []struct {
		name         string
		inputUser    entity.User
		mockBehavior mockBehavior
		expectedErr  error
	}{
		{
			name: "OK",
			inputUser: entity.User{
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
			},
			mockBehavior: func(s *mock_repo.MockUser, ctx context.Context, user entity.User) {
				s.EXPECT().GetByEmail(ctx, user.Email).Return(nil, repo.ErrUserNotFound)
				s.EXPECT().Create(ctx, gomock.AssignableToTypeOf(entity.User{})).
					DoAndReturn(func(ctx context.Context, user entity.User) error {
						createdUser = user
						return nil
					})
				s.EXPECT().GetById(ctx, gomock.AssignableToTypeOf(uuid.UUID{})).
					DoAndReturn(func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
						return &createdUser, nil
					})
			},
			expectedErr: nil,
		},
		{
			name: "user with email already exists",
			inputUser: entity.User{
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
			},
			mockBehavior: func(s *mock_repo.MockUser, ctx context.Context, user entity.User) {
				s.EXPECT().GetByEmail(ctx, user.Email).Return(&entity.User{}, nil)
			},
			expectedErr: ErrEmailAlreadyExists,
		},
		{
			name: "getByEmail error",
			inputUser: entity.User{
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
			},
			mockBehavior: func(s *mock_repo.MockUser, ctx context.Context, user entity.User) {
				s.EXPECT().GetByEmail(ctx, user.Email).Return(nil, errors.New("some error"))
			},
			expectedErr: errors.New("some error"),
		},
		{
			name: "create user error",
			inputUser: entity.User{
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
			},
			mockBehavior: func(s *mock_repo.MockUser, ctx context.Context, user entity.User) {
				s.EXPECT().GetByEmail(ctx, user.Email).Return(nil, repo.ErrUserNotFound)
				s.EXPECT().Create(ctx, gomock.AssignableToTypeOf(entity.User{})).Return(errors.New("some error"))
			},
			expectedErr: errors.New("some error"),
		},
		{
			name: "get user error",
			inputUser: entity.User{
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
			},
			mockBehavior: func(s *mock_repo.MockUser, ctx context.Context, user entity.User) {
				s.EXPECT().GetByEmail(ctx, user.Email).Return(nil, repo.ErrUserNotFound)
				s.EXPECT().Create(ctx, gomock.AssignableToTypeOf(entity.User{})).Return(nil)
				s.EXPECT().GetById(ctx, gomock.AssignableToTypeOf(uuid.UUID{})).Return(nil, errors.New("some error"))
			},
			expectedErr: errors.New("some error"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_repo.NewMockUser(c)
			token := mock_repo.NewMockToken(c)
			testCase.mockBehavior(auth, context.Background(), testCase.inputUser)

			authUsecase := NewAuth(auth, token, 0, 0, "")

			actualUser, err := authUsecase.Register(context.Background(), testCase.inputUser)

			if testCase.expectedErr != nil {
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, actualUser)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, actualUser.Id)
				assert.Equal(t, testCase.inputUser.FirstName, actualUser.FirstName)
				assert.Equal(t, testCase.inputUser.LastName, actualUser.LastName)
				assert.Equal(t, testCase.inputUser.Email, actualUser.Email)
				assert.NotEqual(t, "", actualUser.PasswordHash)
				assert.Equal(t, entity.UserStatusActive, actualUser.Status)
				assert.Equal(t, entity.UserRoleCustomer, actualUser.Role)
				assert.NotEqual(t, 0, actualUser.CreatedAt)
			}
		})
	}
}
