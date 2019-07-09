package conf

type ModuleItem struct {
	ID        string            `json:"id"`
	ProcessID string            `json:"processid"`
	Settings  map[string]string `json:"settings"`
}
