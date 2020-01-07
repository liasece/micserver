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

// ROC对象代理，由 rocutil 创建的ROC对象均由代理来实现roc对象需要的接口，
// 同时实现通过代理对象具有的函数，自动编解码参数列表，实现更为简单的ROC调用。
type ROCObjAgent struct {
	obj interface{}

	opts *options.Options
	typ  roc.ROCObjType
	id   string

	methodMapping sync.Map
}

// 初始化该代理，目标对象不需要实现任何接口，但是需要告知目标对象的类型和ID，以便
// 代理提供ROC对象的接口给ROC系统。
func (this *ROCObjAgent) Init(obj interface{}, rocObjType roc.ROCObjType,
	rocObjID string, ops []*options.Options) error {
	this.obj = obj
	this.typ = rocObjType
	this.id = rocObjID
	this.opts = &options.Options{}
	for _, optItem := range ops {
		this.opts.Merge(optItem)
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
		if this.opts != nil && this.opts.CheckFuncName != nil {
			var use bool
			name, use = this.opts.CheckFuncName(name)
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

// 添加一个ROC对象的方法
func (this *ROCObjAgent) addMethod(m *Method) {
	this.methodMapping.Store(m.GetName(), m)
}

// 获取一个指定名称的方法
func (this *ROCObjAgent) getMethod(name string) *Method {
	vi, ok := this.methodMapping.Load(name)
	if !ok {
		return nil
	}
	return vi.(*Method)
}

// 提供给 roc.Server 的接口，受到ROC调用时调用
func (this *ROCObjAgent) OnROCCall(path *roc.ROCPath, arg []byte) ([]byte, error) {
	funcName := path.Move()
	if method := this.getMethod(funcName); method != nil {
		callArg := &CallArg{}
		err := json.Unmarshal(arg, callArg)
		if err != nil {
			return nil, err
		}

		// 调用前处理
		if this.opts != nil && this.opts.OnBeforeROCCall != nil {
			this.opts.OnBeforeROCCall(this.obj, path, arg)
		}
		// 实际调用该函数
		result, callErr := method.Call(callArg)
		// 调用后处理
		if this.opts != nil && this.opts.OnAfterROCCall != nil {
			this.opts.OnAfterROCCall(this.obj, path, arg)
		}

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

// 提供给 roc.Server 的接口，获取ROC对象的类型
func (this *ROCObjAgent) GetROCObjType() roc.ROCObjType {
	return this.typ
}

// 提供给 roc.Server 的接口，获取ROC对象的ID
func (this *ROCObjAgent) GetROCObjID() string {
	return this.id
}
