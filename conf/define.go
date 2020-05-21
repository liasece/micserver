package conf

// ConfigKey micserver中配置的键类型
type ConfigKey string

// micserver框架内置的配置选项
var (
	// 版本号		string
	Version ConfigKey = "version"
	// 进程ID		string
	ProcessID ConfigKey = "processid"
	// log完整路径		string
	LogWholePath ConfigKey = "logpath"
	// log等级		string
	LogLevel ConfigKey = "loglevel"
	// 服务器TCP子网ip及端口，如 1.0.0.1:80 		string
	SubnetTCPAddr ConfigKey = "subnettcpaddr"
	// 不使用本地chan		bool
	SubnetNoChan ConfigKey = "subnetnochan"
	// 网关TCP地址		string
	GateTCPAddr ConfigKey = "gatetcpaddr"
	// 是否是受保护的进程 		bool
	IsDaemon ConfigKey = "isdaemon"
	// 消息处理并发协程数量		int
	MsgThreadNum ConfigKey = "msgthreadnum"
	// ROC的绑定是否使用异步方式同步到别的module中，会与ROC调用有异步问题 bool
	AsynchronousSyncRocbind ConfigKey = "asynchronous_sync_rocbind"
)
