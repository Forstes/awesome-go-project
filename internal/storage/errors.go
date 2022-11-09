package storage

import (
	"errors"
)

var ErrBucketExists = errors.New("storage: bucket already exists")
