package servercomm

import (
	"encoding/binary"
	"encoding/json"
	"math"
)

const (
	ModuleInfoID              = 36
	STimeTickCommandID        = 37
	STestCommandID            = 38
	SLoginCommandID           = 39
	SLogoutCommandID          = 40
	SSeverStartOKCommandID    = 41
	SLoginRetCommandID        = 42
	SStartRelyNotifyCommandID = 43
	SStartMyNotifyCommandID   = 44
	SNotifyAllInfoID          = 45
	SNotifySafelyQuitID       = 46
	SUpdateSessionID          = 47
	SReqCloseConnectID        = 48
	SForwardToModuleID        = 49
	ModuleMessageID           = 50
	SForwardToClientID        = 51
	SForwardFromGateID        = 52
	ClientMessageID           = 53
	SROCRequestID             = 54
	SROCResponseID            = 55
	SROCBindID                = 56
)

const (
	ModuleInfoName              = "servercomm.ModuleInfo"
	STimeTickCommandName        = "servercomm.STimeTickCommand"
	STestCommandName            = "servercomm.STestCommand"
	SLoginCommandName           = "servercomm.SLoginCommand"
	SLogoutCommandName          = "servercomm.SLogoutCommand"
	SSeverStartOKCommandName    = "servercomm.SSeverStartOKCommand"
	SLoginRetCommandName        = "servercomm.SLoginRetCommand"
	SStartRelyNotifyCommandName = "servercomm.SStartRelyNotifyCommand"
	SStartMyNotifyCommandName   = "servercomm.SStartMyNotifyCommand"
	SNotifyAllInfoName          = "servercomm.SNotifyAllInfo"
	SNotifySafelyQuitName       = "servercomm.SNotifySafelyQuit"
	SUpdateSessionName          = "servercomm.SUpdateSession"
	SReqCloseConnectName        = "servercomm.SReqCloseConnect"
	SForwardToModuleName        = "servercomm.SForwardToModule"
	ModuleMessageName           = "servercomm.ModuleMessage"
	SForwardToClientName        = "servercomm.SForwardToClient"
	SForwardFromGateName        = "servercomm.SForwardFromGate"
	ClientMessageName           = "servercomm.ClientMessage"
	SROCRequestName             = "servercomm.SROCRequest"
	SROCResponseName            = "servercomm.SROCResponse"
	SROCBindName                = "servercomm.SROCBind"
)

func (this *ModuleInfo) WriteBinary(data []byte) int {
	return WriteMsgModuleInfoByObj(data, this)
}

func (this *STimeTickCommand) WriteBinary(data []byte) int {
	return WriteMsgSTimeTickCommandByObj(data, this)
}

func (this *STestCommand) WriteBinary(data []byte) int {
	return WriteMsgSTestCommandByObj(data, this)
}

func (this *SLoginCommand) WriteBinary(data []byte) int {
	return WriteMsgSLoginCommandByObj(data, this)
}

func (this *SLogoutCommand) WriteBinary(data []byte) int {
	return WriteMsgSLogoutCommandByObj(data, this)
}

func (this *SSeverStartOKCommand) WriteBinary(data []byte) int {
	return WriteMsgSSeverStartOKCommandByObj(data, this)
}

func (this *SLoginRetCommand) WriteBinary(data []byte) int {
	return WriteMsgSLoginRetCommandByObj(data, this)
}

func (this *SStartRelyNotifyCommand) WriteBinary(data []byte) int {
	return WriteMsgSStartRelyNotifyCommandByObj(data, this)
}

func (this *SStartMyNotifyCommand) WriteBinary(data []byte) int {
	return WriteMsgSStartMyNotifyCommandByObj(data, this)
}

func (this *SNotifyAllInfo) WriteBinary(data []byte) int {
	return WriteMsgSNotifyAllInfoByObj(data, this)
}

func (this *SNotifySafelyQuit) WriteBinary(data []byte) int {
	return WriteMsgSNotifySafelyQuitByObj(data, this)
}

func (this *SUpdateSession) WriteBinary(data []byte) int {
	return WriteMsgSUpdateSessionByObj(data, this)
}

func (this *SReqCloseConnect) WriteBinary(data []byte) int {
	return WriteMsgSReqCloseConnectByObj(data, this)
}

func (this *SForwardToModule) WriteBinary(data []byte) int {
	return WriteMsgSForwardToModuleByObj(data, this)
}

func (this *ModuleMessage) WriteBinary(data []byte) int {
	return WriteMsgModuleMessageByObj(data, this)
}

func (this *SForwardToClient) WriteBinary(data []byte) int {
	return WriteMsgSForwardToClientByObj(data, this)
}

func (this *SForwardFromGate) WriteBinary(data []byte) int {
	return WriteMsgSForwardFromGateByObj(data, this)
}

func (this *ClientMessage) WriteBinary(data []byte) int {
	return WriteMsgClientMessageByObj(data, this)
}

func (this *SROCRequest) WriteBinary(data []byte) int {
	return WriteMsgSROCRequestByObj(data, this)
}

func (this *SROCResponse) WriteBinary(data []byte) int {
	return WriteMsgSROCResponseByObj(data, this)
}

func (this *SROCBind) WriteBinary(data []byte) int {
	return WriteMsgSROCBindByObj(data, this)
}

func (this *ModuleInfo) ReadBinary(data []byte) int {
	size, _ := ReadMsgModuleInfoByBytes(data, this)
	return size
}

func (this *STimeTickCommand) ReadBinary(data []byte) int {
	size, _ := ReadMsgSTimeTickCommandByBytes(data, this)
	return size
}

func (this *STestCommand) ReadBinary(data []byte) int {
	size, _ := ReadMsgSTestCommandByBytes(data, this)
	return size
}

func (this *SLoginCommand) ReadBinary(data []byte) int {
	size, _ := ReadMsgSLoginCommandByBytes(data, this)
	return size
}

func (this *SLogoutCommand) ReadBinary(data []byte) int {
	size, _ := ReadMsgSLogoutCommandByBytes(data, this)
	return size
}

func (this *SSeverStartOKCommand) ReadBinary(data []byte) int {
	size, _ := ReadMsgSSeverStartOKCommandByBytes(data, this)
	return size
}

