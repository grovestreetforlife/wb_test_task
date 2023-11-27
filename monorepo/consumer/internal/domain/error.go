package domain

import "github.com/pkg/errors"

var (
	// user errors
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrUserDoesNotExists = errors.New("user does not exists")

	// order errors
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrInvalidOrderValue  = errors.New("invalid order value")
	ErrOrderDoesNotExists = errors.New("order does not exists")

	// common errors
	ErrInvalidValue = errors.New("invalid value")
)

var (
	CodeErrDuplicateKey       = "23505"
	CodeErrConstraintLenValue = "22001"
	CodeErrInvalidSyntax      = "22P02"
	CodeErrForeignKey         = "23503"
)
