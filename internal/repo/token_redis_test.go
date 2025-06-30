package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestTokenRedis_SetToken(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	r := NewTokenRedis(rdb)

	type mockBehavior func(userId uuid.UUID, secret string, ttl time.Duration)

	testTable := []struct {
		name         string
		userId       uuid.UUID
		secret       string
		ttl          time.Duration
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name:   "OK",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			secret: "ssJdWpFSBP",
			ttl:    time.Minute * 5,
			mockBehavior: func(userId uuid.UUID, secret string, ttl time.Duration) {
				mock.ExpectSet(tokenPrefix+userId.String(), secret, ttl).SetVal("OK")
			},
			wantErr: false,
		},
		{
			name:   "some error",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			secret: "ssJdWpFSBP",
			ttl:    time.Minute * 5,
			mockBehavior: func(userId uuid.UUID, secret string, ttl time.Duration) {
				mock.ExpectSet(tokenPrefix+userId.String(), secret, ttl).SetErr(someErr)
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userId, testCase.secret, testCase.ttl)
			err := r.SetToken(context.Background(), testCase.userId, testCase.secret, testCase.ttl)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTokenRedis_SetRefreshToken(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	r := NewTokenRedis(rdb)

	type mockBehavior func(userId uuid.UUID, secret string, ttl time.Duration)

	testTable := []struct {
		name         string
		userId       uuid.UUID
		secret       string
		ttl          time.Duration
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name:   "OK",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			secret: "ssJdWpFSBP",
			ttl:    time.Hour * 5,
			mockBehavior: func(userId uuid.UUID, secret string, ttl time.Duration) {
				mock.ExpectSet(refreshTokenPrefix+userId.String(), secret, ttl).SetVal("OK")
			},
			wantErr: false,
		},
		{
			name:   "some error",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			secret: "ssJdWpFSBP",
			ttl:    time.Hour * 5,
			mockBehavior: func(userId uuid.UUID, secret string, ttl time.Duration) {
				mock.ExpectSet(refreshTokenPrefix+userId.String(), secret, ttl).SetErr(someErr)
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userId, testCase.secret, testCase.ttl)
			err := r.SetRefreshToken(context.Background(), testCase.userId, testCase.secret, testCase.ttl)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTokenRedis_DeleteToken(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	r := NewTokenRedis(rdb)

	type mockBehavior func(userId uuid.UUID)

	testTable := []struct {
		name         string
		userId       uuid.UUID
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name:   "OK",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID) {
				mock.ExpectDel(tokenPrefix + userId.String()).SetVal(1)
			},
			wantErr: false,
		},
		{
			name:   "some error",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID) {
				mock.ExpectDel(tokenPrefix + userId.String()).SetErr(someErr)
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userId)
			err := r.DeleteToken(context.Background(), testCase.userId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTokenRedis_DeleteRefreshToken(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	r := NewTokenRedis(rdb)

	type mockBehavior func(userId uuid.UUID)

	testTable := []struct {
		name         string
		userId       uuid.UUID
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name:   "OK",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID) {
				mock.ExpectDel(refreshTokenPrefix + userId.String()).SetVal(1)
			},
			wantErr: false,
		},
		{
			name:   "some error",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID) {
				mock.ExpectDel(refreshTokenPrefix + userId.String()).SetErr(someErr)
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userId)
			err := r.DeleteRefreshToken(context.Background(), testCase.userId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTokenRedis_GetToken(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	r := NewTokenRedis(rdb)

	type mockBehavior func(userId uuid.UUID, secret string)

	testTable := []struct {
		name           string
		userId         uuid.UUID
		mockBehavior   mockBehavior
		expectedSecret string
		expectedErr    error
	}{
		{
			name:   "OK",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID, secret string) {
				mock.ExpectGet(tokenPrefix + userId.String()).SetVal(secret)
			},
			expectedSecret: "ssJdWpFSBP",
			expectedErr:    nil,
		},
		{
			name:   "token not found",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID, secret string) {
				mock.ExpectGet(tokenPrefix + userId.String()).SetErr(redis.Nil)
			},
			expectedSecret: "",
			expectedErr:    ErrTokenNotFound,
		},
		{
			name:   "some error",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID, secret string) {
				mock.ExpectGet(tokenPrefix + userId.String()).SetErr(someErr)
			},
			expectedSecret: "",
			expectedErr:    someErr,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userId, testCase.expectedSecret)
			actualSecret, err := r.GetTokenByUserId(context.Background(), testCase.userId)
			if testCase.expectedErr != nil {
				assert.True(t, errors.Is(testCase.expectedErr, err))
				assert.Equal(t, actualSecret, "")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedSecret, actualSecret)
			}
		})
	}
}

func TestTokenRedis_GetRefreshToken(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	r := NewTokenRedis(rdb)

	type mockBehavior func(userId uuid.UUID, secret string)

	testTable := []struct {
		name           string
		userId         uuid.UUID
		mockBehavior   mockBehavior
		expectedSecret string
		expectedErr    error
	}{
		{
			name:   "OK",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID, secret string) {
				mock.ExpectGet(refreshTokenPrefix + userId.String()).SetVal(secret)
			},
			expectedSecret: "ssJdWpFSBP",
			expectedErr:    nil,
		},
		{
			name:   "token not found",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID, secret string) {
				mock.ExpectGet(refreshTokenPrefix + userId.String()).SetErr(redis.Nil)
			},
			expectedSecret: "",
			expectedErr:    ErrTokenNotFound,
		},
		{
			name:   "some error",
			userId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(userId uuid.UUID, secret string) {
				mock.ExpectGet(refreshTokenPrefix + userId.String()).SetErr(someErr)
			},
			expectedSecret: "",
			expectedErr:    someErr,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userId, testCase.expectedSecret)
			actualSecret, err := r.GetRefreshTokenByUserId(context.Background(), testCase.userId)
			if testCase.expectedErr != nil {
				assert.True(t, errors.Is(testCase.expectedErr, err))
				assert.Equal(t, actualSecret, "")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedSecret, actualSecret)
			}
		})
	}
}