func (this *SLoginRetCommand) ReadBinary(data []byte) int {
	size, _ := ReadMsgSLoginRetCommandByBytes(data, this)
	return size
}

func (this *SStartRelyNotifyCommand) ReadBinary(data []byte) int {
	size, _ := ReadMsgSStartRelyNotifyCommandByBytes(data, this)
	return size
}

func (this *SStartMyNotifyCommand) ReadBinary(data []byte) int {
	size, _ := ReadMsgSStartMyNotifyCommandByBytes(data, this)
	return size
}

func (this *SNotifyAllInfo) ReadBinary(data []byte) int {
	size, _ := ReadMsgSNotifyAllInfoByBytes(data, this)
	return size
}

func (this *SNotifySafelyQuit) ReadBinary(data []byte) int {
	size, _ := ReadMsgSNotifySafelyQuitByBytes(data, this)
	return size
}

func (this *SUpdateSession) ReadBinary(data []byte) int {
	size, _ := ReadMsgSUpdateSessionByBytes(data, this)
	return size
}

func (this *SReqCloseConnect) ReadBinary(data []byte) int {
	size, _ := ReadMsgSReqCloseConnectByBytes(data, this)
	return size
}

func (this *SForwardToModule) ReadBinary(data []byte) int {
	size, _ := ReadMsgSForwardToModuleByBytes(data, this)
	return size
}

func (this *ModuleMessage) ReadBinary(data []byte) int {
	size, _ := ReadMsgModuleMessageByBytes(data, this)
	return size
}

func (this *SForwardToClient) ReadBinary(data []byte) int {
	size, _ := ReadMsgSForwardToClientByBytes(data, this)
	return size
}

func (this *SForwardFromGate) ReadBinary(data []byte) int {
	size, _ := ReadMsgSForwardFromGateByBytes(data, this)
	return size
}

func (this *ClientMessage) ReadBinary(data []byte) int {
	size, _ := ReadMsgClientMessageByBytes(data, this)
	return size
}

func (this *SROCRequest) ReadBinary(data []byte) int {
	size, _ := ReadMsgSROCRequestByBytes(data, this)
	return size
}

func (this *SROCResponse) ReadBinary(data []byte) int {
	size, _ := ReadMsgSROCResponseByBytes(data, this)
	return size
}

func (this *SROCBind) ReadBinary(data []byte) int {
	size, _ := ReadMsgSROCBindByBytes(data, this)
	return size
}

func MsgIdToString(id uint16) string {
	switch id {
	case ModuleInfoID:
		return ModuleInfoName
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
	case SNotifySafelyQuitID:
		return SNotifySafelyQuitName
	case SUpdateSessionID:
		return SUpdateSessionName
	case SReqCloseConnectID:
		return SReqCloseConnectName
	case SForwardToModuleID:
		return SForwardToModuleName
	case ModuleMessageID:
		return ModuleMessageName
	case SForwardToClientID:
		return SForwardToClientName
	case SForwardFromGateID:
		return SForwardFromGateName
	case ClientMessageID:
		return ClientMessageName
	case SROCRequestID:
		return SROCRequestName
	case SROCResponseID:
		return SROCResponseName
	case SROCBindID:
		return SROCBindName
	default:
		return ""
	}
}

func StringToMsgId(msgname string) uint16 {
	switch msgname {
	case ModuleInfoName:
		return ModuleInfoID
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
	case SNotifySafelyQuitName:
		return SNotifySafelyQuitID
	case SUpdateSessionName:
		return SUpdateSessionID
	case SReqCloseConnectName:
		return SReqCloseConnectID
	case SForwardToModuleName:
		return SForwardToModuleID
	case ModuleMessageName:
		return ModuleMessageID
	case SForwardToClientName:
		return SForwardToClientID
	case SForwardFromGateName:
		return SForwardFromGateID
	case ClientMessageName:
		return ClientMessageID
	case SROCRequestName:
		return SROCRequestID
	case SROCResponseName:
		return SROCResponseID
	case SROCBindName:
		return SROCBindID
	default:
		return 0
	}
}

func (this *ModuleInfo) GetMsgId() uint16 {
	return ModuleInfoID
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

func (this *SNotifySafelyQuit) GetMsgId() uint16 {
	return SNotifySafelyQuitID
}

func (this *SUpdateSession) GetMsgId() uint16 {
	return SUpdateSessionID
}

func (this *SReqCloseConnect) GetMsgId() uint16 {
	return SReqCloseConnectID
}

func (this *SForwardToModule) GetMsgId() uint16 {
	return SForwardToModuleID
}

func (this *ModuleMessage) GetMsgId() uint16 {
	return ModuleMessageID
}

func (this *SForwardToClient) GetMsgId() uint16 {
	return SForwardToClientID
}

func (this *SForwardFromGate) GetMsgId() uint16 {
	return SForwardFromGateID
}

func (this *ClientMessage) GetMsgId() uint16 {
	return ClientMessageID
}

func (this *SROCRequest) GetMsgId() uint16 {
	return SROCRequestID
}

func (this *SROCResponse) GetMsgId() uint16 {
	return SROCResponseID
}

func (this *SROCBind) GetMsgId() uint16 {
	return SROCBindID
}

func (this *ModuleInfo) GetMsgName() string {
	return ModuleInfoName
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

func (this *SNotifySafelyQuit) GetMsgName() string {
	return SNotifySafelyQuitName
}

func (this *SUpdateSession) GetMsgName() string {
	return SUpdateSessionName
}

func (this *SReqCloseConnect) GetMsgName() string {
	return SReqCloseConnectName
}

func (this *SForwardToModule) GetMsgName() string {
	return SForwardToModuleName
}

func (this *ModuleMessage) GetMsgName() string {
	return ModuleMessageName
}

func (this *SForwardToClient) GetMsgName() string {
	return SForwardToClientName
}

func (this *SForwardFromGate) GetMsgName() string {
	return SForwardFromGateName
}

func (this *ClientMessage) GetMsgName() string {
	return ClientMessageName
}

func (this *SROCRequest) GetMsgName() string {
	return SROCRequestName
}

func (this *SROCResponse) GetMsgName() string {
	return SROCResponseName
}

func (this *SROCBind) GetMsgName() string {
	return SROCBindName
}

func (this *ModuleInfo) GetSize() int {
	return GetSizeModuleInfo(this)
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

func (this *SNotifySafelyQuit) GetSize() int {
	return GetSizeSNotifySafelyQuit(this)
}

func (this *SUpdateSession) GetSize() int {
	return GetSizeSUpdateSession(this)
}

func (this *SReqCloseConnect) GetSize() int {
	return GetSizeSReqCloseConnect(this)
}

func (this *SForwardToModule) GetSize() int {
	return GetSizeSForwardToModule(this)
}

func (this *ModuleMessage) GetSize() int {
	return GetSizeModuleMessage(this)
}

func (this *SForwardToClient) GetSize() int {
	return GetSizeSForwardToClient(this)
}

func (this *SForwardFromGate) GetSize() int {
	return GetSizeSForwardFromGate(this)
}

func (this *ClientMessage) GetSize() int {
	return GetSizeClientMessage(this)
}

func (this *SROCRequest) GetSize() int {
	return GetSizeSROCRequest(this)
}

func (this *SROCResponse) GetSize() int {
	return GetSizeSROCResponse(this)
}

func (this *SROCBind) GetSize() int {
	return GetSizeSROCBind(this)
}

func (this *ModuleInfo) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *STimeTickCommand) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *STestCommand) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SLoginCommand) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SLogoutCommand) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SSeverStartOKCommand) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SLoginRetCommand) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SStartRelyNotifyCommand) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SStartMyNotifyCommand) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SNotifyAllInfo) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SNotifySafelyQuit) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SUpdateSession) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SReqCloseConnect) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SForwardToModule) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *ModuleMessage) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SForwardToClient) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SForwardFromGate) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *ClientMessage) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SROCRequest) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SROCResponse) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func (this *SROCBind) GetJSON() string {
	json, _ := json.Marshal(this)
	return string(json)
}

