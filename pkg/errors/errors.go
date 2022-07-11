package errors

import (
	"errors"
	perrors "github.com/pkg/errors"
)

//logic error.
var (
	ErrIndexNotFound          = errors.New("index not found")
	ErrInvalidMapping         = errors.New("invalid mapping")
	ErrDocumentNotFound       = errors.New("document not found")
	ErrIndexAlreadyExists     = errors.New("the index already exists")
	ErrIndexCloneNotSupported = errors.New("the index don't support clone")
	ErrBulkDataFormat         = errors.New("error bulk data format")
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
