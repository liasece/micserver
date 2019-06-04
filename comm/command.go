package comm

type SServerInfo struct {
	Serverid      uint32 `xml:"serverid,attr"`
	Servertype    uint32 `xml:"servertype,attr"` // 类型
	Servername    string `xml:"servername,attr"`
	Serverip      string `xml:"serverip,attr"`
	Serverport    uint32 `xml:"serverport,attr"`
	Extip         string `xml:"extip,attr"`
	Httpport      uint32 `xml:"httpport,attr"`
	Httpsport     uint32 `xml:"httpsport,attr"`
	Rpcport       uint32 `xml:"rpcport,attr"`
	Tcpport       uint32 `xml:"tcpport,attr"`
	ClientTcpport uint32 `xml:"tcpsocketport,attr"`

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
	Serverid   uint32
	Servertype uint32 // 类型
	Serverip   string // IP
	Servername string // name
	Serverport uint32 // port

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

// 登录服务器返回
type SLoginRetCommand struct {
	Loginfailed uint32       // 是否连接成功,0成功
	Clientinfo  SServerInfo  // 连接服务器client信息
	Taskinfo    SServerInfo  //	tcptask 所在服务器信息
	Redisinfo   SRedisConfig // Redis连接配置
}

// super通知其它服务器启动成功
type SStartRelyNotifyCommand struct {
	Serverinfos []SServerInfo // 启动成功服务器信息
}

// 启动验证通过通知其它服务器我的新
type SStartMyNotifyCommand struct {
	Serverinfo SServerInfo // 启动成功服务器信息
}

// super 通知所有服务器配置信息
type SNotifyAllInfo struct {
	Serverinfos []SServerInfo // 成功服务器信息
}
type SUpdateGatewayUserAnalysis struct {
	Httpcount         uint32 // 5分钟的量
	Webscoketcount    uint32 // 5分钟的量
	Webscoketcurcount uint32 // 当前的量
}

// 通知super 新添加用户 存入redis
type SAddNewUserToRedisCommand struct {
	Openid       string
	Serverid     uint32
	UUID         uint64
	ClientConnID uint64
}

// gateway 转发过来的消息数据 gateway <<-->> QuizServer
type SGatewayForwardCommand struct {
	Gateserverid uint32
	ClientConnID uint64
	Openid       string
	UUID         uint64
	Cmdid        uint16
	Cmdlen       uint16
	Cmddatas     []byte
}

// gateway 广播转发消息数据 gateway <<-->> clients
type SGatewayForwardBroadcastCommand struct {
	Gateserverid uint32
	// 处理线程哈希，需要提供给Gateway作发送协程选择，
	// 合理配置该项可以保证消息达到顺序且提高gateway性能，否则，
	// 两条广播消息的到达顺序无法保证！！！
	ThreadHash uint32
	UUIDList   []uint64
	Cmdid      uint16
	Cmdlen     uint16
	Cmddatas   []byte
}

// 向Gateway发送GM返回数据，Userserver -> Gatewayserver
type SGatewayForward2HttpCommand struct {
	Gateserverid uint32
	Httptaskid   uint64
	Openid       string
	Cmdname      string
	Cmdlen       uint16
	Cmddatas     []byte
}

//  转发消息给指定的userserver  user->bridge->user 直接给用户
type SBridgeForward2UserCommand struct {
	Fromuuid uint64
	Touuid   uint64
	Cmdid    uint16
	Cmdlen   uint16
	Cmddatas []byte
}

//  广播消息给所有的userserver  userserver->bridge->user 直接给用户
type SBridgeBroadcast2UserCommand struct {
	Fromopenid string
	Toopenid   string
	Cmdid      uint16
	Cmdlen     uint16
	Cmddatas   []byte
}

//  转发消息给指定的userserver  userserver->bridge->userserver 直接给用户所在的server
type SBridgeForward2UserServer struct {
	Fromuuid uint64
	Touuid   uint64
	Cmdid    uint16
	Cmdlen   uint16
	Cmddatas []byte
}

//  转发消息给所有 GatewayServer
//  GatewayServer->bridge->GatewayServer 直接给用户Client
type SBridgeBroadcast2GatewayServer struct {
	Cmdid    uint16
	Cmdlen   uint16
	Cmddatas []byte
}

//  发送消息给指定的userserver  matchserver->userserver 直接给用户所在的server
type SMatchForward2UserServer struct {
	Fromuuid uint64
	Touuid   uint64
	Cmdid    uint16
	Cmdlen   uint16
	Cmddatas []byte
}

//  发送消息给指定的userserver  roomserver->userserver 直接给用户所在的server
type SRoomForward2UserServer struct {
	Fromuuid uint64
	Touuid   uint64
	Cmdid    uint16
	Cmdlen   uint16
	Cmddatas []byte
}

// 网关广播消息发送给每个玩家
type SGatewayBroadcast2UserCommand struct {
	Fromuuid uint64
	Touuid   uint64
	Cmdid    uint16
	Cmdlen   uint16
	Cmddatas []byte
}

// 请求查询指定玩家的好友公开信息 方便搜索好友
type SUserServerSearchFriend struct {
	Fromuuid uint64
	Touuid   uint64
}

// 执行gm指令
type SUserServerGMCommand struct {
	Taskid uint64
	Key    string
	UUID   uint64
	Openid string
	CmdID  uint32
	Param1 string
	Param2 string
}

// 请求远程执行另外一个玩家的操作相关
type SRequestOtherUser struct {
	Fromuuid uint64
	Touuid   uint64
	Cmdid    uint16
	CmdData  []byte
}

// 响应另一个玩家的操作
type SResponseOtherUser struct {
	Fromuuid uint64
	Touuid   uint64
	Cmdid    uint16
	CmdData  []byte
}

// 去bridge获得一个随机的玩家信息 返回User
type SBridgeDialGetUserInfo struct {
	Fromopenid string
	Fromuuid   uint64
	Getuuid    uint64
	Getopenid  string
	Type       uint32
}

// 这个消息之前必定已经经过http login user
// gateway 发起用户登录 Gateway-->UserServer ->QuizServer-->Gateway
type SGatewayWSLoginUser struct {
	Gateserverid uint32
	ClientConnID uint64
	Openid       string
	UUID         uint64
	Token        string
	Tokenendtime uint32
	Sessionkey   string
	Loginappid   string
	Username     string
	Quizid       uint64
	Allmoney     uint64
	Headurl      string
	Female       uint32 //是否女性玩家  1 女玩家 0男玩家
	Retcode      int32  // 0 是成功，非0都表示失败了
	Message      string
	LoginMsg     []byte // 原始登陆消息
	// retcode 1 服务器异常
	// retcode 2 用户登录异常
}

// gateway 发起用户下线 Gateway-->UserServer ->QuizServer
type SGatewayWSOfflineUser struct {
	Openid       string
	UUID         uint64
	Quizid       uint64
	ClientConnID uint64
}

type STemplateMessageKeyWord struct {
	Value string
	Color string
}

// QuizServer通知给玩家发送模板消息
type SQSTemplateMessage struct {
	Openid      string
	Template_id string
	Page        string
	Datalist    []STemplateMessageKeyWord
	Formid      string
}

// 网关通过super广播换accesstoken
type SGatewayChangeAccessToken struct {
	Access_token         string //  微信access_token
	Update_accesstime    uint32 // 需要更新access_token 的时间
	Access_token_QQ      string //  QQaccess_token
	Update_accesstime_QQ uint32 // 需要更新QQ access_token 的时间
}

// 用于在 UserServer MatchServer RoomServer之间互通的玩家信息
type SUserInfo struct {
	UUID            uint64
	Openid          string
	UserServerid    uint32
	GatewayServerid uint32 // 玩家GatewayServer的ID
	ClientConnID    uint64 // 玩家连接ID

	// 玩家段位
	RankID uint32

	// 是否是脚本玩家
	IsScript bool

	// 玩家的GM指令权限等级
	GMLevel uint64

	// 外部携带数据
	ExtraData []byte
}

// UserServer -> MatchServer 玩家请求加入匹配队列
type SJoinMatch struct {
	// 指定房间ID，如果为0，则由系统自动匹配
	RoomID uint64
	// 房间类型
	RoomType uint32
	// 子类型
	SubType uint32
	// 是否仅加入房间
	// 	true: 系统只会加入房间，而不会为其创建房间
	// 	false: 系统检查是否存在合适的房间，如果有，则加入；否则，为其创建一个房间
	OnlyJoin bool
	// 是否仅创建房间
	// 	true: 系统只会创建房间，不进行其他任何操作
	// 	false: 系统检查是否存在合适的房间，如果有，则加入；否则，为其创建一个房间
	OnlyCreate bool
	// 用户基础信息
	Userinfo SUserInfo
}

// MatchServer -> RoomServer 服务器间交互的队伍信息
type STeamInfo struct {
	Members []SUserInfo
}

// MatchServer -> RoomServer 服务器间交互的房间信息
type SRoomInfo struct {
	Teams []STeamInfo
}

// MatchServer -> UserServer 房间的信息
type SUserMatchInfo struct {
	Sec        bool   // 是否成功加入匹配队列
	Done       bool   // 是否已经匹配完成
	Total      uint32 // 当前匹配队列的总人数
	Now        uint32 // 当前匹配队列的人数
	Matchindex uint64 // 匹配队列索引，即RoomID
}

// UserServer -> MatchServer 用户请求退出匹配队列
// MatchServer -> UserServer 用户退出了队列
type SUserQuitMatch struct {
	RoomID     uint64 // 目标房间ID
	Sec        bool   // 操作是否成功
	UUID       uint64
	Type       uint32 // 退出队列的原因： 1主动退出 2用户下线 3匹配完成已销毁匹配队列
	Matchindex uint64 // 匹配队列索引
}

// UserServer -> MatchServer 强制匹配完成
type SMatchDone struct {
	Matchindex uint64 // 匹配队列索引
}

// UserServer -> RoomServer 用户请求退出对局房间
// RoomServer -> UserServer 用户退出了对局房间
type SUserQuitRoom struct {
	RoomID uint64 // 目标房间ID
	Sec    bool   // 操作是否成功
	UUID   uint64
	Type   uint32 // 退出队列的原因： 1主动退出 2用户下线 3游戏已完成销毁房间
	//  	4：用户检查不通过
	Roomindex uint64 // 匹配队列索引
}

// RoomServer -> UserServer 房间的信息
type SUserRoomInfo struct {
	Sec       bool   // 是否成功加入匹配队列
	Roomindex uint64 // 匹配队列索引
}

// UserServer -> MatchServer/RoomServer 通知匹配以及房间服务器，玩家下线了
type SNotifyUserOffline struct {
	UUID       uint64
	Matchindex uint64
	Roomindex  uint64
}

// UserServer -> MatchServer/RoomServer 通知匹配以及房间服务器，玩家重新上线了
type SNotifyUserReonline struct {
	UUID     uint64
	UserInfo SUserInfo
}

// rpc 调用UserServer检查用户的token是否有效
type SRPCCheckUserToken struct {
	Openid  string `json:"openid"`
	Token   string `json:"token"`
	Retcode int32  `json:"retcode"` // 1 表示成功
}

//  广播消息给所有的userserver  userserver->matchserver->userserver
type SMatchBroadcast2UserServerCommand struct {
	Fromuuid   uint64
	Matchindex uint64 // 匹配队列索引，即RoomID
	Cmdid      uint16
	Cmdlen     uint16
	Cmddatas   []byte
}

// 网关通知super自己的信息
type NotifyGatewayInfo struct {
	Serverid   uint32
	Serverip   string
	Serverport uint32
	State      uint16
	TaskSize   uint32
}

type SuperReqGatewayInfo struct {
}

// UserServer 通知 GatewayServer 用户数量
type NotifyGateUserNums struct {
	Usernum uint32
}

// GatewayServer请求UserServer用户数量
type GateReqUserNums struct {
}

// RoomServer 通知 MatchServer 房间数量
type NotifyMatchRoomNums struct {
	Roomnum uint32
}

// MatchServer 请求 RoomServer 房间数量
type MatchReqRoomNums struct {
}

// UserServer -> RoomServer 检查用户是否还存在
// 	如果不存在了，则 RoomServer 会发送 SUserQuitRoom
type SUserCheckEffective struct {
	UUID     uint64
	UserInfo SUserInfo
}

// Redis配置 单项
type SRedisConfigItem struct {
	IP   string
	Port uint32
}

// Redis配置
type SRedisConfig struct {
	RedisList []SRedisConfigItem
}

type SRequestServerInfo struct {
}

// 由 SuperServer 发送给其他服务器
// 通知说明的目标服务器安全退出
// 此消息发送的前提是当前存在可以替代目标服务器的服务器
type SNotifySafelyQuit struct {
	// 目标服务器的信息应该是最新的信息，目标服务器会将该信息改成最新的
	TargetServerInfo SServerInfo
}

// AI 玩家是一种用于计算的特殊玩家，需要在内网端口保证安全行，是可信的
// 可以将一些服务器不便于计算的服务提交到AI玩家进行计算
type SAIUserRegister struct {
	Userinfo SUserInfo
}

// S->C: 请求计算任务
// C->S: 回应计算任务
type STaskCmdForward struct {
	// 消息ID
	CmdID uint16
	// 房间ID
	RoomID uint64
	// 代理UUID
	ProxyUUID uint64
	// 消息Protobuf数据
	ProtoData []byte
}

// UserServer -> MatchServer 开始匹配其他玩家
type SBeginMatch struct {
	Matchindex uint64 // 匹配队列索引
}
