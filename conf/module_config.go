package conf

import (
	"github.com/liasece/micserver/util/conv"
)

type ModuleConfig struct {
	ID          string      `json:"id"`
	Settings    *BaseConfig `json:"settings"`
	AppSettings *BaseConfig `json:"-"`
}

func (this *ModuleConfig) HasModuleSetting(key ConfigKey) bool {
	return this.Settings.Exist(key)
}

func (this *ModuleConfig) get(key ConfigKey) interface{} {
	if this.Settings.Exist(key) {
		return this.Settings.Get(key)
	}
	if this.AppSettings.Exist(key) {
		return this.AppSettings.Get(key)
	}
	return nil
}

func (this *ModuleConfig) exist(key ConfigKey) bool {
	if this.Settings.Exist(key) {
		return true
	}
	if this.AppSettings.Exist(key) {
		return true
	}
	return false
}

func (this *ModuleConfig) existInModule(key ConfigKey) bool {
	if this.Settings.Exist(key) {
		return true
	}
	return false
}

func (this *ModuleConfig) existInApp(key ConfigKey) bool {
	if this.AppSettings.Exist(key) {
		return true
	}
	return false
}

func (this *ModuleConfig) Exist(key ConfigKey) bool {
	if this == nil {
		return false
	}
	return this.exist(key)
}

func (this *ModuleConfig) ExistInModule(key ConfigKey) bool {
	if this == nil {
		return false
	}
	return this.existInModule(key)
}

func (this *ModuleConfig) ExistInApp(key ConfigKey) bool {
	if this == nil {
		return false
	}
	return this.existInApp(key)
}

func (this *ModuleConfig) Get(key ConfigKey) interface{} {
	if this == nil {
		return nil
	}
	return this.get(key)
}

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
