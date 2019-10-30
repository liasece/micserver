package servercomm

type SServerInfo struct {
	ServerID   string
	ServerAddr string
	// 服务器序号 重复不影响正常运行
	// 但是其改动会影响 配置读取/ServerName/Log文件名
	ServerNumber uint32
	// 服务器数字版本
	// 命名规则为： YYYYMMDDhhmm (年月日时分)
	Version uint64
}

type STimeTickCommand struct {
	Testno uint32
}

type STestCommand struct {
	Testno     uint32
	Testttring string // IP
}

type SLoginCommand struct {
	ServerID   string
	ServerAddr string // IP
	// 登录优先级
	ConnectPriority int64
	// 服务器序号 重复不影响正常运行
	// 但是其改动会影响 配置读取/ServerName/Log文件名
	ServerNumber uint32
	// 服务器数字版本
	// 命名规则为： YYYYMMDDhhmm (年月日时分)
	Version uint64
}

// 通知服务器正常退出
type SLogoutCommand struct {
}

// 通知我所连接的服务器启动成功
type SSeverStartOKCommand struct {
	Serverid uint32
}

const (
	// 登录成功
	LOGINRETCODE_SECESS = 0
	// 身份验证错误
	LOGINRETCODE_IDENTITY = 1
	// 重复连接
	LOGINRETCODE_IDENTICAL = 2
)

// 登录服务器返回
type SLoginRetCommand struct {
	Loginfailed uint32       // 是否连接成功,0成功
	Destination *SServerInfo //	tcptask 所在服务器信息
}

// super通知其它服务器启动成功
type SStartRelyNotifyCommand struct {
	Serverinfos []*SServerInfo // 启动成功服务器信息
}

// 启动验证通过通知其它服务器我的新
type SStartMyNotifyCommand struct {
	Serverinfo *SServerInfo // 启动成功服务器信息
}

// super 通知所有服务器配置信息
type SNotifyAllInfo struct {
	Serverinfos []*SServerInfo // 成功服务器信息
}

// 通知说明的目标服务器安全退出
// 此消息发送的前提是当前存在可以替代目标服务器的服务器
type SNotifySafelyQuit struct {
	// 目标服务器的信息应该是最新的信息，目标服务器会将该信息改成最新的
	TargetServerInfo *SServerInfo
}

type SUpdateSession struct {
	ClientConnID string
	Session      map[string]string
}

type SForwardToServer struct {
	FromServerID string
	ToServerID   string
	MsgID        uint16
	Data         []byte
}

type SForwardToClient struct {
	FromServerID string
	ToGateID     string
	ToClientID   string
	MsgID        uint16
	Data         []byte
}

type SForwardFromGate struct {
	FromServerID string
	ToServerID   string
	ClientConnID string
	Session      map[string]string
	MsgID        uint16
	Data         []byte
}

// ROC调用请求
type SROCRequest struct {
	// 请求信息
	FromServerID string
	ToServerID   string
	Seq          int64
	// 调用信息
	CallStr    string
	CallArg    []byte
	NeedReturn bool
}

// ROC调用响应
type SROCResponse struct {
	// 响应信息
	FromServerID string
	ToServerID   string
	ReqSeq       int64
	// 响应数据
	ResData []byte
	Error   string
}

// ROC绑定信息
type SROCBind struct {
	HostServerID string
	IsDelete     bool
	ObjType      string
	ObjIDs       []string
}
