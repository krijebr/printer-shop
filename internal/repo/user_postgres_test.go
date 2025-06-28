package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/stretchr/testify/assert"
)

var someErr = errors.New("some error")

func TestUserRepoPg_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	r := NewUserRepoPg(db)

	type mockBehavior func(ctx context.Context, user entity.User)

	testTable := []struct {
		name         string
		inputUser    entity.User
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			inputUser: entity.User{
				Id:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
				Status:       entity.UserStatusActive,
				Role:         entity.UserRoleCustomer,
				CreatedAt:    time.Date(2025, 6, 25, 0, 0, 0, 0, time.UTC),
			},
			mockBehavior: func(ctx context.Context, user entity.User) {

				mock.ExpectExec("insert into users").
					WithArgs(user.Id, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Status, user.Role, user.CreatedAt).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "user creation error",
			inputUser: entity.User{
				Id:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "12345678910",
				Status:       entity.UserStatusActive,
				Role:         entity.UserRoleCustomer,
				CreatedAt:    time.Date(2025, 6, 25, 0, 0, 0, 0, time.UTC),
			},
			mockBehavior: func(ctx context.Context, user entity.User) {

				mock.ExpectExec("insert into users").
					WithArgs(user.Id, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Status, user.Role, user.CreatedAt).
					WillReturnError(someErr)
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(context.Background(), testCase.inputUser)
			err := r.Create(context.Background(), testCase.inputUser)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestUserRepoPg_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	r := NewUserRepoPg(db)

	type mockBehavior func(ctx context.Context, id uuid.UUID)

	testTable := []struct {
		name         string
		inputUserId  uuid.UUID
		mockBehavior mockBehavior
		expectedUser *entity.User
		expectedErr  error
	}{
		{
			name:        "OK",
			inputUserId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(ctx context.Context, id uuid.UUID) {
				idColumn := sqlmock.NewColumn("id").OfType("uuid", uuid.Nil.String()).Nullable(false)
				firstNameColumn := sqlmock.NewColumn("first_name").OfType("varchar", "Ivan").Nullable(false)
				lastNameColumn := sqlmock.NewColumn("last_name").OfType("varchar", "Ivanov").Nullable(false)
				emailColumn := sqlmock.NewColumn("email").OfType("varchar", "ivan@list.ru").Nullable(false)
				passworHashColumn := sqlmock.NewColumn("password_hash").OfType("varchar", "992320c97d2edc09debf80bc3cd2b770a07a97ecee15771e158a744f38790d2e").Nullable(false)
				statusColumn := sqlmock.NewColumn("status").OfType("varchar", "active").Nullable(false)
				roleColumn := sqlmock.NewColumn("role").OfType("varchar", "customer").Nullable(false)
				createdAtColumn := sqlmock.NewColumn("created_at").OfType("timestamp", "2025-05-19 17:07:13.947").Nullable(false)
				rows := sqlmock.NewRowsWithColumnDefinition(idColumn,
					firstNameColumn,
					lastNameColumn,
					emailColumn,
					passworHashColumn,
					statusColumn,
					roleColumn,
					createdAtColumn).AddRow(id, "Ivan", "Ivanov", "ivan@gmail.com", "992320c97d2edc09debf80bc3cd2b770a07a97ecee15771e158a744f38790d2e", "active", "customer", "2006-01-02T15:04:05Z")
				mock.ExpectQuery(regexp.QuoteMeta("select * from users where id = $1")).
					WithArgs(id).WillReturnRows(rows)
			},
			expectedUser: &entity.User{
				Id:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				FirstName:    "Ivan",
				LastName:     "Ivanov",
				Email:        "ivan@gmail.com",
				PasswordHash: "992320c97d2edc09debf80bc3cd2b770a07a97ecee15771e158a744f38790d2e",
				Status:       entity.UserStatusActive,
				Role:         entity.UserRoleCustomer,
				CreatedAt:    time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
			},
			expectedErr: nil,
		},
		{
			name:        "user not found",
			inputUserId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(ctx context.Context, id uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from users where id = $1")).WithArgs(id).WillReturnError(sql.ErrNoRows)
			},
			expectedUser: nil,
			expectedErr:  ErrUserNotFound,
		},
		{
			name:        "some error",
			inputUserId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			mockBehavior: func(ctx context.Context, id uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from users where id = $1")).WithArgs(id).WillReturnError(someErr)
			},
			expectedUser: nil,
			expectedErr:  someErr,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(context.Background(), testCase.inputUserId)
			actualUser, err := r.GetById(context.Background(), testCase.inputUserId)
			if testCase.expectedErr != nil {
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, actualUser)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedUser, actualUser)
			}
		})
	}

}
