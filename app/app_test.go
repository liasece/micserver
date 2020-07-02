// Package app micserver 最基础的运行单位，app中包含了多个module，app在启动时会初始化所有module，
// 并且根据配置初始化module之间的连接。
// App 是 micserver 中在 "Module" 上一层的概念，使用 micserver 的
// 第一步就是实例化出一个 App 对象，并且向其中插入你的 Modules 。
// 建议一个代码上下文中仅存在一个 App 对象，如果你的需求让你觉得你有
// 必要实例化多个 App 在同一个可执行文件中，那么你应该考虑增加一个
// Module 而不是 App 。
package app

import (
	"testing"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/module"
)

func TestApp_App(t *testing.T) {
	log.SetDefaultLogger(log.NewLogger(&log.Options{NoConsoleColor: true}))

	app := &App{}
	app.Setup(nil)

	modules := make([]module.IModule, 0)
	for i := 0; i < 10; i++ {
		modules = append(modules, &module.BaseModule{})
	}

	app.Init(modules)
	go app.RunAndBlock(modules)
	app.Stop()
}
