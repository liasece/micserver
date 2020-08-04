package log

import (
	"errors"
	"sync"

	"github.com/liasece/micserver/log/core"
)

// error
var (
	ErrNilLogger       = errors.New("logger is nil")
	ErrUnknownLogLevel = errors.New("unknown log level")
)

var _errArrayElemPool = sync.Pool{New: func() interface{} {
	return &errArrayElem{}
}}

// ErrorField is shorthand for the common idiom NamedError("error", err).
func ErrorField(err error) Field {
	return NamedError("error", err)
}

// NamedError constructs a field that lazily stores err.Error() under the
// provided key. Errors which also implement fmt.Formatter (like those produced
// by github.com/liasece/micserver/util/errors) will also have their verbose representation stored
// under key+"Verbose". If passed a nil error, the field is a no-op.
//
// For the common case in which the key is simply "error", the Error function
// is shorter and less repetitive.
func NamedError(key string, err error) Field {
	if err == nil {
		return Skip()
	}
	return Field{Key: key, Type: core.ErrorType, Interface: err}
}

type errArray []error

func (errs errArray) MarshalLogArray(arr core.ArrayEncoder) error {
	for i := range errs {
		if errs[i] == nil {
			continue
		}
		// To represent each error as an object with an "error" attribute and
		// potentially an "errorVerbose" attribute, we need to wrap it in a
		// type that implements LogObjectMarshaler. To prevent this from
		// allocating, pool the wrapper type.
		elem := _errArrayElemPool.Get().(*errArrayElem)
		elem.error = errs[i]
		arr.AppendObject(elem)
		elem.error = nil
		_errArrayElemPool.Put(elem)
	}
	return nil
}

type errArrayElem struct {
	error
}

func (e *errArrayElem) MarshalLogObject(enc core.ObjectEncoder) error {
	// Re-use the error field's logic, which supports non-standard error types.
	ErrorField(e.error).AddTo(enc)
	return nil
}
