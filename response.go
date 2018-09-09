package web

import (
	"bytes"
	"fmt"
	"io"
)

func NewResponse(request *Request) *Response {
	return &Response{
		Base:   request.Base,
		Writer: request.conn.(io.Writer),
	}
}

func (r *Response) write() error {
	// header
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s %d %s\n", r.Protocol, r.Status, StatusText(r.Status)))

	// headers
	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\n", key, value[0]))
	}

	if methodHasBody[r.Method] {
		buf.WriteString("\n")
		buf.Write(r.Body)

		r.conn.Write(buf.Bytes())
	}

	return nil
}
