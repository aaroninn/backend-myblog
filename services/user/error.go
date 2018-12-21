package user

import "errors"

var (
	errEmptyAccount  = errors.New("error account not exists")
	errWrongPassword = errors.New("error password not correct")
)
