package errors

import "errors"

var (
	ErrKeyNotFound      = errors.New("Key not found")
	ErrEmptyKey         = errors.New("Key cannot be empty")
	ErrInvalidIndexName = errors.New("Index name is invalid")
	ErrInvalidDocID     = errors.New("Document ID is invalid")
)
