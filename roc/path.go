package roc

import (
	"fmt"
	"strings"
)

// Path ROC调用的路径
type Path struct {
	strs    []string
	pos     int
	objType ObjType
	objID   string
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
	t := ""
	key := ""
	inKey := false
	for _, k := range pathStr {
		if k == '[' {
			inKey = true
		} else if k == ']' {
			inKey = false
		} else if k == '.' {
			break
		} else {
			if key == "" && !inKey {
				t = t + fmt.Sprintf("%c", k)
			} else if t != "" && inKey {
				key = key + fmt.Sprintf("%c", k)
			} else {
				return "", ""
			}
		}
	}
	return ObjType(t), key
}

// NewPath 根据调用路径字符串，构造一个ROC调用路径
func NewPath(pathstr string) *Path {
	res := &Path{}
	strs := strings.Split(pathstr, ".")
	if len(strs) < 1 {
		return res
	}
	t, id := pathStrDecode(strs[0])
	res.objType = ObjType(t)
	res.objID = id
	res.strs = strs[1:]
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
	return path.pos
}

// Get 获取当前ROC调用路径指定位置的值
func (path *Path) Get(pos int) string {
	if pos < 0 || pos >= len(path.strs) {
		return ""
	}
	return path.strs[pos]
}

// Move 移动当前ROC调用路径到下一个位置，并返回该位置的值
func (path *Path) Move() string {
	if path.pos >= len(path.strs) {
		return ""
	}
	res := path.strs[path.pos]
	path.pos++
	return res
}

// F 添加一个ROC调用路径的函数段，一个ROC调用可以携带多个函数名，中间以 . 号连接
func (path *Path) F(funcName string) *Path {
	res := path
	if res == nil {
		res = &Path{}
	}
	res.strs = append(res.strs, funcName)
	return res
}

// String 获取当前ROC调用路径的字符串描述
func (path *Path) String() string {
	res := string(path.objType) + "[" + string(path.objID) + "]"
	for _, v := range path.strs {
		res += "." + v
	}
	return res
}

// Reset 重置当前ROC调用路径的位置
func (path *Path) Reset() {
	path.pos = 0
}
