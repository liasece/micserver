package conf

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
	// "errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	// "net/http"
	_ "net/http/pprof"
	// "os"
	// "os/signal"
	"github.com/liasece/micserver/comm"
	"path/filepath"
	"strconv"
	"strings"
	// "syscall"
	"sync"
	"time"
)

type DBTableConfig struct {
	// 数据库实例索引
	DBIndex uint32 `json:"dbindex"`
	// 该表对应的表名字
	TableName string `json:"name"`
}

type DBConfig struct {
	// 数据库实例
	// 	key 		---  value
	// 	数据库实例索引    连接字符串
	Dbs map[uint32]string `json:"dbs"`
	// 数据库 表 实例
	// 	key 					---  value
	//  表索引，用于哈希Openid        表配置对象
	Tables map[uint32]*DBTableConfig `json:"tables"`
}

type ServerConfig struct {
	Allprops        map[string]string
	Myserverinfo    comm.SServerInfo // 当前服务器信息
	Myservername    string           // 当前服务器名称
	TerminateServer bool             //  服务器是否需要结束了

	// 服务器数字版本
	// 命名规则为： YYYYMMDDhhmm (年月日时分)
	Version uint64

	dbs         DBConfig          // 数据库配置
	RedisConfig comm.SRedisConfig // Redis配置

	hasAutoConfig  bool
	hasConfigPprof bool
	loadConfigTime uint32
	mutex          sync.Mutex
}

func (this *ServerConfig) AutoConfig(servername string, servertype uint32) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[ServerConfig.AutoConfig] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.loadConfigTime++
	rand.Seed(time.Now().UnixNano())

	this.initParse()
	this.Myservername = servername
	if servername == "SuperServer" || servername == "LoginServer" {
		this.Myservername = servername +
			fmt.Sprintf("%03d", this.getPropUintUnsafe("servernumber"))
	}

	this.loadXMLConfigFile("config.xml")
	err := util.LoadJsonFromFile("dbconfig.json", &this.dbs)
	if err != nil {
		log.Error("[ServerConfig.AutoConfig] 加载数据库配置出错 Err[%s]",
			err.Error())
	}
	this.Myserverinfo.Servertype = servertype

	// 初始化日志文件
	// 重新设置日志文件目录
	logpath := this.getPropUnsafe("logdir")
	if !this.hasAutoConfig {
		logsubname := ""
		if this.getPropUnsafe("servernumber") != "" {
			logsubname = "default_" + strings.ToLower(this.Myservername) +
				fmt.Sprintf("%03d",
					this.getPropUintUnsafe("servernumber")) + ".log"
		} else {
			logsubname = "default_" + strings.ToLower(this.Myservername) +
				".log"
		}
		if servername == "SuperServer" || servername == "LoginServer" {
			logsubname = strings.ToLower(this.Myservername) + ".log"
		}
		logfilename := filepath.Join(logpath, logsubname)
		daemon := this.getPropUnsafe("daemon")
		if daemon == "true" {
			log.AddlogFile(logfilename, true)
			log.RemoveConsoleLog()
			log.Debug("Program is start as a daemon")
		} else {
			log.AddlogFile(logfilename, false)
		}
	} else {
		logsubname := strings.ToLower(this.Myservername) + ".log"
		logfilename := filepath.Join(logpath, logsubname)
		log.ChangelogFile(logfilename)
	}
	// 设置日志级别
	log.SetLogLevel(this.getPropUnsafe("loglevel"))

	// 配置 pprof
	if !this.hasConfigPprof {
		// 外部性能监视
		if this.getPropUnsafe("performance_test") == "true" {
			// 获取 pprof Port
			// pprofport := this.getPropUintUnsafe("pprofport")
			// 获取 pprof IP
			// ifname := this.getPropUnsafe("ifname")
			// localip := util.GetIPv4ByInterface(ifname)
			// if pprofport > 0 && pprofport < 65536 && localip != "" {
			// 	err := BindPprof(localip, pprofport)
			// 	if err == nil {
			// 		this.hasConfigPprof = true
			// 	}
			// } else {
			// 	log.Debug("[ServerConfig.AutoConfig] 未设置 pprof "+
			// 		"IP/Port[%s:%d]",
			// 		localip, pprofport)
			// }
		} else {
			log.Debug("[ServerConfig.AutoConfig] pprof 不启动 "+
				"performance_test[%s]",
				this.getPropUnsafe("performance_test"))
		}
	}

	this.hasAutoConfig = true

	content, _ := json.Marshal(this)
	log.Info("[AutoConfig] 第%d次加载配置完成 配置信息： %s",
		this.loadConfigTime, content)
}

func (this *ServerConfig) ReloadConfig() {
	this.AutoConfig(this.Myservername, this.Myserverinfo.Servertype)
}

