package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrDuplicateName = errors.New("models: duplicate name")
)
