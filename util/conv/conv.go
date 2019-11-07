package conv

func MustInterfaceToString(vi interface{}) string {
	switch vi.(type) {
	case string:
		return vi.(string)
	case []byte:
		return string(vi.([]byte))
	}
	return ""
}

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
	}
	return 0
}

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
	}
	return 0
}
