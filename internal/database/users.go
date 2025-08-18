package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/shekharsikku/go-gin-rest-api/internal/utils"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (em *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING id"

	id := utils.GenerateUniqueID()

	return em.DB.QueryRowContext(ctx, query, id, user.Name, user.Email, user.Password).Scan(&user.Id)
}

func (em *UserModel) getUser(query string, args ...any) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var user User

	err := em.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (em *UserModel) Get(id int) (*User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	return em.getUser(query, id)
}

func (em *UserModel) GetByEmail(email string) (*User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	return em.getUser(query, email)
}

func (em *UserModel) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "SELECT * FROM users"

	rows, err := em.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
