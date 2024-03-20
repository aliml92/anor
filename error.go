package anor

import "errors"

const (
	EINTERNALMSG = "Something went wrong. Please try again later."
)

var (
	ErrNotFound   = errors.New("not found")
	ErrUserExists = errors.New("user exists")
)
