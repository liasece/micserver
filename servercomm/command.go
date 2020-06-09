// Package servercomm micserver中模块间发送的消息的定义，由 *_binary.go 实现二进制协议。
package servercomm

// ModuleInfo 一个模块的信息
type ModuleInfo struct {
	ModuleID   string
	ModuleAddr string
	// 服务器序号 重复不影响正常运行
	// 但是其改动会影响 配置读取/ServerName/Log文件名
	ModuleNumber uint32
	// 服务器数字版本
	// 命名规则为： YYYYMMDDhhmm (年月日时分)
	Version uint64
}

// STimeTickCommand 心跳包请求
type STimeTickCommand struct {
	Testno uint32
}

// STestCommand 测试消息请求
type STestCommand struct {
	Testno     uint32
	Testttring string // IP
}

// SLoginCommand 模块登陆请求
type SLoginCommand struct {
	ModuleID   string
	ModuleAddr string // IP
	// 登录优先级
	ConnectPriority int64
	// 服务器序号 重复不影响正常运行
	// 但是其改动会影响 配置读取/ServerName/Log文件名
	ModuleNumber uint32
	// 服务器数字版本
	// 命名规则为： YYYYMMDDhhmm (年月日时分)
	Version uint64
}

// SLogoutCommand 通知服务器正常退出
type SLogoutCommand struct {
}

// SSeverStartOKCommand 通知我所连接的服务器启动成功
type SSeverStartOKCommand struct {
	ModuleID string
}

// login result
const (
	// 登录成功
	LoginRetCodeSecess = 0
	// 身份验证错误
	LoginRetCodeIdentity = 1
	// 重复连接
	LoginRetCodeIdentical = 2
)

// SLoginRetCommand 登录服务器返回
type SLoginRetCommand struct {
	Loginfailed uint32      // 是否连接成功,0成功
	Destination *ModuleInfo //	tcptask 所在服务器信息
}

// SStartRelyNotifyCommand super通知其它服务器启动成功
type SStartRelyNotifyCommand struct {
	ServerInfos []*ModuleInfo // 启动成功服务器信息
}

// SStartMyNotifyCommand 启动验证通过通知其它服务器我的新
type SStartMyNotifyCommand struct {
	ModuleInfo *ModuleInfo // 启动成功服务器信息
}

// SNotifyAllInfo super 通知所有服务器配置信息
type SNotifyAllInfo struct {
	ServerInfos []*ModuleInfo // 成功服务器信息
}

// SNotifySafelyQuit 通知说明的目标服务器安全退出
// 此消息发送的前提是当前存在可以替代目标服务器的服务器
type SNotifySafelyQuit struct {
	// 目标服务器的信息应该是最新的信息，目标服务器会将该信息改成最新的
	TargetServerInfo *ModuleInfo
}

// SUpdateSession 更新Session的请求
type SUpdateSession struct {
	FromModuleID string
	ToModuleID   string
	ClientConnID string
	SessionUUID  string
	Session      map[string]string
}

// SReqCloseConnect 关闭客户端连接的请求
type SReqCloseConnect struct {
	FromModuleID string
	ToModuleID   string
	ClientConnID string
}

// SForwardToModule 模块间消息转发请求
type SForwardToModule struct {
	FromModuleID string
	ToModuleID   string
	MsgID        uint16
	Data         []byte
}

// ModuleMessage 模块间传递的消息
type ModuleMessage struct {
	FromModule *ModuleInfo
	MsgID      uint16
	Data       []byte
}

// SForwardToClient 请求转发一个客户端消息
type SForwardToClient struct {
	FromModuleID string
	ToGateID     string
	ToClientID   string
	MsgID        uint16
	Data         []byte
}

// SForwardFromGate 网关分发的客户端消息
type SForwardFromGate struct {
	FromModuleID string
	ToModuleID   string
	ClientConnID string
	Session      map[string]string
	MsgID        uint16
	Data         []byte
}

// ClientMessage 客户端消息
type ClientMessage struct {
	FromModule   *ModuleInfo
	ClientConnID string
	MsgID        uint16
	Data         []byte
}

// SROCRequest ROC调用请求
type SROCRequest struct {
	// 请求信息
	FromModuleID string
	ToModuleID   string
	Seq          int64
	// 调用信息
	CallStr    string
	CallArg    []byte
	NeedReturn bool
}

// SROCResponse ROC调用响应
type SROCResponse struct {
	// 响应信息
	FromModuleID string
	ToModuleID   string
	ReqSeq       int64
	// 响应数据
	ResData []byte
	Error   string
}

// SROCBind ROC绑定信息
type SROCBind struct {
	HostModuleID string
	IsDelete     bool
	ObjType      string
	ObjIDs       []string
}
