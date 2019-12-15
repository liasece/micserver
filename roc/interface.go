package roc

type IObj interface {
	GetROCObjType() ROCObjType
	GetROCObjID() string
	OnROCCall(*ROCPath, []byte) ([]byte, error)
}

type IROCObjEventHook interface {
	OnROCObjAdd(IObj)
	OnROCObjDel(IObj)
}

type IROCServer interface {
	GetROC(ROCObjType) *ROC
	NewROC(ROCObjType) *ROC
	ROCCallNR(*ROCPath, []byte) error
	ROCCallBlock(*ROCPath, []byte) ([]byte, error)
	GetROCObjCacheLocation(*ROCPath) string
	RangeROCObjIDByType(ROCObjType, func(id string, location string) bool)
	RangeMyROCObjIDByType(ROCObjType, func(id string, location string) bool)
	RandomROCObjIDByType(ROCObjType) string
	RandomMyROCObjIDByType(ROCObjType) string
}
