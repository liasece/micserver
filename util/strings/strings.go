package strings

import (
	"strconv"
	"strings"
	"unsafe"

	"github.com/liasece/micserver/util/math"
)

// StringToBytes func
func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// StringSplitToUint32 func
func StringSplitToUint32(str string, sli string) []uint32 {
	strlist := strings.Split(str, sli)
	res := make([]uint32, len(strlist))
	for i, str := range strlist {
		tmpint, err := strconv.Atoi(str)
		if err == nil {
			res[i] = uint32(tmpint)
		}
	}
	return res
}

// StringSplitToInt32 func
func StringSplitToInt32(str string, sli string) []int32 {
	strlist := strings.Split(str, sli)
	res := make([]int32, len(strlist))
	for i, str := range strlist {
		tmpint, err := strconv.Atoi(str)
		if err == nil {
			res[i] = int32(tmpint)
		}
	}
	return res
}

// RandString func
func RandString(count int) string {
	var randomstr string
	for r := 0; r < count; r++ {
		i := math.RandBetween(65, 90)
		a := rune(i)
		randomstr += string(a)
	}
	return randomstr
}

// MustInterfaceToString func
func MustInterfaceToString(v interface{}) string {
	if v, ok := v.(string); ok {
		return v
	}
	return ""
}
