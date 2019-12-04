package conf

import (
	"github.com/liasece/micserver/util/conv"
)

type BaseConfig map[string]interface{}

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

func (this *BaseConfig) CopyFrom(src *BaseConfig) {
	if this == nil {
		return
	}
	for k, v := range *src {
		this.set(ConfigKey(k), v)
	}
}

func (this *BaseConfig) Exist(key ConfigKey) bool {
	if this == nil {
		return false
	}
	return this.exist(key)
}

func (this *BaseConfig) Get(key ConfigKey) interface{} {
	if this == nil {
		return nil
	}
	return this.get(key)
}

func (this *BaseConfig) Set(key ConfigKey, v interface{}) {
	if this == nil {
		return
	}
	this.set(key, v)
}

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
