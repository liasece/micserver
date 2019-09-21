package roc

type IObj interface {
	GetObjType() string
	GetObjID() string
	ROCCall([]string, []byte) ([]byte, error)
}
