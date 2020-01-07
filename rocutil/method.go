// Copyright 2019 The Misserver Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

/*
 * @Author: Jansen
 */

package rocutil

import (
	"encoding/json"
	"reflect"
)

// ROC调用的参数列表
type CallArg [][]byte

// 增加一个调用参数
func (this *CallArg) Add(data []byte) {
	(*this) = append(*this, data)
}

// 调用参数的数量
func (this *CallArg) Len() int {
	return len((*this))
}

// 代理对象的方法
type Method struct {
	name  string
	args  []reflect.Type
	value reflect.Value
	typ   reflect.Type
}

// 初始化一个方法
func (this *Method) Init(name string, f reflect.Value) error {
	if f.Kind() != reflect.Func {
		return ErrIsNotFunc
	}
	this.name = name
	this.value = f
	this.typ = f.Type()

	// 遍历所有参数
	numArgs := this.typ.NumIn()
	this.args = make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		argTyp := this.typ.In(i)
		this.args[i] = argTyp
	}
	return nil
}

// 根据远程调用发来的参数，将信息解码为当前方法的参数列表，提供给外层反射调用。
func (this *Method) GetArgValues(data *CallArg) ([]reflect.Value, error) {
	if data.Len() != len(this.args) {
		return nil, ErrArgNumMismatch
	}
	res := make([]reflect.Value, len(this.args))
	for i := 0; i < len(this.args); i++ {
		// new arg value
		v := reflect.New(this.args[i])
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

// 提供编码好的参数二进制流，调用该方法
func (this *Method) Call(data *CallArg) ([]reflect.Value, error) {
	args, err := this.GetArgValues(data)
	if err != nil {
		return nil, err
	}
	result := this.value.Call(args)
	return result, nil
}

// 获取当前方法的名字。
func (this *Method) GetName() string {
	return this.name
}
