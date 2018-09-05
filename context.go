package webserver

import (
	"time"
)

func NewContext(request *Request, response *Response) *Context {
	return &Context{
		StartTime: time.Now(),
		Request:   request,
		Response:  response,
	}
}
