package micserver

import "github.com/liasece/liasece/server"

func CreateApp() module.App {
	return defaultApp.NewApp(Version)
}
