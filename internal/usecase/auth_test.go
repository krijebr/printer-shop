package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/repo"
	mock_repo "github.com/krijebr/printer-shop/internal/repo/mocks"
	"github.com/stretchr/testify/assert"
)

var someErr = errors.New("some error")

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
				s.EXPECT().GetByEmail(ctx, user.Email).Return(nil, someErr)
			},
			expectedErr: someErr,
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
				s.EXPECT().Create(ctx, gomock.AssignableToTypeOf(entity.User{})).Return(someErr)
			},
			expectedErr: someErr,
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
				s.EXPECT().GetById(ctx, gomock.AssignableToTypeOf(uuid.UUID{})).Return(nil, someErr)
			},
			expectedErr: someErr,
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

func TestAuth_HashPassword(t *testing.T) {
	firstSaltWord := "first_salt_word"
	secondSaltWord := "second_salt_word"
	firstAuthUsecase := NewAuth(nil, nil, 0, 0, firstSaltWord)
	secondAuthUsecase := NewAuth(nil, nil, 0, 0, secondSaltWord)
	firstPassword := "firstPassword"
	secondPassword := "secondPassword"
	firstPasswordHash := firstAuthUsecase.HashPassword(firstPassword)
	secondPasswordHash := firstAuthUsecase.HashPassword(secondPassword)
	thirdPasswordHash := secondAuthUsecase.HashPassword(firstPassword)
	t.Run("password hash is not empty string", func(t *testing.T) {
		assert.NotEqual(t, firstPasswordHash, "")
	})
	t.Run("password hash is not equal to password", func(t *testing.T) {
		assert.NotEqual(t, firstPasswordHash, firstPassword)
	})
	t.Run("hashes of two different passwords are not equal", func(t *testing.T) {
		assert.NotEqual(t, firstPasswordHash, secondPasswordHash)
	})
	t.Run("two hashes of the same password, generated with different salt words, are not equal", func(t *testing.T) {
		assert.NotEqual(t, firstPasswordHash, thirdPasswordHash)
	})
}
func TestAuth_generateRandomKey(t *testing.T) {
	var authUsecase auth
	firstRandomKey := authUsecase.generateRandomKey()
	secondRandomKey := authUsecase.generateRandomKey()
	t.Run("random key is not empty string", func(t *testing.T) {
		assert.NotEqual(t, firstRandomKey, "")
	})
	t.Run("two random keys aren't equal", func(t *testing.T) {
		assert.NotEqual(t, firstRandomKey, secondRandomKey)
	})
}

func TestAuth_ValidatePassword(t *testing.T) {
	firstPassword := "firstPassword"
	secondPassword := "secondPassword"
	authUsecase := NewAuth(nil, nil, 0, 0, "salt")
	firstPasswordHash := authUsecase.HashPassword(firstPassword)
	t.Run("validation of password and it's hash passes", func(t *testing.T) {
		assert.True(t, authUsecase.ValidatePassword(firstPassword, firstPasswordHash))
	})
	t.Run("validation of password and hash of another password doesn't pass", func(t *testing.T) {
		assert.False(t, authUsecase.ValidatePassword(secondPassword, firstPasswordHash))
	})
}

