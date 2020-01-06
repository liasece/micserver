/*
一个消息具有的基础数据
*/
package base

// 发送完成回调具备的信息
type SendFinishAgent struct {
	F    func(interface{}) // 用于优化的临时数据指针请注意使用！
	Argv interface{}       // 用于优化的临时数据指针请注意使用！
}

// 基础消息
type MessageBase struct {
	regSendFinish SendFinishAgent // 用于优化的临时数据指针请注意使用！
	msgObject     interface{}     // 用于优化的临时数据指针请注意使用！
}

// 注册当消息发送完成时的回调
func (this *MessageBase) RegSendFinish(cb func(interface{}), argv interface{}) {
	if cb == nil {
		return
	}
	this.regSendFinish.F = cb
	this.regSendFinish.Argv = argv
}

// 当消息发送完成时调用该方法，以执行发送完成回调
func (this *MessageBase) OnSendFinish() {
	if this.regSendFinish.F != nil {
		this.regSendFinish.F(this.regSendFinish.Argv)
	}
}

// 设置该消息的消息对象
func (this *MessageBase) SetObj(obj interface{}) {
	this.msgObject = obj
}

// 获取该消息的消息对象
func (this *MessageBase) GetObj() interface{} {
	return this.msgObject
}

// 重置该消息
func (this *MessageBase) Reset() {
	this.regSendFinish.F = nil
	this.msgObject = nil
}
