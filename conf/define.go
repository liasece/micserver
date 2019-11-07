package conf

type ConfigKey string

var (
	Version       ConfigKey = "version"
	ProcessID     ConfigKey = "processid"
	LogWholePath  ConfigKey = "logpath"
	SubnetTCPAddr ConfigKey = "subnettcpaddr"
	SubnetNoChan  ConfigKey = "subnetnochan"
	GateTCPAddr   ConfigKey = "gatetcpaddr"
	IsDaemon      ConfigKey = "isdaemon"
	MsgThreadNum  ConfigKey = "msgthreadnum"
)
