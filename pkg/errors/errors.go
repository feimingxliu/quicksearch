package errors

import "errors"

//logic error.
var (
	ErrIndexNotFound      = errors.New("Index not found")
	ErrDocumentNotFound   = errors.New("Document not found")
	ErrCloneIndexSameName = errors.New("Cloned index name is same as origin")
)

//underlying db error.
var (
	ErrKeyNotFound      = errors.New("Key not found")
	ErrEmptyKey         = errors.New("Key cannot be empty")
	ErrKeyValueNotMatch = errors.New("Keys and values not match")
)
