package comm
import (
	"encoding/binary"
	"math"
	"encoding/json"
)
const (
	SServerInfoID = 36
	STimeTickCommandID = 37
	STestCommandID = 38
	SLoginCommandID = 39
	SLogoutCommandID = 40
	SSeverStartOKCommandID = 41
	SLoginRetCommandID = 42
	SStartRelyNotifyCommandID = 43
	SStartMyNotifyCommandID = 44
	SNotifyAllInfoID = 45
	SUpdateGatewayUserAnalysisID = 46
	SAddNewUserToRedisCommandID = 47
	SGatewayForwardCommandID = 48
	SGatewayForwardBroadcastCommandID = 49
	SGatewayForward2HttpCommandID = 50
	SBridgeForward2UserCommandID = 51
	SBridgeBroadcast2UserCommandID = 52
	SBridgeForward2UserServerID = 53
	SBridgeBroadcast2GatewayServerID = 54
	SMatchForward2UserServerID = 55
	SRoomForward2UserServerID = 56
	SGatewayBroadcast2UserCommandID = 57
	SUserServerSearchFriendID = 58
	SUserServerGMCommandID = 59
	SRequestOtherUserID = 60
	SResponseOtherUserID = 61
	SBridgeDialGetUserInfoID = 62
	SGatewayWSLoginUserID = 63
	SGatewayWSOfflineUserID = 64
	STemplateMessageKeyWordID = 65
	SQSTemplateMessageID = 66
	SGatewayChangeAccessTokenID = 67
	SUserInfoID = 68
	SJoinMatchID = 69
	STeamInfoID = 70
	SRoomInfoID = 71
	SUserMatchInfoID = 72
	SUserQuitMatchID = 73
	SMatchDoneID = 74
	SUserQuitRoomID = 75
	SUserRoomInfoID = 76
	SNotifyUserOfflineID = 77
	SNotifyUserReonlineID = 78
	SRPCCheckUserTokenID = 79
	SMatchBroadcast2UserServerCommandID = 80
	NotifyGatewayInfoID = 81
	SuperReqGatewayInfoID = 82
	NotifyGateUserNumsID = 83
	GateReqUserNumsID = 84
	NotifyMatchRoomNumsID = 85
	MatchReqRoomNumsID = 86
	SUserCheckEffectiveID = 87
	SRedisConfigItemID = 88
	SRedisConfigID = 89
	SRequestServerInfoID = 90
	SNotifySafelyQuitID = 91
	SAIUserRegisterID = 92
	STaskCmdForwardID = 93
	SBeginMatchID = 94
)
const (
	SServerInfoName = "comm.SServerInfo"
	STimeTickCommandName = "comm.STimeTickCommand"
	STestCommandName = "comm.STestCommand"
	SLoginCommandName = "comm.SLoginCommand"
	SLogoutCommandName = "comm.SLogoutCommand"
	SSeverStartOKCommandName = "comm.SSeverStartOKCommand"
	SLoginRetCommandName = "comm.SLoginRetCommand"
	SStartRelyNotifyCommandName = "comm.SStartRelyNotifyCommand"
	SStartMyNotifyCommandName = "comm.SStartMyNotifyCommand"
	SNotifyAllInfoName = "comm.SNotifyAllInfo"
	SUpdateGatewayUserAnalysisName = "comm.SUpdateGatewayUserAnalysis"
	SAddNewUserToRedisCommandName = "comm.SAddNewUserToRedisCommand"
	SGatewayForwardCommandName = "comm.SGatewayForwardCommand"
	SGatewayForwardBroadcastCommandName = "comm.SGatewayForwardBroadcastCommand"
	SGatewayForward2HttpCommandName = "comm.SGatewayForward2HttpCommand"
	SBridgeForward2UserCommandName = "comm.SBridgeForward2UserCommand"
	SBridgeBroadcast2UserCommandName = "comm.SBridgeBroadcast2UserCommand"
	SBridgeForward2UserServerName = "comm.SBridgeForward2UserServer"
	SBridgeBroadcast2GatewayServerName = "comm.SBridgeBroadcast2GatewayServer"
	SMatchForward2UserServerName = "comm.SMatchForward2UserServer"
	SRoomForward2UserServerName = "comm.SRoomForward2UserServer"
	SGatewayBroadcast2UserCommandName = "comm.SGatewayBroadcast2UserCommand"
	SUserServerSearchFriendName = "comm.SUserServerSearchFriend"
	SUserServerGMCommandName = "comm.SUserServerGMCommand"
	SRequestOtherUserName = "comm.SRequestOtherUser"
	SResponseOtherUserName = "comm.SResponseOtherUser"
	SBridgeDialGetUserInfoName = "comm.SBridgeDialGetUserInfo"
	SGatewayWSLoginUserName = "comm.SGatewayWSLoginUser"
	SGatewayWSOfflineUserName = "comm.SGatewayWSOfflineUser"
	STemplateMessageKeyWordName = "comm.STemplateMessageKeyWord"
	SQSTemplateMessageName = "comm.SQSTemplateMessage"
	SGatewayChangeAccessTokenName = "comm.SGatewayChangeAccessToken"
	SUserInfoName = "comm.SUserInfo"
	SJoinMatchName = "comm.SJoinMatch"
	STeamInfoName = "comm.STeamInfo"
	SRoomInfoName = "comm.SRoomInfo"
	SUserMatchInfoName = "comm.SUserMatchInfo"
	SUserQuitMatchName = "comm.SUserQuitMatch"
	SMatchDoneName = "comm.SMatchDone"
	SUserQuitRoomName = "comm.SUserQuitRoom"
	SUserRoomInfoName = "comm.SUserRoomInfo"
	SNotifyUserOfflineName = "comm.SNotifyUserOffline"
	SNotifyUserReonlineName = "comm.SNotifyUserReonline"
	SRPCCheckUserTokenName = "comm.SRPCCheckUserToken"
	SMatchBroadcast2UserServerCommandName = "comm.SMatchBroadcast2UserServerCommand"
	NotifyGatewayInfoName = "comm.NotifyGatewayInfo"
	SuperReqGatewayInfoName = "comm.SuperReqGatewayInfo"
	NotifyGateUserNumsName = "comm.NotifyGateUserNums"
	GateReqUserNumsName = "comm.GateReqUserNums"
	NotifyMatchRoomNumsName = "comm.NotifyMatchRoomNums"
	MatchReqRoomNumsName = "comm.MatchReqRoomNums"
	SUserCheckEffectiveName = "comm.SUserCheckEffective"
	SRedisConfigItemName = "comm.SRedisConfigItem"
	SRedisConfigName = "comm.SRedisConfig"
	SRequestServerInfoName = "comm.SRequestServerInfo"
	SNotifySafelyQuitName = "comm.SNotifySafelyQuit"
	SAIUserRegisterName = "comm.SAIUserRegister"
	STaskCmdForwardName = "comm.STaskCmdForward"
	SBeginMatchName = "comm.SBeginMatch"
)
func (this *SServerInfo) WriteBinary(data []byte) int {
	return WriteMsgSServerInfoByObj(data,this)
}
func (this *STimeTickCommand) WriteBinary(data []byte) int {
	return WriteMsgSTimeTickCommandByObj(data,this)
}
func (this *STestCommand) WriteBinary(data []byte) int {
	return WriteMsgSTestCommandByObj(data,this)
}
func (this *SLoginCommand) WriteBinary(data []byte) int {
	return WriteMsgSLoginCommandByObj(data,this)
}
func (this *SLogoutCommand) WriteBinary(data []byte) int {
	return WriteMsgSLogoutCommandByObj(data,this)
}
func (this *SSeverStartOKCommand) WriteBinary(data []byte) int {
	return WriteMsgSSeverStartOKCommandByObj(data,this)
}
func (this *SLoginRetCommand) WriteBinary(data []byte) int {
	return WriteMsgSLoginRetCommandByObj(data,this)
}
func (this *SStartRelyNotifyCommand) WriteBinary(data []byte) int {
	return WriteMsgSStartRelyNotifyCommandByObj(data,this)
}
func (this *SStartMyNotifyCommand) WriteBinary(data []byte) int {
	return WriteMsgSStartMyNotifyCommandByObj(data,this)
}
func (this *SNotifyAllInfo) WriteBinary(data []byte) int {
	return WriteMsgSNotifyAllInfoByObj(data,this)
}
func (this *SUpdateGatewayUserAnalysis) WriteBinary(data []byte) int {
	return WriteMsgSUpdateGatewayUserAnalysisByObj(data,this)
}
func (this *SAddNewUserToRedisCommand) WriteBinary(data []byte) int {
	return WriteMsgSAddNewUserToRedisCommandByObj(data,this)
}
func (this *SGatewayForwardCommand) WriteBinary(data []byte) int {
	return WriteMsgSGatewayForwardCommandByObj(data,this)
}
func (this *SGatewayForwardBroadcastCommand) WriteBinary(data []byte) int {
	return WriteMsgSGatewayForwardBroadcastCommandByObj(data,this)
}
func (this *SGatewayForward2HttpCommand) WriteBinary(data []byte) int {
	return WriteMsgSGatewayForward2HttpCommandByObj(data,this)
}
func (this *SBridgeForward2UserCommand) WriteBinary(data []byte) int {
	return WriteMsgSBridgeForward2UserCommandByObj(data,this)
}
func (this *SBridgeBroadcast2UserCommand) WriteBinary(data []byte) int {
	return WriteMsgSBridgeBroadcast2UserCommandByObj(data,this)
}
func (this *SBridgeForward2UserServer) WriteBinary(data []byte) int {
	return WriteMsgSBridgeForward2UserServerByObj(data,this)
}
func (this *SBridgeBroadcast2GatewayServer) WriteBinary(data []byte) int {
	return WriteMsgSBridgeBroadcast2GatewayServerByObj(data,this)
}
func (this *SMatchForward2UserServer) WriteBinary(data []byte) int {
	return WriteMsgSMatchForward2UserServerByObj(data,this)
}
func (this *SRoomForward2UserServer) WriteBinary(data []byte) int {
	return WriteMsgSRoomForward2UserServerByObj(data,this)
}
func (this *SGatewayBroadcast2UserCommand) WriteBinary(data []byte) int {
	return WriteMsgSGatewayBroadcast2UserCommandByObj(data,this)
}
func (this *SUserServerSearchFriend) WriteBinary(data []byte) int {
	return WriteMsgSUserServerSearchFriendByObj(data,this)
}
func (this *SUserServerGMCommand) WriteBinary(data []byte) int {
	return WriteMsgSUserServerGMCommandByObj(data,this)
}
func (this *SRequestOtherUser) WriteBinary(data []byte) int {
	return WriteMsgSRequestOtherUserByObj(data,this)
}
func (this *SResponseOtherUser) WriteBinary(data []byte) int {
	return WriteMsgSResponseOtherUserByObj(data,this)
}
func (this *SBridgeDialGetUserInfo) WriteBinary(data []byte) int {
	return WriteMsgSBridgeDialGetUserInfoByObj(data,this)
}
func (this *SGatewayWSLoginUser) WriteBinary(data []byte) int {
	return WriteMsgSGatewayWSLoginUserByObj(data,this)
}
func (this *SGatewayWSOfflineUser) WriteBinary(data []byte) int {
	return WriteMsgSGatewayWSOfflineUserByObj(data,this)
}
func (this *STemplateMessageKeyWord) WriteBinary(data []byte) int {
	return WriteMsgSTemplateMessageKeyWordByObj(data,this)
}
func (this *SQSTemplateMessage) WriteBinary(data []byte) int {
	return WriteMsgSQSTemplateMessageByObj(data,this)
}
func (this *SGatewayChangeAccessToken) WriteBinary(data []byte) int {
	return WriteMsgSGatewayChangeAccessTokenByObj(data,this)
}
func (this *SUserInfo) WriteBinary(data []byte) int {
	return WriteMsgSUserInfoByObj(data,this)
}
func (this *SJoinMatch) WriteBinary(data []byte) int {
	return WriteMsgSJoinMatchByObj(data,this)
}
func (this *STeamInfo) WriteBinary(data []byte) int {
	return WriteMsgSTeamInfoByObj(data,this)
}
func (this *SRoomInfo) WriteBinary(data []byte) int {
	return WriteMsgSRoomInfoByObj(data,this)
}
func (this *SUserMatchInfo) WriteBinary(data []byte) int {
	return WriteMsgSUserMatchInfoByObj(data,this)
}
func (this *SUserQuitMatch) WriteBinary(data []byte) int {
	return WriteMsgSUserQuitMatchByObj(data,this)
}
func (this *SMatchDone) WriteBinary(data []byte) int {
	return WriteMsgSMatchDoneByObj(data,this)
}
func (this *SUserQuitRoom) WriteBinary(data []byte) int {
	return WriteMsgSUserQuitRoomByObj(data,this)
}
func (this *SUserRoomInfo) WriteBinary(data []byte) int {
	return WriteMsgSUserRoomInfoByObj(data,this)
}
func (this *SNotifyUserOffline) WriteBinary(data []byte) int {
	return WriteMsgSNotifyUserOfflineByObj(data,this)
}
func (this *SNotifyUserReonline) WriteBinary(data []byte) int {
	return WriteMsgSNotifyUserReonlineByObj(data,this)
}
func (this *SRPCCheckUserToken) WriteBinary(data []byte) int {
	return WriteMsgSRPCCheckUserTokenByObj(data,this)
}
func (this *SMatchBroadcast2UserServerCommand) WriteBinary(data []byte) int {
	return WriteMsgSMatchBroadcast2UserServerCommandByObj(data,this)
}
func (this *NotifyGatewayInfo) WriteBinary(data []byte) int {
	return WriteMsgNotifyGatewayInfoByObj(data,this)
}
func (this *SuperReqGatewayInfo) WriteBinary(data []byte) int {
	return WriteMsgSuperReqGatewayInfoByObj(data,this)
}
func (this *NotifyGateUserNums) WriteBinary(data []byte) int {
	return WriteMsgNotifyGateUserNumsByObj(data,this)
}
func (this *GateReqUserNums) WriteBinary(data []byte) int {
	return WriteMsgGateReqUserNumsByObj(data,this)
}
func (this *NotifyMatchRoomNums) WriteBinary(data []byte) int {
	return WriteMsgNotifyMatchRoomNumsByObj(data,this)
}
func (this *MatchReqRoomNums) WriteBinary(data []byte) int {
	return WriteMsgMatchReqRoomNumsByObj(data,this)
}
func (this *SUserCheckEffective) WriteBinary(data []byte) int {
	return WriteMsgSUserCheckEffectiveByObj(data,this)
}
func (this *SRedisConfigItem) WriteBinary(data []byte) int {
	return WriteMsgSRedisConfigItemByObj(data,this)
}
func (this *SRedisConfig) WriteBinary(data []byte) int {
	return WriteMsgSRedisConfigByObj(data,this)
}
func (this *SRequestServerInfo) WriteBinary(data []byte) int {
	return WriteMsgSRequestServerInfoByObj(data,this)
}
func (this *SNotifySafelyQuit) WriteBinary(data []byte) int {
	return WriteMsgSNotifySafelyQuitByObj(data,this)
}
func (this *SAIUserRegister) WriteBinary(data []byte) int {
	return WriteMsgSAIUserRegisterByObj(data,this)
}
func (this *STaskCmdForward) WriteBinary(data []byte) int {
	return WriteMsgSTaskCmdForwardByObj(data,this)
}
func (this *SBeginMatch) WriteBinary(data []byte) int {
	return WriteMsgSBeginMatchByObj(data,this)
}
func (this *SServerInfo) ReadBinary(data []byte) int {
	return ReadMsgSServerInfoByBytes(data, this)
}
func (this *STimeTickCommand) ReadBinary(data []byte) int {
	return ReadMsgSTimeTickCommandByBytes(data, this)
}
func (this *STestCommand) ReadBinary(data []byte) int {
	return ReadMsgSTestCommandByBytes(data, this)
}
func (this *SLoginCommand) ReadBinary(data []byte) int {
	return ReadMsgSLoginCommandByBytes(data, this)
}
func (this *SLogoutCommand) ReadBinary(data []byte) int {
	return ReadMsgSLogoutCommandByBytes(data, this)
}
func (this *SSeverStartOKCommand) ReadBinary(data []byte) int {
	return ReadMsgSSeverStartOKCommandByBytes(data, this)
}
func (this *SLoginRetCommand) ReadBinary(data []byte) int {
	return ReadMsgSLoginRetCommandByBytes(data, this)
}
func (this *SStartRelyNotifyCommand) ReadBinary(data []byte) int {
	return ReadMsgSStartRelyNotifyCommandByBytes(data, this)
}
func (this *SStartMyNotifyCommand) ReadBinary(data []byte) int {
	return ReadMsgSStartMyNotifyCommandByBytes(data, this)
}
func (this *SNotifyAllInfo) ReadBinary(data []byte) int {
	return ReadMsgSNotifyAllInfoByBytes(data, this)
}
func (this *SUpdateGatewayUserAnalysis) ReadBinary(data []byte) int {
	return ReadMsgSUpdateGatewayUserAnalysisByBytes(data, this)
}
func (this *SAddNewUserToRedisCommand) ReadBinary(data []byte) int {
	return ReadMsgSAddNewUserToRedisCommandByBytes(data, this)
}
func (this *SGatewayForwardCommand) ReadBinary(data []byte) int {
	return ReadMsgSGatewayForwardCommandByBytes(data, this)
}
func (this *SGatewayForwardBroadcastCommand) ReadBinary(data []byte) int {
	return ReadMsgSGatewayForwardBroadcastCommandByBytes(data, this)
}
func (this *SGatewayForward2HttpCommand) ReadBinary(data []byte) int {
	return ReadMsgSGatewayForward2HttpCommandByBytes(data, this)
}
func (this *SBridgeForward2UserCommand) ReadBinary(data []byte) int {
	return ReadMsgSBridgeForward2UserCommandByBytes(data, this)
}
func (this *SBridgeBroadcast2UserCommand) ReadBinary(data []byte) int {
	return ReadMsgSBridgeBroadcast2UserCommandByBytes(data, this)
}
func (this *SBridgeForward2UserServer) ReadBinary(data []byte) int {
	return ReadMsgSBridgeForward2UserServerByBytes(data, this)
}
func (this *SBridgeBroadcast2GatewayServer) ReadBinary(data []byte) int {
	return ReadMsgSBridgeBroadcast2GatewayServerByBytes(data, this)
}
func (this *SMatchForward2UserServer) ReadBinary(data []byte) int {
	return ReadMsgSMatchForward2UserServerByBytes(data, this)
}
func (this *SRoomForward2UserServer) ReadBinary(data []byte) int {
	return ReadMsgSRoomForward2UserServerByBytes(data, this)
}
func (this *SGatewayBroadcast2UserCommand) ReadBinary(data []byte) int {
	return ReadMsgSGatewayBroadcast2UserCommandByBytes(data, this)
}
func (this *SUserServerSearchFriend) ReadBinary(data []byte) int {
	return ReadMsgSUserServerSearchFriendByBytes(data, this)
}
func (this *SUserServerGMCommand) ReadBinary(data []byte) int {
	return ReadMsgSUserServerGMCommandByBytes(data, this)
}
func (this *SRequestOtherUser) ReadBinary(data []byte) int {
	return ReadMsgSRequestOtherUserByBytes(data, this)
}
func (this *SResponseOtherUser) ReadBinary(data []byte) int {
	return ReadMsgSResponseOtherUserByBytes(data, this)
}
func (this *SBridgeDialGetUserInfo) ReadBinary(data []byte) int {
	return ReadMsgSBridgeDialGetUserInfoByBytes(data, this)
}
func (this *SGatewayWSLoginUser) ReadBinary(data []byte) int {
	return ReadMsgSGatewayWSLoginUserByBytes(data, this)
}
func (this *SGatewayWSOfflineUser) ReadBinary(data []byte) int {
	return ReadMsgSGatewayWSOfflineUserByBytes(data, this)
}
func (this *STemplateMessageKeyWord) ReadBinary(data []byte) int {
	return ReadMsgSTemplateMessageKeyWordByBytes(data, this)
}
func (this *SQSTemplateMessage) ReadBinary(data []byte) int {
	return ReadMsgSQSTemplateMessageByBytes(data, this)
}
func (this *SGatewayChangeAccessToken) ReadBinary(data []byte) int {
	return ReadMsgSGatewayChangeAccessTokenByBytes(data, this)
}
func (this *SUserInfo) ReadBinary(data []byte) int {
	return ReadMsgSUserInfoByBytes(data, this)
}
func (this *SJoinMatch) ReadBinary(data []byte) int {
	return ReadMsgSJoinMatchByBytes(data, this)
}
func (this *STeamInfo) ReadBinary(data []byte) int {
	return ReadMsgSTeamInfoByBytes(data, this)
}
func (this *SRoomInfo) ReadBinary(data []byte) int {
	return ReadMsgSRoomInfoByBytes(data, this)
}
func (this *SUserMatchInfo) ReadBinary(data []byte) int {
	return ReadMsgSUserMatchInfoByBytes(data, this)
}
func (this *SUserQuitMatch) ReadBinary(data []byte) int {
	return ReadMsgSUserQuitMatchByBytes(data, this)
}
func (this *SMatchDone) ReadBinary(data []byte) int {
	return ReadMsgSMatchDoneByBytes(data, this)
}
func (this *SUserQuitRoom) ReadBinary(data []byte) int {
	return ReadMsgSUserQuitRoomByBytes(data, this)
}
func (this *SUserRoomInfo) ReadBinary(data []byte) int {
	return ReadMsgSUserRoomInfoByBytes(data, this)
}
func (this *SNotifyUserOffline) ReadBinary(data []byte) int {
	return ReadMsgSNotifyUserOfflineByBytes(data, this)
}
func (this *SNotifyUserReonline) ReadBinary(data []byte) int {
	return ReadMsgSNotifyUserReonlineByBytes(data, this)
}
func (this *SRPCCheckUserToken) ReadBinary(data []byte) int {
	return ReadMsgSRPCCheckUserTokenByBytes(data, this)
}
func (this *SMatchBroadcast2UserServerCommand) ReadBinary(data []byte) int {
	return ReadMsgSMatchBroadcast2UserServerCommandByBytes(data, this)
}
func (this *NotifyGatewayInfo) ReadBinary(data []byte) int {
	return ReadMsgNotifyGatewayInfoByBytes(data, this)
}
func (this *SuperReqGatewayInfo) ReadBinary(data []byte) int {
	return ReadMsgSuperReqGatewayInfoByBytes(data, this)
}
func (this *NotifyGateUserNums) ReadBinary(data []byte) int {
	return ReadMsgNotifyGateUserNumsByBytes(data, this)
}
func (this *GateReqUserNums) ReadBinary(data []byte) int {
	return ReadMsgGateReqUserNumsByBytes(data, this)
}
func (this *NotifyMatchRoomNums) ReadBinary(data []byte) int {
	return ReadMsgNotifyMatchRoomNumsByBytes(data, this)
}
func (this *MatchReqRoomNums) ReadBinary(data []byte) int {
	return ReadMsgMatchReqRoomNumsByBytes(data, this)
}
func (this *SUserCheckEffective) ReadBinary(data []byte) int {
	return ReadMsgSUserCheckEffectiveByBytes(data, this)
}
func (this *SRedisConfigItem) ReadBinary(data []byte) int {
	return ReadMsgSRedisConfigItemByBytes(data, this)
}
func (this *SRedisConfig) ReadBinary(data []byte) int {
	return ReadMsgSRedisConfigByBytes(data, this)
}
func (this *SRequestServerInfo) ReadBinary(data []byte) int {
	return ReadMsgSRequestServerInfoByBytes(data, this)
}
func (this *SNotifySafelyQuit) ReadBinary(data []byte) int {
	return ReadMsgSNotifySafelyQuitByBytes(data, this)
}
func (this *SAIUserRegister) ReadBinary(data []byte) int {
	return ReadMsgSAIUserRegisterByBytes(data, this)
}
func (this *STaskCmdForward) ReadBinary(data []byte) int {
	return ReadMsgSTaskCmdForwardByBytes(data, this)
}
func (this *SBeginMatch) ReadBinary(data []byte) int {
	return ReadMsgSBeginMatchByBytes(data, this)
}
func MsgIdToString(id uint16) string {
	switch(id ) {
		case SServerInfoID: 
		return SServerInfoName
		case STimeTickCommandID: 
		return STimeTickCommandName
		case STestCommandID: 
		return STestCommandName
		case SLoginCommandID: 
		return SLoginCommandName
		case SLogoutCommandID: 
		return SLogoutCommandName
		case SSeverStartOKCommandID: 
		return SSeverStartOKCommandName
		case SLoginRetCommandID: 
		return SLoginRetCommandName
		case SStartRelyNotifyCommandID: 
		return SStartRelyNotifyCommandName
		case SStartMyNotifyCommandID: 
		return SStartMyNotifyCommandName
		case SNotifyAllInfoID: 
		return SNotifyAllInfoName
		case SUpdateGatewayUserAnalysisID: 
		return SUpdateGatewayUserAnalysisName
		case SAddNewUserToRedisCommandID: 
		return SAddNewUserToRedisCommandName
		case SGatewayForwardCommandID: 
		return SGatewayForwardCommandName
		case SGatewayForwardBroadcastCommandID: 
		return SGatewayForwardBroadcastCommandName
		case SGatewayForward2HttpCommandID: 
		return SGatewayForward2HttpCommandName
		case SBridgeForward2UserCommandID: 
		return SBridgeForward2UserCommandName
		case SBridgeBroadcast2UserCommandID: 
		return SBridgeBroadcast2UserCommandName
		case SBridgeForward2UserServerID: 
		return SBridgeForward2UserServerName
		case SBridgeBroadcast2GatewayServerID: 
		return SBridgeBroadcast2GatewayServerName
		case SMatchForward2UserServerID: 
		return SMatchForward2UserServerName
		case SRoomForward2UserServerID: 
		return SRoomForward2UserServerName
		case SGatewayBroadcast2UserCommandID: 
		return SGatewayBroadcast2UserCommandName
		case SUserServerSearchFriendID: 
		return SUserServerSearchFriendName
		case SUserServerGMCommandID: 
		return SUserServerGMCommandName
		case SRequestOtherUserID: 
		return SRequestOtherUserName
		case SResponseOtherUserID: 
		return SResponseOtherUserName
		case SBridgeDialGetUserInfoID: 
		return SBridgeDialGetUserInfoName
		case SGatewayWSLoginUserID: 
		return SGatewayWSLoginUserName
		case SGatewayWSOfflineUserID: 
		return SGatewayWSOfflineUserName
		case STemplateMessageKeyWordID: 
		return STemplateMessageKeyWordName
		case SQSTemplateMessageID: 
		return SQSTemplateMessageName
		case SGatewayChangeAccessTokenID: 
		return SGatewayChangeAccessTokenName
		case SUserInfoID: 
		return SUserInfoName
		case SJoinMatchID: 
		return SJoinMatchName
		case STeamInfoID: 
		return STeamInfoName
		case SRoomInfoID: 
		return SRoomInfoName
		case SUserMatchInfoID: 
		return SUserMatchInfoName
		case SUserQuitMatchID: 
		return SUserQuitMatchName
		case SMatchDoneID: 
		return SMatchDoneName
		case SUserQuitRoomID: 
		return SUserQuitRoomName
		case SUserRoomInfoID: 
		return SUserRoomInfoName
		case SNotifyUserOfflineID: 
		return SNotifyUserOfflineName
		case SNotifyUserReonlineID: 
		return SNotifyUserReonlineName
		case SRPCCheckUserTokenID: 
		return SRPCCheckUserTokenName
		case SMatchBroadcast2UserServerCommandID: 
		return SMatchBroadcast2UserServerCommandName
		case NotifyGatewayInfoID: 
		return NotifyGatewayInfoName
		case SuperReqGatewayInfoID: 
		return SuperReqGatewayInfoName
		case NotifyGateUserNumsID: 
		return NotifyGateUserNumsName
		case GateReqUserNumsID: 
		return GateReqUserNumsName
		case NotifyMatchRoomNumsID: 
		return NotifyMatchRoomNumsName
		case MatchReqRoomNumsID: 
		return MatchReqRoomNumsName
		case SUserCheckEffectiveID: 
		return SUserCheckEffectiveName
		case SRedisConfigItemID: 
		return SRedisConfigItemName
		case SRedisConfigID: 
		return SRedisConfigName
		case SRequestServerInfoID: 
		return SRequestServerInfoName
		case SNotifySafelyQuitID: 
		return SNotifySafelyQuitName
		case SAIUserRegisterID: 
		return SAIUserRegisterName
		case STaskCmdForwardID: 
		return STaskCmdForwardName
		case SBeginMatchID: 
		return SBeginMatchName
		default:
		return ""
	}
}
func StringToMsgId(msgname string) uint16 {
	switch(msgname ) {
		case SServerInfoName: 
		return SServerInfoID
		case STimeTickCommandName: 
		return STimeTickCommandID
		case STestCommandName: 
		return STestCommandID
		case SLoginCommandName: 
		return SLoginCommandID
		case SLogoutCommandName: 
		return SLogoutCommandID
		case SSeverStartOKCommandName: 
		return SSeverStartOKCommandID
		case SLoginRetCommandName: 
		return SLoginRetCommandID
		case SStartRelyNotifyCommandName: 
		return SStartRelyNotifyCommandID
		case SStartMyNotifyCommandName: 
		return SStartMyNotifyCommandID
		case SNotifyAllInfoName: 
		return SNotifyAllInfoID
		case SUpdateGatewayUserAnalysisName: 
		return SUpdateGatewayUserAnalysisID
		case SAddNewUserToRedisCommandName: 
		return SAddNewUserToRedisCommandID
		case SGatewayForwardCommandName: 
		return SGatewayForwardCommandID
		case SGatewayForwardBroadcastCommandName: 
		return SGatewayForwardBroadcastCommandID
		case SGatewayForward2HttpCommandName: 
		return SGatewayForward2HttpCommandID
		case SBridgeForward2UserCommandName: 
		return SBridgeForward2UserCommandID
		case SBridgeBroadcast2UserCommandName: 
		return SBridgeBroadcast2UserCommandID
		case SBridgeForward2UserServerName: 
		return SBridgeForward2UserServerID
		case SBridgeBroadcast2GatewayServerName: 
		return SBridgeBroadcast2GatewayServerID
		case SMatchForward2UserServerName: 
		return SMatchForward2UserServerID
		case SRoomForward2UserServerName: 
		return SRoomForward2UserServerID
		case SGatewayBroadcast2UserCommandName: 
		return SGatewayBroadcast2UserCommandID
		case SUserServerSearchFriendName: 
		return SUserServerSearchFriendID
		case SUserServerGMCommandName: 
		return SUserServerGMCommandID
		case SRequestOtherUserName: 
		return SRequestOtherUserID
		case SResponseOtherUserName: 
		return SResponseOtherUserID
		case SBridgeDialGetUserInfoName: 
		return SBridgeDialGetUserInfoID
		case SGatewayWSLoginUserName: 
		return SGatewayWSLoginUserID
		case SGatewayWSOfflineUserName: 
		return SGatewayWSOfflineUserID
		case STemplateMessageKeyWordName: 
		return STemplateMessageKeyWordID
		case SQSTemplateMessageName: 
		return SQSTemplateMessageID
		case SGatewayChangeAccessTokenName: 
		return SGatewayChangeAccessTokenID
		case SUserInfoName: 
		return SUserInfoID
		case SJoinMatchName: 
		return SJoinMatchID
		case STeamInfoName: 
		return STeamInfoID
		case SRoomInfoName: 
		return SRoomInfoID
		case SUserMatchInfoName: 
		return SUserMatchInfoID
		case SUserQuitMatchName: 
		return SUserQuitMatchID
		case SMatchDoneName: 
		return SMatchDoneID
		case SUserQuitRoomName: 
		return SUserQuitRoomID
		case SUserRoomInfoName: 
		return SUserRoomInfoID
		case SNotifyUserOfflineName: 
		return SNotifyUserOfflineID
		case SNotifyUserReonlineName: 
		return SNotifyUserReonlineID
		case SRPCCheckUserTokenName: 
		return SRPCCheckUserTokenID
		case SMatchBroadcast2UserServerCommandName: 
		return SMatchBroadcast2UserServerCommandID
		case NotifyGatewayInfoName: 
		return NotifyGatewayInfoID
		case SuperReqGatewayInfoName: 
		return SuperReqGatewayInfoID
		case NotifyGateUserNumsName: 
		return NotifyGateUserNumsID
		case GateReqUserNumsName: 
		return GateReqUserNumsID
		case NotifyMatchRoomNumsName: 
		return NotifyMatchRoomNumsID
		case MatchReqRoomNumsName: 
		return MatchReqRoomNumsID
		case SUserCheckEffectiveName: 
		return SUserCheckEffectiveID
		case SRedisConfigItemName: 
		return SRedisConfigItemID
		case SRedisConfigName: 
		return SRedisConfigID
		case SRequestServerInfoName: 
		return SRequestServerInfoID
		case SNotifySafelyQuitName: 
		return SNotifySafelyQuitID
		case SAIUserRegisterName: 
		return SAIUserRegisterID
		case STaskCmdForwardName: 
		return STaskCmdForwardID
		case SBeginMatchName: 
		return SBeginMatchID
		default:
		return 0
	}
}
func MsgIdToType(id uint16) rune {
	switch(id ) {
		case SServerInfoID: 
		return rune('S')
		case STimeTickCommandID: 
		return rune('S')
		case STestCommandID: 
		return rune('S')
		case SLoginCommandID: 
		return rune('S')
		case SLogoutCommandID: 
		return rune('S')
		case SSeverStartOKCommandID: 
		return rune('S')
		case SLoginRetCommandID: 
		return rune('S')
		case SStartRelyNotifyCommandID: 
		return rune('S')
		case SStartMyNotifyCommandID: 
		return rune('S')
		case SNotifyAllInfoID: 
		return rune('S')
		case SUpdateGatewayUserAnalysisID: 
		return rune('S')
		case SAddNewUserToRedisCommandID: 
		return rune('S')
		case SGatewayForwardCommandID: 
		return rune('S')
		case SGatewayForwardBroadcastCommandID: 
		return rune('S')
		case SGatewayForward2HttpCommandID: 
		return rune('S')
		case SBridgeForward2UserCommandID: 
		return rune('S')
		case SBridgeBroadcast2UserCommandID: 
		return rune('S')
		case SBridgeForward2UserServerID: 
		return rune('S')
		case SBridgeBroadcast2GatewayServerID: 
		return rune('S')
		case SMatchForward2UserServerID: 
		return rune('S')
		case SRoomForward2UserServerID: 
		return rune('S')
		case SGatewayBroadcast2UserCommandID: 
		return rune('S')
		case SUserServerSearchFriendID: 
		return rune('S')
		case SUserServerGMCommandID: 
		return rune('S')
		case SRequestOtherUserID: 
		return rune('S')
		case SResponseOtherUserID: 
		return rune('S')
		case SBridgeDialGetUserInfoID: 
		return rune('S')
		case SGatewayWSLoginUserID: 
		return rune('S')
		case SGatewayWSOfflineUserID: 
		return rune('S')
		case STemplateMessageKeyWordID: 
		return rune('S')
		case SQSTemplateMessageID: 
		return rune('S')
		case SGatewayChangeAccessTokenID: 
		return rune('S')
		case SUserInfoID: 
		return rune('S')
		case SJoinMatchID: 
		return rune('S')
		case STeamInfoID: 
		return rune('S')
		case SRoomInfoID: 
		return rune('S')
		case SUserMatchInfoID: 
		return rune('S')
		case SUserQuitMatchID: 
		return rune('S')
		case SMatchDoneID: 
		return rune('S')
		case SUserQuitRoomID: 
		return rune('S')
		case SUserRoomInfoID: 
		return rune('S')
		case SNotifyUserOfflineID: 
		return rune('S')
		case SNotifyUserReonlineID: 
		return rune('S')
		case SRPCCheckUserTokenID: 
		return rune('S')
		case SMatchBroadcast2UserServerCommandID: 
		return rune('S')
		case NotifyGatewayInfoID: 
		return rune('N')
		case SuperReqGatewayInfoID: 
		return rune('S')
		case NotifyGateUserNumsID: 
		return rune('N')
		case GateReqUserNumsID: 
		return rune('G')
		case NotifyMatchRoomNumsID: 
		return rune('N')
		case MatchReqRoomNumsID: 
		return rune('M')
		case SUserCheckEffectiveID: 
		return rune('S')
		case SRedisConfigItemID: 
		return rune('S')
		case SRedisConfigID: 
		return rune('S')
		case SRequestServerInfoID: 
		return rune('S')
		case SNotifySafelyQuitID: 
		return rune('S')
		case SAIUserRegisterID: 
		return rune('S')
		case STaskCmdForwardID: 
		return rune('S')
		case SBeginMatchID: 
		return rune('S')
		default:
		return rune(0)
	}
}
func (this *SServerInfo) GetMsgId() uint16 {
	return SServerInfoID
}
func (this *STimeTickCommand) GetMsgId() uint16 {
	return STimeTickCommandID
}
func (this *STestCommand) GetMsgId() uint16 {
	return STestCommandID
}
func (this *SLoginCommand) GetMsgId() uint16 {
	return SLoginCommandID
}
func (this *SLogoutCommand) GetMsgId() uint16 {
	return SLogoutCommandID
}
func (this *SSeverStartOKCommand) GetMsgId() uint16 {
	return SSeverStartOKCommandID
}
func (this *SLoginRetCommand) GetMsgId() uint16 {
	return SLoginRetCommandID
}
func (this *SStartRelyNotifyCommand) GetMsgId() uint16 {
	return SStartRelyNotifyCommandID
}
func (this *SStartMyNotifyCommand) GetMsgId() uint16 {
	return SStartMyNotifyCommandID
}
func (this *SNotifyAllInfo) GetMsgId() uint16 {
	return SNotifyAllInfoID
}
func (this *SUpdateGatewayUserAnalysis) GetMsgId() uint16 {
	return SUpdateGatewayUserAnalysisID
}
func (this *SAddNewUserToRedisCommand) GetMsgId() uint16 {
	return SAddNewUserToRedisCommandID
}
func (this *SGatewayForwardCommand) GetMsgId() uint16 {
	return SGatewayForwardCommandID
}
func (this *SGatewayForwardBroadcastCommand) GetMsgId() uint16 {
	return SGatewayForwardBroadcastCommandID
}
func (this *SGatewayForward2HttpCommand) GetMsgId() uint16 {
	return SGatewayForward2HttpCommandID
}
func (this *SBridgeForward2UserCommand) GetMsgId() uint16 {
	return SBridgeForward2UserCommandID
}
func (this *SBridgeBroadcast2UserCommand) GetMsgId() uint16 {
	return SBridgeBroadcast2UserCommandID
}
func (this *SBridgeForward2UserServer) GetMsgId() uint16 {
	return SBridgeForward2UserServerID
}
func (this *SBridgeBroadcast2GatewayServer) GetMsgId() uint16 {
	return SBridgeBroadcast2GatewayServerID
}
func (this *SMatchForward2UserServer) GetMsgId() uint16 {
	return SMatchForward2UserServerID
}
func (this *SRoomForward2UserServer) GetMsgId() uint16 {
	return SRoomForward2UserServerID
}
func (this *SGatewayBroadcast2UserCommand) GetMsgId() uint16 {
	return SGatewayBroadcast2UserCommandID
}
func (this *SUserServerSearchFriend) GetMsgId() uint16 {
	return SUserServerSearchFriendID
}
func (this *SUserServerGMCommand) GetMsgId() uint16 {
	return SUserServerGMCommandID
}
func (this *SRequestOtherUser) GetMsgId() uint16 {
	return SRequestOtherUserID
}
func (this *SResponseOtherUser) GetMsgId() uint16 {
	return SResponseOtherUserID
}
func (this *SBridgeDialGetUserInfo) GetMsgId() uint16 {
	return SBridgeDialGetUserInfoID
}
func (this *SGatewayWSLoginUser) GetMsgId() uint16 {
	return SGatewayWSLoginUserID
}
func (this *SGatewayWSOfflineUser) GetMsgId() uint16 {
	return SGatewayWSOfflineUserID
}
func (this *STemplateMessageKeyWord) GetMsgId() uint16 {
	return STemplateMessageKeyWordID
}
func (this *SQSTemplateMessage) GetMsgId() uint16 {
	return SQSTemplateMessageID
}
func (this *SGatewayChangeAccessToken) GetMsgId() uint16 {
	return SGatewayChangeAccessTokenID
}
func (this *SUserInfo) GetMsgId() uint16 {
	return SUserInfoID
}
func (this *SJoinMatch) GetMsgId() uint16 {
	return SJoinMatchID
}
func (this *STeamInfo) GetMsgId() uint16 {
	return STeamInfoID
}
func (this *SRoomInfo) GetMsgId() uint16 {
	return SRoomInfoID
}
func (this *SUserMatchInfo) GetMsgId() uint16 {
	return SUserMatchInfoID
}
func (this *SUserQuitMatch) GetMsgId() uint16 {
	return SUserQuitMatchID
}
func (this *SMatchDone) GetMsgId() uint16 {
	return SMatchDoneID
}
func (this *SUserQuitRoom) GetMsgId() uint16 {
	return SUserQuitRoomID
}
func (this *SUserRoomInfo) GetMsgId() uint16 {
	return SUserRoomInfoID
}
func (this *SNotifyUserOffline) GetMsgId() uint16 {
	return SNotifyUserOfflineID
}
func (this *SNotifyUserReonline) GetMsgId() uint16 {
	return SNotifyUserReonlineID
}
func (this *SRPCCheckUserToken) GetMsgId() uint16 {
	return SRPCCheckUserTokenID
}
func (this *SMatchBroadcast2UserServerCommand) GetMsgId() uint16 {
	return SMatchBroadcast2UserServerCommandID
}
func (this *NotifyGatewayInfo) GetMsgId() uint16 {
	return NotifyGatewayInfoID
}
func (this *SuperReqGatewayInfo) GetMsgId() uint16 {
	return SuperReqGatewayInfoID
}
func (this *NotifyGateUserNums) GetMsgId() uint16 {
	return NotifyGateUserNumsID
}
func (this *GateReqUserNums) GetMsgId() uint16 {
	return GateReqUserNumsID
}
func (this *NotifyMatchRoomNums) GetMsgId() uint16 {
	return NotifyMatchRoomNumsID
}
func (this *MatchReqRoomNums) GetMsgId() uint16 {
	return MatchReqRoomNumsID
}
func (this *SUserCheckEffective) GetMsgId() uint16 {
	return SUserCheckEffectiveID
}
func (this *SRedisConfigItem) GetMsgId() uint16 {
	return SRedisConfigItemID
}
func (this *SRedisConfig) GetMsgId() uint16 {
	return SRedisConfigID
}
func (this *SRequestServerInfo) GetMsgId() uint16 {
	return SRequestServerInfoID
}
func (this *SNotifySafelyQuit) GetMsgId() uint16 {
	return SNotifySafelyQuitID
}
func (this *SAIUserRegister) GetMsgId() uint16 {
	return SAIUserRegisterID
}
func (this *STaskCmdForward) GetMsgId() uint16 {
	return STaskCmdForwardID
}
func (this *SBeginMatch) GetMsgId() uint16 {
	return SBeginMatchID
}
func (this *SServerInfo) GetMsgName() string {
	return SServerInfoName
}
func (this *STimeTickCommand) GetMsgName() string {
	return STimeTickCommandName
}
func (this *STestCommand) GetMsgName() string {
	return STestCommandName
}
func (this *SLoginCommand) GetMsgName() string {
	return SLoginCommandName
}
func (this *SLogoutCommand) GetMsgName() string {
	return SLogoutCommandName
}
func (this *SSeverStartOKCommand) GetMsgName() string {
	return SSeverStartOKCommandName
}
func (this *SLoginRetCommand) GetMsgName() string {
	return SLoginRetCommandName
}
func (this *SStartRelyNotifyCommand) GetMsgName() string {
	return SStartRelyNotifyCommandName
}
func (this *SStartMyNotifyCommand) GetMsgName() string {
	return SStartMyNotifyCommandName
}
func (this *SNotifyAllInfo) GetMsgName() string {
	return SNotifyAllInfoName
}
func (this *SUpdateGatewayUserAnalysis) GetMsgName() string {
	return SUpdateGatewayUserAnalysisName
}
func (this *SAddNewUserToRedisCommand) GetMsgName() string {
	return SAddNewUserToRedisCommandName
}
func (this *SGatewayForwardCommand) GetMsgName() string {
	return SGatewayForwardCommandName
}
func (this *SGatewayForwardBroadcastCommand) GetMsgName() string {
	return SGatewayForwardBroadcastCommandName
}
func (this *SGatewayForward2HttpCommand) GetMsgName() string {
	return SGatewayForward2HttpCommandName
}
func (this *SBridgeForward2UserCommand) GetMsgName() string {
	return SBridgeForward2UserCommandName
}
func (this *SBridgeBroadcast2UserCommand) GetMsgName() string {
	return SBridgeBroadcast2UserCommandName
}
func (this *SBridgeForward2UserServer) GetMsgName() string {
	return SBridgeForward2UserServerName
}
func (this *SBridgeBroadcast2GatewayServer) GetMsgName() string {
	return SBridgeBroadcast2GatewayServerName
}
func (this *SMatchForward2UserServer) GetMsgName() string {
	return SMatchForward2UserServerName
}
func (this *SRoomForward2UserServer) GetMsgName() string {
	return SRoomForward2UserServerName
}
func (this *SGatewayBroadcast2UserCommand) GetMsgName() string {
	return SGatewayBroadcast2UserCommandName
}
func (this *SUserServerSearchFriend) GetMsgName() string {
	return SUserServerSearchFriendName
}
func (this *SUserServerGMCommand) GetMsgName() string {
	return SUserServerGMCommandName
}
func (this *SRequestOtherUser) GetMsgName() string {
	return SRequestOtherUserName
}
func (this *SResponseOtherUser) GetMsgName() string {
	return SResponseOtherUserName
}
func (this *SBridgeDialGetUserInfo) GetMsgName() string {
	return SBridgeDialGetUserInfoName
}
func (this *SGatewayWSLoginUser) GetMsgName() string {
	return SGatewayWSLoginUserName
}
func (this *SGatewayWSOfflineUser) GetMsgName() string {
	return SGatewayWSOfflineUserName
}
func (this *STemplateMessageKeyWord) GetMsgName() string {
	return STemplateMessageKeyWordName
}
func (this *SQSTemplateMessage) GetMsgName() string {
	return SQSTemplateMessageName
}
func (this *SGatewayChangeAccessToken) GetMsgName() string {
	return SGatewayChangeAccessTokenName
}
func (this *SUserInfo) GetMsgName() string {
	return SUserInfoName
}
func (this *SJoinMatch) GetMsgName() string {
	return SJoinMatchName
}
func (this *STeamInfo) GetMsgName() string {
	return STeamInfoName
}
func (this *SRoomInfo) GetMsgName() string {
	return SRoomInfoName
}
func (this *SUserMatchInfo) GetMsgName() string {
	return SUserMatchInfoName
}
func (this *SUserQuitMatch) GetMsgName() string {
	return SUserQuitMatchName
}
func (this *SMatchDone) GetMsgName() string {
	return SMatchDoneName
}
func (this *SUserQuitRoom) GetMsgName() string {
	return SUserQuitRoomName
}
func (this *SUserRoomInfo) GetMsgName() string {
	return SUserRoomInfoName
}
func (this *SNotifyUserOffline) GetMsgName() string {
	return SNotifyUserOfflineName
}
func (this *SNotifyUserReonline) GetMsgName() string {
	return SNotifyUserReonlineName
}
func (this *SRPCCheckUserToken) GetMsgName() string {
	return SRPCCheckUserTokenName
}
func (this *SMatchBroadcast2UserServerCommand) GetMsgName() string {
	return SMatchBroadcast2UserServerCommandName
}
func (this *NotifyGatewayInfo) GetMsgName() string {
	return NotifyGatewayInfoName
}
func (this *SuperReqGatewayInfo) GetMsgName() string {
	return SuperReqGatewayInfoName
}
func (this *NotifyGateUserNums) GetMsgName() string {
	return NotifyGateUserNumsName
}
func (this *GateReqUserNums) GetMsgName() string {
	return GateReqUserNumsName
}
func (this *NotifyMatchRoomNums) GetMsgName() string {
	return NotifyMatchRoomNumsName
}
func (this *MatchReqRoomNums) GetMsgName() string {
	return MatchReqRoomNumsName
}
func (this *SUserCheckEffective) GetMsgName() string {
	return SUserCheckEffectiveName
}
func (this *SRedisConfigItem) GetMsgName() string {
	return SRedisConfigItemName
}
func (this *SRedisConfig) GetMsgName() string {
	return SRedisConfigName
}
func (this *SRequestServerInfo) GetMsgName() string {
	return SRequestServerInfoName
}
func (this *SNotifySafelyQuit) GetMsgName() string {
	return SNotifySafelyQuitName
}
func (this *SAIUserRegister) GetMsgName() string {
	return SAIUserRegisterName
}
func (this *STaskCmdForward) GetMsgName() string {
	return STaskCmdForwardName
}
func (this *SBeginMatch) GetMsgName() string {
	return SBeginMatchName
}
func (this *SServerInfo) GetSize() int {
	return GetSizeSServerInfo(this)
}
func (this *STimeTickCommand) GetSize() int {
	return GetSizeSTimeTickCommand(this)
}
func (this *STestCommand) GetSize() int {
	return GetSizeSTestCommand(this)
}
func (this *SLoginCommand) GetSize() int {
	return GetSizeSLoginCommand(this)
}
func (this *SLogoutCommand) GetSize() int {
	return GetSizeSLogoutCommand(this)
}
func (this *SSeverStartOKCommand) GetSize() int {
	return GetSizeSSeverStartOKCommand(this)
}
func (this *SLoginRetCommand) GetSize() int {
	return GetSizeSLoginRetCommand(this)
}
func (this *SStartRelyNotifyCommand) GetSize() int {
	return GetSizeSStartRelyNotifyCommand(this)
}
func (this *SStartMyNotifyCommand) GetSize() int {
	return GetSizeSStartMyNotifyCommand(this)
}
func (this *SNotifyAllInfo) GetSize() int {
	return GetSizeSNotifyAllInfo(this)
}
func (this *SUpdateGatewayUserAnalysis) GetSize() int {
	return GetSizeSUpdateGatewayUserAnalysis(this)
}
func (this *SAddNewUserToRedisCommand) GetSize() int {
	return GetSizeSAddNewUserToRedisCommand(this)
}
func (this *SGatewayForwardCommand) GetSize() int {
	return GetSizeSGatewayForwardCommand(this)
}
func (this *SGatewayForwardBroadcastCommand) GetSize() int {
	return GetSizeSGatewayForwardBroadcastCommand(this)
}
func (this *SGatewayForward2HttpCommand) GetSize() int {
	return GetSizeSGatewayForward2HttpCommand(this)
}
func (this *SBridgeForward2UserCommand) GetSize() int {
	return GetSizeSBridgeForward2UserCommand(this)
}
func (this *SBridgeBroadcast2UserCommand) GetSize() int {
	return GetSizeSBridgeBroadcast2UserCommand(this)
}
func (this *SBridgeForward2UserServer) GetSize() int {
	return GetSizeSBridgeForward2UserServer(this)
}
func (this *SBridgeBroadcast2GatewayServer) GetSize() int {
	return GetSizeSBridgeBroadcast2GatewayServer(this)
}
func (this *SMatchForward2UserServer) GetSize() int {
	return GetSizeSMatchForward2UserServer(this)
}
func (this *SRoomForward2UserServer) GetSize() int {
	return GetSizeSRoomForward2UserServer(this)
}
func (this *SGatewayBroadcast2UserCommand) GetSize() int {
	return GetSizeSGatewayBroadcast2UserCommand(this)
}
func (this *SUserServerSearchFriend) GetSize() int {
	return GetSizeSUserServerSearchFriend(this)
}
func (this *SUserServerGMCommand) GetSize() int {
	return GetSizeSUserServerGMCommand(this)
}
func (this *SRequestOtherUser) GetSize() int {
	return GetSizeSRequestOtherUser(this)
}
func (this *SResponseOtherUser) GetSize() int {
	return GetSizeSResponseOtherUser(this)
}
func (this *SBridgeDialGetUserInfo) GetSize() int {
	return GetSizeSBridgeDialGetUserInfo(this)
}
func (this *SGatewayWSLoginUser) GetSize() int {
	return GetSizeSGatewayWSLoginUser(this)
}
func (this *SGatewayWSOfflineUser) GetSize() int {
	return GetSizeSGatewayWSOfflineUser(this)
}
func (this *STemplateMessageKeyWord) GetSize() int {
	return GetSizeSTemplateMessageKeyWord(this)
}
func (this *SQSTemplateMessage) GetSize() int {
	return GetSizeSQSTemplateMessage(this)
}
func (this *SGatewayChangeAccessToken) GetSize() int {
	return GetSizeSGatewayChangeAccessToken(this)
}
func (this *SUserInfo) GetSize() int {
	return GetSizeSUserInfo(this)
}
func (this *SJoinMatch) GetSize() int {
	return GetSizeSJoinMatch(this)
}
func (this *STeamInfo) GetSize() int {
	return GetSizeSTeamInfo(this)
}
func (this *SRoomInfo) GetSize() int {
	return GetSizeSRoomInfo(this)
}
func (this *SUserMatchInfo) GetSize() int {
	return GetSizeSUserMatchInfo(this)
}
func (this *SUserQuitMatch) GetSize() int {
	return GetSizeSUserQuitMatch(this)
}
func (this *SMatchDone) GetSize() int {
	return GetSizeSMatchDone(this)
}
func (this *SUserQuitRoom) GetSize() int {
	return GetSizeSUserQuitRoom(this)
}
func (this *SUserRoomInfo) GetSize() int {
	return GetSizeSUserRoomInfo(this)
}
func (this *SNotifyUserOffline) GetSize() int {
	return GetSizeSNotifyUserOffline(this)
}
func (this *SNotifyUserReonline) GetSize() int {
	return GetSizeSNotifyUserReonline(this)
}
func (this *SRPCCheckUserToken) GetSize() int {
	return GetSizeSRPCCheckUserToken(this)
}
func (this *SMatchBroadcast2UserServerCommand) GetSize() int {
	return GetSizeSMatchBroadcast2UserServerCommand(this)
}
func (this *NotifyGatewayInfo) GetSize() int {
	return GetSizeNotifyGatewayInfo(this)
}
func (this *SuperReqGatewayInfo) GetSize() int {
	return GetSizeSuperReqGatewayInfo(this)
}
func (this *NotifyGateUserNums) GetSize() int {
	return GetSizeNotifyGateUserNums(this)
}
func (this *GateReqUserNums) GetSize() int {
	return GetSizeGateReqUserNums(this)
}
func (this *NotifyMatchRoomNums) GetSize() int {
	return GetSizeNotifyMatchRoomNums(this)
}
func (this *MatchReqRoomNums) GetSize() int {
	return GetSizeMatchReqRoomNums(this)
}
func (this *SUserCheckEffective) GetSize() int {
	return GetSizeSUserCheckEffective(this)
}
func (this *SRedisConfigItem) GetSize() int {
	return GetSizeSRedisConfigItem(this)
}
func (this *SRedisConfig) GetSize() int {
	return GetSizeSRedisConfig(this)
}
func (this *SRequestServerInfo) GetSize() int {
	return GetSizeSRequestServerInfo(this)
}
func (this *SNotifySafelyQuit) GetSize() int {
	return GetSizeSNotifySafelyQuit(this)
}
func (this *SAIUserRegister) GetSize() int {
	return GetSizeSAIUserRegister(this)
}
func (this *STaskCmdForward) GetSize() int {
	return GetSizeSTaskCmdForward(this)
}
func (this *SBeginMatch) GetSize() int {
	return GetSizeSBeginMatch(this)
}
func (this *SServerInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *STimeTickCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *STestCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SLoginCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SLogoutCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SSeverStartOKCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SLoginRetCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SStartRelyNotifyCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SStartMyNotifyCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SNotifyAllInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUpdateGatewayUserAnalysis) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SAddNewUserToRedisCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SGatewayForwardCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SGatewayForwardBroadcastCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SGatewayForward2HttpCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SBridgeForward2UserCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SBridgeBroadcast2UserCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SBridgeForward2UserServer) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SBridgeBroadcast2GatewayServer) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SMatchForward2UserServer) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SRoomForward2UserServer) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SGatewayBroadcast2UserCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUserServerSearchFriend) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUserServerGMCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SRequestOtherUser) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SResponseOtherUser) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SBridgeDialGetUserInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SGatewayWSLoginUser) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SGatewayWSOfflineUser) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *STemplateMessageKeyWord) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SQSTemplateMessage) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SGatewayChangeAccessToken) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUserInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SJoinMatch) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *STeamInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SRoomInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUserMatchInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUserQuitMatch) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SMatchDone) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUserQuitRoom) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUserRoomInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SNotifyUserOffline) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SNotifyUserReonline) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SRPCCheckUserToken) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SMatchBroadcast2UserServerCommand) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *NotifyGatewayInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SuperReqGatewayInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *NotifyGateUserNums) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *GateReqUserNums) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *NotifyMatchRoomNums) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *MatchReqRoomNums) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUserCheckEffective) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SRedisConfigItem) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SRedisConfig) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SRequestServerInfo) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SNotifySafelyQuit) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SAIUserRegister) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *STaskCmdForward) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SBeginMatch) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func readBinaryString(data []byte) string {
	strfunclen := binary.BigEndian.Uint16(data[:2])
	if int(strfunclen) + 2 > len(data ) {
		return ""
	}
	return string(data[2:2+strfunclen])
}
func writeBinaryString(data []byte,obj string) int {
	objlen := len(obj)
	binary.BigEndian.PutUint16(data[:2],uint16(objlen))
	copy(data[2:2+objlen], obj)
	return 2+objlen
}
func bool2int(value bool) int {
	if value {
		return 1
	}
	return 0
}
func readBinaryInt64(data []byte) int64 {
	// 大端模式
	num := int64(0)
	num |= int64(data[7]) << 0
	num |= int64(data[6]) << 8
	num |= int64(data[5]) << 16
	num |= int64(data[4]) << 24
	num |= int64(data[3]) << 32
	num |= int64(data[2]) << 40
	num |= int64(data[1]) << 48
	num |= int64(data[0]) << 56
	return num
}
func writeBinaryInt64(data []byte, num int64 ) {
	// 大端模式
	data[7] = byte((num >> 0) & 0xff)
	data[6] = byte((num >> 8) & 0xff)
	data[5] = byte((num >> 16) & 0xff)
	data[4] = byte((num >> 24) & 0xff)
	data[3] = byte((num >> 32) & 0xff)
	data[2] = byte((num >> 40) & 0xff)
	data[1] = byte((num >> 48) & 0xff)
	data[0] = byte((num >> 56) & 0xff)
}
func readBinaryInt32(data []byte) int32 {
	// 大端模式
	num := int32(0)
	num |= int32(data[3]) << 0
	num |= int32(data[2]) << 8
	num |= int32(data[1]) << 16
	num |= int32(data[0]) << 24
	return num
}
func writeBinaryInt32(data []byte, num int32 ) {
	// 大端模式
	data[3] = byte((num >> 0) & 0xff)
	data[2] = byte((num >> 8) & 0xff)
	data[1] = byte((num >> 16) & 0xff)
	data[0] = byte((num >> 24) & 0xff)
}
func readBinaryInt(data []byte) int {
	return int(readBinaryInt32(data))
}
func writeBinaryInt(data []byte, num int ) {
	writeBinaryInt32(data,int32(num))
}
func readBinaryInt16(data []byte) int16 {
	// 大端模式
	num := int16(0)
	num |= int16(data[1]) << 0
	num |= int16(data[0]) << 8
	return num
}
func writeBinaryInt16(data []byte, num int16 ) {
	// 大端模式
	data[1] = byte((num >> 0) & 0xff)
	data[0] = byte((num >> 8) & 0xff)
}
func readBinaryInt8(data []byte) int8 {
	// 大端模式
	num := int8(0)
	num |= int8(data[0]) << 0
	return num
}
func writeBinaryInt8(data []byte, num int8 ) {
	// 大端模式
	data[0] = byte(num)
}
func readBinaryBool(data []byte) bool {
	// 大端模式
	num := int8(0)
	num |= int8(data[0]) << 0
	return num>0
}
func writeBinaryBool(data []byte, num bool ) {
	// 大端模式
	if num == true {
		data[0] = byte(1)
	} else {
		data[0] = byte(0)
	}
}
func readBinaryUint8(data []byte) uint8 {
	return uint8(data[0])
}
func writeBinaryUint8(data []byte, num uint8 ) {
	data[0] = byte(num)
}
func writeBinaryFloat32(data []byte, num float32 ) {
	bits := math.Float32bits(num)
	binary.BigEndian.PutUint32(data,bits)
}
func readBinaryFloat32(data []byte) float32 {
	bits := binary.BigEndian.Uint32(data)
	return math.Float32frombits(bits)
}
func writeBinaryFloat64(data []byte, num float64 ) {
	bits := math.Float64bits(num)
	binary.BigEndian.PutUint64(data,bits)
}
func readBinaryFloat64(data []byte) float64 {
	bits := binary.BigEndian.Uint64(data)
	return math.Float64frombits(bits)
}
func ReadMsgSServerInfoByBytes(indata []byte, obj *SServerInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Serverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.Servertype = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 2 + len(obj.Servername) > data__len{
		return endpos
	}
	obj.Servername = readBinaryString(data[offset:])
	offset += 2 + len(obj.Servername)
	if offset + 2 + len(obj.Serverip) > data__len{
		return endpos
	}
	obj.Serverip = readBinaryString(data[offset:])
	offset += 2 + len(obj.Serverip)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Serverport = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 2 + len(obj.Extip) > data__len{
		return endpos
	}
	obj.Extip = readBinaryString(data[offset:])
	offset += 2 + len(obj.Extip)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Httpport = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.Httpsport = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.Rpcport = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.Tcpport = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.ClientTcpport = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.ServerNumber = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.Version = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSServerInfoByObj(data []byte, obj *SServerInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Serverid)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Servertype)
	offset+=4
	writeBinaryString(data[offset:],obj.Servername)
	offset += 2 + len(obj.Servername)
	writeBinaryString(data[offset:],obj.Serverip)
	offset += 2 + len(obj.Serverip)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Serverport)
	offset+=4
	writeBinaryString(data[offset:],obj.Extip)
	offset += 2 + len(obj.Extip)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Httpport)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Httpsport)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Rpcport)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Tcpport)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.ClientTcpport)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.ServerNumber)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Version)
	offset+=8
	return offset
}
func GetSizeSServerInfo(obj *SServerInfo) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + 4 + 2 + len(obj.Servername) + 2 + len(obj.Serverip) + 
	4 + 2 + len(obj.Extip) + 4 + 4 + 4 + 
	4 + 4 + 4 + 8
}
func ReadMsgSTimeTickCommandByBytes(indata []byte, obj *STimeTickCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Testno = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgSTimeTickCommandByObj(data []byte, obj *STimeTickCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Testno)
	offset+=4
	return offset
}
func GetSizeSTimeTickCommand(obj *STimeTickCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 4
}
func ReadMsgSTestCommandByBytes(indata []byte, obj *STestCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Testno = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 2 + len(obj.Testttring) > data__len{
		return endpos
	}
	obj.Testttring = readBinaryString(data[offset:])
	offset += 2 + len(obj.Testttring)
	return endpos
}
func WriteMsgSTestCommandByObj(data []byte, obj *STestCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Testno)
	offset+=4
	writeBinaryString(data[offset:],obj.Testttring)
	offset += 2 + len(obj.Testttring)
	return offset
}
func GetSizeSTestCommand(obj *STestCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + 2 + len(obj.Testttring)
}
func ReadMsgSLoginCommandByBytes(indata []byte, obj *SLoginCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Serverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.Servertype = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 2 + len(obj.Serverip) > data__len{
		return endpos
	}
	obj.Serverip = readBinaryString(data[offset:])
	offset += 2 + len(obj.Serverip)
	if offset + 2 + len(obj.Servername) > data__len{
		return endpos
	}
	obj.Servername = readBinaryString(data[offset:])
	offset += 2 + len(obj.Servername)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Serverport = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.ServerNumber = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.Version = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSLoginCommandByObj(data []byte, obj *SLoginCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Serverid)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Servertype)
	offset+=4
	writeBinaryString(data[offset:],obj.Serverip)
	offset += 2 + len(obj.Serverip)
	writeBinaryString(data[offset:],obj.Servername)
	offset += 2 + len(obj.Servername)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Serverport)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.ServerNumber)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Version)
	offset+=8
	return offset
}
func GetSizeSLoginCommand(obj *SLoginCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + 4 + 2 + len(obj.Serverip) + 2 + len(obj.Servername) + 
	4 + 4 + 8
}
func ReadMsgSLogoutCommandByBytes(indata []byte, obj *SLogoutCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	return endpos
}
func WriteMsgSLogoutCommandByObj(data []byte, obj *SLogoutCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	return offset
}
func GetSizeSLogoutCommand(obj *SLogoutCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 0
}
func ReadMsgSSeverStartOKCommandByBytes(indata []byte, obj *SSeverStartOKCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Serverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgSSeverStartOKCommandByObj(data []byte, obj *SSeverStartOKCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Serverid)
	offset+=4
	return offset
}
func GetSizeSSeverStartOKCommand(obj *SSeverStartOKCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 4
}
func ReadMsgSLoginRetCommandByBytes(indata []byte, obj *SLoginRetCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Loginfailed = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + obj.Clientinfo.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSServerInfoByBytes(data[offset:], &obj.Clientinfo)
	if offset + obj.Taskinfo.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSServerInfoByBytes(data[offset:], &obj.Taskinfo)
	if offset + obj.Redisinfo.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSRedisConfigByBytes(data[offset:], &obj.Redisinfo)
	return endpos
}
func WriteMsgSLoginRetCommandByObj(data []byte, obj *SLoginRetCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Loginfailed)
	offset+=4
	offset += WriteMsgSServerInfoByObj(data[offset:], &obj.Clientinfo)
	offset += WriteMsgSServerInfoByObj(data[offset:], &obj.Taskinfo)
	offset += WriteMsgSRedisConfigByObj(data[offset:], &obj.Redisinfo)
	return offset
}
func GetSizeSLoginRetCommand(obj *SLoginRetCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + obj.Clientinfo.GetSize() + obj.Taskinfo.GetSize() + obj.Redisinfo.GetSize()
}
func ReadMsgSStartRelyNotifyCommandByBytes(indata []byte, obj *SStartRelyNotifyCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	Serverinfos_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Serverinfos_slen := 0
	Serverinfos_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Serverinfos_slen = int(Serverinfos_slent)
	obj.Serverinfos = make([]SServerInfo,Serverinfos_slen)
	i1i := 0
	for Serverinfos_slen > i1i {
		offset += ReadMsgSServerInfoByBytes(data[offset:],&obj.Serverinfos[i1i])
		i1i++
	}
	return endpos
}
func WriteMsgSStartRelyNotifyCommandByObj(data []byte, obj *SStartRelyNotifyCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Serverinfos)))
	offset += 2
	i1i := 0
	Serverinfos_slen := len(obj.Serverinfos)
	for Serverinfos_slen > i1i {
		offset += WriteMsgSServerInfoByObj(data[offset:],&obj.Serverinfos[i1i])
		i1i++
	}
	return offset
}
func GetSizeSStartRelyNotifyCommand(obj *SStartRelyNotifyCommand) int {
	if obj == nil {
		return 2
	}
	sizerelySServerInfo1 := func()int{
		resnum := 0
		i1i := 0
		Serverinfos_slen := len(obj.Serverinfos)
		for Serverinfos_slen > i1i {
			resnum += obj.Serverinfos[i1i].GetSize()
			i1i++
		}
		return resnum
	}
	return 2 + 2 + sizerelySServerInfo1()
}
func ReadMsgSStartMyNotifyCommandByBytes(indata []byte, obj *SStartMyNotifyCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + obj.Serverinfo.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSServerInfoByBytes(data[offset:], &obj.Serverinfo)
	return endpos
}
func WriteMsgSStartMyNotifyCommandByObj(data []byte, obj *SStartMyNotifyCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	offset += WriteMsgSServerInfoByObj(data[offset:], &obj.Serverinfo)
	return offset
}
func GetSizeSStartMyNotifyCommand(obj *SStartMyNotifyCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + obj.Serverinfo.GetSize()
}
func ReadMsgSNotifyAllInfoByBytes(indata []byte, obj *SNotifyAllInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	Serverinfos_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Serverinfos_slen := 0
	Serverinfos_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Serverinfos_slen = int(Serverinfos_slent)
	obj.Serverinfos = make([]SServerInfo,Serverinfos_slen)
	i1i := 0
	for Serverinfos_slen > i1i {
		offset += ReadMsgSServerInfoByBytes(data[offset:],&obj.Serverinfos[i1i])
		i1i++
	}
	return endpos
}
func WriteMsgSNotifyAllInfoByObj(data []byte, obj *SNotifyAllInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Serverinfos)))
	offset += 2
	i1i := 0
	Serverinfos_slen := len(obj.Serverinfos)
	for Serverinfos_slen > i1i {
		offset += WriteMsgSServerInfoByObj(data[offset:],&obj.Serverinfos[i1i])
		i1i++
	}
	return offset
}
func GetSizeSNotifyAllInfo(obj *SNotifyAllInfo) int {
	if obj == nil {
		return 2
	}
	sizerelySServerInfo1 := func()int{
		resnum := 0
		i1i := 0
		Serverinfos_slen := len(obj.Serverinfos)
		for Serverinfos_slen > i1i {
			resnum += obj.Serverinfos[i1i].GetSize()
			i1i++
		}
		return resnum
	}
	return 2 + 2 + sizerelySServerInfo1()
}
func ReadMsgSUpdateGatewayUserAnalysisByBytes(indata []byte, obj *SUpdateGatewayUserAnalysis) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Httpcount = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.Webscoketcount = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.Webscoketcurcount = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgSUpdateGatewayUserAnalysisByObj(data []byte, obj *SUpdateGatewayUserAnalysis) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Httpcount)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Webscoketcount)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Webscoketcurcount)
	offset+=4
	return offset
}
func GetSizeSUpdateGatewayUserAnalysis(obj *SUpdateGatewayUserAnalysis) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + 4 + 4
}
func ReadMsgSAddNewUserToRedisCommandByBytes(indata []byte, obj *SAddNewUserToRedisCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 + len(obj.Openid) > data__len{
		return endpos
	}
	obj.Openid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Openid)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Serverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.ClientConnID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSAddNewUserToRedisCommandByObj(data []byte, obj *SAddNewUserToRedisCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	writeBinaryString(data[offset:],obj.Openid)
	offset += 2 + len(obj.Openid)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Serverid)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.ClientConnID)
	offset+=8
	return offset
}
func GetSizeSAddNewUserToRedisCommand(obj *SAddNewUserToRedisCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + len(obj.Openid) + 4 + 8 + 8
}
func ReadMsgSGatewayForwardCommandByBytes(indata []byte, obj *SGatewayForwardCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Gateserverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.ClientConnID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 + len(obj.Openid) > data__len{
		return endpos
	}
	obj.Openid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Openid)
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i7i := 0
	for Cmddatas_slen > i7i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i7i] = tmpCmddatasvalue
		offset += 1
		i7i++
	}
	return endpos
}
func WriteMsgSGatewayForwardCommandByObj(data []byte, obj *SGatewayForwardCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Gateserverid)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.ClientConnID)
	offset+=8
	writeBinaryString(data[offset:],obj.Openid)
	offset += 2 + len(obj.Openid)
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i7i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i7i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i7i])
		offset += 1
		i7i++
	}
	return offset
}
func GetSizeSGatewayForwardCommand(obj *SGatewayForwardCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + 8 + 2 + len(obj.Openid) + 8 + 
	2 + 2 + 2 + len(obj.Cmddatas) * 1
}
func ReadMsgSGatewayForwardBroadcastCommandByBytes(indata []byte, obj *SGatewayForwardBroadcastCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Gateserverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.ThreadHash = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	UUIDList_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	UUIDList_slen := 0
	UUIDList_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	UUIDList_slen = int(UUIDList_slent)
	obj.UUIDList = make([]uint64,UUIDList_slen)
	i3i := 0
	for UUIDList_slen > i3i {
		if offset + 8 > data__len{
			return endpos
		}
		tmpUUIDListvalue := binary.BigEndian.Uint64(data[offset:offset+8])
		obj.UUIDList[i3i] = tmpUUIDListvalue
		offset += 8
		i3i++
	}
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i6i := 0
	for Cmddatas_slen > i6i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i6i] = tmpCmddatasvalue
		offset += 1
		i6i++
	}
	return endpos
}
func WriteMsgSGatewayForwardBroadcastCommandByObj(data []byte, obj *SGatewayForwardBroadcastCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Gateserverid)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.ThreadHash)
	offset+=4
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.UUIDList)))
	offset += 2
	i3i := 0
	UUIDList_slen := len(obj.UUIDList)
	for UUIDList_slen > i3i {
		binary.BigEndian.PutUint64(data[offset:offset+8],obj.UUIDList[i3i])
		offset += 8
		i3i++
	}
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i6i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i6i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i6i])
		offset += 1
		i6i++
	}
	return offset
}
func GetSizeSGatewayForwardBroadcastCommand(obj *SGatewayForwardBroadcastCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + 4 + 2 + len(obj.UUIDList) * 8 + 2 + 
	2 + 2 + len(obj.Cmddatas) * 1
}
func ReadMsgSGatewayForward2HttpCommandByBytes(indata []byte, obj *SGatewayForward2HttpCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Gateserverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.Httptaskid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 + len(obj.Openid) > data__len{
		return endpos
	}
	obj.Openid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Openid)
	if offset + 2 + len(obj.Cmdname) > data__len{
		return endpos
	}
	obj.Cmdname = readBinaryString(data[offset:])
	offset += 2 + len(obj.Cmdname)
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i6i := 0
	for Cmddatas_slen > i6i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i6i] = tmpCmddatasvalue
		offset += 1
		i6i++
	}
	return endpos
}
func WriteMsgSGatewayForward2HttpCommandByObj(data []byte, obj *SGatewayForward2HttpCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Gateserverid)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Httptaskid)
	offset+=8
	writeBinaryString(data[offset:],obj.Openid)
	offset += 2 + len(obj.Openid)
	writeBinaryString(data[offset:],obj.Cmdname)
	offset += 2 + len(obj.Cmdname)
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i6i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i6i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i6i])
		offset += 1
		i6i++
	}
	return offset
}
func GetSizeSGatewayForward2HttpCommand(obj *SGatewayForward2HttpCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + 8 + 2 + len(obj.Openid) + 2 + len(obj.Cmdname) + 
	2 + 2 + len(obj.Cmddatas) * 1
}
func ReadMsgSBridgeForward2UserCommandByBytes(indata []byte, obj *SBridgeForward2UserCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Touuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i5i := 0
	for Cmddatas_slen > i5i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i5i] = tmpCmddatasvalue
		offset += 1
		i5i++
	}
	return endpos
}
func WriteMsgSBridgeForward2UserCommandByObj(data []byte, obj *SBridgeForward2UserCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Touuid)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i5i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i5i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i5i])
		offset += 1
		i5i++
	}
	return offset
}
func GetSizeSBridgeForward2UserCommand(obj *SBridgeForward2UserCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8 + 2 + 2 + 
	2 + len(obj.Cmddatas) * 1
}
func ReadMsgSBridgeBroadcast2UserCommandByBytes(indata []byte, obj *SBridgeBroadcast2UserCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 + len(obj.Fromopenid) > data__len{
		return endpos
	}
	obj.Fromopenid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Fromopenid)
	if offset + 2 + len(obj.Toopenid) > data__len{
		return endpos
	}
	obj.Toopenid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Toopenid)
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i5i := 0
	for Cmddatas_slen > i5i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i5i] = tmpCmddatasvalue
		offset += 1
		i5i++
	}
	return endpos
}
func WriteMsgSBridgeBroadcast2UserCommandByObj(data []byte, obj *SBridgeBroadcast2UserCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	writeBinaryString(data[offset:],obj.Fromopenid)
	offset += 2 + len(obj.Fromopenid)
	writeBinaryString(data[offset:],obj.Toopenid)
	offset += 2 + len(obj.Toopenid)
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i5i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i5i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i5i])
		offset += 1
		i5i++
	}
	return offset
}
func GetSizeSBridgeBroadcast2UserCommand(obj *SBridgeBroadcast2UserCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + len(obj.Fromopenid) + 2 + len(obj.Toopenid) + 2 + 2 + 
	2 + len(obj.Cmddatas) * 1
}
func ReadMsgSBridgeForward2UserServerByBytes(indata []byte, obj *SBridgeForward2UserServer) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Touuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i5i := 0
	for Cmddatas_slen > i5i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i5i] = tmpCmddatasvalue
		offset += 1
		i5i++
	}
	return endpos
}
func WriteMsgSBridgeForward2UserServerByObj(data []byte, obj *SBridgeForward2UserServer) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Touuid)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i5i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i5i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i5i])
		offset += 1
		i5i++
	}
	return offset
}
func GetSizeSBridgeForward2UserServer(obj *SBridgeForward2UserServer) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8 + 2 + 2 + 
	2 + len(obj.Cmddatas) * 1
}
func ReadMsgSBridgeBroadcast2GatewayServerByBytes(indata []byte, obj *SBridgeBroadcast2GatewayServer) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i3i := 0
	for Cmddatas_slen > i3i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i3i] = tmpCmddatasvalue
		offset += 1
		i3i++
	}
	return endpos
}
func WriteMsgSBridgeBroadcast2GatewayServerByObj(data []byte, obj *SBridgeBroadcast2GatewayServer) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i3i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i3i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i3i])
		offset += 1
		i3i++
	}
	return offset
}
func GetSizeSBridgeBroadcast2GatewayServer(obj *SBridgeBroadcast2GatewayServer) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + 2 + 2 + len(obj.Cmddatas) * 1
}
func ReadMsgSMatchForward2UserServerByBytes(indata []byte, obj *SMatchForward2UserServer) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Touuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i5i := 0
	for Cmddatas_slen > i5i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i5i] = tmpCmddatasvalue
		offset += 1
		i5i++
	}
	return endpos
}
func WriteMsgSMatchForward2UserServerByObj(data []byte, obj *SMatchForward2UserServer) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Touuid)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i5i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i5i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i5i])
		offset += 1
		i5i++
	}
	return offset
}
func GetSizeSMatchForward2UserServer(obj *SMatchForward2UserServer) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8 + 2 + 2 + 
	2 + len(obj.Cmddatas) * 1
}
func ReadMsgSRoomForward2UserServerByBytes(indata []byte, obj *SRoomForward2UserServer) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Touuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i5i := 0
	for Cmddatas_slen > i5i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i5i] = tmpCmddatasvalue
		offset += 1
		i5i++
	}
	return endpos
}
func WriteMsgSRoomForward2UserServerByObj(data []byte, obj *SRoomForward2UserServer) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Touuid)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i5i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i5i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i5i])
		offset += 1
		i5i++
	}
	return offset
}
func GetSizeSRoomForward2UserServer(obj *SRoomForward2UserServer) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8 + 2 + 2 + 
	2 + len(obj.Cmddatas) * 1
}
func ReadMsgSGatewayBroadcast2UserCommandByBytes(indata []byte, obj *SGatewayBroadcast2UserCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Touuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i5i := 0
	for Cmddatas_slen > i5i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i5i] = tmpCmddatasvalue
		offset += 1
		i5i++
	}
	return endpos
}
func WriteMsgSGatewayBroadcast2UserCommandByObj(data []byte, obj *SGatewayBroadcast2UserCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Touuid)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i5i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i5i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i5i])
		offset += 1
		i5i++
	}
	return offset
}
func GetSizeSGatewayBroadcast2UserCommand(obj *SGatewayBroadcast2UserCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8 + 2 + 2 + 
	2 + len(obj.Cmddatas) * 1
}
func ReadMsgSUserServerSearchFriendByBytes(indata []byte, obj *SUserServerSearchFriend) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Touuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSUserServerSearchFriendByObj(data []byte, obj *SUserServerSearchFriend) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Touuid)
	offset+=8
	return offset
}
func GetSizeSUserServerSearchFriend(obj *SUserServerSearchFriend) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8
}
func ReadMsgSUserServerGMCommandByBytes(indata []byte, obj *SUserServerGMCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Taskid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 + len(obj.Key) > data__len{
		return endpos
	}
	obj.Key = readBinaryString(data[offset:])
	offset += 2 + len(obj.Key)
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 + len(obj.Openid) > data__len{
		return endpos
	}
	obj.Openid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Openid)
	if offset + 4 > data__len{
		return endpos
	}
	obj.CmdID = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 2 + len(obj.Param1) > data__len{
		return endpos
	}
	obj.Param1 = readBinaryString(data[offset:])
	offset += 2 + len(obj.Param1)
	if offset + 2 + len(obj.Param2) > data__len{
		return endpos
	}
	obj.Param2 = readBinaryString(data[offset:])
	offset += 2 + len(obj.Param2)
	return endpos
}
func WriteMsgSUserServerGMCommandByObj(data []byte, obj *SUserServerGMCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Taskid)
	offset+=8
	writeBinaryString(data[offset:],obj.Key)
	offset += 2 + len(obj.Key)
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	writeBinaryString(data[offset:],obj.Openid)
	offset += 2 + len(obj.Openid)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.CmdID)
	offset+=4
	writeBinaryString(data[offset:],obj.Param1)
	offset += 2 + len(obj.Param1)
	writeBinaryString(data[offset:],obj.Param2)
	offset += 2 + len(obj.Param2)
	return offset
}
func GetSizeSUserServerGMCommand(obj *SUserServerGMCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 2 + len(obj.Key) + 8 + 2 + len(obj.Openid) + 
	4 + 2 + len(obj.Param1) + 2 + len(obj.Param2)
}
func ReadMsgSRequestOtherUserByBytes(indata []byte, obj *SRequestOtherUser) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Touuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	CmdData_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	CmdData_slen := 0
	CmdData_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	CmdData_slen = int(CmdData_slent)
	obj.CmdData = make([]byte,CmdData_slen)
	i4i := 0
	for CmdData_slen > i4i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmdDatavalue := readBinaryUint8(data[offset:offset+1])
		obj.CmdData[i4i] = tmpCmdDatavalue
		offset += 1
		i4i++
	}
	return endpos
}
func WriteMsgSRequestOtherUserByObj(data []byte, obj *SRequestOtherUser) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Touuid)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.CmdData)))
	offset += 2
	i4i := 0
	CmdData_slen := len(obj.CmdData)
	for CmdData_slen > i4i {
		writeBinaryUint8(data[offset:offset+1],obj.CmdData[i4i])
		offset += 1
		i4i++
	}
	return offset
}
func GetSizeSRequestOtherUser(obj *SRequestOtherUser) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8 + 2 + 2 + len(obj.CmdData) * 1
}
func ReadMsgSResponseOtherUserByBytes(indata []byte, obj *SResponseOtherUser) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Touuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	CmdData_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	CmdData_slen := 0
	CmdData_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	CmdData_slen = int(CmdData_slent)
	obj.CmdData = make([]byte,CmdData_slen)
	i4i := 0
	for CmdData_slen > i4i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmdDatavalue := readBinaryUint8(data[offset:offset+1])
		obj.CmdData[i4i] = tmpCmdDatavalue
		offset += 1
		i4i++
	}
	return endpos
}
func WriteMsgSResponseOtherUserByObj(data []byte, obj *SResponseOtherUser) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Touuid)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.CmdData)))
	offset += 2
	i4i := 0
	CmdData_slen := len(obj.CmdData)
	for CmdData_slen > i4i {
		writeBinaryUint8(data[offset:offset+1],obj.CmdData[i4i])
		offset += 1
		i4i++
	}
	return offset
}
func GetSizeSResponseOtherUser(obj *SResponseOtherUser) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8 + 2 + 2 + len(obj.CmdData) * 1
}
func ReadMsgSBridgeDialGetUserInfoByBytes(indata []byte, obj *SBridgeDialGetUserInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 + len(obj.Fromopenid) > data__len{
		return endpos
	}
	obj.Fromopenid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Fromopenid)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Getuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 + len(obj.Getopenid) > data__len{
		return endpos
	}
	obj.Getopenid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Getopenid)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Type = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgSBridgeDialGetUserInfoByObj(data []byte, obj *SBridgeDialGetUserInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	writeBinaryString(data[offset:],obj.Fromopenid)
	offset += 2 + len(obj.Fromopenid)
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Getuuid)
	offset+=8
	writeBinaryString(data[offset:],obj.Getopenid)
	offset += 2 + len(obj.Getopenid)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Type)
	offset+=4
	return offset
}
func GetSizeSBridgeDialGetUserInfo(obj *SBridgeDialGetUserInfo) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + len(obj.Fromopenid) + 8 + 8 + 2 + len(obj.Getopenid) + 
	4
}
func ReadMsgSGatewayWSLoginUserByBytes(indata []byte, obj *SGatewayWSLoginUser) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Gateserverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.ClientConnID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 + len(obj.Openid) > data__len{
		return endpos
	}
	obj.Openid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Openid)
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 + len(obj.Token) > data__len{
		return endpos
	}
	obj.Token = readBinaryString(data[offset:])
	offset += 2 + len(obj.Token)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Tokenendtime = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 2 + len(obj.Sessionkey) > data__len{
		return endpos
	}
	obj.Sessionkey = readBinaryString(data[offset:])
	offset += 2 + len(obj.Sessionkey)
	if offset + 2 + len(obj.Loginappid) > data__len{
		return endpos
	}
	obj.Loginappid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Loginappid)
	if offset + 2 + len(obj.Username) > data__len{
		return endpos
	}
	obj.Username = readBinaryString(data[offset:])
	offset += 2 + len(obj.Username)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Quizid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Allmoney = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 + len(obj.Headurl) > data__len{
		return endpos
	}
	obj.Headurl = readBinaryString(data[offset:])
	offset += 2 + len(obj.Headurl)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Female = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.Retcode = readBinaryInt32(data[offset:offset+4])
	offset+=4
	if offset + 2 + len(obj.Message) > data__len{
		return endpos
	}
	obj.Message = readBinaryString(data[offset:])
	offset += 2 + len(obj.Message)
	LoginMsg_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	LoginMsg_slen := 0
	LoginMsg_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	LoginMsg_slen = int(LoginMsg_slent)
	obj.LoginMsg = make([]byte,LoginMsg_slen)
	i16i := 0
	for LoginMsg_slen > i16i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpLoginMsgvalue := readBinaryUint8(data[offset:offset+1])
		obj.LoginMsg[i16i] = tmpLoginMsgvalue
		offset += 1
		i16i++
	}
	return endpos
}
func WriteMsgSGatewayWSLoginUserByObj(data []byte, obj *SGatewayWSLoginUser) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Gateserverid)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.ClientConnID)
	offset+=8
	writeBinaryString(data[offset:],obj.Openid)
	offset += 2 + len(obj.Openid)
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	writeBinaryString(data[offset:],obj.Token)
	offset += 2 + len(obj.Token)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Tokenendtime)
	offset+=4
	writeBinaryString(data[offset:],obj.Sessionkey)
	offset += 2 + len(obj.Sessionkey)
	writeBinaryString(data[offset:],obj.Loginappid)
	offset += 2 + len(obj.Loginappid)
	writeBinaryString(data[offset:],obj.Username)
	offset += 2 + len(obj.Username)
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Quizid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Allmoney)
	offset+=8
	writeBinaryString(data[offset:],obj.Headurl)
	offset += 2 + len(obj.Headurl)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Female)
	offset+=4
	writeBinaryInt32(data[offset:offset+4], obj.Retcode)
	offset+=4
	writeBinaryString(data[offset:],obj.Message)
	offset += 2 + len(obj.Message)
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.LoginMsg)))
	offset += 2
	i16i := 0
	LoginMsg_slen := len(obj.LoginMsg)
	for LoginMsg_slen > i16i {
		writeBinaryUint8(data[offset:offset+1],obj.LoginMsg[i16i])
		offset += 1
		i16i++
	}
	return offset
}
func GetSizeSGatewayWSLoginUser(obj *SGatewayWSLoginUser) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + 8 + 2 + len(obj.Openid) + 8 + 
	2 + len(obj.Token) + 4 + 2 + len(obj.Sessionkey) + 2 + len(obj.Loginappid) + 2 + len(obj.Username) + 
	8 + 8 + 2 + len(obj.Headurl) + 4 + 4 + 
	2 + len(obj.Message) + 2 + len(obj.LoginMsg) * 1
}
func ReadMsgSGatewayWSOfflineUserByBytes(indata []byte, obj *SGatewayWSOfflineUser) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 + len(obj.Openid) > data__len{
		return endpos
	}
	obj.Openid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Openid)
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Quizid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.ClientConnID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSGatewayWSOfflineUserByObj(data []byte, obj *SGatewayWSOfflineUser) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	writeBinaryString(data[offset:],obj.Openid)
	offset += 2 + len(obj.Openid)
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Quizid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.ClientConnID)
	offset+=8
	return offset
}
func GetSizeSGatewayWSOfflineUser(obj *SGatewayWSOfflineUser) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + len(obj.Openid) + 8 + 8 + 8
}
func ReadMsgSTemplateMessageKeyWordByBytes(indata []byte, obj *STemplateMessageKeyWord) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 + len(obj.Value) > data__len{
		return endpos
	}
	obj.Value = readBinaryString(data[offset:])
	offset += 2 + len(obj.Value)
	if offset + 2 + len(obj.Color) > data__len{
		return endpos
	}
	obj.Color = readBinaryString(data[offset:])
	offset += 2 + len(obj.Color)
	return endpos
}
func WriteMsgSTemplateMessageKeyWordByObj(data []byte, obj *STemplateMessageKeyWord) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	writeBinaryString(data[offset:],obj.Value)
	offset += 2 + len(obj.Value)
	writeBinaryString(data[offset:],obj.Color)
	offset += 2 + len(obj.Color)
	return offset
}
func GetSizeSTemplateMessageKeyWord(obj *STemplateMessageKeyWord) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + len(obj.Value) + 2 + len(obj.Color)
}
func ReadMsgSQSTemplateMessageByBytes(indata []byte, obj *SQSTemplateMessage) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 + len(obj.Openid) > data__len{
		return endpos
	}
	obj.Openid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Openid)
	if offset + 2 + len(obj.Template_id) > data__len{
		return endpos
	}
	obj.Template_id = readBinaryString(data[offset:])
	offset += 2 + len(obj.Template_id)
	if offset + 2 + len(obj.Page) > data__len{
		return endpos
	}
	obj.Page = readBinaryString(data[offset:])
	offset += 2 + len(obj.Page)
	Datalist_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Datalist_slen := 0
	Datalist_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Datalist_slen = int(Datalist_slent)
	obj.Datalist = make([]STemplateMessageKeyWord,Datalist_slen)
	i4i := 0
	for Datalist_slen > i4i {
		offset += ReadMsgSTemplateMessageKeyWordByBytes(data[offset:],&obj.Datalist[i4i])
		i4i++
	}
	if offset + 2 + len(obj.Formid) > data__len{
		return endpos
	}
	obj.Formid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Formid)
	return endpos
}
func WriteMsgSQSTemplateMessageByObj(data []byte, obj *SQSTemplateMessage) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	writeBinaryString(data[offset:],obj.Openid)
	offset += 2 + len(obj.Openid)
	writeBinaryString(data[offset:],obj.Template_id)
	offset += 2 + len(obj.Template_id)
	writeBinaryString(data[offset:],obj.Page)
	offset += 2 + len(obj.Page)
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Datalist)))
	offset += 2
	i4i := 0
	Datalist_slen := len(obj.Datalist)
	for Datalist_slen > i4i {
		offset += WriteMsgSTemplateMessageKeyWordByObj(data[offset:],&obj.Datalist[i4i])
		i4i++
	}
	writeBinaryString(data[offset:],obj.Formid)
	offset += 2 + len(obj.Formid)
	return offset
}
func GetSizeSQSTemplateMessage(obj *SQSTemplateMessage) int {
	if obj == nil {
		return 2
	}
	sizerelySTemplateMessageKeyWord4 := func()int{
		resnum := 0
		i4i := 0
		Datalist_slen := len(obj.Datalist)
		for Datalist_slen > i4i {
			resnum += obj.Datalist[i4i].GetSize()
			i4i++
		}
		return resnum
	}
	return 2 + 2 + len(obj.Openid) + 2 + len(obj.Template_id) + 2 + len(obj.Page) + 2 + sizerelySTemplateMessageKeyWord4() + 
	2 + len(obj.Formid)
}
func ReadMsgSGatewayChangeAccessTokenByBytes(indata []byte, obj *SGatewayChangeAccessToken) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 + len(obj.Access_token) > data__len{
		return endpos
	}
	obj.Access_token = readBinaryString(data[offset:])
	offset += 2 + len(obj.Access_token)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Update_accesstime = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 2 + len(obj.Access_token_QQ) > data__len{
		return endpos
	}
	obj.Access_token_QQ = readBinaryString(data[offset:])
	offset += 2 + len(obj.Access_token_QQ)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Update_accesstime_QQ = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgSGatewayChangeAccessTokenByObj(data []byte, obj *SGatewayChangeAccessToken) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	writeBinaryString(data[offset:],obj.Access_token)
	offset += 2 + len(obj.Access_token)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Update_accesstime)
	offset+=4
	writeBinaryString(data[offset:],obj.Access_token_QQ)
	offset += 2 + len(obj.Access_token_QQ)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Update_accesstime_QQ)
	offset+=4
	return offset
}
func GetSizeSGatewayChangeAccessToken(obj *SGatewayChangeAccessToken) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + len(obj.Access_token) + 4 + 2 + len(obj.Access_token_QQ) + 4
}
func ReadMsgSUserInfoByBytes(indata []byte, obj *SUserInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 + len(obj.Openid) > data__len{
		return endpos
	}
	obj.Openid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Openid)
	if offset + 4 > data__len{
		return endpos
	}
	obj.UserServerid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.GatewayServerid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.ClientConnID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 4 > data__len{
		return endpos
	}
	obj.RankID = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 1 > data__len{
		return endpos
	}
	obj.IsScript = uint8(data[offset]) != 0
	offset += 1
	if offset + 8 > data__len{
		return endpos
	}
	obj.GMLevel = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	ExtraData_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	ExtraData_slen := 0
	ExtraData_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	ExtraData_slen = int(ExtraData_slent)
	obj.ExtraData = make([]byte,ExtraData_slen)
	i9i := 0
	for ExtraData_slen > i9i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpExtraDatavalue := readBinaryUint8(data[offset:offset+1])
		obj.ExtraData[i9i] = tmpExtraDatavalue
		offset += 1
		i9i++
	}
	return endpos
}
func WriteMsgSUserInfoByObj(data []byte, obj *SUserInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	writeBinaryString(data[offset:],obj.Openid)
	offset += 2 + len(obj.Openid)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.UserServerid)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.GatewayServerid)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.ClientConnID)
	offset+=8
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.RankID)
	offset+=4
	data[offset] = uint8(bool2int(obj.IsScript))
	offset += 1
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.GMLevel)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.ExtraData)))
	offset += 2
	i9i := 0
	ExtraData_slen := len(obj.ExtraData)
	for ExtraData_slen > i9i {
		writeBinaryUint8(data[offset:offset+1],obj.ExtraData[i9i])
		offset += 1
		i9i++
	}
	return offset
}
func GetSizeSUserInfo(obj *SUserInfo) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 2 + len(obj.Openid) + 4 + 4 + 
	8 + 4 + 1 + 8 + 2 + len(obj.ExtraData) * 1
}
func ReadMsgSJoinMatchByBytes(indata []byte, obj *SJoinMatch) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.RoomID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 4 > data__len{
		return endpos
	}
	obj.RoomType = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.SubType = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 1 > data__len{
		return endpos
	}
	obj.OnlyJoin = uint8(data[offset]) != 0
	offset += 1
	if offset + 1 > data__len{
		return endpos
	}
	obj.OnlyCreate = uint8(data[offset]) != 0
	offset += 1
	if offset + obj.Userinfo.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSUserInfoByBytes(data[offset:], &obj.Userinfo)
	return endpos
}
func WriteMsgSJoinMatchByObj(data []byte, obj *SJoinMatch) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.RoomID)
	offset+=8
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.RoomType)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.SubType)
	offset+=4
	data[offset] = uint8(bool2int(obj.OnlyJoin))
	offset += 1
	data[offset] = uint8(bool2int(obj.OnlyCreate))
	offset += 1
	offset += WriteMsgSUserInfoByObj(data[offset:], &obj.Userinfo)
	return offset
}
func GetSizeSJoinMatch(obj *SJoinMatch) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 4 + 4 + 1 + 
	1 + obj.Userinfo.GetSize()
}
func ReadMsgSTeamInfoByBytes(indata []byte, obj *STeamInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	Members_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Members_slen := 0
	Members_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Members_slen = int(Members_slent)
	obj.Members = make([]SUserInfo,Members_slen)
	i1i := 0
	for Members_slen > i1i {
		offset += ReadMsgSUserInfoByBytes(data[offset:],&obj.Members[i1i])
		i1i++
	}
	return endpos
}
func WriteMsgSTeamInfoByObj(data []byte, obj *STeamInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Members)))
	offset += 2
	i1i := 0
	Members_slen := len(obj.Members)
	for Members_slen > i1i {
		offset += WriteMsgSUserInfoByObj(data[offset:],&obj.Members[i1i])
		i1i++
	}
	return offset
}
func GetSizeSTeamInfo(obj *STeamInfo) int {
	if obj == nil {
		return 2
	}
	sizerelySUserInfo1 := func()int{
		resnum := 0
		i1i := 0
		Members_slen := len(obj.Members)
		for Members_slen > i1i {
			resnum += obj.Members[i1i].GetSize()
			i1i++
		}
		return resnum
	}
	return 2 + 2 + sizerelySUserInfo1()
}
func ReadMsgSRoomInfoByBytes(indata []byte, obj *SRoomInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	Teams_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Teams_slen := 0
	Teams_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Teams_slen = int(Teams_slent)
	obj.Teams = make([]STeamInfo,Teams_slen)
	i1i := 0
	for Teams_slen > i1i {
		offset += ReadMsgSTeamInfoByBytes(data[offset:],&obj.Teams[i1i])
		i1i++
	}
	return endpos
}
func WriteMsgSRoomInfoByObj(data []byte, obj *SRoomInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Teams)))
	offset += 2
	i1i := 0
	Teams_slen := len(obj.Teams)
	for Teams_slen > i1i {
		offset += WriteMsgSTeamInfoByObj(data[offset:],&obj.Teams[i1i])
		i1i++
	}
	return offset
}
func GetSizeSRoomInfo(obj *SRoomInfo) int {
	if obj == nil {
		return 2
	}
	sizerelySTeamInfo1 := func()int{
		resnum := 0
		i1i := 0
		Teams_slen := len(obj.Teams)
		for Teams_slen > i1i {
			resnum += obj.Teams[i1i].GetSize()
			i1i++
		}
		return resnum
	}
	return 2 + 2 + sizerelySTeamInfo1()
}
func ReadMsgSUserMatchInfoByBytes(indata []byte, obj *SUserMatchInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 1 > data__len{
		return endpos
	}
	obj.Sec = uint8(data[offset]) != 0
	offset += 1
	if offset + 1 > data__len{
		return endpos
	}
	obj.Done = uint8(data[offset]) != 0
	offset += 1
	if offset + 4 > data__len{
		return endpos
	}
	obj.Total = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 > data__len{
		return endpos
	}
	obj.Now = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.Matchindex = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSUserMatchInfoByObj(data []byte, obj *SUserMatchInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	data[offset] = uint8(bool2int(obj.Sec))
	offset += 1
	data[offset] = uint8(bool2int(obj.Done))
	offset += 1
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Total)
	offset+=4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Now)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Matchindex)
	offset+=8
	return offset
}
func GetSizeSUserMatchInfo(obj *SUserMatchInfo) int {
	if obj == nil {
		return 2
	}
	return 2 + 1 + 1 + 4 + 4 + 
	8
}
func ReadMsgSUserQuitMatchByBytes(indata []byte, obj *SUserQuitMatch) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.RoomID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 1 > data__len{
		return endpos
	}
	obj.Sec = uint8(data[offset]) != 0
	offset += 1
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 4 > data__len{
		return endpos
	}
	obj.Type = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.Matchindex = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSUserQuitMatchByObj(data []byte, obj *SUserQuitMatch) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.RoomID)
	offset+=8
	data[offset] = uint8(bool2int(obj.Sec))
	offset += 1
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Type)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Matchindex)
	offset+=8
	return offset
}
func GetSizeSUserQuitMatch(obj *SUserQuitMatch) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 1 + 8 + 4 + 
	8
}
func ReadMsgSMatchDoneByBytes(indata []byte, obj *SMatchDone) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Matchindex = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSMatchDoneByObj(data []byte, obj *SMatchDone) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Matchindex)
	offset+=8
	return offset
}
func GetSizeSMatchDone(obj *SMatchDone) int {
	if obj == nil {
		return 2
	}
	return 2 + 8
}
func ReadMsgSUserQuitRoomByBytes(indata []byte, obj *SUserQuitRoom) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.RoomID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 1 > data__len{
		return endpos
	}
	obj.Sec = uint8(data[offset]) != 0
	offset += 1
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 4 > data__len{
		return endpos
	}
	obj.Type = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos
	}
	obj.Roomindex = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSUserQuitRoomByObj(data []byte, obj *SUserQuitRoom) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.RoomID)
	offset+=8
	data[offset] = uint8(bool2int(obj.Sec))
	offset += 1
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Type)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Roomindex)
	offset+=8
	return offset
}
func GetSizeSUserQuitRoom(obj *SUserQuitRoom) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 1 + 8 + 4 + 
	8
}
func ReadMsgSUserRoomInfoByBytes(indata []byte, obj *SUserRoomInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 1 > data__len{
		return endpos
	}
	obj.Sec = uint8(data[offset]) != 0
	offset += 1
	if offset + 8 > data__len{
		return endpos
	}
	obj.Roomindex = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSUserRoomInfoByObj(data []byte, obj *SUserRoomInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	data[offset] = uint8(bool2int(obj.Sec))
	offset += 1
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Roomindex)
	offset+=8
	return offset
}
func GetSizeSUserRoomInfo(obj *SUserRoomInfo) int {
	if obj == nil {
		return 2
	}
	return 2 + 1 + 8
}
func ReadMsgSNotifyUserOfflineByBytes(indata []byte, obj *SNotifyUserOffline) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Matchindex = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Roomindex = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSNotifyUserOfflineByObj(data []byte, obj *SNotifyUserOffline) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Matchindex)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Roomindex)
	offset+=8
	return offset
}
func GetSizeSNotifyUserOffline(obj *SNotifyUserOffline) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8 + 8
}
func ReadMsgSNotifyUserReonlineByBytes(indata []byte, obj *SNotifyUserReonline) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + obj.UserInfo.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSUserInfoByBytes(data[offset:], &obj.UserInfo)
	return endpos
}
func WriteMsgSNotifyUserReonlineByObj(data []byte, obj *SNotifyUserReonline) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	offset += WriteMsgSUserInfoByObj(data[offset:], &obj.UserInfo)
	return offset
}
func GetSizeSNotifyUserReonline(obj *SNotifyUserReonline) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + obj.UserInfo.GetSize()
}
func ReadMsgSRPCCheckUserTokenByBytes(indata []byte, obj *SRPCCheckUserToken) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 + len(obj.Openid) > data__len{
		return endpos
	}
	obj.Openid = readBinaryString(data[offset:])
	offset += 2 + len(obj.Openid)
	if offset + 2 + len(obj.Token) > data__len{
		return endpos
	}
	obj.Token = readBinaryString(data[offset:])
	offset += 2 + len(obj.Token)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Retcode = readBinaryInt32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgSRPCCheckUserTokenByObj(data []byte, obj *SRPCCheckUserToken) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	writeBinaryString(data[offset:],obj.Openid)
	offset += 2 + len(obj.Openid)
	writeBinaryString(data[offset:],obj.Token)
	offset += 2 + len(obj.Token)
	writeBinaryInt32(data[offset:offset+4], obj.Retcode)
	offset+=4
	return offset
}
func GetSizeSRPCCheckUserToken(obj *SRPCCheckUserToken) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + len(obj.Openid) + 2 + len(obj.Token) + 4
}
func ReadMsgSMatchBroadcast2UserServerCommandByBytes(indata []byte, obj *SMatchBroadcast2UserServerCommand) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Fromuuid = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.Matchindex = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdid = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 2 > data__len{
		return endpos
	}
	obj.Cmdlen = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	Cmddatas_slen := 0
	Cmddatas_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	Cmddatas_slen = int(Cmddatas_slent)
	obj.Cmddatas = make([]byte,Cmddatas_slen)
	i5i := 0
	for Cmddatas_slen > i5i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpCmddatasvalue := readBinaryUint8(data[offset:offset+1])
		obj.Cmddatas[i5i] = tmpCmddatasvalue
		offset += 1
		i5i++
	}
	return endpos
}
func WriteMsgSMatchBroadcast2UserServerCommandByObj(data []byte, obj *SMatchBroadcast2UserServerCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Fromuuid)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Matchindex)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdid)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.Cmdlen)
	offset+=2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.Cmddatas)))
	offset += 2
	i5i := 0
	Cmddatas_slen := len(obj.Cmddatas)
	for Cmddatas_slen > i5i {
		writeBinaryUint8(data[offset:offset+1],obj.Cmddatas[i5i])
		offset += 1
		i5i++
	}
	return offset
}
func GetSizeSMatchBroadcast2UserServerCommand(obj *SMatchBroadcast2UserServerCommand) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + 8 + 2 + 2 + 
	2 + len(obj.Cmddatas) * 1
}
func ReadMsgNotifyGatewayInfoByBytes(indata []byte, obj *NotifyGatewayInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Serverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 2 + len(obj.Serverip) > data__len{
		return endpos
	}
	obj.Serverip = readBinaryString(data[offset:])
	offset += 2 + len(obj.Serverip)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Serverport = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 2 > data__len{
		return endpos
	}
	obj.State = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 4 > data__len{
		return endpos
	}
	obj.TaskSize = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgNotifyGatewayInfoByObj(data []byte, obj *NotifyGatewayInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Serverid)
	offset+=4
	writeBinaryString(data[offset:],obj.Serverip)
	offset += 2 + len(obj.Serverip)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Serverport)
	offset+=4
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.State)
	offset+=2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.TaskSize)
	offset+=4
	return offset
}
func GetSizeNotifyGatewayInfo(obj *NotifyGatewayInfo) int {
	if obj == nil {
		return 2
	}
	return 2 + 4 + 2 + len(obj.Serverip) + 4 + 2 + 
	4
}
func ReadMsgSuperReqGatewayInfoByBytes(indata []byte, obj *SuperReqGatewayInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	return endpos
}
func WriteMsgSuperReqGatewayInfoByObj(data []byte, obj *SuperReqGatewayInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	return offset
}
func GetSizeSuperReqGatewayInfo(obj *SuperReqGatewayInfo) int {
	if obj == nil {
		return 2
	}
	return 2 + 0
}
func ReadMsgNotifyGateUserNumsByBytes(indata []byte, obj *NotifyGateUserNums) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Usernum = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgNotifyGateUserNumsByObj(data []byte, obj *NotifyGateUserNums) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Usernum)
	offset+=4
	return offset
}
func GetSizeNotifyGateUserNums(obj *NotifyGateUserNums) int {
	if obj == nil {
		return 2
	}
	return 2 + 4
}
func ReadMsgGateReqUserNumsByBytes(indata []byte, obj *GateReqUserNums) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	return endpos
}
func WriteMsgGateReqUserNumsByObj(data []byte, obj *GateReqUserNums) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	return offset
}
func GetSizeGateReqUserNums(obj *GateReqUserNums) int {
	if obj == nil {
		return 2
	}
	return 2 + 0
}
func ReadMsgNotifyMatchRoomNumsByBytes(indata []byte, obj *NotifyMatchRoomNums) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Roomnum = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgNotifyMatchRoomNumsByObj(data []byte, obj *NotifyMatchRoomNums) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Roomnum)
	offset+=4
	return offset
}
func GetSizeNotifyMatchRoomNums(obj *NotifyMatchRoomNums) int {
	if obj == nil {
		return 2
	}
	return 2 + 4
}
func ReadMsgMatchReqRoomNumsByBytes(indata []byte, obj *MatchReqRoomNums) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	return endpos
}
func WriteMsgMatchReqRoomNumsByObj(data []byte, obj *MatchReqRoomNums) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	return offset
}
func GetSizeMatchReqRoomNums(obj *MatchReqRoomNums) int {
	if obj == nil {
		return 2
	}
	return 2 + 0
}
func ReadMsgSUserCheckEffectiveByBytes(indata []byte, obj *SUserCheckEffective) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.UUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + obj.UserInfo.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSUserInfoByBytes(data[offset:], &obj.UserInfo)
	return endpos
}
func WriteMsgSUserCheckEffectiveByObj(data []byte, obj *SUserCheckEffective) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.UUID)
	offset+=8
	offset += WriteMsgSUserInfoByObj(data[offset:], &obj.UserInfo)
	return offset
}
func GetSizeSUserCheckEffective(obj *SUserCheckEffective) int {
	if obj == nil {
		return 2
	}
	return 2 + 8 + obj.UserInfo.GetSize()
}
func ReadMsgSRedisConfigItemByBytes(indata []byte, obj *SRedisConfigItem) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 + len(obj.IP) > data__len{
		return endpos
	}
	obj.IP = readBinaryString(data[offset:])
	offset += 2 + len(obj.IP)
	if offset + 4 > data__len{
		return endpos
	}
	obj.Port = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos
}
func WriteMsgSRedisConfigItemByObj(data []byte, obj *SRedisConfigItem) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	writeBinaryString(data[offset:],obj.IP)
	offset += 2 + len(obj.IP)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Port)
	offset+=4
	return offset
}
func GetSizeSRedisConfigItem(obj *SRedisConfigItem) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + len(obj.IP) + 4
}
func ReadMsgSRedisConfigByBytes(indata []byte, obj *SRedisConfig) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	RedisList_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	RedisList_slen := 0
	RedisList_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	RedisList_slen = int(RedisList_slent)
	obj.RedisList = make([]SRedisConfigItem,RedisList_slen)
	i1i := 0
	for RedisList_slen > i1i {
		offset += ReadMsgSRedisConfigItemByBytes(data[offset:],&obj.RedisList[i1i])
		i1i++
	}
	return endpos
}
func WriteMsgSRedisConfigByObj(data []byte, obj *SRedisConfig) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.RedisList)))
	offset += 2
	i1i := 0
	RedisList_slen := len(obj.RedisList)
	for RedisList_slen > i1i {
		offset += WriteMsgSRedisConfigItemByObj(data[offset:],&obj.RedisList[i1i])
		i1i++
	}
	return offset
}
func GetSizeSRedisConfig(obj *SRedisConfig) int {
	if obj == nil {
		return 2
	}
	sizerelySRedisConfigItem1 := func()int{
		resnum := 0
		i1i := 0
		RedisList_slen := len(obj.RedisList)
		for RedisList_slen > i1i {
			resnum += obj.RedisList[i1i].GetSize()
			i1i++
		}
		return resnum
	}
	return 2 + 2 + sizerelySRedisConfigItem1()
}
func ReadMsgSRequestServerInfoByBytes(indata []byte, obj *SRequestServerInfo) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	return endpos
}
func WriteMsgSRequestServerInfoByObj(data []byte, obj *SRequestServerInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	return offset
}
func GetSizeSRequestServerInfo(obj *SRequestServerInfo) int {
	if obj == nil {
		return 2
	}
	return 2 + 0
}
func ReadMsgSNotifySafelyQuitByBytes(indata []byte, obj *SNotifySafelyQuit) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + obj.TargetServerInfo.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSServerInfoByBytes(data[offset:], &obj.TargetServerInfo)
	return endpos
}
func WriteMsgSNotifySafelyQuitByObj(data []byte, obj *SNotifySafelyQuit) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	offset += WriteMsgSServerInfoByObj(data[offset:], &obj.TargetServerInfo)
	return offset
}
func GetSizeSNotifySafelyQuit(obj *SNotifySafelyQuit) int {
	if obj == nil {
		return 2
	}
	return 2 + obj.TargetServerInfo.GetSize()
}
func ReadMsgSAIUserRegisterByBytes(indata []byte, obj *SAIUserRegister) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + obj.Userinfo.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSUserInfoByBytes(data[offset:], &obj.Userinfo)
	return endpos
}
func WriteMsgSAIUserRegisterByObj(data []byte, obj *SAIUserRegister) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	offset += WriteMsgSUserInfoByObj(data[offset:], &obj.Userinfo)
	return offset
}
func GetSizeSAIUserRegister(obj *SAIUserRegister) int {
	if obj == nil {
		return 2
	}
	return 2 + obj.Userinfo.GetSize()
}
func ReadMsgSTaskCmdForwardByBytes(indata []byte, obj *STaskCmdForward) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 2 > data__len{
		return endpos
	}
	obj.CmdID = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 8 > data__len{
		return endpos
	}
	obj.RoomID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	if offset + 8 > data__len{
		return endpos
	}
	obj.ProxyUUID = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	ProtoData_slent := uint16(0)
	if offset + 2 > data__len{
		return endpos
	}
	ProtoData_slen := 0
	ProtoData_slent = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	ProtoData_slen = int(ProtoData_slent)
	obj.ProtoData = make([]byte,ProtoData_slen)
	i4i := 0
	for ProtoData_slen > i4i {
		if offset + 1 > data__len{
			return endpos
		}
		tmpProtoDatavalue := readBinaryUint8(data[offset:offset+1])
		obj.ProtoData[i4i] = tmpProtoDatavalue
		offset += 1
		i4i++
	}
	return endpos
}
func WriteMsgSTaskCmdForwardByObj(data []byte, obj *STaskCmdForward) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.CmdID)
	offset+=2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.RoomID)
	offset+=8
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.ProxyUUID)
	offset+=8
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.ProtoData)))
	offset += 2
	i4i := 0
	ProtoData_slen := len(obj.ProtoData)
	for ProtoData_slen > i4i {
		writeBinaryUint8(data[offset:offset+1],obj.ProtoData[i4i])
		offset += 1
		i4i++
	}
	return offset
}
func GetSizeSTaskCmdForward(obj *STaskCmdForward) int {
	if obj == nil {
		return 2
	}
	return 2 + 2 + 8 + 8 + 2 + len(obj.ProtoData) * 1
}
func ReadMsgSBeginMatchByBytes(indata []byte, obj *SBeginMatch) int {
	offset := 0
	if len(indata) < 2 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))
	offset += 2
	if objsize == 0 {
		return 2
	}
	if offset + objsize > len(indata ) {
		return 2
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 8 > data__len{
		return endpos
	}
	obj.Matchindex = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos
}
func WriteMsgSBeginMatchByObj(data []byte, obj *SBeginMatch) int {
	if obj == nil {
		binary.BigEndian.PutUint16(data[0:2],0)
		return 2
	}
	objsize := obj.GetSize() - 2
	offset := 0
	binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))
	offset += 2
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Matchindex)
	offset+=8
	return offset
}
func GetSizeSBeginMatch(obj *SBeginMatch) int {
	if obj == nil {
		return 2
	}
	return 2 + 8
}
