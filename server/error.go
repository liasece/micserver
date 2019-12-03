package server

import (
	"errors"
)

var (
	ErrTargetClientDontExist = errors.New("target client does not exist")
)
