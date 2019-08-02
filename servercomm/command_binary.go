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
	SForwardToServerID = 47
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
	SForwardToServerName = "servercomm.SForwardToServer"
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
func (this *SForwardToServer) WriteBinary(data []byte) int {
	return WriteMsgSForwardToServerByObj(data,this)
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
func (this *SNotifySafelyQuit) ReadBinary(data []byte) int {
	return ReadMsgSNotifySafelyQuitByBytes(data, this)
}
func (this *SForwardToServer) ReadBinary(data []byte) int {
	return ReadMsgSForwardToServerByBytes(data, this)
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
		case SForwardToServerID: 
		return SForwardToServerName
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
		case SForwardToServerName: 
		return SForwardToServerID
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
		case SNotifySafelyQuitID: 
		return rune('S')
		case SForwardToServerID: 
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
func (this *SNotifySafelyQuit) GetMsgId() uint16 {
	return SNotifySafelyQuitID
}
func (this *SForwardToServer) GetMsgId() uint16 {
	return SForwardToServerID
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
func (this *SForwardToServer) GetMsgName() string {
	return SForwardToServerName
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
func (this *SForwardToServer) GetSize() int {
	return GetSizeSForwardToServer(this)
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
func (this *SForwardToServer) GetJson() string {
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
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 + len(obj.ServerID) > data__len{
		return endpos
	}
	obj.ServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ServerID)
	if offset + 4 + len(obj.ServerAddr) > data__len{
		return endpos
	}
	obj.ServerAddr = readBinaryString(data[offset:])
	offset += 4 + len(obj.ServerAddr)
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
func ReadMsgSTimeTickCommandByBytes(indata []byte, obj *STimeTickCommand) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
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
func ReadMsgSTestCommandByBytes(indata []byte, obj *STestCommand) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
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
	if offset + 4 + len(obj.Testttring) > data__len{
		return endpos
	}
	obj.Testttring = readBinaryString(data[offset:])
	offset += 4 + len(obj.Testttring)
	return endpos
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
func ReadMsgSLoginCommandByBytes(indata []byte, obj *SLoginCommand) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 + len(obj.ServerID) > data__len{
		return endpos
	}
	obj.ServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ServerID)
	if offset + 4 + len(obj.ServerAddr) > data__len{
		return endpos
	}
	obj.ServerAddr = readBinaryString(data[offset:])
	offset += 4 + len(obj.ServerAddr)
	if offset + 8 > data__len{
		return endpos
	}
	obj.ConnectPriority = readBinaryInt64(data[offset:offset+8])
	offset+=8
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
func ReadMsgSLogoutCommandByBytes(indata []byte, obj *SLogoutCommand) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
	}
	endpos := offset+objsize
	return endpos
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
func ReadMsgSSeverStartOKCommandByBytes(indata []byte, obj *SSeverStartOKCommand) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
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
func ReadMsgSLoginRetCommandByBytes(indata []byte, obj *SLoginRetCommand) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
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
	if offset + obj.Destination.GetSize() > data__len{
		return endpos
	}
	offset += ReadMsgSServerInfoByBytes(data[offset:], &obj.Destination)
	return endpos
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
	offset += WriteMsgSServerInfoByObj(data[offset:], &obj.Destination)
	return offset
}
func GetSizeSLoginRetCommand(obj *SLoginRetCommand) int {
	if obj == nil {
		return 4
	}
	return 4 + 4 + obj.Destination.GetSize()
}
func ReadMsgSStartRelyNotifyCommandByBytes(indata []byte, obj *SStartRelyNotifyCommand) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	Serverinfos_slent := uint32(0)
	if offset + 4 > data__len{
		return endpos
	}
	Serverinfos_slent = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	Serverinfos_slen := int(Serverinfos_slent)
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
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Serverinfos)))
	offset += 4
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
		return 4
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
	return 4 + 4 + sizerelySServerInfo1()
}
func ReadMsgSStartMyNotifyCommandByBytes(indata []byte, obj *SStartMyNotifyCommand) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
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
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	offset += WriteMsgSServerInfoByObj(data[offset:], &obj.Serverinfo)
	return offset
}
func GetSizeSStartMyNotifyCommand(obj *SStartMyNotifyCommand) int {
	if obj == nil {
		return 4
	}
	return 4 + obj.Serverinfo.GetSize()
}
func ReadMsgSNotifyAllInfoByBytes(indata []byte, obj *SNotifyAllInfo) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	Serverinfos_slent := uint32(0)
	if offset + 4 > data__len{
		return endpos
	}
	Serverinfos_slent = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	Serverinfos_slen := int(Serverinfos_slent)
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
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Serverinfos)))
	offset += 4
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
		return 4
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
	return 4 + 4 + sizerelySServerInfo1()
}
func ReadMsgSNotifySafelyQuitByBytes(indata []byte, obj *SNotifySafelyQuit) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
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
		binary.BigEndian.PutUint32(data[0:4],0)
		return 4
	}
	objsize := obj.GetSize() - 4
	offset := 0
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))
	offset += 4
	offset += WriteMsgSServerInfoByObj(data[offset:], &obj.TargetServerInfo)
	return offset
}
func GetSizeSNotifySafelyQuit(obj *SNotifySafelyQuit) int {
	if obj == nil {
		return 4
	}
	return 4 + obj.TargetServerInfo.GetSize()
}
func ReadMsgSForwardToServerByBytes(indata []byte, obj *SForwardToServer) int {
	offset := 0
	if len(indata) < 4 {
		return 0
	}
	objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))
	offset += 4
	if objsize == 0 {
		return 4
	}
	if offset + objsize > len(indata ) {
		return offset
	}
	endpos := offset+objsize
	data := indata[offset:offset+objsize]
	offset = 0
	data__len := len(data)
	if offset + 4 + len(obj.FromServerID) > data__len{
		return endpos
	}
	obj.FromServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.FromServerID)
	if offset + 4 + len(obj.ToServerID) > data__len{
		return endpos
	}
	obj.ToServerID = readBinaryString(data[offset:])
	offset += 4 + len(obj.ToServerID)
	if offset + 4 + len(obj.MsgName) > data__len{
		return endpos
	}
	obj.MsgName = readBinaryString(data[offset:])
	offset += 4 + len(obj.MsgName)
	Data_slent := uint32(0)
	if offset + 4 > data__len{
		return endpos
	}
	Data_slent = binary.BigEndian.Uint32(data[offset:offset+4])
	offset+=4
	Data_slen := int(Data_slent)
	obj.Data = make([]byte,Data_slen)
	if offset+(Data_slen*1) > data__len {
		return endpos
	}
	copy(obj.Data, data[offset:offset+Data_slen])
	return endpos
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
	writeBinaryString(data[offset:],obj.MsgName)
	offset += 4 + len(obj.MsgName)
	binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.Data)))
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
	return 4 + 4 + len(obj.FromServerID) + 4 + len(obj.ToServerID) + 4 + len(obj.MsgName) + 4 + len(obj.Data) * 1
}
