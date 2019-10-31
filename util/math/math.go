package math

import ()

func Abs(n int32) uint32 {
	if n < 0 {
		return uint32(-n)
	}
	return uint32(n)
}
