package dusk

import (
	"errors"
)

var (
	ErrDoesNotExist     = errors.New("the item does not exist")
	ErrNoRows           = errors.New("no items found")
	ErrUniqueConstraint = errors.New("the item already exists")
	ErrIsbnExists       = errors.New("isbn already exists")
	ErrNoChange         = errors.New("no change executed")
)
