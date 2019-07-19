package conf

type AppConfig struct {
	AppSettings map[string]string        `json:"settings"`
	Modules     map[string]*ModuleConfig `json:"modules"`
}

func (this *AppConfig) BuildModuleIDFromMapkey() {
	for k, m := range this.Modules {
		m.ID = k
	}
}

func (this *AppConfig) HasSetting(key string) bool {
	if this.AppSettings == nil {
		return false
	}
	if _, ok := this.AppSettings[key]; ok {
		return true
	}
	return false
}

func (this *AppConfig) GetSetting(key string) string {
	if this.AppSettings == nil {
		return ""
	}
	if v, ok := this.AppSettings[key]; ok {
		return v
	}
	return ""
}

func (this *AppConfig) GetModuleConfig(moduleid string) *ModuleConfig {
	var res ModuleConfig
	if v, ok := this.Modules[moduleid]; ok {
		res = *v
	}
	if res.Settings == nil {
		res.Settings = make(map[string]string)
	}
	if res.AppSettings == nil {
		res.AppSettings = make(map[string]string)
	}
	for k, v := range this.AppSettings {
		if _, ok := res.AppSettings[k]; !ok {
			res.AppSettings[k] = v
		}
	}
	return &res
}

func (this *AppConfig) GetModuleConfigList() map[string]*ModuleConfig {
	res := make(map[string]*ModuleConfig)
	for k, _ := range this.Modules {
		res[k] = this.GetModuleConfig(k)
	}
	return res
}

func (this *AppConfig) GetSubnetTCPAddrMap() map[string]string {
	res := make(map[string]string)
	for k, m := range this.Modules {
		if addr, ok := m.Settings["subnettcpaddr"]; ok {
			res[k] = addr
		}
	}
	return res
}
