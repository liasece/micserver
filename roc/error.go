package roc

import (
	"errors"
)

var (
	ErrUnregisterRoc = errors.New("unregistered roc")
	ErrUnknowObj     = errors.New("unknow roc obj")
)
