package micserver

import "github.com/liasece/micserver"

func CreateApp() module.App {
	return defaultApp.NewApp(Version)
}
