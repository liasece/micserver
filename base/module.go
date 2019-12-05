package base

type IModule interface {
	GetModuleID() string
	GetModuleType() string
	GetModuleNum() int
	GetModuleIDHash() uint32
}
