package servercomm
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
	SNotifySafelyQuitID = 46
	SUpdateSessionID = 47
	SForwardToServerID = 48
	SForwardToClientID = 49
	SForwardFromGateID = 50
)
const (
	SServerInfoName = "servercomm.SServerInfo"
	STimeTickCommandName = "servercomm.STimeTickCommand"
	STestCommandName = "servercomm.STestCommand"
	SLoginCommandName = "servercomm.SLoginCommand"
	SLogoutCommandName = "servercomm.SLogoutCommand"
	SSeverStartOKCommandName = "servercomm.SSeverStartOKCommand"
	SLoginRetCommandName = "servercomm.SLoginRetCommand"
	SStartRelyNotifyCommandName = "servercomm.SStartRelyNotifyCommand"
	SStartMyNotifyCommandName = "servercomm.SStartMyNotifyCommand"
	SNotifyAllInfoName = "servercomm.SNotifyAllInfo"
	SNotifySafelyQuitName = "servercomm.SNotifySafelyQuit"
	SUpdateSessionName = "servercomm.SUpdateSession"
	SForwardToServerName = "servercomm.SForwardToServer"
	SForwardToClientName = "servercomm.SForwardToClient"
	SForwardFromGateName = "servercomm.SForwardFromGate"
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
func (this *SNotifySafelyQuit) WriteBinary(data []byte) int {
	return WriteMsgSNotifySafelyQuitByObj(data,this)
}
func (this *SUpdateSession) WriteBinary(data []byte) int {
	return WriteMsgSUpdateSessionByObj(data,this)
}
func (this *SForwardToServer) WriteBinary(data []byte) int {
	return WriteMsgSForwardToServerByObj(data,this)
}
func (this *SForwardToClient) WriteBinary(data []byte) int {
	return WriteMsgSForwardToClientByObj(data,this)
}
func (this *SForwardFromGate) WriteBinary(data []byte) int {
	return WriteMsgSForwardFromGateByObj(data,this)
}
func (this *SServerInfo) ReadBinary(data []byte) int {
	size,_ := ReadMsgSServerInfoByBytes(data, this)
	return size
}
func (this *STimeTickCommand) ReadBinary(data []byte) int {
	size,_ := ReadMsgSTimeTickCommandByBytes(data, this)
	return size
}
func (this *STestCommand) ReadBinary(data []byte) int {
	size,_ := ReadMsgSTestCommandByBytes(data, this)
	return size
}
func (this *SLoginCommand) ReadBinary(data []byte) int {
	size,_ := ReadMsgSLoginCommandByBytes(data, this)
	return size
}
func (this *SLogoutCommand) ReadBinary(data []byte) int {
	size,_ := ReadMsgSLogoutCommandByBytes(data, this)
	return size
}
func (this *SSeverStartOKCommand) ReadBinary(data []byte) int {
	size,_ := ReadMsgSSeverStartOKCommandByBytes(data, this)
	return size
}
func (this *SLoginRetCommand) ReadBinary(data []byte) int {
	size,_ := ReadMsgSLoginRetCommandByBytes(data, this)
	return size
}
func (this *SStartRelyNotifyCommand) ReadBinary(data []byte) int {
	size,_ := ReadMsgSStartRelyNotifyCommandByBytes(data, this)
	return size
}
func (this *SStartMyNotifyCommand) ReadBinary(data []byte) int {
	size,_ := ReadMsgSStartMyNotifyCommandByBytes(data, this)
	return size
}
func (this *SNotifyAllInfo) ReadBinary(data []byte) int {
	size,_ := ReadMsgSNotifyAllInfoByBytes(data, this)
	return size
}
func (this *SNotifySafelyQuit) ReadBinary(data []byte) int {
	size,_ := ReadMsgSNotifySafelyQuitByBytes(data, this)
	return size
}
func (this *SUpdateSession) ReadBinary(data []byte) int {
	size,_ := ReadMsgSUpdateSessionByBytes(data, this)
	return size
}
func (this *SForwardToServer) ReadBinary(data []byte) int {
	size,_ := ReadMsgSForwardToServerByBytes(data, this)
	return size
}
func (this *SForwardToClient) ReadBinary(data []byte) int {
	size,_ := ReadMsgSForwardToClientByBytes(data, this)
	return size
}
func (this *SForwardFromGate) ReadBinary(data []byte) int {
	size,_ := ReadMsgSForwardFromGateByBytes(data, this)
	return size
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
		case SNotifySafelyQuitID: 
		return SNotifySafelyQuitName
		case SUpdateSessionID: 
		return SUpdateSessionName
		case SForwardToServerID: 
		return SForwardToServerName
		case SForwardToClientID: 
		return SForwardToClientName
		case SForwardFromGateID: 
		return SForwardFromGateName
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
		case SNotifySafelyQuitName: 
		return SNotifySafelyQuitID
		case SUpdateSessionName: 
		return SUpdateSessionID
		case SForwardToServerName: 
		return SForwardToServerID
		case SForwardToClientName: 
		return SForwardToClientID
		case SForwardFromGateName: 
		return SForwardFromGateID
		default:
		return 0
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
func (this *SNotifySafelyQuit) GetMsgId() uint16 {
	return SNotifySafelyQuitID
}
func (this *SUpdateSession) GetMsgId() uint16 {
	return SUpdateSessionID
}
func (this *SForwardToServer) GetMsgId() uint16 {
	return SForwardToServerID
}
func (this *SForwardToClient) GetMsgId() uint16 {
	return SForwardToClientID
}
func (this *SForwardFromGate) GetMsgId() uint16 {
	return SForwardFromGateID
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
func (this *SNotifySafelyQuit) GetMsgName() string {
	return SNotifySafelyQuitName
}
func (this *SUpdateSession) GetMsgName() string {
	return SUpdateSessionName
}
func (this *SForwardToServer) GetMsgName() string {
	return SForwardToServerName
}
func (this *SForwardToClient) GetMsgName() string {
	return SForwardToClientName
}
func (this *SForwardFromGate) GetMsgName() string {
	return SForwardFromGateName
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
func (this *SNotifySafelyQuit) GetSize() int {
	return GetSizeSNotifySafelyQuit(this)
}
func (this *SUpdateSession) GetSize() int {
	return GetSizeSUpdateSession(this)
}
func (this *SForwardToServer) GetSize() int {
	return GetSizeSForwardToServer(this)
}
func (this *SForwardToClient) GetSize() int {
	return GetSizeSForwardToClient(this)
}
func (this *SForwardFromGate) GetSize() int {
	return GetSizeSForwardFromGate(this)
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
func (this *SNotifySafelyQuit) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SUpdateSession) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SForwardToServer) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SForwardToClient) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func (this *SForwardFromGate) GetJson() string {
	json,_ := json.Marshal(this)
	return string(json)
}
func readBinaryString(data []byte) string {
	strfunclen := binary.BigEndian.Uint32(data[:4])
	if int(strfunclen) + 4 > len(data ) {
		return ""
	}
	return string(data[4:4+strfunclen])
}
func writeBinaryString(data []byte,obj string) int {
	objlen := len(obj)
	binary.BigEndian.PutUint32(data[:4],uint32(objlen))
	copy(data[4:4+objlen], obj)
	return 4+objlen
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
func readBinaryUint(data []byte) uint {
	return uint(binary.BigEndian.Uint32(data))
}
func writeBinaryUint(data []byte, num uint ) {
	binary.BigEndian.PutUint32(data,uint32(num))
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
func ReadMsgSServerInfoByBytes(indata []byte, obj *SServerInfo) (int,*SServerInfo ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SServerInfo{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 + len(obj.ServerID) > data__len{
		return endpos,obj
	}
	obj.ServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ServerID)
	if offset + 4 + len(obj.ServerAddr) > data__len{
		return endpos,obj
	}
	obj.ServerAddr = readBinaryString(data[offset:])
	offset += 4 + len(obj.ServerAddr)
	if offset + 4 > data__len{
		return endpos,obj
	}
	obj.ServerNumber = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos,obj
	}
	obj.Version = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos,obj
}
func WriteMsgSServerInfoByObj(data []byte, obj *SServerInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:],obj.ServerID)
	offset += 4 + len(obj.ServerID)
	writeBinaryString(data[offset:],obj.ServerAddr)
	offset += 4 + len(obj.ServerAddr)
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.ServerNumber)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Version)
	offset+=8
	return offset
}
func GetSizeSServerInfo(obj *SServerInfo) int {
	if obj == nil {
		return 4
	}
	return 4 + 4 + len(obj.ServerID) + 4 + len(obj.ServerAddr) + 4 + 8
}
func ReadMsgSTimeTickCommandByBytes(indata []byte, obj *STimeTickCommand) (int,*STimeTickCommand ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&STimeTickCommand{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos,obj
	}
	obj.Testno = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos,obj
}
func WriteMsgSTimeTickCommandByObj(data []byte, obj *STimeTickCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Testno)
	offset+=4
	return offset
}
func GetSizeSTimeTickCommand(obj *STimeTickCommand) int {
	if obj == nil {
		return 4
	}
	return 4 + 4
}
func ReadMsgSTestCommandByBytes(indata []byte, obj *STestCommand) (int,*STestCommand ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&STestCommand{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos,obj
	}
	obj.Testno = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 4 + len(obj.Testttring) > data__len{
		return endpos,obj
	}
	obj.Testttring = readBinaryString(data[offset:])
	offset += 4 + len(obj.Testttring)
	return endpos,obj
}
func WriteMsgSTestCommandByObj(data []byte, obj *STestCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Testno)
	offset+=4
	writeBinaryString(data[offset:],obj.Testttring)
	offset += 4 + len(obj.Testttring)
	return offset
}
func GetSizeSTestCommand(obj *STestCommand) int {
	if obj == nil {
		return 4
	}
	return 4 + 4 + 4 + len(obj.Testttring)
}
func ReadMsgSLoginCommandByBytes(indata []byte, obj *SLoginCommand) (int,*SLoginCommand ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SLoginCommand{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 + len(obj.ServerID) > data__len{
		return endpos,obj
	}
	obj.ServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ServerID)
	if offset + 4 + len(obj.ServerAddr) > data__len{
		return endpos,obj
	}
	obj.ServerAddr = readBinaryString(data[offset:])
	offset += 4 + len(obj.ServerAddr)
	if offset + 8 > data__len{
		return endpos,obj
	}
	obj.ConnectPriority = readBinaryInt64(data[offset:offset+8])
	offset+=8
	if offset + 4 > data__len{
		return endpos,obj
	}
	obj.ServerNumber = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + 8 > data__len{
		return endpos,obj
	}
	obj.Version = binary.BigEndian.Uint64(data[offset:offset+8])
	offset+=8
	return endpos,obj
}
func WriteMsgSLoginCommandByObj(data []byte, obj *SLoginCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:],obj.ServerID)
	offset += 4 + len(obj.ServerID)
	writeBinaryString(data[offset:],obj.ServerAddr)
	offset += 4 + len(obj.ServerAddr)
	writeBinaryInt64(data[offset:offset+8], obj.ConnectPriority)
	offset+=8
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.ServerNumber)
	offset+=4
	binary.BigEndian.PutUint64(data[offset:offset+8], obj.Version)
	offset+=8
	return offset
}
func GetSizeSLoginCommand(obj *SLoginCommand) int {
	if obj == nil {
		return 4
	}
	return 4 + 4 + len(obj.ServerID) + 4 + len(obj.ServerAddr) + 8 + 4 + 
	8
}
func ReadMsgSLogoutCommandByBytes(indata []byte, obj *SLogoutCommand) (int,*SLogoutCommand ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SLogoutCommand{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	return endpos,obj
}
func WriteMsgSLogoutCommandByObj(data []byte, obj *SLogoutCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	return offset
}
func GetSizeSLogoutCommand(obj *SLogoutCommand) int {
	if obj == nil {
		return 4
	}
	return 4 + 0
}
func ReadMsgSSeverStartOKCommandByBytes(indata []byte, obj *SSeverStartOKCommand) (int,*SSeverStartOKCommand ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SSeverStartOKCommand{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos,obj
	}
	obj.Serverid = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	return endpos,obj
}
func WriteMsgSSeverStartOKCommandByObj(data []byte, obj *SSeverStartOKCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Serverid)
	offset+=4
	return offset
}
func GetSizeSSeverStartOKCommand(obj *SSeverStartOKCommand) int {
	if obj == nil {
		return 4
	}
	return 4 + 4
}
func ReadMsgSLoginRetCommandByBytes(indata []byte, obj *SLoginRetCommand) (int,*SLoginRetCommand ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SLoginRetCommand{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos,obj
	}
	obj.Loginfailed = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	if offset + obj.Destination.GetSize() > data__len{
		return endpos,obj
	}
	rsize_Destination := 0
	rsize_Destination,obj.Destination = ReadMsgSServerInfoByBytes(data[offset:], nil)
	offset += rsize_Destination
	return endpos,obj
}
func WriteMsgSLoginRetCommandByObj(data []byte, obj *SLoginRetCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	binary.BigEndian.PutUint32(data[offset:offset+4], obj.Loginfailed)
	offset+=4
	offset += WriteMsgSServerInfoByObj(data[offset:], obj.Destination)
	return offset
}
func GetSizeSLoginRetCommand(obj *SLoginRetCommand) int {
	if obj == nil {
		return 4
	}
	return 4 + 4 + obj.Destination.GetSize()
}
func ReadMsgSStartRelyNotifyCommandByBytes(indata []byte, obj *SStartRelyNotifyCommand) (int,*SStartRelyNotifyCommand ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SStartRelyNotifyCommand{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos,obj
	}
	Serverinfos_slen := int(binary.BigEndian.Uint32(data[offset:offset+4]))
	offset += 4
	if Serverinfos_slen != 0xffffffff {
		obj.Serverinfos = make([]*SServerInfo,Serverinfos_slen)
		for i1i := 0; Serverinfos_slen > i1i; i1i++ {
			rsize_Serverinfos := 0
			rsize_Serverinfos,obj.Serverinfos[i1i] = ReadMsgSServerInfoByBytes(data[offset:],nil)
			offset += rsize_Serverinfos
		}
	}
	return endpos,obj
}
func WriteMsgSStartRelyNotifyCommandByObj(data []byte, obj *SStartRelyNotifyCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	if obj.Serverinfos == nil {
		binary.BigEndian.PutUint32(data[offset:offset+4],0xffffffff)
	} else {
		binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Serverinfos)))
	}
	offset += 4
	i1i := 0
	Serverinfos_slen := len(obj.Serverinfos)
	for Serverinfos_slen > i1i {
		offset += WriteMsgSServerInfoByObj(data[offset:],obj.Serverinfos[i1i])
		i1i++
	}
	return offset
}
func GetSizeSStartRelyNotifyCommand(obj *SStartRelyNotifyCommand) int {
	if obj == nil {
		return 4
	}
	sizerelySServerInfo1 := 0
	i1i := 0
	Serverinfos_slen := len(obj.Serverinfos)
	for Serverinfos_slen > i1i {
		sizerelySServerInfo1 += obj.Serverinfos[i1i].GetSize()
		i1i++
	}
	return 4 + 4 + sizerelySServerInfo1
}
func ReadMsgSStartMyNotifyCommandByBytes(indata []byte, obj *SStartMyNotifyCommand) (int,*SStartMyNotifyCommand ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SStartMyNotifyCommand{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + obj.Serverinfo.GetSize() > data__len{
		return endpos,obj
	}
	rsize_Serverinfo := 0
	rsize_Serverinfo,obj.Serverinfo = ReadMsgSServerInfoByBytes(data[offset:], nil)
	offset += rsize_Serverinfo
	return endpos,obj
}
func WriteMsgSStartMyNotifyCommandByObj(data []byte, obj *SStartMyNotifyCommand) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	offset += WriteMsgSServerInfoByObj(data[offset:], obj.Serverinfo)
	return offset
}
func GetSizeSStartMyNotifyCommand(obj *SStartMyNotifyCommand) int {
	if obj == nil {
		return 4
	}
	return 4 + obj.Serverinfo.GetSize()
}
func ReadMsgSNotifyAllInfoByBytes(indata []byte, obj *SNotifyAllInfo) (int,*SNotifyAllInfo ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SNotifyAllInfo{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 > data__len{
		return endpos,obj
	}
	Serverinfos_slen := int(binary.BigEndian.Uint32(data[offset:offset+4]))
	offset += 4
	if Serverinfos_slen != 0xffffffff {
		obj.Serverinfos = make([]*SServerInfo,Serverinfos_slen)
		for i1i := 0; Serverinfos_slen > i1i; i1i++ {
			rsize_Serverinfos := 0
			rsize_Serverinfos,obj.Serverinfos[i1i] = ReadMsgSServerInfoByBytes(data[offset:],nil)
			offset += rsize_Serverinfos
		}
	}
	return endpos,obj
}
func WriteMsgSNotifyAllInfoByObj(data []byte, obj *SNotifyAllInfo) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	if obj.Serverinfos == nil {
		binary.BigEndian.PutUint32(data[offset:offset+4],0xffffffff)
	} else {
		binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Serverinfos)))
	}
	offset += 4
	i1i := 0
	Serverinfos_slen := len(obj.Serverinfos)
	for Serverinfos_slen > i1i {
		offset += WriteMsgSServerInfoByObj(data[offset:],obj.Serverinfos[i1i])
		i1i++
	}
	return offset
}
func GetSizeSNotifyAllInfo(obj *SNotifyAllInfo) int {
	if obj == nil {
		return 4
	}
	sizerelySServerInfo1 := 0
	i1i := 0
	Serverinfos_slen := len(obj.Serverinfos)
	for Serverinfos_slen > i1i {
		sizerelySServerInfo1 += obj.Serverinfos[i1i].GetSize()
		i1i++
	}
	return 4 + 4 + sizerelySServerInfo1
}
func ReadMsgSNotifySafelyQuitByBytes(indata []byte, obj *SNotifySafelyQuit) (int,*SNotifySafelyQuit ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SNotifySafelyQuit{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + obj.TargetServerInfo.GetSize() > data__len{
		return endpos,obj
	}
	rsize_TargetServerInfo := 0
	rsize_TargetServerInfo,obj.TargetServerInfo = ReadMsgSServerInfoByBytes(data[offset:], nil)
	offset += rsize_TargetServerInfo
	return endpos,obj
}
func WriteMsgSNotifySafelyQuitByObj(data []byte, obj *SNotifySafelyQuit) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	offset += WriteMsgSServerInfoByObj(data[offset:], obj.TargetServerInfo)
	return offset
}
func GetSizeSNotifySafelyQuit(obj *SNotifySafelyQuit) int {
	if obj == nil {
		return 4
	}
	return 4 + obj.TargetServerInfo.GetSize()
}
func ReadMsgSUpdateSessionByBytes(indata []byte, obj *SUpdateSession) (int,*SUpdateSession ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SUpdateSession{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 + len(obj.ClientConnID) > data__len{
		return endpos,obj
	}
	obj.ClientConnID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ClientConnID)
	if offset + 4 > data__len{
		return endpos,obj
	}
	Session_slen := binary.BigEndian.Uint32(data[offset:offset+4])
	offset += 4
	if Session_slen != 0xffffffff {
		obj.Session = make(map[string]string)
		for i2i := uint32(0); i2i < Session_slen; i2i++ {
			if offset + 0 > data__len{
				return endpos,obj
			}
			keySession := readBinaryString(data[offset:])
			Session_kcatlen := len(keySession)
			offset += Session_kcatlen + 4
			if offset + 2 > data__len{
				return endpos,obj
			}
			valueSession := readBinaryString(data[offset:])
			Session_vcatlen := len(valueSession)
			offset += Session_vcatlen + 4
			obj.Session[keySession] = valueSession
		}
	}
	return endpos,obj
}
func WriteMsgSUpdateSessionByObj(data []byte, obj *SUpdateSession) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:],obj.ClientConnID)
	offset += 4 + len(obj.ClientConnID)
	if obj.Session == nil {
		binary.BigEndian.PutUint32(data[offset:offset+4],0xffffffff)
	} else {
		binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Session)))
	}
	offset += 4
	for Sessionkey,Sessionvalue := range obj.Session {
		Session_kcatlen := writeBinaryString(data[offset:],Sessionkey)
		offset += Session_kcatlen
		Session_vcatlen := writeBinaryString(data[offset:],Sessionvalue)
		offset += Session_vcatlen
	}
	return offset
}
func GetSizeSUpdateSession(obj *SUpdateSession) int {
	if obj == nil {
		return 4
	}
	sizerelystring2 := 0
	for Sessionvalue,Sessionkey := range obj.Session {
		sizerelystring2 += len(Sessionvalue) + 4
		sizerelystring2 += len(Sessionkey) + 4
	}
	return 4 + 4 + len(obj.ClientConnID) + 4 + sizerelystring2
}
func ReadMsgSForwardToServerByBytes(indata []byte, obj *SForwardToServer) (int,*SForwardToServer ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SForwardToServer{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 + len(obj.FromServerID) > data__len{
		return endpos,obj
	}
	obj.FromServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromServerID)
	if offset + 4 + len(obj.ToServerID) > data__len{
		return endpos,obj
	}
	obj.ToServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToServerID)
	if offset + 2 > data__len{
		return endpos,obj
	}
	obj.MsgID = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 4 > data__len{
		return endpos,obj
	}
	Data_slen := int(binary.BigEndian.Uint32(data[offset:offset+4]))
	offset += 4
	if Data_slen != 0xffffffff {
		if offset + Data_slen > data__len {
			return endpos,obj
		}
		obj.Data = make([]byte,Data_slen)
		copy(obj.Data, data[offset:offset+Data_slen])
		offset += Data_slen
	}
	return endpos,obj
}
func WriteMsgSForwardToServerByObj(data []byte, obj *SForwardToServer) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:],obj.FromServerID)
	offset += 4 + len(obj.FromServerID)
	writeBinaryString(data[offset:],obj.ToServerID)
	offset += 4 + len(obj.ToServerID)
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.MsgID)
	offset+=2
	if obj.Data == nil {
		binary.BigEndian.PutUint32(data[offset:offset+4],0xffffffff)
	} else {
		binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Data)))
	}
	offset += 4
	Data_slen := len(obj.Data)
	copy(data[offset:offset+Data_slen], obj.Data)
	offset += Data_slen
	return offset
}
func GetSizeSForwardToServer(obj *SForwardToServer) int {
	if obj == nil {
		return 4
	}
	return 4 + 4 + len(obj.FromServerID) + 4 + len(obj.ToServerID) + 2 + 4 + len(obj.Data) * 1
}
func ReadMsgSForwardToClientByBytes(indata []byte, obj *SForwardToClient) (int,*SForwardToClient ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SForwardToClient{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 + len(obj.FromServerID) > data__len{
		return endpos,obj
	}
	obj.FromServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromServerID)
	if offset + 4 + len(obj.ToGateID) > data__len{
		return endpos,obj
	}
	obj.ToGateID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToGateID)
	if offset + 4 + len(obj.ToClientID) > data__len{
		return endpos,obj
	}
	obj.ToClientID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToClientID)
	if offset + 2 > data__len{
		return endpos,obj
	}
	obj.MsgID = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 4 > data__len{
		return endpos,obj
	}
	Data_slen := int(binary.BigEndian.Uint32(data[offset:offset+4]))
	offset += 4
	if Data_slen != 0xffffffff {
		if offset + Data_slen > data__len {
			return endpos,obj
		}
		obj.Data = make([]byte,Data_slen)
		copy(obj.Data, data[offset:offset+Data_slen])
		offset += Data_slen
	}
	return endpos,obj
}
func WriteMsgSForwardToClientByObj(data []byte, obj *SForwardToClient) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:],obj.FromServerID)
	offset += 4 + len(obj.FromServerID)
	writeBinaryString(data[offset:],obj.ToGateID)
	offset += 4 + len(obj.ToGateID)
	writeBinaryString(data[offset:],obj.ToClientID)
	offset += 4 + len(obj.ToClientID)
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.MsgID)
	offset+=2
	if obj.Data == nil {
		binary.BigEndian.PutUint32(data[offset:offset+4],0xffffffff)
	} else {
		binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Data)))
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
	return 4 + 4 + len(obj.FromServerID) + 4 + len(obj.ToGateID) + 4 + len(obj.ToClientID) + 2 + 
	4 + len(obj.Data) * 1
}
func ReadMsgSForwardFromGateByBytes(indata []byte, obj *SForwardFromGate) (int,*SForwardFromGate ) {
	offset := 0
	if len(indata) < 4 {
		return 0,nil
	}
	objsize := int(binary.BigEndian.Uint32(indata))
	offset += 4
	if objsize == 0 {
		return 4,nil
	}
	if obj == nil{
		obj=&SForwardFromGate{}
	}
	if offset + objsize > len(indata ) {
		return offset,obj
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 + len(obj.FromServerID) > data__len{
		return endpos,obj
	}
	obj.FromServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromServerID)
	if offset + 4 + len(obj.ToServerID) > data__len{
		return endpos,obj
	}
	obj.ToServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToServerID)
	if offset + 4 + len(obj.ClientConnID) > data__len{
		return endpos,obj
	}
	obj.ClientConnID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ClientConnID)
	if offset + 4 > data__len{
		return endpos,obj
	}
	Session_slen := binary.BigEndian.Uint32(data[offset:offset+4])
	offset += 4
	if Session_slen != 0xffffffff {
		obj.Session = make(map[string]string)
		for i4i := uint32(0); i4i < Session_slen; i4i++ {
			if offset + 0 > data__len{
				return endpos,obj
			}
			keySession := readBinaryString(data[offset:])
			Session_kcatlen := len(keySession)
			offset += Session_kcatlen + 4
			if offset + 2 > data__len{
				return endpos,obj
			}
			valueSession := readBinaryString(data[offset:])
			Session_vcatlen := len(valueSession)
			offset += Session_vcatlen + 4
			obj.Session[keySession] = valueSession
		}
	}
	if offset + 2 > data__len{
		return endpos,obj
	}
	obj.MsgID = binary.BigEndian.Uint16(data[offset:offset+2])
	offset+=2
	if offset + 4 > data__len{
		return endpos,obj
	}
	Data_slen := int(binary.BigEndian.Uint32(data[offset:offset+4]))
	offset += 4
	if Data_slen != 0xffffffff {
		if offset + Data_slen > data__len {
			return endpos,obj
		}
		obj.Data = make([]byte,Data_slen)
		copy(obj.Data, data[offset:offset+Data_slen])
		offset += Data_slen
	}
	return endpos,obj
}
func WriteMsgSForwardFromGateByObj(data []byte, obj *SForwardFromGate) int {
	if obj == nil {
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	writeBinaryString(data[offset:],obj.FromServerID)
	offset += 4 + len(obj.FromServerID)
	writeBinaryString(data[offset:],obj.ToServerID)
	offset += 4 + len(obj.ToServerID)
	writeBinaryString(data[offset:],obj.ClientConnID)
	offset += 4 + len(obj.ClientConnID)
	if obj.Session == nil {
		binary.BigEndian.PutUint32(data[offset:offset+4],0xffffffff)
	} else {
		binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Session)))
	}
	offset += 4
	for Sessionkey,Sessionvalue := range obj.Session {
		Session_kcatlen := writeBinaryString(data[offset:],Sessionkey)
		offset += Session_kcatlen
		Session_vcatlen := writeBinaryString(data[offset:],Sessionvalue)
		offset += Session_vcatlen
	}
	binary.BigEndian.PutUint16(data[offset:offset+2], obj.MsgID)
	offset+=2
	if obj.Data == nil {
		binary.BigEndian.PutUint32(data[offset:offset+4],0xffffffff)
	} else {
		binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Data)))
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
	for Sessionvalue,Sessionkey := range obj.Session {
		sizerelystring4 += len(Sessionvalue) + 4
		sizerelystring4 += len(Sessionkey) + 4
	}
	return 4 + 4 + len(obj.FromServerID) + 4 + len(obj.ToServerID) + 4 + len(obj.ClientConnID) + 4 + sizerelystring4 + 
	2 + 4 + len(obj.Data) * 1
}
