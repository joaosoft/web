package errors

import "encoding/json"

func (el *ErrorList) Len() int {
	return len(*el)
}

func (el *ErrorList) IsEmpty() bool {
	return len(*el) == 0
}

func (el *ErrorList) Add(err *Error) *ErrorList {
	*el = append(*el, err)
	return el
}

func (el *ErrorList) String() string {
	b, _ := json.Marshal(el)
	return string(b)
}
