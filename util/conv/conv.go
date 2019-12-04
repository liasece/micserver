package conv

import (
	"fmt"
	"strconv"
)

// interface{} -->> string
func MustInterfaceToString(vi interface{}) string {
	switch vi.(type) {
	case string:
		return vi.(string)
	case []byte:
		return string(vi.([]byte))
	default:
		return fmt.Sprint(vi)
	}
	return ""
}

// interface{} -->> bool
func MustInterfaceToBool(vi interface{}) bool {
	if MustInterfaceToInt64(vi) != 0 {
		return true
	}
	switch vi.(type) {
	case bool:
		return vi.(bool)
	case string:
		v := vi.(string)
		return v == "true" || v == "TRUE" || v == "True"
	}
	return false
}

// interface{} -->> int64
func MustInterfaceToInt64(vi interface{}) int64 {
	switch vi.(type) {
	case int64:
		return vi.(int64)
	case int32:
		return int64(vi.(int32))
	case uint64:
		return int64(vi.(uint64))
	case uint32:
		return int64(vi.(uint32))
	case int:
		return int64(vi.(int))
	case float32:
		return int64(vi.(float32))
	case float64:
		return int64(vi.(float64))
	case string:
		v, err := strconv.Atoi(vi.(string))
		if err == nil {
			return int64(v)
		}
	}
	return 0
}

// interface{} -->> float64
func MustInterfaceToFloat64(vi interface{}) float64 {
	switch vi.(type) {
	case int64:
		return float64(vi.(int64))
	case int32:
		return float64(vi.(int32))
	case uint64:
		return float64(vi.(uint64))
	case uint32:
		return float64(vi.(uint32))
	case int:
		return float64(vi.(int))
	case float32:
		return float64(vi.(float32))
	case float64:
		return float64(vi.(float64))
	case string:
		v, err := strconv.Atoi(vi.(string))
		if err == nil {
			return float64(v)
		}
	}
	return 0
}

// []interface{} -->> []int64
func MustInterfaceSliceToInt64Slice(vsi []interface{}) []int64 {
	if vsi == nil {
		return nil
	}
	res := make([]int64, len(vsi))
	for i, vi := range vsi {
		res[i] = MustInterfaceToInt64(vi)
	}
	return res
}

// interface{} -->> []int64
func MustInterfaceToInt64Slice(vi interface{}) []int64 {
	switch vi.(type) {
	case []interface{}:
		return MustInterfaceSliceToInt64Slice(vi.([]interface{}))
	case int64:
		return vi.([]int64)
	case []int32:
		slice := vi.([]int32)
		res := make([]int64, len(slice))
		for i, v := range slice {
			res[i] = int64(v)
		}
		return res
	case []uint64:
		slice := vi.([]uint64)
		res := make([]int64, len(slice))
		for i, v := range slice {
			res[i] = int64(v)
		}
		return res
	case []uint32:
		slice := vi.([]uint32)
		res := make([]int64, len(slice))
		for i, v := range slice {
			res[i] = int64(v)
		}
		return res
	case []int:
		slice := vi.([]int)
		res := make([]int64, len(slice))
		for i, v := range slice {
			res[i] = int64(v)
		}
		return res
	case []float32:
		slice := vi.([]float32)
		res := make([]int64, len(slice))
		for i, v := range slice {
			res[i] = int64(v)
		}
		return res
	case []float64:
		slice := vi.([]float64)
		res := make([]int64, len(slice))
		for i, v := range slice {
			res[i] = int64(v)
		}
		return res
	case []string:
		slice := vi.([]string)
		res := make([]int64, len(slice))
		for i, v := range slice {
			res[i] = MustInterfaceToInt64(v)
		}
		return res
	}
	return nil
}