func readBinaryString(data []byte) string {
	strfunclen := binary.LittleEndian.Uint32(data[:4])
	if int(strfunclen)+4 > len(data) {
		return ""
	}
	return string(data[4 : 4+strfunclen])
}

func writeBinaryString(data []byte, obj string) int {
	objlen := len(obj)
	binary.LittleEndian.PutUint32(data[:4], uint32(objlen))
	copy(data[4:4+objlen], obj)
	return 4 + objlen
}

func bool2int(value bool) int {
	if value {
		return 1
	}
	return 0
}

func readBinaryInt(data []byte) int {
	return int(int32(binary.LittleEndian.Uint32(data)))
}

func writeBinaryInt(data []byte, num int) {
	binary.LittleEndian.PutUint32(data, uint32(int32(num)))
}

func readBinaryInt8(data []byte) int8 {
	// 大端模式
	num := int8(0)
	num |= int8(data[0]) << 0
	return num
}

func writeBinaryInt8(data []byte, num int8) {
	// 大端模式
	data[0] = byte(num)
}

func readBinaryBool(data []byte) bool {
	// 大端模式
	num := int8(0)
	num |= int8(data[0]) << 0
	return num > 0
}

func writeBinaryBool(data []byte, num bool) {
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

func writeBinaryUint8(data []byte, num uint8) {
	data[0] = byte(num)
}

func readBinaryUint(data []byte) uint {
	return uint(binary.LittleEndian.Uint32(data))
}

func writeBinaryUint(data []byte, num uint) {
	binary.LittleEndian.PutUint32(data, uint32(num))
}

func writeBinaryFloat32(data []byte, num float32) {
	bits := math.Float32bits(num)
	binary.LittleEndian.PutUint32(data, bits)
}

func readBinaryFloat32(data []byte) float32 {
	bits := binary.LittleEndian.Uint32(data)
	return math.Float32frombits(bits)
}

func writeBinaryFloat64(data []byte, num float64) {
	bits := math.Float64bits(num)
	binary.LittleEndian.PutUint64(data, bits)
}

func readBinaryFloat64(data []byte) float64 {
	bits := binary.LittleEndian.Uint64(data)
	return math.Float64frombits(bits)
}

func ReadMsgModuleInfoByBytes(indata []byte, obj *ModuleInfo) (int, *ModuleInfo) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &ModuleInfo{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.ModuleID) > data__len {
		return endpos, obj
	}
	obj.ModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ModuleID)
	if offset+4+len(obj.ModuleAddr) > data__len {
		return endpos, obj
	}
	obj.ModuleAddr = readBinaryString(data[offset:])
	offset += 4 + len(obj.ModuleAddr)
	if offset+4 > data__len {
		return endpos, obj
	}
	obj.ModuleNumber = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4
	if offset+8 > data__len {
		return endpos, obj
	}
	obj.Version = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8

	return endpos, obj
}

func WriteMsgModuleInfoByObj(data []byte, obj *ModuleInfo) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.ModuleID)
	offset += 4 + len(obj.ModuleID)
	writeBinaryString(data[offset:], obj.ModuleAddr)
	offset += 4 + len(obj.ModuleAddr)
	binary.LittleEndian.PutUint32(data[offset:offset+4], obj.ModuleNumber)
	offset += 4
	binary.LittleEndian.PutUint64(data[offset:offset+8], obj.Version)
	offset += 8

	return offset
}

func GetSizeModuleInfo(obj *ModuleInfo) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + len(obj.ModuleID) + 4 + len(obj.ModuleAddr) + 4 + 8
}

func ReadMsgSTimeTickCommandByBytes(indata []byte, obj *STimeTickCommand) (int, *STimeTickCommand) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &STimeTickCommand{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4 > data__len {
		return endpos, obj
	}
	obj.TestNO = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	return endpos, obj
}

func WriteMsgSTimeTickCommandByObj(data []byte, obj *STimeTickCommand) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	binary.LittleEndian.PutUint32(data[offset:offset+4], obj.TestNO)
	offset += 4

	return offset
}

func GetSizeSTimeTickCommand(obj *STimeTickCommand) int {
	if obj == nil {
		return 4
	}

	return 4 + 4
}

func ReadMsgSTestCommandByBytes(indata []byte, obj *STestCommand) (int, *STestCommand) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &STestCommand{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4 > data__len {
		return endpos, obj
	}
	obj.TestNO = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4
	if offset+4+len(obj.TestString) > data__len {
		return endpos, obj
	}
	obj.TestString = readBinaryString(data[offset:])
	offset += 4 + len(obj.TestString)

	return endpos, obj
}

