// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

/*
 * @Author: Jansen
 */

package options

import (
	"github.com/liasece/micserver/roc"
)

// Options options
type Options struct {
	// 检查ROC调用的函数名称，并且返回一个最终的名称以及是否使用它
	CheckFuncName func(method string) (string, bool)
	// 在ROC被调用前，会执行该函数，你可以在这里完成一些保证同步的操作，如加锁
	OnBeforeROCCall func(obj interface{}, callpath *roc.Path,
		arg []byte)
	// 在ROC被调用后，会执行该函数，你可以在这里完成一些保证同步的操作，如加锁
	OnAfterROCCall func(obj interface{}, callpath *roc.Path,
		arg []byte)
}

// Merge merge options
func (o *Options) Merge(opt *Options) {
	if opt == nil {
		return
	}
	if opt.CheckFuncName != nil {
		o.CheckFuncName = opt.CheckFuncName
	}
	if opt.OnBeforeROCCall != nil {
		o.OnBeforeROCCall = opt.OnBeforeROCCall
	}
	if opt.OnAfterROCCall != nil {
		o.OnAfterROCCall = opt.OnAfterROCCall
	}
}
