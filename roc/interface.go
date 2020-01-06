package roc

// ROC对象需要实现的接口
type IObj interface {
	GetROCObjType() ROCObjType
	GetROCObjID() string
	OnROCCall(*ROCPath, []byte) ([]byte, error)
}

// ROC事件处理需要实现的接口
type IROCObjEventHook interface {
	OnROCObjAdd(IObj)
	OnROCObjDel(IObj)
}

// ROC服务器实现的接口
type IROCServer interface {
	GetROC(ROCObjType) *ROC
	NewROC(ROCObjType) *ROC
	ROCCallNR(*ROCPath, []byte) error
	ROCCallBlock(*ROCPath, []byte) ([]byte, error)
	GetROCCachedLocation(ROCObjType, string) string
	RangeROCCachedByType(ROCObjType, func(id string, location string) bool)
	RandomROCCachedByType(ROCObjType) string
}