func TestAuth_Login(t *testing.T) {
	type mockUserBehavior func(s *mock_repo.MockUser, ctx context.Context, email string, password string, user entity.User)
	type mockTokenBehavior func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID)
	testTable := []struct {
		name              string
		inputEmail        string
		inputPassword     string
		outputUser        entity.User
		tokenTTL          time.Duration
		refreshTokenTTL   time.Duration
		mockUserBehavior  mockUserBehavior
		mockTokenBehavior mockTokenBehavior
		expectedErr       error
	}{
		{
			name:          "OK",
			inputEmail:    "ivan@gmail.com",
			inputPassword: "12345678910",
			outputUser: entity.User{
				Id:           uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "992320c97d2edc09debf80bc3cd2b770a07a97ecee15771e158a744f38790d2e",
				Status:       entity.UserStatusActive,
				Role:         entity.UserRoleCustomer,
				CreatedAt:    time.Date(2025, 6, 25, 0, 0, 0, 0, time.UTC),
			},
			tokenTTL:        time.Minute * 5,
			refreshTokenTTL: time.Hour * 4,
			mockUserBehavior: func(s *mock_repo.MockUser, ctx context.Context, email string, password string, user entity.User) {
				s.EXPECT().GetByEmail(ctx, email).Return(&user, nil)
			},
			mockTokenBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID) {
				s.EXPECT().SetToken(ctx, userId, gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(time.Duration(0))).Return(nil)
				s.EXPECT().SetRefreshToken(ctx, userId, gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(time.Duration(0))).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:          "user not found",
			inputEmail:    "ivan1@gmail.com",
			inputPassword: "12345678910",
			outputUser: entity.User{
				Id:           uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "992320c97d2edc09debf80bc3cd2b770a07a97ecee15771e158a744f38790d2e",
				Status:       entity.UserStatusActive,
				Role:         entity.UserRoleCustomer,
				CreatedAt:    time.Date(2025, 6, 25, 0, 0, 0, 0, time.UTC),
			},
			tokenTTL:        time.Minute * 5,
			refreshTokenTTL: time.Hour * 4,
			mockUserBehavior: func(s *mock_repo.MockUser, ctx context.Context, email string, password string, user entity.User) {
				s.EXPECT().GetByEmail(ctx, email).Return(nil, repo.ErrUserNotFound)
			},
			mockTokenBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID) {
			},
			expectedErr: ErrUserNotFound,
		},
		{
			name:          "wrong password",
			inputEmail:    "ivan@gmail.com",
			inputPassword: "1234567891",
			outputUser: entity.User{
				Id:           uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "992320c97d2edc09debf80bc3cd2b770a07a97ecee15771e158a744f38790d2e",
				Status:       entity.UserStatusActive,
				Role:         entity.UserRoleCustomer,
				CreatedAt:    time.Date(2025, 6, 25, 0, 0, 0, 0, time.UTC),
			},
			tokenTTL:        time.Minute * 5,
			refreshTokenTTL: time.Hour * 4,
			mockUserBehavior: func(s *mock_repo.MockUser, ctx context.Context, email string, password string, user entity.User) {
				s.EXPECT().GetByEmail(ctx, email).Return(&user, nil)
			},
			mockTokenBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID) {
			},
			expectedErr: ErrWrongPassword,
		},
		{
			name:          "some error",
			inputEmail:    "ivan@gmail.com",
			inputPassword: "12345678910",
			outputUser: entity.User{
				Id:           uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "992320c97d2edc09debf80bc3cd2b770a07a97ecee15771e158a744f38790d2e",
				Status:       entity.UserStatusActive,
				Role:         entity.UserRoleCustomer,
				CreatedAt:    time.Date(2025, 6, 25, 0, 0, 0, 0, time.UTC),
			},
			tokenTTL:        time.Minute * 5,
			refreshTokenTTL: time.Hour * 4,
			mockUserBehavior: func(s *mock_repo.MockUser, ctx context.Context, email string, password string, user entity.User) {
				s.EXPECT().GetByEmail(ctx, email).Return(&user, someErr)
			},
			mockTokenBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID) {
			},
			expectedErr: someErr,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userMock := mock_repo.NewMockUser(c)
			tokenMock := mock_repo.NewMockToken(c)
			testCase.mockUserBehavior(userMock, context.Background(), testCase.inputEmail, testCase.inputPassword, testCase.outputUser)
			testCase.mockTokenBehavior(tokenMock, context.Background(), testCase.outputUser.Id)
			authUsecase := NewAuth(userMock, tokenMock, testCase.tokenTTL, testCase.refreshTokenTTL, "qwerty")

			token, refreshToken, err := authUsecase.Login(context.Background(), testCase.inputEmail, testCase.inputPassword)

			if testCase.expectedErr != nil {
				assert.Equal(t, testCase.expectedErr, err)
				assert.Equal(t, token, "")
				assert.Equal(t, refreshToken, "")
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, token, "")
				assert.NotEqual(t, refreshToken, "")
				assert.NotEqual(t, token, refreshToken)
			}
		})
	}
}

