package miniofs

import (
	"errors"
	"os"
)

var (
	ErrNoBucketInName  = errors.New("no bucket name found in the name")
	ErrEmptyObjectName = errors.New("storage: object name is empty")
	ErrFileClosed      = os.ErrClosed
	ErrOutOfRange      = errors.New("out of range")
	ErrNotSupported    = errors.New("doesn't support this operation")
)
