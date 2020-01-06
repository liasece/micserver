package base

// module提供给外部的接口
type IModule interface {
	GetModuleID() string
	GetModuleType() string
	GetModuleNum() int
	GetModuleIDHash() uint32
}
