package crepe

import (
	"errors"
)

var (
	ErrInvalidValue    = errors.New("invalid value")
	ErrUnexportedField = errors.New("unexported field")
	ErrInvalidTag      = errors.New("invalid tag")
	ErrUnknownAction   = errors.New("unknown action")
	ErrUnimplemented   = errors.New("unimplemented")
	ErrUnsupportedType = errors.New("unsupported type")
)
