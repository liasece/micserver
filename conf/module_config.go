package conf

type ModuleConfig struct {
	ID          string            `json:"id"`
	Settings    map[string]string `json:"settings"`
	AppSettings map[string]string `json:"-"`
}

func (this *ModuleConfig) HasModuleSetting(key string) bool {
	if this.Settings == nil {
		return false
	}
	if _, ok := this.Settings[key]; ok {
		return true
	}
	return false
}

func (this *ModuleConfig) GetModuleSetting(key string) string {
	if this.Settings == nil {
		return ""
	}
	if v, ok := this.Settings[key]; ok {
		return v
	}
	return ""
}

func (this *ModuleConfig) GetModuleSettingMap() map[string]string {
	res := make(map[string]string)
	for k, v := range this.AppSettings {
		res[k] = v
	}
	for k, v := range this.Settings {
		res[k] = v
	}
	return res
}

func (this *ModuleConfig) HasSetting(key string) bool {
	if this.Settings != nil {
		if _, ok := this.Settings[key]; ok {
			return true
		}
	}
	if this.AppSettings != nil {
		if _, ok := this.AppSettings[key]; ok {
			return true
		}
	}
	return false
}

func (this *ModuleConfig) GetSetting(key string) string {
	if this.Settings != nil {
		if v, ok := this.Settings[key]; ok {
			return v
		}
	}
	if this.AppSettings != nil {
		if v, ok := this.AppSettings[key]; ok {
			return v
		}
	}
	return ""
}
