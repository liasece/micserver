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

func LoadConfig(filepath string) (*TopConfig, error) {
	res := &TopConfig{}
	err := res.loadJsonConfigFile(filepath)
	if err != nil {
		log.Error("loadJsonConfigFile(filepath) err:%v", err)
		return nil, err
	}
	res.InitParse()
	res.AppConfig.BuildModuleIDFromMapkey()
	return res, nil
}

func (this *TopConfig) GetArgvModuleList() []string {
	res := make([]string, len(this.argvModuleList))
	for i, v := range this.argvModuleList {
		res[i] = v
	}
	return res
}

func (this *TopConfig) HasGlobalConfig(key string) bool {
	return this.hasPropUnsafe(key)
}

func (this *TopConfig) GetGlobalConfig(key string) string {
	return this.getPropUnsafe(key)
}

func (this *TopConfig) GetProp(propname string) string {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getPropUnsafe(propname)
}

func (this *TopConfig) GetPropInt(propname string) int32 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getPropIntUnsafe(propname)
}

func (this *TopConfig) GetPropUint(propname string) uint32 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getPropUintUnsafe(propname)
}

func (this *TopConfig) GetPropBool(propname string) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getPropBoolUnsafe(propname)
}

func (this *TopConfig) getPropUnsafe(propname string) string {
	if propvalue, found := this.globalProp[propname+"_s"]; found {
		return propvalue
	}
	if propvalue, found := this.globalProp[propname]; found {
		return propvalue
	}
	return ""
}

func (this *TopConfig) hasPropUnsafe(propname string) bool {
	if _, found := this.globalProp[propname+"_s"]; found {
		return true
	}
	if _, found := this.globalProp[propname]; found {
		return true
	}
	return false
}

func (this *TopConfig) getPropIntUnsafe(propname string) int32 {
	retvalue, _ := strconv.Atoi(this.getPropUnsafe(propname))
	return int32(retvalue)
}

func (this *TopConfig) getPropUintUnsafe(propname string) uint32 {
	retvalue, _ := strconv.Atoi(this.getPropUnsafe(propname))
	return uint32(retvalue)
}

func (this *TopConfig) getPropBoolUnsafe(propname string) bool {
	retvalue := this.getPropUnsafe(propname)
	if retvalue == "true" || retvalue == "True" || retvalue == "TRUE" {
		return true
	}
	return false
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

// 参数解析相关
func (this *TopConfig) InitParse() {
	if this.globalProp == nil {
		this.globalProp = make(map[string]string)
	}
	for k, v := range cmdLineArgv {
		this.globalProp[k] = v
	}
	if cmdLineArgv != nil {
		if v, ok := cmdLineArgv["isdaemon"]; ok {
			this.AppConfig.Set(IsDaemon, v)
		}
		if v, ok := cmdLineArgv["logpath"]; ok {
			this.AppConfig.Set(LogWholePath, v)
		}
		if v, ok := cmdLineArgv["processid"]; ok {
			this.AppConfig.Set(ProcessID, v)
			if v != "development" {
				this.argvModuleList = strings.Split(v, ",")
			}
		}

		if v, ok := cmdLineArgv["version"]; ok {
			tint, err := strconv.ParseUint(v, 10, 64)
			if err == nil {
				this.Version = tint
			}
			this.AppConfig.Set(Version, v)
		}
	}
}

var cmdLineArgv map[string]string

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
