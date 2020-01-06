package conf

import (
	"github.com/liasece/micserver/util/conv"
)

// 模块的配置，也包括了该模块所属的App配置
type ModuleConfig struct {
	ID          string      `json:"id"`
	Settings    *BaseConfig `json:"settings"`
	AppSettings *BaseConfig `json:"-"`
}

// 获取配置值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (this *ModuleConfig) get(key ConfigKey) interface{} {
	if this.Settings.Exist(key) {
		return this.Settings.Get(key)
	}
	if this.AppSettings.Exist(key) {
		return this.AppSettings.Get(key)
	}
	return nil
}

// 判断配置是否存在，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (this *ModuleConfig) exist(key ConfigKey) bool {
	if this.Settings.Exist(key) {
		return true
	}
	if this.AppSettings.Exist(key) {
		return true
	}
	return false
}

// 判断目标配置是否在 Module 配置中存在，无视 App 的配置
func (this *ModuleConfig) existInModule(key ConfigKey) bool {
	if this.Settings.Exist(key) {
		return true
	}
	return false
}

// 判断目标配置是否在 App 配置中存在，无视 Module 的配置
func (this *ModuleConfig) existInApp(key ConfigKey) bool {
	if this.AppSettings.Exist(key) {
		return true
	}
	return false
}

// 判断配置是否存在，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (this *ModuleConfig) Exist(key ConfigKey) bool {
	if this == nil {
		return false
	}
	return this.exist(key)
}

// 判断目标配置是否在 Module 配置中存在，无视 App 的配置
func (this *ModuleConfig) ExistInModule(key ConfigKey) bool {
	if this == nil {
		return false
	}
	return this.existInModule(key)
}

// 判断目标配置是否在 App 配置中存在，无视 Module 的配置
func (this *ModuleConfig) ExistInApp(key ConfigKey) bool {
	if this == nil {
		return false
	}
	return this.existInApp(key)
}

// 获取配置值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (this *ModuleConfig) Get(key ConfigKey) interface{} {
	if this == nil {
		return nil
	}
	return this.get(key)
}

// 获取配置 bool 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (this *ModuleConfig) GetBool(key ConfigKey) bool {
	if this == nil {
		return false
	}
	v := this.get(key)
	if v == nil {
		return false
	}
	return conv.MustInterfaceToBool(v)
}

// 获取配置 string 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (this *ModuleConfig) GetString(key ConfigKey) string {
	if this == nil {
		return ""
	}
	v := this.get(key)
	if v == nil {
		return ""
	}
	return conv.MustInterfaceToString(v)
}

// 获取配置 int64 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (this *ModuleConfig) GetInt64(key ConfigKey) int64 {
	if this == nil {
		return 0
	}
	v := this.get(key)
	if v == nil {
		return 0
	}
	return conv.MustInterfaceToInt64(v)
}

// 获取配置 []int64 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (this *ModuleConfig) GetInt64Slice(key ConfigKey) []int64 {
	if this == nil {
		return nil
	}
	v := this.get(key)
	if v == nil {
		return nil
	}
	return conv.MustInterfaceToInt64Slice(v)
}

// 获取配置 []string 值，先尝试获取 Module 中的配置，如果 Module 中不存在该配置，
// 再尝试获取 App 中的配置，如果都不存在则返回 nil
func (this *ModuleConfig) GetStringSlice(key ConfigKey) []string {
	if this == nil {
		return nil
	}
	v := this.get(key)
	if v == nil {
		return nil
	}
	return conv.MustInterfaceToStringSlice(v)
}
