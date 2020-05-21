package conf

import (
	"github.com/liasece/micserver/util/conv"
)

// BaseConfig 配置由 encoding/json 读入为 map[string]interface{} 类型的结构。
// 不可写，所以不对外提供写接口
type BaseConfig map[string]interface{}

// NewBaseConfig 获取一个空的 BaseConfig 结构
func NewBaseConfig() *BaseConfig {
	res := make(BaseConfig)
	return &res
}

func (bf *BaseConfig) get(key ConfigKey) interface{} {
	if v, ok := (*bf)[string(key)]; ok {
		return v
	}
	return nil
}

func (bf *BaseConfig) set(key ConfigKey, v interface{}) {
	(*bf)[string(key)] = v
}

func (bf *BaseConfig) exist(key ConfigKey) bool {
	_, ok := (*bf)[string(key)]
	return ok
}

// Exist 判断键是否存在
func (bf *BaseConfig) Exist(key ConfigKey) bool {
	if bf == nil {
		return false
	}
	return bf.exist(key)
}

// Get 获取指定键的值
func (bf *BaseConfig) Get(key ConfigKey) interface{} {
	if bf == nil {
		return nil
	}
	return bf.get(key)
}

// GetBool 获取指定键的bool类型值
func (bf *BaseConfig) GetBool(key ConfigKey) bool {
	if bf == nil {
		return false
	}
	v := bf.get(key)
	if v == nil {
		return false
	}
	return conv.MustInterfaceToBool(v)
}

// GetString 获取指定键的 string 类型值
func (bf *BaseConfig) GetString(key ConfigKey) string {
	if bf == nil {
		return ""
	}
	v := bf.get(key)
	if v == nil {
		return ""
	}
	return conv.MustInterfaceToString(v)
}

// GetInt64 获取指定键的 int64 类型值
func (bf *BaseConfig) GetInt64(key ConfigKey) int64 {
	if bf == nil {
		return 0
	}
	v := bf.get(key)
	if v == nil {
		return 0
	}
	return conv.MustInterfaceToInt64(v)
}

// GetInt64Slice 获取指定键的 []int64 类型值
func (bf *BaseConfig) GetInt64Slice(key ConfigKey) []int64 {
	if bf == nil {
		return nil
	}
	v := bf.get(key)
	if v == nil {
		return nil
	}
	return conv.MustInterfaceToInt64Slice(v)
}

// GetStringSlice 获取指定键的 []string 类型值
func (bf *BaseConfig) GetStringSlice(key ConfigKey) []string {
	if bf == nil {
		return nil
	}
	v := bf.get(key)
	if v == nil {
		return nil
	}
	return conv.MustInterfaceToStringSlice(v)
}
