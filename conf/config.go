// Package conf 包含micserver中使用的 配置文件/命令行参数 的配置
package conf

import (
	// pprof
	_ "net/http/pprof"

	"encoding/json"
	"flag"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/liasece/micserver/log"
)

// TopConfig 所有配置，包括从文件以进命令行输入的所有配置内容
type TopConfig struct {
	AppConfig AppConfig `json:"app"` // 进程配置信息

	// 全局配置字符串，一般是从进程启动时携带的参数提供的
	globalProp     map[string]string
	argvModuleList []string

	// 服务器数字版本
	// 建议命名规则为： YYYYMMDDhhmmss (年月日时分秒)
	Version uint64 `json:"-"`

	hasConfigPprof bool
	loadConfigTime uint32
	mutex          sync.Mutex
}

// LoadConfig 从配置文件以及命令行参数中加载配置
func LoadConfig(filepath string) (*TopConfig, error) {
	res := &TopConfig{}
	err := res.loadJSONConfigFile(filepath)
	if err != nil {
		log.Error("loadJSONConfigFile(filepath) err:%v", err)
		return nil, err
	}
	res.InitParse()
	res.AppConfig.buildModuleIDFromMapkey()
	return res, nil
}

func (tc *TopConfig) loadJSONConfigFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &tc)
	if err != nil {
		return err
	}
	return nil
}

// GetArgvModuleList 命令行参数中的模块名列表，可使用命令行参数的 -p 选项携带，例如：
// ./foo -p gate001,gate002,logic001
// 表明该进程中启动的模块包括 gate001 gate002 logic001
func (tc *TopConfig) GetArgvModuleList() []string {
	res := make([]string, len(tc.argvModuleList))
	for i, v := range tc.argvModuleList {
		res[i] = v
	}
	return res
}

// HasProp 判断命令行参数是否携带名为key的参数
func (tc *TopConfig) HasProp(key string) bool {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()
	return tc.hasProp(key)
}

// GetProp 获取名为key的命令行参数
func (tc *TopConfig) GetProp(propname string) string {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()
	return tc.getProp(propname)
}

// GetPropInt64 获取int类型，名为key的命令行参数
func (tc *TopConfig) GetPropInt64(propname string) int64 {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()
	return int64(tc.getPropInt(propname))
}

// GetPropBool 获取bool类型，名为key的命令行参数
func (tc *TopConfig) GetPropBool(propname string) bool {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()
	return tc.getPropBool(propname)
}

// getProp 获取命令行参数，无锁
func (tc *TopConfig) getProp(propname string) string {
	if propvalue, found := tc.globalProp[propname+"_s"]; found {
		return propvalue
	}
	if propvalue, found := tc.globalProp[propname]; found {
		return propvalue
	}
	return ""
}

func (tc *TopConfig) hasProp(propname string) bool {
	if _, found := tc.globalProp[propname+"_s"]; found {
		return true
	}
	if _, found := tc.globalProp[propname]; found {
		return true
	}
	return false
}

func (tc *TopConfig) getPropInt(propname string) int {
	retvalue, _ := strconv.Atoi(tc.getProp(propname))
	return retvalue
}

func (tc *TopConfig) getPropBool(propname string) bool {
	retvalue := tc.getProp(propname)
	if retvalue == "true" || retvalue == "True" || retvalue == "TRUE" {
		return true
	}
	return false
}

// InitParse 初始化命令行参数
func (tc *TopConfig) InitParse() {
	if tc.globalProp == nil {
		tc.globalProp = make(map[string]string)
	}
	if tc.AppConfig.BaseConfig == nil {
		tc.AppConfig.BaseConfig = &BaseConfig{}
	}
	for k, v := range cmdLineArgv {
		tc.globalProp[k] = v
	}
	if cmdLineArgv != nil {
		if v, ok := cmdLineArgv["isdaemon"]; ok {
			tc.AppConfig.set(IsDaemon, v)
		}
		if v, ok := cmdLineArgv["logpath"]; ok {
			tc.AppConfig.set(LogWholePath, v)
		}
		if v, ok := cmdLineArgv["processid"]; ok {
			tc.AppConfig.set(ProcessID, v)
			if v != "development" {
				tc.argvModuleList = strings.Split(v, ",")
			}
		}

		if v, ok := cmdLineArgv["version"]; ok {
			tint, err := strconv.ParseUint(v, 10, 64)
			if err == nil {
				tc.Version = tint
			}
			tc.AppConfig.set(Version, v)
		}
	}
}

// 命令行参数列表
var cmdLineArgv map[string]string

// init 读取命令行参数
func init() {
	if cmdLineArgv == nil {
		cmdLineArgv = make(map[string]string)
	}
	var daemonflag string
	flag.StringVar(&daemonflag, "d", "", "as a daemon true or false")

	var processflag string
	flag.StringVar(&processflag, "p", "", "process id as gate001")

	var logpathflag string
	flag.StringVar(&logpathflag, "l", "", "log path as /log/")

	var serverversion string
	flag.StringVar(&serverversion, "v", "", "server version as [0-9]{14}")

	flag.Parse()

	if len(daemonflag) > 0 {
		if daemonflag == "true" {
			cmdLineArgv["isdaemon"] = "true"
		} else {
			cmdLineArgv["isdaemon"] = "false"
		}
	}
	if len(processflag) > 0 {
		cmdLineArgv["processid"] = processflag
	} else {
		cmdLineArgv["processid"] = "development"
	}
	if len(logpathflag) > 0 {
		cmdLineArgv["logpath"] = logpathflag
	}
	if len(serverversion) > 0 {
		cmdLineArgv["version"] = serverversion
	}
}
