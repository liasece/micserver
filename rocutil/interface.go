// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

/*
 * @Author: Jansen
 */

package rocutil

import (
	"github.com/liasece/micserver/roc"
)

type IROCObjBase interface {
	GetROCObjType() roc.ROCObjType
	GetROCObjID() string
}
