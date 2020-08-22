package conf

import (
	"github.com/liasece/micserver/util/conv"
)

// BaseModuleConfig origin module config
type BaseModuleConfig struct {
	ID       string      `json:"id"`
	Settings *BaseConfig `json:"settings"`
}

// ModuleConfig 模块的配置，也包括了该模块所属的App配置
type ModuleConfig struct {
	BaseModuleConfig
	AppSettings *BaseConfig `json:"-"`
}

// get 获取配置值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (mc *ModuleConfig) get(key ConfigKey) interface{} {
	if mc.Settings.Exist(key) {
		return mc.Settings.Get(key)
	}
	if mc.AppSettings.Exist(key) {
		return mc.AppSettings.Get(key)
	}
	return nil
}

// exist 判断配置是否存在，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (mc *ModuleConfig) exist(key ConfigKey) bool {
	if mc.Settings.Exist(key) {
		return true
	}
	if mc.AppSettings.Exist(key) {
		return true
	}
	return false
}

// existInModule 判断目标配置是否在 Module 配置中存在，无视 App 的配置
func (mc *ModuleConfig) existInModule(key ConfigKey) bool {
	if mc.Settings.Exist(key) {
		return true
	}
	return false
}

// existInApp 判断目标配置是否在 App 配置中存在，无视 Module 的配置
func (mc *ModuleConfig) existInApp(key ConfigKey) bool {
	if mc.AppSettings.Exist(key) {
		return true
	}
	return false
}

// Exist 判断配置是否存在，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (mc *ModuleConfig) Exist(key ConfigKey) bool {
	if mc == nil {
		return false
	}
	return mc.exist(key)
}

// ExistInModule 判断目标配置是否在 Module 配置中存在，无视 App 的配置
func (mc *ModuleConfig) ExistInModule(key ConfigKey) bool {
	if mc == nil {
		return false
	}
	return mc.existInModule(key)
}

// ExistInApp 判断目标配置是否在 App 配置中存在，无视 Module 的配置
func (mc *ModuleConfig) ExistInApp(key ConfigKey) bool {
	if mc == nil {
		return false
	}
	return mc.existInApp(key)
}

// Get 获取配置值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (mc *ModuleConfig) Get(key ConfigKey) interface{} {
	if mc == nil {
		return nil
	}
	return mc.get(key)
}

// GetBool 获取配置 bool 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (mc *ModuleConfig) GetBool(key ConfigKey) bool {
	if mc == nil {
		return false
	}
	v := mc.get(key)
	if v == nil {
		return false
	}
	return conv.MustInterfaceToBool(v)
}

// GetString 获取配置 string 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (mc *ModuleConfig) GetString(key ConfigKey) string {
	if mc == nil {
		return ""
	}
	v := mc.get(key)
	if v == nil {
		return ""
	}
	return conv.MustInterfaceToString(v)
}

// GetInt64 获取配置 int64 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (mc *ModuleConfig) GetInt64(key ConfigKey) int64 {
	if mc == nil {
		return 0
	}
	v := mc.get(key)
	if v == nil {
		return 0
	}
	return conv.MustInterfaceToInt64(v)
}

// GetInt64Slice 获取配置 []int64 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (mc *ModuleConfig) GetInt64Slice(key ConfigKey) []int64 {
	if mc == nil {
		return nil
	}
	v := mc.get(key)
	if v == nil {
		return nil
	}
	return conv.MustInterfaceToInt64Slice(v)
}

// GetStringSlice 获取配置 []string 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (mc *ModuleConfig) GetStringSlice(key ConfigKey) []string {
	if mc == nil {
		return nil
	}
	v := mc.get(key)
	if v == nil {
		return nil
	}
	return conv.MustInterfaceToStringSlice(v)
}
