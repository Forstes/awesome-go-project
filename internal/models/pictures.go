package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Picture struct {
	ID      int
	Owner   string
	Title   string
	Path    string
	Created time.Time
	Expires time.Time
}

type PictureModel struct {
	DB *pgxpool.Pool
}

func (m *PictureModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO pictures (title, path, created, expires)
	VALUES($1, $2, NOW(), NOW() + INTERVAL '1 DAY' * $3) RETURNING ID`

	var id int
	err := m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PictureModel) Get(id int) (*Picture, error) {

	stmt := `SELECT title, path, created, expires FROM pictures
	WHERE expires > NOW() AND id = $1`

	s := &Picture{}

	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&s.Title, &s.Path, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *PictureModel) Latest() ([]*Picture, error) {

	stmt := `SELECT u.name, title, path, created, expires FROM pictures
	JOIN users u ON pictures.owner_id = u.id
	WHERE expires > NOW()
	ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pictures := []*Picture{}

	for rows.Next() {
		s := &Picture{}
		err = rows.Scan(&s.Owner, &s.Title, &s.Path, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		pictures = append(pictures, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pictures, nil
}
