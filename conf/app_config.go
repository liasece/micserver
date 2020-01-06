package conf

import (
	"strings"
)

// App的配置，包括App全局配置以及配置文件中的模块配置
type AppConfig struct {
	*BaseConfig `json:"settings"`
	Modules     map[string]*ModuleConfig `json:"modules"`
}

// 将模块配置中的模块ID设置为模块配置json块的键值
func (this *AppConfig) buildModuleIDFromMapkey() {
	for k, m := range this.Modules {
		m.ID = k
		m.AppSettings = this.BaseConfig
		if m.Settings == nil {
			m.Settings = NewBaseConfig()
		}
	}
}

// 获取目标模块ID的模块配置，返回模块的指针，如果目标模块不存在，返回nil。
// 为了性能考虑不进行拷贝，
// 调用方不允许使用代码修改 ModuleConfig 的内容，
// 你应该修改配置文件而不是用代码设置配置值。
func (this *AppConfig) GetModuleConfig(moduleid string) *ModuleConfig {
	var res *ModuleConfig
	if v, ok := this.Modules[moduleid]; ok {
		res = v
	} else {
		return nil
	}
	return res
}

// 获取所有的模块配置
func (this *AppConfig) GetModuleConfigList() map[string]*ModuleConfig {
	return this.Modules
}

// 获取配置中存在的所有模块的subnet地址
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
