package webserver

import (
	"net"
)

func NewResponse(request *Request) *Response {
	return &Response{
		Base: Base{
			Method:   request.Method,
			Url:      request.Url,
			Protocol: request.Protocol,
			Headers:  request.Headers,
			Cookies:  request.Cookies,
		},
	}
}

func (r *Response) write(conn net.Conn) error {
	return nil
}
