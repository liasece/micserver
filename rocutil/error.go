// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

/*
 * @Author: Jansen
 */

package rocutil

import (
	"errors"
)

// 错误定义
var (
	ErrIsNotFunc      = errors.New("target funcname isn't func")
	ErrUnknownFunc    = errors.New("unknown function name")
	ErrArgNumMismatch = errors.New("call arg num mismatch")
)
