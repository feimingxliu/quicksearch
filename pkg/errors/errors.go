package errors

import (
	"errors"
	perrors "github.com/pkg/errors"
)

//logic error.
var (
	ErrIndexNotFound          = errors.New("Index not found")
	ErrInvalidMapping         = errors.New("invalid mapping")
	ErrDocumentNotFound       = errors.New("Document not found")
	ErrIndexAlreadyExists     = errors.New("The index already exists")
	ErrIndexCloneNotSupported = errors.New("The index don't support clone")
)

//underlying db error.
var (
	ErrKeyNotFound      = errors.New("Key not found")
	ErrEmptyKey         = errors.New("Key cannot be empty")
	ErrKeyValueNotMatch = errors.New("Keys and values not match")
)

func WithStack(err error) error {
	return perrors.WithStack(err)
}