func (this *ServerConfig) GetTablesSum() uint32 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return uint32(len(this.dbs.Tables))
}

func (this *ServerConfig) GetTableInfo(
	tableindex uint32) (*DBTableConfig, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, finded := this.dbs.Tables[tableindex]; !finded {
		return nil,
			fmt.Errorf("tableindex %d dose't exit", tableindex)
	}
	return this.dbs.Tables[tableindex], nil
}

func (this *ServerConfig) GetDBsDBConfigs() map[uint32]string {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.dbs.Dbs
}

func (this *ServerConfig) GetProp(propname string) string {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getPropUnsafe(propname)
}

func (this *ServerConfig) GetPropInt(propname string) int32 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getPropIntUnsafe(propname)
}

func (this *ServerConfig) GetPropUint(propname string) uint32 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getPropUintUnsafe(propname)
}

func (this *ServerConfig) GetPropBool(propname string) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getPropBoolUnsafe(propname)
}

func (this *ServerConfig) getPropUnsafe(propname string) string {
	if propvalue, found := this.Allprops[propname+"_s"]; found {
		return propvalue
	}
	if propvalue, found := this.Allprops[propname]; found {
		return propvalue
	}
	return ""
}

func (this *ServerConfig) getPropIntUnsafe(propname string) int32 {
	retvalue, _ := strconv.Atoi(this.getPropUnsafe(propname))
	return int32(retvalue)
}

func (this *ServerConfig) getPropUintUnsafe(propname string) uint32 {
	retvalue, _ := strconv.Atoi(this.getPropUnsafe(propname))
	return uint32(retvalue)
}

func (this *ServerConfig) getPropBoolUnsafe(propname string) bool {
	retvalue := this.getPropUnsafe(propname)
	if retvalue == "true" || retvalue == "True" || retvalue == "TRUE" {
		return true
	}
	return false
}

func (this *ServerConfig) SetProp(propname string, value string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.setProp(propname, value)
}

func (this *ServerConfig) setProp(propname string, value string) {
	this.Allprops[propname] = value
}

func (this *ServerConfig) loadXMLConfigFile(filename string) bool {
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
	if len(this.getPropUnsafe("servername")) > 0 {
		this.Myservername = this.getPropUnsafe("servername")
	}
	if len(this.Myserverinfo.Servername) > 0 {
		this.Myservername = this.Myserverinfo.Servername
	}
	log.Debug("加载配置文件,servername:%s", this.Myservername)
	this.parse_token(decoder, t)

	// ifname := this.getPropUnsafe("ifname")
	// localip := util.GetIPv4ByInterface(ifname)
	// this.Myserverinfo.Serverip = localip
	log.Debug("加载配置文件完成,servername:%s", this.Myservername)
	return true
}

func (this *ServerConfig) parse_token(decoder *xml.Decoder,
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
			nodename := token.Name.Local
			if (len(this.Myservername) > 0 && nodename == this.Myservername) ||
				nodename == "global" {
				checkservername = true
			}
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
				this.Allprops["superserverport"] = superserverport
				superserverid := attrs["serverid"]
				this.Allprops["superserverid"] = superserverid
			}
			continue
			// 处理元素结束（标签）
		case xml.EndElement:
			token := t.(xml.EndElement)
			nodename := token.Name.Local
			if (len(this.Myservername) > 0 && nodename == this.Myservername) ||
				nodename == "global" {
				checkservername = false
			}
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
				this.Allprops[propname] = content
				if propname == "superserver" {
					this.Allprops["superserverip"] = content
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
func (this *ServerConfig) initParse() {
	if this.hasAutoConfig {
		return
	}
	var daemonflag string
	flag.StringVar(&daemonflag, "d", "", "as a daemon true or false")

	var serverflag string
	flag.StringVar(&serverflag, "s", "", "server name info as UserServer01")

	var lognameflag string
	flag.StringVar(&lognameflag, "l", "", "log name  as /log/gatewayserver.log")

	var servernumber string
	flag.StringVar(&servernumber, "n", "", "server number  as 0 1 2...")

	var serverversion string
	flag.StringVar(&serverversion, "v", "", "server version  as [0-9]{12}")

	flag.Parse()

	if len(daemonflag) > 0 {
		if daemonflag == "true" {
			this.Allprops["daemon_s"] = "true"
		} else {
			this.Allprops["daemon_s"] = "false"
		}
	}
	if len(serverflag) > 0 {
		this.Allprops["servername_s"] = serverflag
	}
	if len(lognameflag) > 0 {
		this.Allprops["logfilename_s"] = lognameflag
	}
	if len(servernumber) > 0 {
		this.Allprops["servernumber_s"] = servernumber
	}
	if len(serverversion) > 0 {
		tint, err := strconv.ParseUint(serverversion, 10, 64)
		if err == nil {
			this.Version = tint
		}
	}
}
