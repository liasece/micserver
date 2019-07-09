package conf

type AppConfig struct {
	Modules map[string][]*ModuleItem `json:"module"`
}
