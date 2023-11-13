package minio

import (
	"errors"
	"os"
)

var (
	ErrNoBucketInName     = errors.New("no bucket name found in the name")
	ErrObjectDoesNotExist = errors.New("storage: object doesn't exist")
	ErrEmptyObjectName    = errors.New("storage: object name is empty")
	ErrFileClosed         = os.ErrClosed
	ErrOutOfRange         = errors.New("out of range")
	ErrFileNotFound       = os.ErrNotExist
	ErrNotSupported       = errors.New("doesn't support this operation")
)
