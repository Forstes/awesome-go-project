package models

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	HashedPassword []byte
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name string, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (name, password)
VALUES($1, $2) RETURNING ID`

	var id int
	err1 := m.DB.QueryRow(context.Background(), stmt, name, string(hashedPassword)).Scan(&id)

	if err1 != nil {
		return 0, err1
	}

	return int(id), nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil

}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
