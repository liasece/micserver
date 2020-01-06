package roc

import (
	"fmt"
	"strings"
)

// ROC调用的路径
type ROCPath struct {
	strs    []string
	pos     int
	objType ROCObjType
	objID   string
}

// 根据目标ROC的类型及ID，构造一个ROC调用路径
func O(objType ROCObjType, objID string) *ROCPath {
	res := &ROCPath{}
	res.objType = objType
	res.objID = objID
	return res
}

// kstr的格式必须为 ROC 远程对象调用那样定义的格式，如：
// 对象类型[对象的键]
func kstrDecode(kstr string) (ROCObjType, string) {
	t := ""
	key := ""
	inkey := false
	for _, k := range kstr {
		if k == '[' {
			inkey = true
		} else if k == ']' {
			inkey = false
		} else if k == '.' {
			break
		} else {
			if key == "" && !inkey {
				t = t + fmt.Sprintf("%c", k)
			} else if t != "" && inkey {
				key = key + fmt.Sprintf("%c", k)
			} else {
				return "", ""
			}
		}
	}
	return ROCObjType(t), key
}

// 根据调用路径字符串，构造一个ROC调用路径
func NewROCPath(pathstr string) *ROCPath {
	res := &ROCPath{}
	strs := strings.Split(pathstr, ".")
	if len(strs) < 1 {
		return res
	}
	t, id := kstrDecode(strs[0])
	res.objType = ROCObjType(t)
	res.objID = id
	res.strs = strs[1:]
	return res
}

// 获取ROC调用路径的目标ROC对象类型
func (this *ROCPath) GetObjType() ROCObjType {
	return this.objType
}

// 获取ROC调用路径的目标ROC对象ID
func (this *ROCPath) GetObjID() string {
	return this.objID
}

// 获取ROC调用路径当前行进到的位置
func (this *ROCPath) GetPos() int {
	return this.pos
}

// 获取当前ROC调用路径指定位置的值
func (this *ROCPath) Get(pos int) string {
	if pos < 0 || pos >= len(this.strs) {
		return ""
	}
	return this.strs[pos]
}

// 移动当前ROC调用路径到下一个位置，并返回该位置的值
func (this *ROCPath) Move() string {
	if this.pos >= len(this.strs) {
		return ""
	}
	res := this.strs[this.pos]
	this.pos++
	return res
}

// 添加一个ROC调用路径的函数段，一个ROC调用可以携带多个函数名，中间以 . 号连接
func (this *ROCPath) F(funcName string) *ROCPath {
	res := this
	if res == nil {
		res = &ROCPath{}
	}
	res.strs = append(res.strs, funcName)
	return res
}

// 获取当前ROC调用路径的字符串描述
func (this *ROCPath) String() string {
	res := string(this.objType) + "[" + string(this.objID) + "]"
	for _, v := range this.strs {
		res += "." + v
	}
	return res
}

// 重置当前ROC调用路径的位置
func (this *ROCPath) Reset() {
	this.pos = 0
}
