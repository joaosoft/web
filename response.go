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
		Base:        request.Base,
		Attachments: make(map[string]Attachment),
		Boundary:    RandomBoundary(),
		Writer:      request.conn.(io.Writer),
	}
}

func (r *Response) write() error {
	hasAttachments := len(r.Attachments) > 0
	// header
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Protocol, r.Status, StatusText(r.Status)))

	// headers
	r.Headers[HeaderServer] = []string{"webserver"}
	r.Headers[HeaderDate] = []string{time.Now().Format(TimeFormat)}

	if hasAttachments {
		r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; boundary=%s", ContentMultipartFormData, r.Boundary)}
		r.Headers[HeaderContentLength] = []string{strconv.Itoa(len(r.Body))}
	} else {
		r.Headers[HeaderContentType] = []string{string(r.ContentType)}

		size := len(r.Body)
		for _, attachment := range r.Attachments {
			size += len(attachment.Body)
		}
		r.Headers[HeaderContentLength] = []string{strconv.Itoa(size)}
	}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	buf.WriteString("\r\n")

	if hasAttachments {
		for _, attachment := range r.Attachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
			buf.WriteString(fmt.Sprintf("%s: %s; name=%q; filename=%q\r\n", HeaderContentDisposition, attachment.ContentDisposition, attachment.Name, attachment.File))
			buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", HeaderContentType, attachment.ContentType))
			buf.Write(attachment.Body)
			buf.WriteString("\r\n")
		}
	}

	if methodHasBody[r.Method] {
		if hasAttachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", HeaderContentDisposition, ContentDispositionFormData))
			buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", HeaderContentType, r.ContentType))
			buf.Write(r.Body)
		} else {
			buf.Write(r.Body)
		}
	}

	if hasAttachments {
		buf.WriteString(fmt.Sprintf("\r\n--%s--", r.Boundary))
	}

	fmt.Println(buf.String())
	r.conn.Write(buf.Bytes())

	return nil
}
