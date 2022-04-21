package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

func (e *Error) Add(err *Error) *Error {
	prevError := &Error{
		Previous: e.Previous,
		Level:    e.Level,
		Code:     e.Code,
		Message:  e.Message,
		Stack:    e.Stack,
	}

	var stack string
	if err.Level <= LevelError {
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		function := runtime.FuncForPC(pc[0])
		stack = string(debug.Stack())
		index := strings.Index(stack, function.Name())
		if index < 0 {
			index = 0
		}
		stack = stack[index:]
	}

	return &Error{
		Previous: prevError,
		Level:    err.Level,
		Code:     err.Code,
		Message:  err.Message,
		Stack:    stack,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Cause() string {
	str := fmt.Sprintf("'%s'", e.Message)

	prevErr := e.Previous
	for prevErr != nil {
		str += fmt.Sprintf(", caused by '%s'", prevErr.Message)
		prevErr = prevErr.Previous
	}
	return str
}

func (e *Error) Errors() []*Error {
	errors := make([]*Error, 0)
	errors = append(errors, e)

	nextErr := e.Previous
	for nextErr != nil {
		errors = append(errors, e.Previous)
		nextErr = nextErr.Previous
	}

	return errors
}

func (e *Error) Format(values ...interface{}) *Error {
	e.Message = fmt.Sprintf(e.Error(), values...)
	return e
}

func (e *Error) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