func WriteMsgSTestCommandByObj(data []byte, obj *STestCommand) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	binary.LittleEndian.PutUint32(data[offset:offset+4], obj.TestNO)
	offset += 4
	writeBinaryString(data[offset:], obj.TestString)
	offset += 4 + len(obj.TestString)

	return offset
}

func GetSizeSTestCommand(obj *STestCommand) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + 4 + len(obj.TestString)
}

func ReadMsgSLoginCommandByBytes(indata []byte, obj *SLoginCommand) (int, *SLoginCommand) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SLoginCommand{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.ModuleID) > data__len {
		return endpos, obj
	}
	obj.ModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ModuleID)
	if offset+4+len(obj.ModuleAddr) > data__len {
		return endpos, obj
	}
	obj.ModuleAddr = readBinaryString(data[offset:])
	offset += 4 + len(obj.ModuleAddr)
	if offset+8 > data__len {
		return endpos, obj
	}
	obj.ConnectPriority = int64(binary.LittleEndian.Uint64(data[offset : offset+8]))
	offset += 8
	if offset+4 > data__len {
		return endpos, obj
	}
	obj.ModuleNumber = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4
	if offset+8 > data__len {
		return endpos, obj
	}
	obj.Version = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8

	return endpos, obj
}

func WriteMsgSLoginCommandByObj(data []byte, obj *SLoginCommand) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.ModuleID)
	offset += 4 + len(obj.ModuleID)
	writeBinaryString(data[offset:], obj.ModuleAddr)
	offset += 4 + len(obj.ModuleAddr)
	binary.LittleEndian.PutUint64(data[offset:offset+8], uint64(obj.ConnectPriority))
	offset += 8
	binary.LittleEndian.PutUint32(data[offset:offset+4], obj.ModuleNumber)
	offset += 4
	binary.LittleEndian.PutUint64(data[offset:offset+8], obj.Version)
	offset += 8

	return offset
}

func GetSizeSLoginCommand(obj *SLoginCommand) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + len(obj.ModuleID) + 4 + len(obj.ModuleAddr) + 8 + 4 +
		8
}

func ReadMsgSLogoutCommandByBytes(indata []byte, obj *SLogoutCommand) (int, *SLogoutCommand) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SLogoutCommand{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize

	return endpos, obj
}

func WriteMsgSLogoutCommandByObj(data []byte, obj *SLogoutCommand) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4

	return offset
}

func GetSizeSLogoutCommand(obj *SLogoutCommand) int {
	if obj == nil {
		return 4
	}

	return 4 + 0
}

func ReadMsgSSeverStartOKCommandByBytes(indata []byte, obj *SSeverStartOKCommand) (int, *SSeverStartOKCommand) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SSeverStartOKCommand{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.ModuleID) > data__len {
		return endpos, obj
	}
	obj.ModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ModuleID)

	return endpos, obj
}

func WriteMsgSSeverStartOKCommandByObj(data []byte, obj *SSeverStartOKCommand) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.ModuleID)
	offset += 4 + len(obj.ModuleID)

	return offset
}

func GetSizeSSeverStartOKCommand(obj *SSeverStartOKCommand) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + len(obj.ModuleID)
}

func ReadMsgSLoginRetCommandByBytes(indata []byte, obj *SLoginRetCommand) (int, *SLoginRetCommand) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SLoginRetCommand{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4 > data__len {
		return endpos, obj
	}
	obj.Loginfailed = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4
	if offset+obj.Destination.GetSize() > data__len {
		return endpos, obj
	}
	rsize_Destination := 0
	rsize_Destination, obj.Destination = ReadMsgModuleInfoByBytes(data[offset:], nil)
	offset += rsize_Destination

	return endpos, obj
}

func WriteMsgSLoginRetCommandByObj(data []byte, obj *SLoginRetCommand) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	binary.LittleEndian.PutUint32(data[offset:offset+4], obj.Loginfailed)
	offset += 4
	offset += WriteMsgModuleInfoByObj(data[offset:], obj.Destination)

	return offset
}

func GetSizeSLoginRetCommand(obj *SLoginRetCommand) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + obj.Destination.GetSize()
}

func ReadMsgSStartRelyNotifyCommandByBytes(indata []byte, obj *SStartRelyNotifyCommand) (int, *SStartRelyNotifyCommand) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SStartRelyNotifyCommand{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4 > data__len {
		return endpos, obj
	}
	ServerInfos_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if ServerInfos_slen != 0xffffffff {
		obj.ServerInfos = make([]*ModuleInfo, ServerInfos_slen)

		for i1i := 0; ServerInfos_slen > i1i; i1i++ {
			rsize_ServerInfos := 0
			rsize_ServerInfos, obj.ServerInfos[i1i] = ReadMsgModuleInfoByBytes(data[offset:], nil)
			offset += rsize_ServerInfos
		}
	}

	return endpos, obj
}

func WriteMsgSStartRelyNotifyCommandByObj(data []byte, obj *SStartRelyNotifyCommand) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	if obj.ServerInfos == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.ServerInfos)))
	}
	offset += 4
	i1i := 0
	ServerInfos_slen := len(obj.ServerInfos)
	for ServerInfos_slen > i1i {
		offset += WriteMsgModuleInfoByObj(data[offset:], obj.ServerInfos[i1i])
		i1i++
	}

	return offset
}

func GetSizeSStartRelyNotifyCommand(obj *SStartRelyNotifyCommand) int {
	if obj == nil {
		return 4
	}
	sizerelyModuleInfo1 := 0
	i1i := 0
	ServerInfos_slen := len(obj.ServerInfos)
	for ServerInfos_slen > i1i {
		sizerelyModuleInfo1 += obj.ServerInfos[i1i].GetSize()
		i1i++
	}

	return 4 + 4 + sizerelyModuleInfo1
}

