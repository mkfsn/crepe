package crepe

import (
	"errors"
)

var (
	ErrInvalidValue    = errors.New("invalid value")
	ErrUnexportedField = errors.New("unexported field")
	ErrUnimplemented   = errors.New("unimplemented")
	ErrUnsupportedType = errors.New("unsupported type")
)
