package errors

import "errors"

var ErrKeyNotFound = errors.New("Key not found")
var ErrEmptyKey = errors.New("Key cannot be empty")
