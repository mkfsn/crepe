package internal

import (
	"errors"
)

var (
	ErrUnknownTarget = errors.New("unknown target")
	ErrEmptyTag      = errors.New("empty tag")
)
