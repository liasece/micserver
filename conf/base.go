package conf

import (
	"github.com/liasece/micserver/util/conv"
)

// 配置由 encoding/json 读入为 map[string]interface{} 类型的结构。
// 不可写，所以不对外提供写接口
type BaseConfig map[string]interface{}

// 获取一个空的 BaseConfig 结构
func NewBaseConfig() *BaseConfig {
	res := make(BaseConfig)
	return &res
}

func (this *BaseConfig) get(key ConfigKey) interface{} {
	if v, ok := (*this)[string(key)]; ok {
		return v
	}
	return nil
}

func (this *BaseConfig) set(key ConfigKey, v interface{}) {
	(*this)[string(key)] = v
}

func (this *BaseConfig) exist(key ConfigKey) bool {
	_, ok := (*this)[string(key)]
	return ok
}

// 判断键是否存在
func (this *BaseConfig) Exist(key ConfigKey) bool {
	if this == nil {
		return false
	}
	return this.exist(key)
}

// 获取指定键的值
func (this *BaseConfig) Get(key ConfigKey) interface{} {
	if this == nil {
		return nil
	}
	return this.get(key)
}

// 获取指定键的bool类型值
func (this *BaseConfig) GetBool(key ConfigKey) bool {
	if this == nil {
		return false
	}
	v := this.get(key)
	if v == nil {
		return false
	}
	return conv.MustInterfaceToBool(v)
}

// 获取指定键的 string 类型值
func (this *BaseConfig) GetString(key ConfigKey) string {
	if this == nil {
		return ""
	}
	v := this.get(key)
	if v == nil {
		return ""
	}
	return conv.MustInterfaceToString(v)
}

// 获取指定键的 int64 类型值
func (this *BaseConfig) GetInt64(key ConfigKey) int64 {
	if this == nil {
		return 0
	}
	v := this.get(key)
	if v == nil {
		return 0
	}
	return conv.MustInterfaceToInt64(v)
}

// 获取指定键的 []int64 类型值
func (this *BaseConfig) GetInt64Slice(key ConfigKey) []int64 {
	if this == nil {
		return nil
	}
	v := this.get(key)
	if v == nil {
		return nil
	}
	return conv.MustInterfaceToInt64Slice(v)
}

// 获取指定键的 []string 类型值
func (this *BaseConfig) GetStringSlice(key ConfigKey) []string {
	if this == nil {
		return nil
	}
	v := this.get(key)
	if v == nil {
		return nil
	}
	return conv.MustInterfaceToStringSlice(v)
}
