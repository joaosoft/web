package client

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"web"
)

func (c *Client) NewRequest(method web.Method, url string) (*Request, error) {

	// validate url
	regx := regexp.MustCompile(web.RegexForURL)
	if !regx.MatchString(url) {
		return nil, fmt.Errorf("invalid url [%s]", url)
	}

	return &Request{
		Base: Base{
			Protocol:  web.ProtocolHttp1p1,
			Method:    method,
			Url:       url,
			Headers:   make(web.Headers),
			Cookies:   make(web.Cookies),
			Params:    make(web.Params),
			UrlParams: make(web.UrlParams),
			Charset:   web.CharsetUTF8,
		},
		Attachments:         make(map[string]Attachment),
		MultiAttachmentMode: c.multiAttachmentMode,
		Boundary:            web.RandomBoundary(),
	}, nil
}

func (r *Request) build() ([]byte, error) {
	var buf bytes.Buffer
	var lenAttachments = len(r.Attachments)

	if headers, err := r.handleHeaders(); err != nil {
		return nil, err
	} else {
		buf.Write(headers)
	}

	if lenAttachments > 0 {
		switch r.MultiAttachmentMode {
		case web.MultiAttachmentModeBoundary:
			if body, err := r.handleBody(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
			if body, err := r.handleBoundaryAttachments(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
		case web.MultiAttachmentModeZip:
			if lenAttachments > 1 {
				if body, err := r.handleZippedAttachments(); err != nil {
					return nil, err
				} else {
					buf.Write(body)
				}
			} else {
				if body, err := r.handleSingleAttachment(); err != nil {
					return nil, err
				} else {
					buf.Write(body)
				}
			}
		}
	} else {
		if body, err := r.handleBody(); err != nil {
			return nil, err
		} else {
			buf.Write(body)
		}
	}

	return buf.Bytes(), nil
}

func (r *Request) handleHeaders() ([]byte, error) {
	var buf bytes.Buffer
	lenAttachments := len(r.Attachments)

	// header
	split := strings.SplitN(r.Url, "/", 2)
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid url [%s]", r.Url)
	}

	buf.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Method, fmt.Sprintf("/%s", split[1]), r.Protocol))

	// headers
	r.Headers[web.HeaderServer] = []string{"client"}
	r.Headers[web.HeaderDate] = []string{time.Now().Format(web.TimeFormat)}

	if lenAttachments > 0 {

		switch r.MultiAttachmentMode {
		case web.MultiAttachmentModeBoundary:
			r.Headers[web.HeaderContentType] = []string{fmt.Sprintf("%s; boundary=%s; charset=%s", web.ContentTypeMultipartFormData, r.Boundary, r.Charset)}
		case web.MultiAttachmentModeZip:
			var name = "attachments"
			var fileName = "attachments.zip"
			var contentType = web.ContentTypeApplicationZip
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
			r.Headers[web.HeaderContentType] = []string{fmt.Sprintf("%s; attachment; name=%q; filename=%q; charset=%s", contentType, name, fileName, charset)}
		}
	} else {
		r.Headers[web.HeaderContentType] = []string{string(r.ContentType)}
		r.Headers[web.HeaderContentLength] = []string{strconv.Itoa(len(r.Body))}
	}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

func (r *Request) handleBody() ([]byte, error) {
	var buf bytes.Buffer

	if web.MethodHasBody[r.Method] {
		buf.Write(r.Body)
		if r.MultiAttachmentMode == web.MultiAttachmentModeBoundary && len(r.Attachments) > 0 {
			buf.WriteString("\r\n\r\n")
		}
	}

	return buf.Bytes(), nil
}

func (r *Request) handleSingleAttachment() ([]byte, error) {
	for _, attachment := range r.Attachments {
		return attachment.Body, nil
	}
	return []byte{}, nil
}

func (r *Request) handleBoundaryAttachments() ([]byte, error) {
	var buf bytes.Buffer

	if len(r.Attachments) == 0 {
		return buf.Bytes(), nil
	}

	for _, attachment := range r.Attachments {
		buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
		buf.WriteString(fmt.Sprintf("%s: %s; name=%q; filename=%q\r\n", web.HeaderContentDisposition, attachment.ContentDisposition, attachment.Name, attachment.File))
		buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", web.HeaderContentType, attachment.ContentType))
		buf.Write(attachment.Body)
		buf.WriteString("\r\n")
	}

	buf.WriteString(fmt.Sprintf("\r\n--%s--", r.Boundary))

	return buf.Bytes(), nil
}

func (r *Request) handleZippedAttachments() ([]byte, error) {
	// create a buffer to write our archive
	buf := new(bytes.Buffer)

	if len(r.Attachments) == 0 {
		return buf.Bytes(), nil
	}

	// create a new zip archive
	w := zip.NewWriter(buf)

	// register a custom deflate compressor to override the default Deflate compressor with a higher compression level
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

	err := w.Close()
	if err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}
