package termrender

import "errors"

var (
	ErrInvalidRenderFormat = errors.New("invalid rendering format")
	ErrInvalidInputType    = errors.New("invalid input type")
)