func ReadMsgSStartMyNotifyCommandByBytes(indata []byte, obj *SStartMyNotifyCommand) (int, *SStartMyNotifyCommand) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SStartMyNotifyCommand{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+obj.ModuleInfo.GetSize() > data__len {
		return endpos, obj
	}
	rsize_ModuleInfo := 0
	rsize_ModuleInfo, obj.ModuleInfo = ReadMsgModuleInfoByBytes(data[offset:], nil)
	offset += rsize_ModuleInfo

	return endpos, obj
}

func WriteMsgSStartMyNotifyCommandByObj(data []byte, obj *SStartMyNotifyCommand) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	offset += WriteMsgModuleInfoByObj(data[offset:], obj.ModuleInfo)

	return offset
}

func GetSizeSStartMyNotifyCommand(obj *SStartMyNotifyCommand) int {
	if obj == nil {
		return 4
	}

	return 4 + obj.ModuleInfo.GetSize()
}

func ReadMsgSNotifyAllInfoByBytes(indata []byte, obj *SNotifyAllInfo) (int, *SNotifyAllInfo) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SNotifyAllInfo{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4 > data__len {
		return endpos, obj
	}
	ServerInfos_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if ServerInfos_slen != 0xffffffff {
		obj.ServerInfos = make([]*ModuleInfo, ServerInfos_slen)

		for i1i := 0; ServerInfos_slen > i1i; i1i++ {
			rsize_ServerInfos := 0
			rsize_ServerInfos, obj.ServerInfos[i1i] = ReadMsgModuleInfoByBytes(data[offset:], nil)
			offset += rsize_ServerInfos
		}
	}

	return endpos, obj
}

func WriteMsgSNotifyAllInfoByObj(data []byte, obj *SNotifyAllInfo) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	if obj.ServerInfos == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.ServerInfos)))
	}
	offset += 4
	i1i := 0
	ServerInfos_slen := len(obj.ServerInfos)
	for ServerInfos_slen > i1i {
		offset += WriteMsgModuleInfoByObj(data[offset:], obj.ServerInfos[i1i])
		i1i++
	}

	return offset
}

func GetSizeSNotifyAllInfo(obj *SNotifyAllInfo) int {
	if obj == nil {
		return 4
	}
	sizerelyModuleInfo1 := 0
	i1i := 0
	ServerInfos_slen := len(obj.ServerInfos)
	for ServerInfos_slen > i1i {
		sizerelyModuleInfo1 += obj.ServerInfos[i1i].GetSize()
		i1i++
	}

	return 4 + 4 + sizerelyModuleInfo1
}

func ReadMsgSNotifySafelyQuitByBytes(indata []byte, obj *SNotifySafelyQuit) (int, *SNotifySafelyQuit) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SNotifySafelyQuit{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+obj.TargetServerInfo.GetSize() > data__len {
		return endpos, obj
	}
	rsize_TargetServerInfo := 0
	rsize_TargetServerInfo, obj.TargetServerInfo = ReadMsgModuleInfoByBytes(data[offset:], nil)
	offset += rsize_TargetServerInfo

	return endpos, obj
}

func WriteMsgSNotifySafelyQuitByObj(data []byte, obj *SNotifySafelyQuit) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	offset += WriteMsgModuleInfoByObj(data[offset:], obj.TargetServerInfo)

	return offset
}

func GetSizeSNotifySafelyQuit(obj *SNotifySafelyQuit) int {
	if obj == nil {
		return 4
	}

	return 4 + obj.TargetServerInfo.GetSize()
}

func ReadMsgSUpdateSessionByBytes(indata []byte, obj *SUpdateSession) (int, *SUpdateSession) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SUpdateSession{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.FromModuleID) > data__len {
		return endpos, obj
	}
	obj.FromModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromModuleID)
	if offset+4+len(obj.ToModuleID) > data__len {
		return endpos, obj
	}
	obj.ToModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToModuleID)
	if offset+4+len(obj.ClientConnID) > data__len {
		return endpos, obj
	}
	obj.ClientConnID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ClientConnID)
	if offset+4+len(obj.SessionUUID) > data__len {
		return endpos, obj
	}
	obj.SessionUUID = readBinaryString(data[offset:])
	offset += 4 + len(obj.SessionUUID)
	if offset+4 > data__len {
		return endpos, obj
	}
	Session_slen := binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4
	if Session_slen != 0xffffffff {
		obj.Session = make(map[string]string)
		for i5i := uint32(0); i5i < Session_slen; i5i++ {
			if offset+0 > data__len {
				return endpos, obj
			}
			keySession := readBinaryString(data[offset:])
			Session_kcatlen := len(keySession)
			offset += Session_kcatlen + 4
			if offset+2 > data__len {
				return endpos, obj
			}
			valueSession := readBinaryString(data[offset:])
			Session_vcatlen := len(valueSession)
			offset += Session_vcatlen + 4
			obj.Session[keySession] = valueSession
		}
	}

	return endpos, obj
}

func WriteMsgSUpdateSessionByObj(data []byte, obj *SUpdateSession) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.FromModuleID)
	offset += 4 + len(obj.FromModuleID)
	writeBinaryString(data[offset:], obj.ToModuleID)
	offset += 4 + len(obj.ToModuleID)
	writeBinaryString(data[offset:], obj.ClientConnID)
	offset += 4 + len(obj.ClientConnID)
	writeBinaryString(data[offset:], obj.SessionUUID)
	offset += 4 + len(obj.SessionUUID)
	if obj.Session == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.Session)))
	}
	offset += 4
	for Sessionkey, Sessionvalue := range obj.Session {
		Session_kcatlen := writeBinaryString(data[offset:], Sessionkey)
		offset += Session_kcatlen
		Session_vcatlen := writeBinaryString(data[offset:], Sessionvalue)
		offset += Session_vcatlen
	}

	return offset
}

func GetSizeSUpdateSession(obj *SUpdateSession) int {
	if obj == nil {
		return 4
	}
	sizerelystring5 := 0
	for Sessionvalue, Sessionkey := range obj.Session {
		sizerelystring5 += len(Sessionvalue) + 4
		sizerelystring5 += len(Sessionkey) + 4
	}

	return 4 + 4 + len(obj.FromModuleID) + 4 + len(obj.ToModuleID) + 4 + len(obj.ClientConnID) + 4 + len(obj.SessionUUID) +
		4 + sizerelystring5
}

