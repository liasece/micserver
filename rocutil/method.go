// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of method source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

/*
 * @Author: Jansen
 */

package rocutil

import (
	"encoding/json"
	"reflect"
)

// CallArg ROC调用的参数列表
type CallArg [][]byte

// Add 增加一个调用参数
func (arg *CallArg) Add(data []byte) {
	(*arg) = append(*arg, data)
}

// Len 调用参数的数量
func (arg *CallArg) Len() int {
	return len((*arg))
}

// Method 代理对象的方法
type Method struct {
	name  string
	args  []reflect.Type
	value reflect.Value
	typ   reflect.Type
}

// Init 初始化一个方法
func (method *Method) Init(name string, f reflect.Value) error {
	if f.Kind() != reflect.Func {
		return ErrIsNotFunc
	}
	method.name = name
	method.value = f
	method.typ = f.Type()

	// 遍历所有参数
	numArgs := method.typ.NumIn()
	method.args = make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		argTyp := method.typ.In(i)
		method.args[i] = argTyp
	}
	return nil
}

// GetArgValues 根据远程调用发来的参数，将信息解码为当前方法的参数列表，提供给外层反射调用。
func (method *Method) GetArgValues(data *CallArg) ([]reflect.Value, error) {
	if data.Len() != len(method.args) {
		return nil, ErrArgNumMismatch
	}
	res := make([]reflect.Value, len(method.args))
	for i := 0; i < len(method.args); i++ {
		// new arg value
		v := reflect.New(method.args[i])
		// unmarshal arg value
		err := json.Unmarshal((*data)[i], v.Interface())
		if err != nil {
			return nil, err
		}
		// add to result
		res[i] = v.Elem()
	}
	return res, nil
}

// Call 提供编码好的参数二进制流，调用该方法
func (method *Method) Call(data *CallArg) ([]reflect.Value, error) {
	args, err := method.GetArgValues(data)
	if err != nil {
		return nil, err
	}
	result := method.value.Call(args)
	return result, nil
}

// GetName 获取当前方法的名字。
func (method *Method) GetName() string {
	return method.name
}
