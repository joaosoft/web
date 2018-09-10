package webserver

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"
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
	buf.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Protocol, r.Status, StatusText(r.Status)))

	// headers
	r.Headers[HeaderContentLength] = []string{strconv.Itoa(len(r.Body))}
	r.Headers[HeaderServer] = []string{"webserver"}
	r.Headers[HeaderDate] = []string{time.Now().Format(TimeFormat)}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	if methodHasBody[r.Method] {
		buf.WriteString("\r\n")
		buf.Write(r.Body)

		r.conn.Write(buf.Bytes())
	}

	return nil
}
