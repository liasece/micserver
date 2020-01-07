package server

import (
	"errors"
)

// 服务的错误定义
var (
	ErrTargetClientDontExist = errors.New("target client does not exist")
)
