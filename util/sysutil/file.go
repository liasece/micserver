package sysutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// 读取一个json文件到指定对象中
func LoadJsonFromFile(filename string, v interface{}) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("[LoadJsonFromFile] Load %s failed \n%s ",
			filename, err.Error())
	}
	err = json.Unmarshal([]byte(content), v)
	if err != nil {
		return fmt.Errorf(
			"[LoadJsonFromFile] Load %s failed, Unmarshal failed :\n%s ",
			filename, err.Error())
	}
	return nil
}
