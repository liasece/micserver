package roc

type IObj interface {
	GetObjType() string
	GetObjID() string
	ROCCall(*ROCPath, []byte) ([]byte, error)
}

type IROCObjEventHook interface {
	OnROCObjAdd(IObj)
	OnROCObjDel(IObj)
}
