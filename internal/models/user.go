package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx"
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

	stmt := `INSERT INTO users (name, password) VALUES($1, $2) RETURNING ID`

	var id int
	err = m.DB.QueryRow(context.Background(), stmt, name, string(hashedPassword)).Scan(&id)

	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *UserModel) Authenticate(name, password string) (int, error) {

	stmt := `SELECT id, password FROM users WHERE name = $1`
	s := User{}

	err := m.DB.QueryRow(context.Background(), stmt, name).Scan(&s.ID, &s.HashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(s.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return s.ID, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