func ReadMsgSReqCloseConnectByBytes(indata []byte, obj *SReqCloseConnect) (int, *SReqCloseConnect) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SReqCloseConnect{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.FromModuleID) > data__len {
		return endpos, obj
	}
	obj.FromModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromModuleID)
	if offset+4+len(obj.ToModuleID) > data__len {
		return endpos, obj
	}
	obj.ToModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToModuleID)
	if offset+4+len(obj.ClientConnID) > data__len {
		return endpos, obj
	}
	obj.ClientConnID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ClientConnID)

	return endpos, obj
}

func WriteMsgSReqCloseConnectByObj(data []byte, obj *SReqCloseConnect) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.FromModuleID)
	offset += 4 + len(obj.FromModuleID)
	writeBinaryString(data[offset:], obj.ToModuleID)
	offset += 4 + len(obj.ToModuleID)
	writeBinaryString(data[offset:], obj.ClientConnID)
	offset += 4 + len(obj.ClientConnID)

	return offset
}

func GetSizeSReqCloseConnect(obj *SReqCloseConnect) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + len(obj.FromModuleID) + 4 + len(obj.ToModuleID) + 4 + len(obj.ClientConnID)
}

func ReadMsgSForwardToModuleByBytes(indata []byte, obj *SForwardToModule) (int, *SForwardToModule) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SForwardToModule{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.FromModuleID) > data__len {
		return endpos, obj
	}
	obj.FromModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromModuleID)
	if offset+4+len(obj.ToModuleID) > data__len {
		return endpos, obj
	}
	obj.ToModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToModuleID)
	if offset+2 > data__len {
		return endpos, obj
	}
	obj.MsgID = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2
	if offset+4 > data__len {
		return endpos, obj
	}
	Data_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if Data_slen != 0xffffffff {
		if offset+Data_slen > data__len {
			return endpos, obj
		}
		obj.Data = make([]byte, Data_slen)
		copy(obj.Data, data[offset:offset+Data_slen])
		offset += Data_slen
	}

	return endpos, obj
}

func WriteMsgSForwardToModuleByObj(data []byte, obj *SForwardToModule) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.FromModuleID)
	offset += 4 + len(obj.FromModuleID)
	writeBinaryString(data[offset:], obj.ToModuleID)
	offset += 4 + len(obj.ToModuleID)
	binary.LittleEndian.PutUint16(data[offset:offset+2], obj.MsgID)
	offset += 2
	if obj.Data == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.Data)))
	}
	offset += 4
	Data_slen := len(obj.Data)
	copy(data[offset:offset+Data_slen], obj.Data)
	offset += Data_slen

	return offset
}

func GetSizeSForwardToModule(obj *SForwardToModule) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + len(obj.FromModuleID) + 4 + len(obj.ToModuleID) + 2 + 4 + len(obj.Data)*1
}

func ReadMsgModuleMessageByBytes(indata []byte, obj *ModuleMessage) (int, *ModuleMessage) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &ModuleMessage{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+obj.FromModule.GetSize() > data__len {
		return endpos, obj
	}
	rsize_FromModule := 0
	rsize_FromModule, obj.FromModule = ReadMsgModuleInfoByBytes(data[offset:], nil)
	offset += rsize_FromModule
	if offset+2 > data__len {
		return endpos, obj
	}
	obj.MsgID = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2
	if offset+4 > data__len {
		return endpos, obj
	}
	Data_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if Data_slen != 0xffffffff {
		if offset+Data_slen > data__len {
			return endpos, obj
		}
		obj.Data = make([]byte, Data_slen)
		copy(obj.Data, data[offset:offset+Data_slen])
		offset += Data_slen
	}

	return endpos, obj
}

func WriteMsgModuleMessageByObj(data []byte, obj *ModuleMessage) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	offset += WriteMsgModuleInfoByObj(data[offset:], obj.FromModule)
	binary.LittleEndian.PutUint16(data[offset:offset+2], obj.MsgID)
	offset += 2
	if obj.Data == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.Data)))
	}
	offset += 4
	Data_slen := len(obj.Data)
	copy(data[offset:offset+Data_slen], obj.Data)
	offset += Data_slen

	return offset
}

func GetSizeModuleMessage(obj *ModuleMessage) int {
	if obj == nil {
		return 4
	}

	return 4 + obj.FromModule.GetSize() + 2 + 4 + len(obj.Data)*1
}

func ReadMsgSForwardToClientByBytes(indata []byte, obj *SForwardToClient) (int, *SForwardToClient) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SForwardToClient{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.FromModuleID) > data__len {
		return endpos, obj
	}
	obj.FromModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromModuleID)
	if offset+4+len(obj.ToGateID) > data__len {
		return endpos, obj
	}
	obj.ToGateID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToGateID)
	if offset+4+len(obj.ToClientID) > data__len {
		return endpos, obj
	}
	obj.ToClientID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToClientID)
	if offset+2 > data__len {
		return endpos, obj
	}
	obj.MsgID = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2
	if offset+4 > data__len {
		return endpos, obj
	}
	Data_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if Data_slen != 0xffffffff {
		if offset+Data_slen > data__len {
			return endpos, obj
		}
		obj.Data = make([]byte, Data_slen)
		copy(obj.Data, data[offset:offset+Data_slen])
		offset += Data_slen
	}

	return endpos, obj
}

func WriteMsgSForwardToClientByObj(data []byte, obj *SForwardToClient) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.FromModuleID)
	offset += 4 + len(obj.FromModuleID)
	writeBinaryString(data[offset:], obj.ToGateID)
	offset += 4 + len(obj.ToGateID)
	writeBinaryString(data[offset:], obj.ToClientID)
	offset += 4 + len(obj.ToClientID)
	binary.LittleEndian.PutUint16(data[offset:offset+2], obj.MsgID)
	offset += 2
	if obj.Data == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.Data)))
	}
	offset += 4
	Data_slen := len(obj.Data)
	copy(data[offset:offset+Data_slen], obj.Data)
	offset += Data_slen

	return offset
}

func GetSizeSForwardToClient(obj *SForwardToClient) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + len(obj.FromModuleID) + 4 + len(obj.ToGateID) + 4 + len(obj.ToClientID) + 2 +
		4 + len(obj.Data)*1
}

