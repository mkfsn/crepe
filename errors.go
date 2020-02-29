package crepe

import (
	"errors"
)

var (
	ErrInvalidValue    = errors.New("invalid value")
	ErrUnexportedField = errors.New("unexported field")
	ErrEmptyTag        = errors.New("empty tag")
	ErrUnknownAction   = errors.New("unknown action")
	ErrUnknownTarget   = errors.New("unknown target")
	ErrUnimplemented   = errors.New("unimplemented")
	ErrUnsupportedType = errors.New("unsupported type")
)
