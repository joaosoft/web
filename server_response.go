package web

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"strconv"
	"time"
)

func (w *Server) NewResponse(request *Request) *Response {
	return &Response{
		Base:                request.Base,
		FormData:            make(map[string]*FormData),
		MultiAttachmentMode: w.multiAttachmentMode,
		Boundary:            RandomBoundary(),
		Writer:              request.conn.(io.Writer),
		Status:              StatusNoContent,
		StatusText:          StatusText(StatusNoContent),
	}
}

func (r *Response) write() error {
	var buf bytes.Buffer
	var lenFormData = len(r.FormData)

	if headers, err := r.handleHeaders(); err != nil {
		return err
	} else {
		buf.Write(headers)
	}

	if lenFormData > 0 {
		switch r.MultiAttachmentMode {
		case MultiAttachmentModeBoundary:
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
		case MultiAttachmentModeZip:
			if lenFormData > 1 {
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

			if r.Server.logger.IsDebugEnabled() {
				r.Server.logger.Infof("[RESPONSE BODY] [%s]", string(body))
			}
		}
	}

	r.conn.Write(buf.Bytes())

	return nil
}

func (r *Response) handleHeaders() ([]byte, error) {
	var buf bytes.Buffer
	lenFormData := len(r.FormData)

	r.Headers[HeaderServer] = []string{"Server"}
	r.Headers[HeaderDate] = []string{time.Now().Format(TimeFormat)}
	r.Headers[HeaderAccessControlAllowCredentials] = []string{"true"}
	if val, ok := r.Headers[HeaderOrigin]; ok {
		r.Headers[HeaderAccessControlAllowOrigin] = val
	}

	// header
	buf.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Protocol, r.Status, StatusText(r.Status)))

	if lenFormData > 0 {

		switch r.MultiAttachmentMode {
		case MultiAttachmentModeBoundary:
			r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; boundary=%s; charset=%s", ContentTypeMultipartFormData, r.Boundary, r.Charset)}
		case MultiAttachmentModeZip:
			var name = "attachments"
			var fileName = "attachments.zip"
			var contentType = ContentTypeApplicationZip
			var charset = r.Charset

			if lenFormData == 1 {
				for _, attachment := range r.FormData {
					name = attachment.Name
					fileName = attachment.FileName
					contentType = attachment.ContentType
					if attachment.Charset != "" {
						charset = attachment.Charset
					}
					break
				}
			}
			r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; attachment; name=%q; filename=%q; charset=%s", contentType, name, fileName, charset)}
		}
	} else {
		r.Headers[HeaderContentType] = []string{string(r.ContentType)}
		r.Headers[HeaderContentLength] = []string{strconv.Itoa(len(r.Body))}
	}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

func (r *Response) handleBody() ([]byte, error) {
	var buf bytes.Buffer

	if MethodHasBody[r.Method] {
		buf.Write(r.Body)
		if r.MultiAttachmentMode == MultiAttachmentModeBoundary && len(r.FormData) > 0 {
			buf.WriteString("\r\n\r\n")
		}
	}

	return buf.Bytes(), nil
}

func (r *Response) handleSingleAttachment() ([]byte, error) {
	for _, formData := range r.FormData {
		return formData.Body, nil
	}
	return []byte{}, nil
}

func (r *Response) handleBoundaryAttachments() ([]byte, error) {
	var buf bytes.Buffer

	if len(r.FormData) == 0 {
		return buf.Bytes(), nil
	}

	lenF := len(r.FormData)
	i := 0

	for _, formData := range r.FormData {
		buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
		if formData.IsAttachment {
			buf.WriteString(fmt.Sprintf("%s: %s; name=%q; filename=%q\r\n", HeaderContentDisposition, formData.ContentDisposition, formData.Name, formData.FileName))
		} else {
			buf.WriteString(fmt.Sprintf("%s: %s; name=%q\r\n", HeaderContentDisposition, formData.ContentDisposition, formData.Name))

		}
		buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", HeaderContentType, formData.ContentType))
		buf.Write(formData.Body)

		if i < lenF-1 {
			buf.WriteString("\r\n")
		}
		i++
	}

	buf.WriteString(fmt.Sprintf("\r\n--%s--", r.Boundary))

	return buf.Bytes(), nil
}

func (r *Response) handleZippedAttachments() ([]byte, error) {
	// create a buffer to write our archive
	buf := new(bytes.Buffer)

	if len(r.FormData) == 0 {
		return buf.Bytes(), nil
	}

	// create a new zip archive
	w := zip.NewWriter(buf)

	// register a custom deflate compressor to override the default Deflate compressor with a higher compression level
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	for _, attachment := range r.FormData {
		f, err := w.Create(attachment.FileName)
		if err != nil {
			return buf.Bytes(), err
		}
		_, err = f.Write([]byte(attachment.Body))
		if err != nil {
			return buf.Bytes(), err
		}
	}

	err := w.Close()
	if err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}
