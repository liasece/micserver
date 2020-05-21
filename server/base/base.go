/*
Package base 服务的基本接口
*/
package base

import (
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/session"
)

// ServerHook 上层服务(模块)需要处理模块消息时需要实现的接口
type ServerHook interface {
	OnModuleMessage(msg *servercomm.ModuleMessage)
	OnClientMessage(se *session.Session, msg *servercomm.ClientMessage)
}
