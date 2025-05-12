package repo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	_ "github.com/lib/pq"
)

type UserRepoPg struct {
	db *sql.DB
}

func NewUserRepoPg(db *sql.DB) User {
	return &UserRepoPg{
		db: db,
	}
}
func (u *UserRepoPg) GetAll(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	log.Println(filter)
	where := ""
	if filter != nil {
		where = where + " where "
		first := true
		if filter.UserStatus != nil {
			where = where + "status = '" + string(*filter.UserStatus) + "'"
			first = false
		}
		if filter.UserRole != nil {
			if !first {
				where = where + " and "
			}
			where = where + "role = '" + string(*filter.UserRole) + "'"
		}
	}
	rows, err := u.db.QueryContext(ctx, "select * from users"+where)
	if err != nil {
		return nil, err
	}
	users := []*entity.User{}
	for rows.Next() {
		var dateStr string
		var firstNameStr, lastNameStr sql.NullString
		user := new(entity.User)
		err := rows.Scan(&user.Id, &firstNameStr, &lastNameStr, &user.Email, &user.PasswordHash, &user.Status, &user.Role, &dateStr)
		if err != nil {
			log.Println("Ошибка чтения строки", err)
			continue
		}
		if firstNameStr.Valid {
			user.FirstName = firstNameStr.String
		}
		if lastNameStr.Valid {
			user.LastName = lastNameStr.String
		}
		user.CreatedAt, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			log.Println("Ошибка преобразования времени")
			continue
		}
		users = append(users, user)
	}
	return users, nil
}
func (u *UserRepoPg) GetById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return nil, nil
}
func (u *UserRepoPg) Create(ctx context.Context, user entity.User) (*entity.User, error) {

	_, err := u.db.ExecContext(ctx, "insert into users (id, first_name, last_name, email, password_hash, status, role, created_at) values ($1,$2,$3,$4,$5,$6,$7,$8)", user.Id, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Status, user.Role, user.CreatedAt)
	if err != nil {
		return nil, err
	}
	row := u.db.QueryRowContext(ctx, "select * from users where id = $1", user.Id)
	var dateStr string
	var firstNameStr, lastNameStr sql.NullString
	newUser := new(entity.User)
	err = row.Scan(&newUser.Id, &firstNameStr, &lastNameStr, &newUser.Email, &newUser.PasswordHash, &newUser.Status, &newUser.Role, &dateStr)
	if err != nil {
		return nil, err
	}
	if firstNameStr.Valid {
		newUser.FirstName = firstNameStr.String
	}
	if lastNameStr.Valid {
		newUser.LastName = lastNameStr.String
	}
	newUser.CreatedAt, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
func (u *UserRepoPg) UpdateById(ctx context.Context, id uuid.UUID, user entity.User) (*entity.User, error) {
	return nil, nil
}
func (u *UserRepoPg) DeleteById(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (u *UserRepoPg) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
