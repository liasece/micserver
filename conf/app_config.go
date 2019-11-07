package conf

import (
	"strings"
)

type AppConfig struct {
	*BaseConfig `json:"settings"`
	Modules     map[string]*ModuleConfig `json:"modules"`
}

func (this *AppConfig) BuildModuleIDFromMapkey() {
	for k, m := range this.Modules {
		m.ID = k
	}
}

func (this *AppConfig) GetModuleConfig(moduleid string) *ModuleConfig {
	var res ModuleConfig
	if v, ok := this.Modules[moduleid]; ok {
		res = *v
	} else {
		res.ID = moduleid
	}
	if res.Settings == nil {
		res.Settings = NewBaseConfig()
	}
	if res.AppSettings == nil {
		res.AppSettings = NewBaseConfig()
	}
	res.AppSettings.CopyFrom(this.BaseConfig)
	return &res
}

func (this *AppConfig) GetModuleConfigList() map[string]*ModuleConfig {
	res := make(map[string]*ModuleConfig)
	for k, _ := range this.Modules {
		if !strings.HasPrefix(k, "//") {
			res[k] = this.GetModuleConfig(k)
		}
	}
	return res
}

func (this *AppConfig) GetSubnetTCPAddrMap() map[string]string {
	res := make(map[string]string)
	for k, m := range this.Modules {
		if !strings.HasPrefix(k, "//") {
			if m.Settings.Exist(SubnetTCPAddr) {
				res[k] = m.Settings.GetString(SubnetTCPAddr)
			}
		}
	}
	return res
}
