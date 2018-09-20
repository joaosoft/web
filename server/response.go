package server

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"strconv"
	"time"
	"webserver"
)

func (w *WebServer) NewResponse(request *Request) *Response {
	return &Response{
		Base:                request.Base,
		Attachments:         make(map[string]Attachment),
		MultiAttachmentMode: w.multiAttachmentMode,
		Boundary:            webserver.RandomBoundary(),
		Writer:              request.conn.(io.Writer),
	}
}

func (r *Response) write() error {
	var buf bytes.Buffer
	var lenAttachments = len(r.Attachments)

	if headers, err := r.handleHeaders(); err != nil {
		return err
	} else {
		buf.Write(headers)
	}

	if lenAttachments > 0 {
		switch r.MultiAttachmentMode {
		case webserver.MultiAttachmentModeBoundary:
			if body, err := r.handleBody(); err != nil {
				return err
			} else {
				buf.Write(body)
			}
			if body, err := r.handleBoundaryAttachments(); err != nil {
				return err
			} else {
				buf.Write(body)
			}
		case webserver.MultiAttachmentModeZip:
			if lenAttachments > 1 {
				if body, err := r.handleZippedAttachments(); err != nil {
					return err
				} else {
					buf.Write(body)
				}
			} else {
				if body, err := r.handleSingleAttachment(); err != nil {
					return err
				} else {
					buf.Write(body)
				}
			}
		}
	} else {
		if body, err := r.handleBody(); err != nil {
			return err
		} else {
			buf.Write(body)
		}
	}

	r.conn.Write(buf.Bytes())

	return nil
}

func (r *Response) handleHeaders() ([]byte, error) {
	var buf bytes.Buffer
	lenAttachments := len(r.Attachments)

	// header
	buf.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Protocol, r.Status, webserver.StatusText(r.Status)))

	// headers
	r.Headers[webserver.HeaderServer] = []string{"webserver"}
	r.Headers[webserver.HeaderDate] = []string{time.Now().Format(webserver.TimeFormat)}

	if lenAttachments > 0 {

		switch r.MultiAttachmentMode {
		case webserver.MultiAttachmentModeBoundary:
			r.Headers[webserver.HeaderContentType] = []string{fmt.Sprintf("%s; boundary=%s; charset=%s", webserver.ContentTypeMultipartFormData, r.Boundary, r.Charset)}
		case webserver.MultiAttachmentModeZip:
			var name = "attachments"
			var fileName = "attachments.zip"
			var contentType = webserver.ContentTypeApplicationZip
			var charset = r.Charset

			if lenAttachments == 1 {
				for _, attachment := range r.Attachments {
					name = attachment.Name
					fileName = attachment.File
					contentType = attachment.ContentType
					if attachment.Charset != "" {
						charset = attachment.Charset
					}
					break
				}
			}
			r.Headers[webserver.HeaderContentType] = []string{fmt.Sprintf("%s; attachment; name=%q; filename=%q; charset=%s", contentType, name, fileName, charset)}
		}
	} else {
		r.Headers[webserver.HeaderContentType] = []string{string(r.ContentType)}
		r.Headers[webserver.HeaderContentLength] = []string{strconv.Itoa(len(r.Body))}
	}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

func (r *Response) handleBody() ([]byte, error) {
	var buf bytes.Buffer

	if webserver.MethodHasBody[r.Method] {
		buf.Write(r.Body)
		if r.MultiAttachmentMode == webserver.MultiAttachmentModeBoundary && len(r.Attachments) > 0 {
			buf.WriteString("\r\n\r\n")
		}
	}

	return buf.Bytes(), nil
}

func (r *Response) handleSingleAttachment() ([]byte, error) {
	for _, attachment := range r.Attachments {
		return attachment.Body, nil
	}
	return []byte{}, nil
}

func (r *Response) handleBoundaryAttachments() ([]byte, error) {
	var buf bytes.Buffer

	if len(r.Attachments) == 0 {
		return buf.Bytes(), nil
	}

	for _, attachment := range r.Attachments {
		buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
		buf.WriteString(fmt.Sprintf("%s: %s; name=%q; filename=%q\r\n", webserver.HeaderContentDisposition, attachment.ContentDisposition, attachment.Name, attachment.File))
		buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", webserver.HeaderContentType, attachment.ContentType))
		buf.Write(attachment.Body)
		buf.WriteString("\r\n")
	}

	buf.WriteString(fmt.Sprintf("\r\n--%s--", r.Boundary))

	return buf.Bytes(), nil
}

func (r *Response) handleZippedAttachments() ([]byte, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	if len(r.Attachments) == 0 {
		return buf.Bytes(), nil
	}

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Register a custom Deflate compressor to override the default Deflate compressor with a higher compression level.
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	for _, attachment := range r.Attachments {
		f, err := w.Create(attachment.File)
		if err != nil {
			return buf.Bytes(), err
		}
		_, err = f.Write([]byte(attachment.Body))
		if err != nil {
			return buf.Bytes(), err
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}
