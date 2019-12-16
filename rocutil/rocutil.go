// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

/*
 * @Author: Jansen
 */

package rocutil

import (
	"encoding/json"

	"github.com/liasece/micserver/roc"
	"github.com/liasece/micserver/rocutil/options"
)

func NewROCObj(obj interface{}, rocObjType roc.ROCObjType,
	rocObjID string, ops ...*options.Options) (roc.IObj, error) {
	agent := &ROCObjAgent{}
	agent.Init(obj, rocObjType, rocObjID, ops)
	return agent, nil
}

func ServerROCObj(rocserver roc.IROCServer, obj interface{},
	rocObjType roc.ROCObjType, rocObjID string,
	opts ...*options.Options) (roc.IObj, error) {
	agent, err := NewROCObj(obj, rocObjType, rocObjID, opts...)
	if err != nil {
		return nil, err
	}
	res, _ := rocserver.NewROC(rocObjType).
		GetOrRegObj(rocObjID, agent)
	return res, nil
}

// call remote object function
func CallNR(rocServer roc.IROCServer, typ roc.ROCObjType, objID string,
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
