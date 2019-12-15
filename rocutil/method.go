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

type CallArg [][]byte

func (this *CallArg) Add(data []byte) {
	(*this) = append(*this, data)
}

func (this *CallArg) Len() int {
	return len((*this))
}

type Method struct {
	name  string
	args  []reflect.Type
	value reflect.Value
	typ   reflect.Type
}

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

// get call args use of Call()
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

// do call this method
func (this *Method) Call(data *CallArg) ([]reflect.Value, error) {
	args, err := this.GetArgValues(data)
	if err != nil {
		return nil, err
	}
	result := this.value.Call(args)
	return result, nil
}

// get this method name in rocutil register
func (this *Method) GetName() string {
	return this.name
}