func TestAuth_ValidateToken(t *testing.T) {
	type mockUserBehavior func(s *mock_repo.MockUser, ctx context.Context, user *entity.User)
	type mockTokenBehavior func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID, secret string)
	testTable := []struct {
		name              string
		userId            uuid.UUID
		expectedUser      *entity.User
		mockUserBehavior  mockUserBehavior
		mockTokenBehavior mockTokenBehavior
		expectedErr       error
	}{
		{
			name:   "OK",
			userId: uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
			expectedUser: &entity.User{
				Id:           uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "992320c97d2edc09debf80bc3cd2b770a07a97ecee15771e158a744f38790d2e",
				Status:       entity.UserStatusActive,
				Role:         entity.UserRoleCustomer,
				CreatedAt:    time.Date(2025, 6, 25, 0, 0, 0, 0, time.UTC),
			},
			mockUserBehavior: func(s *mock_repo.MockUser, ctx context.Context, user *entity.User) {
				s.EXPECT().GetById(ctx, user.Id).Return(user, nil)
			},
			mockTokenBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID, secret string) {
				s.EXPECT().GetTokenByUserId(ctx, userId).Return(secret, nil)
			},
			expectedErr: nil,
		},
		{
			name:         "some error",
			userId:       uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
			expectedUser: nil,
			mockUserBehavior: func(s *mock_repo.MockUser, ctx context.Context, user *entity.User) {
			},
			mockTokenBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID, secret string) {
				s.EXPECT().GetTokenByUserId(ctx, userId).Return("", someErr)
			},
			expectedErr: someErr,
		},
		{
			name:         "invalid token",
			userId:       uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
			expectedUser: nil,
			mockUserBehavior: func(s *mock_repo.MockUser, ctx context.Context, user *entity.User) {
			},
			mockTokenBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID, secret string) {
				secretSlice := []byte(secret)
				secretSlice[0] = secretSlice[0] + 1
				secret = string(secretSlice)
				s.EXPECT().GetTokenByUserId(ctx, userId).Return(secret, nil)
			},
			expectedErr: ErrInvalidToken,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userMock := mock_repo.NewMockUser(c)
			tokenMock := mock_repo.NewMockToken(c)

			var authUsecase auth
			authUsecase.userRepo = userMock
			authUsecase.tokenRepo = tokenMock
			authUsecase.tokenTTL = time.Minute * 5
			secret := authUsecase.generateRandomKey()
			expTime := time.Now().Add(authUsecase.tokenTTL)
			tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iss": testCase.userId.String(),
				"exp": expTime.Unix(),
			})

			token, err := tokenObj.SignedString([]byte(secret))

			testCase.mockTokenBehavior(tokenMock, context.Background(), testCase.userId, secret)
			testCase.mockUserBehavior(userMock, context.Background(), testCase.expectedUser)

			actualUser, err := authUsecase.ValidateToken(context.Background(), token)

			if testCase.expectedErr != nil {
				assert.True(t, errors.Is(err, testCase.expectedErr))
				assert.Equal(t, actualUser, testCase.expectedUser)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedUser, actualUser)
			}
		})
	}
}

func TestAuth_RefreshToken(t *testing.T) {
	type mockUserBehavior func(s *mock_repo.MockUser, ctx context.Context)
	type mockTokenBehavior func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID, secret string)
	testTable := []struct {
		name         string
		userId       uuid.UUID
		mockBehavior mockTokenBehavior
		expectedErr  error
	}{
		{
			name:   "OK",
			userId: uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
			mockBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID, secret string) {
				s.EXPECT().GetRefreshTokenByUserId(ctx, userId).Return(secret, nil)
				s.EXPECT().SetToken(ctx, userId, gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(time.Duration(0))).Return(nil)
				s.EXPECT().SetRefreshToken(ctx, userId, gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(time.Duration(0))).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "some error",
			userId: uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
			mockBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID, secret string) {
				s.EXPECT().GetRefreshTokenByUserId(ctx, userId).Return("", someErr)
			},
			expectedErr: someErr,
		},
		{
			name:   "invalid token",
			userId: uuid.MustParse("8be456fa-aa6b-4310-b321-2cacfb8193a9"),
			mockBehavior: func(s *mock_repo.MockToken, ctx context.Context, userId uuid.UUID, secret string) {
				secretSlice := []byte(secret)
				secretSlice[0] = secretSlice[0] + 1
				secret = string(secretSlice)
				s.EXPECT().GetRefreshTokenByUserId(ctx, userId).Return(secret, nil)
			},
			expectedErr: ErrInvalidToken,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userMock := mock_repo.NewMockUser(c)
			tokenMock := mock_repo.NewMockToken(c)

			var authUsecase auth
			authUsecase.userRepo = userMock
			authUsecase.tokenRepo = tokenMock
			authUsecase.refreshTokenTTL = time.Hour * 5
			secret := authUsecase.generateRandomKey()
			expTime := time.Now().Add(authUsecase.refreshTokenTTL)
			refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iss": testCase.userId.String(),
				"exp": expTime.Unix(),
			})

			refreshToken, err := refreshTokenObj.SignedString([]byte(secret))

			testCase.mockBehavior(tokenMock, context.Background(), testCase.userId, secret)

			token, newRefreshToken, err := authUsecase.RefreshToken(context.Background(), refreshToken)

			if testCase.expectedErr != nil {
				assert.True(t, errors.Is(err, testCase.expectedErr))
				assert.Equal(t, token, "")
				assert.Equal(t, newRefreshToken, "")
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, token, refreshToken)
				assert.NotEqual(t, token, "")
				assert.NotEqual(t, refreshToken, "")
			}
		})
	}
}
