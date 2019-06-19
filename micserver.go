package micserver

import "github.com/liasece/micserver/module"

var Version string = "0.0.1"

var defaultApp module.App

func CreateApp() *module.App {
	return defaultApp.NewApp(Version)
}
