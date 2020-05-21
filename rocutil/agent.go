// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of agent source code is governed by a MIT-style
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

// ROCObjAgent ROC对象代理，由 rocutil 创建的ROC对象均由代理来实现roc对象需要的接口，
// 同时实现通过代理对象具有的函数，自动编解码参数列表，实现更为简单的ROC调用。
type ROCObjAgent struct {
	obj interface{}

	opts *options.Options
	typ  roc.ObjType
	id   string

	methodMapping sync.Map
}

// Init 初始化该代理，目标对象不需要实现任何接口，但是需要告知目标对象的类型和ID，以便
// 代理提供ROC对象的接口给ROC系统。
func (agent *ROCObjAgent) Init(obj interface{}, rocObjType roc.ObjType,
	rocObjID string, ops []*options.Options) error {
	agent.obj = obj
	agent.typ = rocObjType
	agent.id = rocObjID
	agent.opts = &options.Options{}
	for _, optItem := range ops {
		agent.opts.Merge(optItem)
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
		if agent.opts != nil && agent.opts.CheckFuncName != nil {
			var use bool
			name, use = agent.opts.CheckFuncName(name)
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
		agent.addMethod(newMethod)
	}
	return nil
}

// addMethod 添加一个ROC对象的方法
func (agent *ROCObjAgent) addMethod(m *Method) {
	agent.methodMapping.Store(m.GetName(), m)
}

// getMethod 获取一个指定名称的方法
func (agent *ROCObjAgent) getMethod(name string) *Method {
	vi, ok := agent.methodMapping.Load(name)
	if !ok {
		return nil
	}
	return vi.(*Method)
}

// OnROCCall 提供给 roc.Server 的接口，受到ROC调用时调用
func (agent *ROCObjAgent) OnROCCall(path *roc.Path, arg []byte) ([]byte, error) {
	funcName := path.Move()
	if method := agent.getMethod(funcName); method != nil {
		callArg := &CallArg{}
		err := json.Unmarshal(arg, callArg)
		if err != nil {
			return nil, err
		}

		// 调用前处理
		if agent.opts != nil && agent.opts.OnBeforeROCCall != nil {
			agent.opts.OnBeforeROCCall(agent.obj, path, arg)
		}
		// 实际调用该函数
		result, callErr := method.Call(callArg)
		// 调用后处理
		if agent.opts != nil && agent.opts.OnAfterROCCall != nil {
			agent.opts.OnAfterROCCall(agent.obj, path, arg)
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

// GetROCObjType 提供给 roc.Server 的接口，获取ROC对象的类型
func (agent *ROCObjAgent) GetROCObjType() roc.ObjType {
	return agent.typ
}

// GetROCObjID 提供给 roc.Server 的接口，获取ROC对象的ID
func (agent *ROCObjAgent) GetROCObjID() string {
	return agent.id
}
