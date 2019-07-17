package conf

type ModuleConfig struct {
	ID        string            `json:"id"`
	ProcessID string            `json:"processid"`
	Settings  map[string]string `json:"settings"`
}

func (this *ModuleConfig) HasSetting(key string) bool {
	if this.Settings == nil {
		return false
	}
	if _, ok := this.Settings[key]; ok {
		return true
	}
	return false
}

func (this *ModuleConfig) GetSetting(key string) string {
	if this.Settings == nil {
		return ""
	}
	if v, ok := this.Settings[key]; ok {
		return v
	}
	return ""
}
