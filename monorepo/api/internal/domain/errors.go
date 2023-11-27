package domain

import "github.com/pkg/errors"

var (
	ErrOrderNotExists = errors.New("order does not exists")
	ErrItemsNotExists = errors.New("items not exists")
	ErrInvalidSyntax  = errors.New("invalid syntax value")
)

var (
	CodeInvalidSyntax = "22P02"
)
