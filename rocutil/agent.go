// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

/*
 * @Author: Jansen
 */

package rocutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/liasece/micserver/roc"
	"github.com/liasece/micserver/rocutil/options"
)

type ROCObjAgent struct {
	IROCObjBase

	methodMapping sync.Map
}

func (this *ROCObjAgent) Init(obj IROCObjBase, ops []*options.Options) error {
	this.IROCObjBase = obj
	opt := &options.Options{}
	for _, optItem := range ops {
		opt.Merge(optItem)
	}

	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	numMethod := objType.NumMethod()
	for i := 0; i < numMethod; i++ {
		method := objValue.Method(i)
		if method.Kind() != reflect.Func {
			continue
		}
		name := objType.Method(i).Name
		// check func name
		if opt != nil && opt.CheckFuncName != nil {
			var use bool
			name, use = opt.CheckFuncName(name)
			if !use {
				continue
			}
		}
		// "" can't be roc call function's name
		if name == "" {
			continue
		}
		// init new method
		newMethod := &Method{}
		newMethod.Init(name, method)
		this.addMethod(newMethod)
	}
	return nil
}

func (this *ROCObjAgent) addMethod(m *Method) {
	this.methodMapping.Store(m.GetName(), m)
}

func (this *ROCObjAgent) getMethod(name string) *Method {
	vi, ok := this.methodMapping.Load(name)
	if !ok {
		return nil
	}
	return vi.(*Method)
}

func (this *ROCObjAgent) OnROCCall(path *roc.ROCPath, arg []byte) ([]byte, error) {
	funcName := path.Move()
	if method := this.getMethod(funcName); method != nil {
		callArg := &CallArg{}
		err := json.Unmarshal(arg, callArg)
		if err != nil {
			return nil, err
		}
		result, callErr := method.Call(callArg)
		if callErr != nil {
			return nil, callErr
		}
		if len(result) > 0 {
			sendb, err := json.Marshal(result)
			return sendb, err
		}
	} else {
		return nil, fmt.Errorf("%s:%s", ErrUnknownFunc.Error(), funcName)
	}
	return nil, nil
}
