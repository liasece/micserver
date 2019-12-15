// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

/*
 * @Author: Jansen
 */

package options

type Options struct {
	// 检查ROC调用的函数名称，并且返回一个最终的名称以及是否使用它
	CheckFuncName func(method string) (string, bool)
}

func (this *Options) Merge(opt *Options) {
	if opt == nil {
		return
	}
	if opt.CheckFuncName != nil {
		this.CheckFuncName = opt.CheckFuncName
	}
}
