package errors

import (
	"encoding/json"
	"fmt"
)

func (e *Err) Add(newErr *Err) {
	prevErr := &Err{
		Previous: e.Previous,
		Code:     e.Code,
		Err:      e.Err,
	}

	e.Previous = prevErr
	e.Code = newErr.Code
	e.Err = newErr.Err
}

func (e *Err) Error() string {
	return e.Err
}

func (e *Err) Cause() string {
	str := fmt.Sprintf("'%s'", e.Err)

	prevErr := e.Previous
	for prevErr != nil {
		str += fmt.Sprintf(", caused by '%s'", prevErr.Err)
		prevErr = prevErr.Previous
	}
	return str
}

func (e *Err) SetErr(newErr *Err) {
	*e = *newErr
}

func (e *Err) GetErr() *Err {
	return e
}

func (e *Err) GetPrevious() *Err {
	return e.Previous
}

func (e *Err) GetErrors() []*Err {
	errors := make([]*Err, 0)
	errors = append(errors, e)

	nextErr := e.Previous
	for nextErr != nil {
		errors = append(errors, e.Previous)
		nextErr = nextErr.Previous
	}

	return errors
}

func (e *Err) SetCode(code string) {
	e.Code = code

}
func (e *Err) GetCode() string {
	return e.Code
}

func (e *Err) Format(values ...interface{}) *Err {
	e.SetErr(New(e.Code, fmt.Sprintf(e.Error(), values)))
	return e
}

func (e *Err) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