func ReadMsgSForwardFromGateByBytes(indata []byte, obj *SForwardFromGate) (int, *SForwardFromGate) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SForwardFromGate{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.FromModuleID) > data__len {
		return endpos, obj
	}
	obj.FromModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromModuleID)
	if offset+4+len(obj.ToModuleID) > data__len {
		return endpos, obj
	}
	obj.ToModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToModuleID)
	if offset+4+len(obj.ClientConnID) > data__len {
		return endpos, obj
	}
	obj.ClientConnID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ClientConnID)
	if offset+4 > data__len {
		return endpos, obj
	}
	Session_slen := binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4
	if Session_slen != 0xffffffff {
		obj.Session = make(map[string]string)
		for i4i := uint32(0); i4i < Session_slen; i4i++ {
			if offset+0 > data__len {
				return endpos, obj
			}
			keySession := readBinaryString(data[offset:])
			Session_kcatlen := len(keySession)
			offset += Session_kcatlen + 4
			if offset+2 > data__len {
				return endpos, obj
			}
			valueSession := readBinaryString(data[offset:])
			Session_vcatlen := len(valueSession)
			offset += Session_vcatlen + 4
			obj.Session[keySession] = valueSession
		}
	}
	if offset+2 > data__len {
		return endpos, obj
	}
	obj.MsgID = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2
	if offset+4 > data__len {
		return endpos, obj
	}
	Data_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if Data_slen != 0xffffffff {
		if offset+Data_slen > data__len {
			return endpos, obj
		}
		obj.Data = make([]byte, Data_slen)
		copy(obj.Data, data[offset:offset+Data_slen])
		offset += Data_slen
	}

	return endpos, obj
}

func WriteMsgSForwardFromGateByObj(data []byte, obj *SForwardFromGate) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.FromModuleID)
	offset += 4 + len(obj.FromModuleID)
	writeBinaryString(data[offset:], obj.ToModuleID)
	offset += 4 + len(obj.ToModuleID)
	writeBinaryString(data[offset:], obj.ClientConnID)
	offset += 4 + len(obj.ClientConnID)
	if obj.Session == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.Session)))
	}
	offset += 4
	for Sessionkey, Sessionvalue := range obj.Session {
		Session_kcatlen := writeBinaryString(data[offset:], Sessionkey)
		offset += Session_kcatlen
		Session_vcatlen := writeBinaryString(data[offset:], Sessionvalue)
		offset += Session_vcatlen
	}
	binary.LittleEndian.PutUint16(data[offset:offset+2], obj.MsgID)
	offset += 2
	if obj.Data == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.Data)))
	}
	offset += 4
	Data_slen := len(obj.Data)
	copy(data[offset:offset+Data_slen], obj.Data)
	offset += Data_slen

	return offset
}

func GetSizeSForwardFromGate(obj *SForwardFromGate) int {
	if obj == nil {
		return 4
	}
	sizerelystring4 := 0
	for Sessionvalue, Sessionkey := range obj.Session {
		sizerelystring4 += len(Sessionvalue) + 4
		sizerelystring4 += len(Sessionkey) + 4
	}

	return 4 + 4 + len(obj.FromModuleID) + 4 + len(obj.ToModuleID) + 4 + len(obj.ClientConnID) + 4 + sizerelystring4 +
		2 + 4 + len(obj.Data)*1
}

func ReadMsgClientMessageByBytes(indata []byte, obj *ClientMessage) (int, *ClientMessage) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &ClientMessage{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+obj.FromModule.GetSize() > data__len {
		return endpos, obj
	}
	rsize_FromModule := 0
	rsize_FromModule, obj.FromModule = ReadMsgModuleInfoByBytes(data[offset:], nil)
	offset += rsize_FromModule
	if offset+4+len(obj.ClientConnID) > data__len {
		return endpos, obj
	}
	obj.ClientConnID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ClientConnID)
	if offset+2 > data__len {
		return endpos, obj
	}
	obj.MsgID = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2
	if offset+4 > data__len {
		return endpos, obj
	}
	Data_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if Data_slen != 0xffffffff {
		if offset+Data_slen > data__len {
			return endpos, obj
		}
		obj.Data = make([]byte, Data_slen)
		copy(obj.Data, data[offset:offset+Data_slen])
		offset += Data_slen
	}

	return endpos, obj
}

func WriteMsgClientMessageByObj(data []byte, obj *ClientMessage) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	offset += WriteMsgModuleInfoByObj(data[offset:], obj.FromModule)
	writeBinaryString(data[offset:], obj.ClientConnID)
	offset += 4 + len(obj.ClientConnID)
	binary.LittleEndian.PutUint16(data[offset:offset+2], obj.MsgID)
	offset += 2
	if obj.Data == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.Data)))
	}
	offset += 4
	Data_slen := len(obj.Data)
	copy(data[offset:offset+Data_slen], obj.Data)
	offset += Data_slen

	return offset
}

func GetSizeClientMessage(obj *ClientMessage) int {
	if obj == nil {
		return 4
	}

	return 4 + obj.FromModule.GetSize() + 4 + len(obj.ClientConnID) + 2 + 4 + len(obj.Data)*1
}

func ReadMsgSROCRequestByBytes(indata []byte, obj *SROCRequest) (int, *SROCRequest) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SROCRequest{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.FromModuleID) > data__len {
		return endpos, obj
	}
	obj.FromModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromModuleID)
	if offset+4+len(obj.ToModuleID) > data__len {
		return endpos, obj
	}
	obj.ToModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToModuleID)
	if offset+8 > data__len {
		return endpos, obj
	}
	obj.Seq = int64(binary.LittleEndian.Uint64(data[offset : offset+8]))
	offset += 8
	if offset+4+len(obj.CallStr) > data__len {
		return endpos, obj
	}
	obj.CallStr = readBinaryString(data[offset:])
	offset += 4 + len(obj.CallStr)
	if offset+4 > data__len {
		return endpos, obj
	}
	CallArg_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if CallArg_slen != 0xffffffff {
		if offset+CallArg_slen > data__len {
			return endpos, obj
		}
		obj.CallArg = make([]byte, CallArg_slen)
		copy(obj.CallArg, data[offset:offset+CallArg_slen])
		offset += CallArg_slen
	}
	if offset+1 > data__len {
		return endpos, obj
	}
	obj.NeedReturn = uint8(data[offset]) != 0
	offset += 1

	return endpos, obj
}

