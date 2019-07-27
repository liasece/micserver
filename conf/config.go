package conf

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
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
	globalProp map[string]string `json:"-"`

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

func (this *TopConfig) SetProp(propname string, value string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.setProp(propname, value)
}

func (this *TopConfig) setProp(propname string, value string) {
	this.globalProp[propname] = value
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

func (this *TopConfig) loadXMLConfigFile(filename string) bool {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return false
	}
	str := string(content)
	// 必须把换行符去掉，不然解析不出chardata数据
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	decoder := xml.NewDecoder(bytes.NewBuffer([]byte(str)))
	t, decerr := decoder.Token()
	if decerr != nil {
		fmt.Println(decerr)
		return false
	}

	// 优先使用命令行的
	// if len(this.getPropUnsafe("modulename")) > 0 {
	// 	this.Myservername = this.getPropUnsafe("modulename")
	// }
	// if len(this.Myserverinfo.Servername) > 0 {
	// 	this.Myservername = this.Myserverinfo.Servername
	// }
	// log.Debug("加载配置文件,modulename:%s", this.Myservername)
	this.parse_token(decoder, t)

	// ifname := this.getPropUnsafe("ifname")
	// localip := util.GetIPv4ByInterface(ifname)
	// this.Myserverinfo.Serverip = localip
	log.Debug("加载配置文件完成")
	return true
}

func (this *TopConfig) parse_token(decoder *xml.Decoder,
	xmltoken xml.Token) {
	var t xml.Token
	var err error
	var propname string

	checkservername := false
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:
			token := t.(xml.StartElement)
			// nodename := token.Name.Local
			// if (len(this.Myservername) > 0 && nodename == this.Myservername) ||
			// 	nodename == "global" {
			// 	checkservername = true
			// }
			if !checkservername {
				continue
			}
			propname = token.Name.Local
			var attrs = make(map[string]string)
			for _, attr := range token.Attr {
				attrname := attr.Name.Local
				attrvalue := attr.Value
				attrs[attrname] = attrvalue
				//log.Debug("读取服务器配置信息222,%s,%s", attrname, attrvalue)
			}
			if propname == "superserver" {
				superserverport := attrs["port"]
				this.globalProp["superserverport"] = superserverport
				superserverid := attrs["serverid"]
				this.globalProp["superserverid"] = superserverid
			}
			continue
			// 处理元素结束（标签）
		case xml.EndElement:
			// token := t.(xml.EndElement)
			// nodename := token.Name.Local
			// if (len(this.Myservername) > 0 && nodename == this.Myservername) ||
			// 	nodename == "global" {
			// 	checkservername = false
			// }
			continue
		case xml.Comment:
			continue
		case xml.CharData:
			if !checkservername {
				continue
			}
			token := t.(xml.CharData)
			content := string([]byte(token))
			if propname != "" && len(content) > 0 {
				this.globalProp[propname] = content
				if propname == "superserver" {
					this.globalProp["superserverip"] = content
				}
			}
			propname = ""
			continue
		case xml.Directive:
		default:
			continue
		}
	}
}

// 参数解析相关
func (this *TopConfig) InitParse() {
	if this.globalProp == nil {
		this.globalProp = make(map[string]string)
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
			this.globalProp["isdaemon_s"] = "true"
		} else {
			this.globalProp["daemon_s"] = "false"
		}
		if this.AppConfig.AppSettings != nil {
			this.AppConfig.AppSettings["isdaemon"] = daemonflag
		}
	}
	if len(processflag) > 0 {
		this.globalProp["processid_s"] = processflag
	} else {
		this.globalProp["processid_s"] = "development"
	}
	if len(logpathflag) > 0 {
		this.globalProp["logpath_s"] = logpathflag
		if this.AppConfig.AppSettings != nil {
			this.AppConfig.AppSettings["logpath"] = logpathflag
		}
	}
	if len(serverversion) > 0 {
		tint, err := strconv.ParseUint(serverversion, 10, 64)
		if err == nil {
			this.Version = tint
		}
	}
}
