// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

/*
 * @Author: Jansen
 */

// Package rocutil roc对象构造及调用的简单使用工具，基于反射实现，性能敏感处不宜使用。
package rocutil

import (
	"encoding/json"

	"github.com/liasece/micserver/roc"
	"github.com/liasece/micserver/rocutil/options"
)

// NewROCObj 基于源对象构造一个 ROC Util 对象，该对象的远程调用必须使用 rocutil.CallNR 调用，
// 本质上，是实例化出一个源对象的代理，由该代理实现 roc.IObj 接口，然后返回该代理，
// 这简化了源对象的构造，并且通过反射来将源对象中具有的方法映射到ROC调用的路径中，
// 使得源对象的 OnROCCall 方法也省略了。
func NewROCObj(obj interface{}, rocObjType roc.ObjType,
	rocObjID string, ops ...*options.Options) (roc.IObj, error) {
	agent := &ROCObjAgent{}
	agent.Init(obj, rocObjType, rocObjID, ops)
	return agent, nil
}

// ServerROCObj 服务一个对象，并且将该对象注册到
// 源对象的 roc.IObj 接口由 NewROCObj 实现，
// 实际注册到ROC系统中的是代理对象，如果目标类型的ROC没被创建则自动创建。
func ServerROCObj(rocserver roc.IROCServer, obj interface{},
	rocObjType roc.ObjType, rocObjID string,
	opts ...*options.Options) (roc.IObj, error) {
	agent, err := NewROCObj(obj, rocObjType, rocObjID, opts...)
	if err != nil {
		return nil, err
	}
	res, _ := rocserver.NewROC(rocObjType).
		GetOrRegObj(rocObjID, agent)
	return res, nil
}

// CallNR call remote object function
// 调用一个由 rocutil 创建的ROC对象，由于该调用会将参数编码，所以对端的参数解码也
// 必须由 rocutil 实现。
func CallNR(rocServer roc.IROCServer, typ roc.ObjType, objID string,
	funcName string, args ...interface{}) error {
	callArg := &CallArg{}
	for _, arg := range args {
		// default encoder is encoding/json
		b, err := json.Marshal(arg)
		if err != nil {
			return err
		}
		callArg.Add(b)
	}
	// outermost arg marshal
	b, err := json.Marshal(callArg)
	if err != nil {
		return err
	}
	return rocServer.ROCCallNR(roc.O(typ, objID).F(funcName), b)
}
