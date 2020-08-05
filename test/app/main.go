package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/liasece/micserver"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/module"
	"github.com/liasece/micserver/util"
	"github.com/liasece/micserver/util/monitor"
)

// GatewayModule module
type GatewayModule struct {
	module.BaseModule

	testSeqTimes    int64
	testCheckTimeNS int64
	testSwitch      bool
	// 模块的负载
	clientMsgLoad          monitor.Load
	lastCheckClientMsgLoad int64
}

// NewGatewayModule func
func NewGatewayModule(moduleid string) *GatewayModule {
	res := &GatewayModule{}
	res.BaseModule.SetModuleID(moduleid)
	return res
}

// PlayerModule module
type PlayerModule struct {
	module.BaseModule
}

// NewPlayerModule func
func NewPlayerModule(moduleid string) *PlayerModule {
	res := &PlayerModule{}
	res.BaseModule.SetModuleID(moduleid)
	return res
}

// LoginModule module
type LoginModule struct {
	module.BaseModule
}

// NewLoginModule func
func NewLoginModule(moduleid string) *LoginModule {
	res := &LoginModule{}
	res.BaseModule.SetModuleID(moduleid)
	return res
}

// InitManager init manager
type InitManager struct {
	modules    map[string]module.IModule
	configPath string

	hasInit bool
	mutex   sync.Mutex
}

var mInitManager *InitManager
var mInitManagerLock sync.Once

// GetInitManger func
func GetInitManger() *InitManager {
	mInitManagerLock.Do(func() {
		mInitManager = &InitManager{}
	})
	return mInitManager
}

// GetProgramModuleList func
func (im *InitManager) GetProgramModuleList() []module.IModule {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	if !im.hasInit {
		im.hasInit = true
		im.modules = make(map[string]module.IModule)
		config := conf.TopConfig{}
		config.InitFlag()
		isDevelopment := true

		// 遍历所有的参数指定的模块名
		for _, pid := range config.GetArgvModuleList() {
			isDevelopment = false
			basepid := pid
			for i := 0; i < 1; i++ {
				if i != 0 {
					pid = fmt.Sprint(basepid, "_", i)
				}
				stype := util.GetModuleIDType(pid)
				log.Debug("App initialize ServerType[%s] ServerID[%s]", stype, pid)
				switch stype {
				case "gate":
					im.addModule(NewGatewayModule(pid))
				case "player":
					im.addModule(NewPlayerModule(pid))
				case "login":
					im.addModule(NewLoginModule(pid))
				default:
					panic(fmt.Sprintf("无法解析的模块 %s:%s", stype, pid))
				}
			}
		}

		// 如果当前是开发模式，添加如下的列表
		if isDevelopment {
			im.addModule(NewGatewayModule("gate001"))
			im.addModule(NewGatewayModule("gate002"))
		}
	}
	return im.getModuleSlice()
}

// 添加一个模块
func (im *InitManager) addModule(module module.IModule) bool {
	if im.modules == nil {
		return false
	}
	im.modules[module.GetModuleID()] = module
	return true
}

// 获取模块列表的切片形式
func (im *InitManager) getModuleSlice() []module.IModule {
	res := make([]module.IModule, 0)
	for _, m := range im.modules {
		res = append(res, m)
	}
	return res
}

// GetConfigPath 获取配置文件的路径
func (im *InitManager) GetConfigPath() string {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	if im.configPath == "" {
		pwd, err := os.Getwd()
		if err == nil {
			im.configPath = filepath.Join(pwd, "config", "config.json")
		} else {
			log.Error("os.Getwd() err:%v", err)
			return ""
		}
	}
	return im.configPath
}

func main() {
	// 初始化 MicServer
	app, err := micserver.SetupApp(GetInitManger().GetConfigPath())
	if err != nil {
		log.Fatal("Create app fatal: %v", err)
		time.Sleep(time.Second * 1)
		return
	}

	log.Info("即将运行")
	// 初始化性能监控
	monitor.BindPprof("", 8888)

	// app 开始运行 阻塞
	app.RunAndBlock(GetInitManger().GetProgramModuleList())
}
