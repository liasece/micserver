package server

import (
	"errors"
)

// 服务的错误定义
var (
	ErrTargetClientNoExist = errors.New("target client does not exist")
)
