package server

import (
	"fmt"
	"webserver"
)

var (
	ErrorNotFound = NewError(webserver.StatusNotFound, "route not found")
)

type Error struct {
	Status   webserver.Status
	Messages interface{}
}

func NewError(status webserver.Status, errors ...interface{}) *Error {
	err := &Error{
		Status: status,
	}

	if len(errors) > 0 {
		err.Messages = errors
	} else {
		err.Messages = []string{webserver.StatusText(status)}
	}

	return err
}

func (e *Error) Error() string {
	return fmt.Sprintf("status=%d, messages=%v", e.Status, e.Messages)
}

func (w *WebServer) DefaultErrorHandler(ctx *Context, err error) error {
	w.logger.Infof("handling error: %s", err)

	if e, ok := err.(*Error); ok {
		return ctx.Response.JSON(e.Status, e)
	}

	return ctx.Response.JSON(webserver.StatusInternalServerError, NewError(webserver.StatusInternalServerError, err))
}
