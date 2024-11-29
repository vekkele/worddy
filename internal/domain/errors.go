package domain

import "errors"

var ErrNoUserFound = errors.New("auth: no user found")
var ErrDuplicateEmail = errors.New("auth: duplicate email")
var ErrInvalidCredentials = errors.New("auth: invalid credentials")
