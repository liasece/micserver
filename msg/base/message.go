package base

type SendFinishAgent struct {
	F    func(interface{}) // 用于优化的临时数据指针请注意使用！
	Argv interface{}       // 用于优化的临时数据指针请注意使用！
}

type MessageBase struct {
	regSendFinish SendFinishAgent // 用于优化的临时数据指针请注意使用！
	msgObject     interface{}     // 用于优化的临时数据指针请注意使用！
}

func (this *MessageBase) RegSendFinish(cb func(interface{}), argv interface{}) {
	if cb == nil {
		return
	}
	this.regSendFinish.F = cb
	this.regSendFinish.Argv = argv
}

func (this *MessageBase) OnSendFinish() {
	if this.regSendFinish.F != nil {
		this.regSendFinish.F(this.regSendFinish.Argv)
	}
}

func (this *MessageBase) SetObj(obj interface{}) {
	this.msgObject = obj
}

func (this *MessageBase) GetObj() interface{} {
	return this.msgObject
}

func (this *MessageBase) Reset() {
	this.regSendFinish.F = nil
	this.msgObject = nil
}
