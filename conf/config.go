/*
包含micserver中使用的 配置文件/命令行参数 的配置
*/
package conf

import (
	"encoding/json"
	"flag"
	"github.com/liasece/micserver/log"
	"io/ioutil"
	_ "net/http/pprof"
	"strconv"
	"strings"
	"sync"
)

// 所有配置，包括从文件以进命令行输入的所有配置内容
type TopConfig struct {
	AppConfig AppConfig `json:"app"` // 进程配置信息

	// 全局配置字符串，一般是从进程启动时携带的参数提供的
	globalProp     map[string]string `json:"-"`
	argvModuleList []string

	// 服务器数字版本
	// 建议命名规则为： YYYYMMDDhhmmss (年月日时分秒)
	Version uint64 `json:"-"`

	hasConfigPprof bool       `json:"-"`
	loadConfigTime uint32     `json:"-"`
	mutex          sync.Mutex `json:"-"`
}

// 从配置文件以及命令行参数中加载配置
func LoadConfig(filepath string) (*TopConfig, error) {
	res := &TopConfig{}
	err := res.loadJsonConfigFile(filepath)
	if err != nil {
		log.Error("loadJsonConfigFile(filepath) err:%v", err)
		return nil, err
	}
	res.InitParse()
	res.AppConfig.buildModuleIDFromMapkey()
	return res, nil
}

func (this *TopConfig) loadJsonConfigFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &this)
	if err != nil {
		return err
	}
	return nil
}

// 命令行参数中的模块名列表，可使用命令行参数的 -p 选项携带，例如：
// ./foo -p gate001,gate002,logic001
// 表明该进程中启动的模块包括 gate001 gate002 logic001
func (this *TopConfig) GetArgvModuleList() []string {
	res := make([]string, len(this.argvModuleList))
	for i, v := range this.argvModuleList {
		res[i] = v
	}
	return res
}

// 判断命令行参数是否携带名为key的参数
func (this *TopConfig) HasProp(key string) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.hasProp(key)
}

// 获取名为key的命令行参数
func (this *TopConfig) GetProp(propname string) string {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getProp(propname)
}

// 获取int类型，名为key的命令行参数
func (this *TopConfig) GetPropInt64(propname string) int64 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return int64(this.getPropInt(propname))
}

// 获取bool类型，名为key的命令行参数
func (this *TopConfig) GetPropBool(propname string) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getPropBool(propname)
}

// 获取命令行参数，无锁
func (this *TopConfig) getProp(propname string) string {
	if propvalue, found := this.globalProp[propname+"_s"]; found {
		return propvalue
	}
	if propvalue, found := this.globalProp[propname]; found {
		return propvalue
	}
	return ""
}

func (this *TopConfig) hasProp(propname string) bool {
	if _, found := this.globalProp[propname+"_s"]; found {
		return true
	}
	if _, found := this.globalProp[propname]; found {
		return true
	}
	return false
}

func (this *TopConfig) getPropInt(propname string) int {
	retvalue, _ := strconv.Atoi(this.getProp(propname))
	return retvalue
}

func (this *TopConfig) getPropBool(propname string) bool {
	retvalue := this.getProp(propname)
	if retvalue == "true" || retvalue == "True" || retvalue == "TRUE" {
		return true
	}
	return false
}

// 初始化命令行参数
func (this *TopConfig) InitParse() {
	if this.globalProp == nil {
		this.globalProp = make(map[string]string)
	}
	if this.AppConfig.BaseConfig == nil {
		this.AppConfig.BaseConfig = &BaseConfig{}
	}
	for k, v := range cmdLineArgv {
		this.globalProp[k] = v
	}
	if cmdLineArgv != nil {
		if v, ok := cmdLineArgv["isdaemon"]; ok {
			this.AppConfig.set(IsDaemon, v)
		}
		if v, ok := cmdLineArgv["logpath"]; ok {
			this.AppConfig.set(LogWholePath, v)
		}
		if v, ok := cmdLineArgv["processid"]; ok {
			this.AppConfig.set(ProcessID, v)
			if v != "development" {
				this.argvModuleList = strings.Split(v, ",")
			}
		}

		if v, ok := cmdLineArgv["version"]; ok {
			tint, err := strconv.ParseUint(v, 10, 64)
			if err == nil {
				this.Version = tint
			}
			this.AppConfig.set(Version, v)
		}
	}
}

// 命令行参数列表
var cmdLineArgv map[string]string

// 读取命令行参数
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
