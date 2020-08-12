package roc

import (
	"errors"
)

// ROC错误定义
var (
	ErrUnregisterROC = errors.New("unregistered roc")
	ErrUnknownObj    = errors.New("unknow roc obj")
)
