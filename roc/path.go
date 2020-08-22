package roc

import (
	"fmt"
	"strings"
)

// Path ROC调用的路径
type Path struct {
	funcStrings []string
	funcPos     int
	objType     ObjType
	objID       string
}

// O 根据目标ROC的类型及ID，构造一个ROC调用路径
func O(objType ObjType, objID string) *Path {
	res := &Path{}
	res.objType = objType
	res.objID = objID
	return res
}

// pathStr的格式必须为 ROC 远程对象调用那样定义的格式，如：
// pathStrDecode 对象类型[对象的键]
func pathStrDecode(pathStr string) (ObjType, string) {
	objType := ""
	objID := ""
	inKey := false
	for _, k := range pathStr {
		if k == '[' && !inKey {
			inKey = true
		} else if k == ']' && inKey {
			break
		} else if k == '.' {
			break
		} else {
			if !inKey {
				objType = objType + fmt.Sprintf("%c", k)
			} else {
				objID = objID + fmt.Sprintf("%c", k)
			}
		}
	}
	return ObjType(objType), objID
}

// NewPath 根据调用路径字符串，构造一个ROC调用路径
func NewPath(pathStr string) *Path {
	res := &Path{}
	funcStrings := strings.Split(pathStr, ".")
	if len(funcStrings) < 1 {
		return res
	}
	t, id := pathStrDecode(funcStrings[0])
	res.objType = ObjType(t)
	res.objID = id
	if len(funcStrings) > 1 {
		res.funcStrings = funcStrings[1:]
	}
	return res
}

// GetObjType 获取ROC调用路径的目标ROC对象类型
func (path *Path) GetObjType() ObjType {
	return path.objType
}

// GetObjID 获取ROC调用路径的目标ROC对象ID
func (path *Path) GetObjID() string {
	return path.objID
}

// GetPos 获取ROC调用路径当前行进到的位置
func (path *Path) GetPos() int {
	return path.funcPos
}

// Get 获取当前ROC调用路径指定位置的值
func (path *Path) Get(funcPos int) string {
	if funcPos < 0 || funcPos >= len(path.funcStrings) {
		return ""
	}
	return path.funcStrings[funcPos]
}

// Move 移动当前ROC调用路径到下一个位置，并返回该位置的值
func (path *Path) Move() string {
	if path.funcPos >= len(path.funcStrings) {
		return ""
	}
	res := path.funcStrings[path.funcPos]
	path.funcPos++
	return res
}

// F 添加一个ROC调用路径的函数段，一个ROC调用可以携带多个函数名，中间以 . 号连接
func (path *Path) F(funcName string) *Path {
	res := path
	if res == nil {
		res = &Path{}
	}
	res.funcStrings = append(res.funcStrings, funcName)
	return res
}

// String 获取当前ROC调用路径的字符串描述
func (path *Path) String() string {
	res := string(path.objType) + "[" + string(path.objID) + "]"
	for _, v := range path.funcStrings {
		res += "." + v
	}
	return res
}

// Reset 重置当前ROC调用路径的位置
func (path *Path) Reset() {
	path.funcPos = 0
}
