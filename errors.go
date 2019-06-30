package wgcreate

import (
	"errors"
)

var (
	ErrInterfaceNotFound = errors.New(`interface not found`)
	ErrNoFreeInterface   = errors.New(`no free interface`)
)
