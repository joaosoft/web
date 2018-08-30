package errors

import (
	"fmt"
)

func New(code string, err interface{}, params ...interface{}) *Err {

	switch v := err.(type) {
	case error:
		return &Err{Code: code, Err: v.Error()}

	case string:
		return &Err{Code: code, Err: fmt.Sprintf(v, params...)}

	default:
		return &Err{Code: code, Err: fmt.Sprint(v)}

	}
}
