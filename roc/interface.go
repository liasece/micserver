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
	GetROCCachedLocation(ROCObjType, string) string
	RangeROCCachedByType(ROCObjType, func(id string, location string) bool)
	RandomROCCachedByType(ROCObjType) string
}