func WriteMsgSROCRequestByObj(data []byte, obj *SROCRequest) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.FromModuleID)
	offset += 4 + len(obj.FromModuleID)
	writeBinaryString(data[offset:], obj.ToModuleID)
	offset += 4 + len(obj.ToModuleID)
	binary.LittleEndian.PutUint64(data[offset:offset+8], uint64(obj.Seq))
	offset += 8
	writeBinaryString(data[offset:], obj.CallStr)
	offset += 4 + len(obj.CallStr)
	if obj.CallArg == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.CallArg)))
	}
	offset += 4
	CallArg_slen := len(obj.CallArg)
	copy(data[offset:offset+CallArg_slen], obj.CallArg)
	offset += CallArg_slen
	data[offset] = uint8(bool2int(obj.NeedReturn))
	offset += 1

	return offset
}

func GetSizeSROCRequest(obj *SROCRequest) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + len(obj.FromModuleID) + 4 + len(obj.ToModuleID) + 8 + 4 + len(obj.CallStr) +
		4 + len(obj.CallArg)*1 + 1
}

func ReadMsgSROCResponseByBytes(indata []byte, obj *SROCResponse) (int, *SROCResponse) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SROCResponse{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.FromModuleID) > data__len {
		return endpos, obj
	}
	obj.FromModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromModuleID)
	if offset+4+len(obj.ToModuleID) > data__len {
		return endpos, obj
	}
	obj.ToModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToModuleID)
	if offset+8 > data__len {
		return endpos, obj
	}
	obj.ReqSeq = int64(binary.LittleEndian.Uint64(data[offset : offset+8]))
	offset += 8
	if offset+4 > data__len {
		return endpos, obj
	}
	ResData_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if ResData_slen != 0xffffffff {
		if offset+ResData_slen > data__len {
			return endpos, obj
		}
		obj.ResData = make([]byte, ResData_slen)
		copy(obj.ResData, data[offset:offset+ResData_slen])
		offset += ResData_slen
	}
	if offset+4+len(obj.Error) > data__len {
		return endpos, obj
	}
	obj.Error = readBinaryString(data[offset:])
	offset += 4 + len(obj.Error)

	return endpos, obj
}

func WriteMsgSROCResponseByObj(data []byte, obj *SROCResponse) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.FromModuleID)
	offset += 4 + len(obj.FromModuleID)
	writeBinaryString(data[offset:], obj.ToModuleID)
	offset += 4 + len(obj.ToModuleID)
	binary.LittleEndian.PutUint64(data[offset:offset+8], uint64(obj.ReqSeq))
	offset += 8
	if obj.ResData == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.ResData)))
	}
	offset += 4
	ResData_slen := len(obj.ResData)
	copy(data[offset:offset+ResData_slen], obj.ResData)
	offset += ResData_slen
	writeBinaryString(data[offset:], obj.Error)
	offset += 4 + len(obj.Error)

	return offset
}

func GetSizeSROCResponse(obj *SROCResponse) int {
	if obj == nil {
		return 4
	}

	return 4 + 4 + len(obj.FromModuleID) + 4 + len(obj.ToModuleID) + 8 + 4 + len(obj.ResData)*1 +
		4 + len(obj.Error)
}

func ReadMsgSROCBindByBytes(indata []byte, obj *SROCBind) (int, *SROCBind) {
	offset := 0
	if len(indata) < 4 {
		return 0, nil
	}
	objsize := int(binary.LittleEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4, nil
	}
	if obj == nil {
		obj = &SROCBind{}
	}
	if offset+objsize > len(indata) {
		return offset, obj
	}
	endpos := offset + objsize
	data := indata[offset : offset+objsize]
	offset = 0
	data__len := len(data)
	if offset+4+len(obj.HostModuleID) > data__len {
		return endpos, obj
	}
	obj.HostModuleID = readBinaryString(data[offset:])
	offset += 4 + len(obj.HostModuleID)
	if offset+1 > data__len {
		return endpos, obj
	}
	obj.IsDelete = uint8(data[offset]) != 0
	offset += 1
	if offset+4+len(obj.ObjType) > data__len {
		return endpos, obj
	}
	obj.ObjType = readBinaryString(data[offset:])
	offset += 4 + len(obj.ObjType)
	if offset+4 > data__len {
		return endpos, obj
	}
	ObjIDs_slen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if ObjIDs_slen != 0xffffffff {
		if offset+(ObjIDs_slen*0) > data__len {
			return endpos, obj
		}
		obj.ObjIDs = make([]string, ObjIDs_slen)
		for i4i := 0; ObjIDs_slen > i4i; i4i++ {
			obj.ObjIDs[i4i] = string(readBinaryString(data[offset:]))
			offset += 0
		}
	}

	return endpos, obj
}

func WriteMsgSROCBindByObj(data []byte, obj *SROCBind) int {
	if obj == nil {
		binary.LittleEndian.PutUint32(data[0:4], 0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:], obj.HostModuleID)
	offset += 4 + len(obj.HostModuleID)
	data[offset] = uint8(bool2int(obj.IsDelete))
	offset += 1
	writeBinaryString(data[offset:], obj.ObjType)
	offset += 4 + len(obj.ObjType)
	if obj.ObjIDs == nil {
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xffffffff)
	} else {
		binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(obj.ObjIDs)))
	}
	offset += 4
	ObjIDs_slen := len(obj.ObjIDs)
	for i4i := 0; ObjIDs_slen > i4i; i4i++ {
		writeBinaryString(data[offset:offset+0], obj.ObjIDs[i4i])
		offset += 0
	}

	return offset
}

func GetSizeSROCBind(obj *SROCBind) int {
	if obj == nil {
		return 4
	}
	sizerelystring4 := 0
	i4i := 0
	ObjIDs_slen := len(obj.ObjIDs)
	for ObjIDs_slen > i4i {
		sizerelystring4 += len(obj.ObjIDs[i4i]) + 4
		i4i++
	}

	return 4 + 4 + len(obj.HostModuleID) + 1 + 4 + len(obj.ObjType) + 4 + sizerelystring4
}
