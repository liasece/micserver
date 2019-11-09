package roc

import (
	"fmt"
)

type ROCPath struct {
	strs    []string
	pos     int
	objType ROCObjType
	objID   string
}

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

func NewROCPath(strs []string) *ROCPath {
	res := &ROCPath{
		strs: strs,
	}
	if len(strs) > 0 {
		t, id := kstrDecode(strs[0])
		res.objType = ROCObjType(t)
		res.objID = id
	}
	return res
}

func (this *ROCPath) GetObjType() ROCObjType {
	return this.objType
}

func (this *ROCPath) GetObjID() string {
	return this.objID
}

func (this *ROCPath) GetPos() int {
	return this.pos
}

func (this *ROCPath) Get(pos int) string {
	if pos < 0 || pos >= len(this.strs) {
		return ""
	}
	return this.strs[pos]
}

func (this *ROCPath) Move() string {
	if this.pos >= len(this.strs) {
		return ""
	}
	this.pos++
	return this.strs[this.pos]
}

func (this *ROCPath) F(funcName string) *ROCPath {
	res := this
	if res == nil {
		res = &ROCPath{}
	}
	res.strs = append(res.strs, funcName)
	return res
}

func (this *ROCPath) String() string {
	res := string(this.objType) + "[" + string(this.objID) + "]"
	for _, v := range this.strs {
		res += "." + v
	}
	return res
}

func (this *ROCPath) Reset() {
	this.pos = 0
}
