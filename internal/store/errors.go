package store

import "errors"

var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrDuplicateEmail = errors.New("models: duplicate email")
