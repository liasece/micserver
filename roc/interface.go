package roc

// IObj ROC对象需要实现的接口
type IObj interface {
	GetROCObjType() ObjType
	GetROCObjID() string
	OnROCCall(*Path, []byte) ([]byte, error)
}

// IROCObjEventHook ROC事件处理需要实现的接口
type IROCObjEventHook interface {
	OnROCObjAdd(IObj)
	OnROCObjDel(IObj)
}

// IROCServer ROC服务器实现的接口
type IROCServer interface {
	GetROC(ObjType) *ROC
	NewROC(ObjType) *ROC
	ROCCallNR(*Path, []byte) error
	ROCCallBlock(*Path, []byte) ([]byte, error)
	GetROCCachedLocation(ObjType, string) string
	RangeROCCachedByType(ObjType, func(id string, location string) bool)
	RandomROCCachedByType(ObjType) string
}
