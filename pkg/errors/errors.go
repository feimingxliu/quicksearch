package errors

import "errors"

//logic error.
var ()

//underlying db error.
var (
	ErrKeyNotFound      = errors.New("Key not found")
	ErrEmptyKey         = errors.New("Key cannot be empty")
	ErrKeyValueNotMatch = errors.New("Keys and values not match")
)
