package roc

import (
	"errors"
)

// ROC错误定义
var (
	ErrUnregisteredROC = errors.New("unregistered roc")
	ErrUnknownObj      = errors.New("unknown roc obj")
)
