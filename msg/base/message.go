// Package base 一个消息具有的基础数据
package base

// SendFinishAgent 发送完成回调具备的信息
type SendFinishAgent struct {
	F    func(interface{}) // 用于优化的临时数据指针请注意使用！
	Argv interface{}       // 用于优化的临时数据指针请注意使用！
}

// MessageBase 基础消息
type MessageBase struct {
	regSendFinish SendFinishAgent // 用于优化的临时数据指针请注意使用！
	msgObject     interface{}     // 用于优化的临时数据指针请注意使用！
}

// RegSendFinish 注册当消息发送完成时的回调
func (mb *MessageBase) RegSendFinish(cb func(interface{}), argv interface{}) {
	if cb == nil {
		return
	}
	mb.regSendFinish.F = cb
	mb.regSendFinish.Argv = argv
}

// OnSendFinish 当消息发送完成时调用该方法，以执行发送完成回调
func (mb *MessageBase) OnSendFinish() {
	if mb.regSendFinish.F != nil {
		mb.regSendFinish.F(mb.regSendFinish.Argv)
	}
}

// SetObj 设置该消息的消息对象
func (mb *MessageBase) SetObj(obj interface{}) {
	mb.msgObject = obj
}

// GetObj 获取该消息的消息对象
func (mb *MessageBase) GetObj() interface{} {
	return mb.msgObject
}

// Reset 重置该消息
func (mb *MessageBase) Reset() {
	mb.regSendFinish.F = nil
	mb.msgObject = nil
}
