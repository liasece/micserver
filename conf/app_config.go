package conf

type AppConfig struct {
	AppSettings map[string]string        `json:"settings"`
	Modules     map[string]*ModuleConfig `json:"modules"`
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
	for k, v := range this.AppSettings {
		if _, ok := res.Settings[k]; !ok {
			res.Settings[k] = v
		}
	}
	return &res
}
