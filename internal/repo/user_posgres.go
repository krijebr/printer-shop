package repo

import (
	"context"
	"database/sql"
	"log"
	"strings"
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
	where := ""
	if filter != nil {
		whereS := []string{}
		if filter.UserStatus != nil {
			whereS = append(whereS, "status = '"+string(*filter.UserStatus)+"'")
		}
		if filter.UserRole != nil {
			whereS = append(whereS, "role = '"+string(*filter.UserRole)+"'")
		}
		where = " where " + strings.Join(whereS, " and ")
	}
	rows, err := u.db.QueryContext(ctx, "select * from users"+where)
	if err != nil {
		return nil, err
	}
	users := []*entity.User{}
	for rows.Next() {
		var dateStr string
		user := new(entity.User)
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.Status, &user.Role, &dateStr)
		if err != nil {
			log.Println("Ошибка чтения строки", err)
			continue
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
	row := u.db.QueryRowContext(ctx, "select * from users where id = $1", id)
	var dateStr string
	user := new(entity.User)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.Status, &user.Role, &dateStr)
	if err != nil {
		return nil, err
	}
	user.CreatedAt, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UserRepoPg) Create(ctx context.Context, user entity.User) error {
	_, err := u.db.ExecContext(ctx,
		"insert into users (id, first_name, last_name, email, password_hash, status, role, created_at) values ($1,$2,$3,$4,$5,$6,$7,$8)",
		user.Id, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Status, user.Role, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepoPg) Update(ctx context.Context, user entity.User) error {
	set := []string{}
	if user.FirstName != "" {
		set = append(set, "first_name = '"+user.FirstName+"'")
	}
	if user.LastName != "" {
		set = append(set, "last_name = '"+user.LastName+"'")
	}
	if user.PasswordHash != "" {
		set = append(set, "password_hash = '"+user.PasswordHash+"'")
	}
	if user.Status != "" {
		set = append(set, "status = '"+string(user.Status)+"'")
	}
	if user.Role != "" {
		set = append(set, "role = '"+string(user.Role)+"'")
	}

	_, err := u.db.ExecContext(ctx, "update users set "+strings.Join(set, ", ")+" where id = $1", user.Id)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepoPg) DeleteById(ctx context.Context, id uuid.UUID) error {
	_, err := u.db.ExecContext(ctx, "delete from users where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepoPg) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := u.db.QueryRowContext(ctx, "select * from users where email = $1", email)
	var dateStr string
	user := new(entity.User)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.Status, &user.Role, &dateStr)
	if err != nil {
		return nil, err
	}
	user.CreatedAt, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, err
	}
	return user, nil
}
