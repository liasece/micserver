package roc

type IObj interface {
	GetObjType() string
	GetObjID() string
	ROCCall(*ROCPath, []byte) ([]byte, error)
}
