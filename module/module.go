package module

import (
	"github.com/liasece/micserver/log"
)

type App struct {
}

func (this *App) NewApp(version string) *App {
	log.AutoConfig("/home/jansen/logs/main.log", "Main", true)
	return &App{}
}

func (this *App) Init() {
	log.Debug("hello world!")
}
