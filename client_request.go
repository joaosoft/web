package web

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/joaosoft/auth-types/basic"
	"github.com/joaosoft/auth-types/jwt"
)

func (c *Client) NewRequest(method Method, url string) (*Request, error) {

	// validate url
	regx := regexp.MustCompile(RegexForURL)
	if !regx.MatchString(url) {
		return nil, fmt.Errorf("invalid url [%s]", url)
	}

	address := NewAddress(url)
	params := address.Params

	return &Request{
		Base: Base{
			Client:      c,
			Protocol:    ProtocolHttp1p1,
			Method:      method,
			Address:     address,
			Headers:     make(Headers),
			Cookies:     make(Cookies),
			Params:      params,
			UrlParams:   make(UrlParams),
			Charset:     CharsetUTF8,
			ContentType: ContentTypeApplicationJSON,
		},
		FormData:            make(map[string]*FormData),
		Attachments:         make(map[string]*Attachment),
		MultiAttachmentMode: c.multiAttachmentMode,
		Boundary:            RandomBoundary(),
	}, nil
}

func (r *Request) WithBody(body []byte, contentType ContentType) *Request {
	r.Body = body
	r.ContentType = contentType

	return r
}

func (r *Request) WithAuthBasic(username, password string) (*Request, error) {
	r.SetHeader(HeaderAuthorization, []string{basic.Generate(username, password)})

	return r, nil
}

func (r *Request) WithAuthJwt(claims jwt.Claims, key interface{}) (*Request, error) {
	token, err := jwt.New(jwt.SignatureHS384).Generate(claims, key)
	if err != nil {
		return r, err
	}

	r.SetHeader(HeaderAuthorization, []string{token})

	return r, nil
}

func (r *Request) WithContentType(contentType ContentType) *Request {
	r.ContentType = contentType

	return r
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
		case MultiAttachmentModeBoundary:
			if body, err := r.handleBody(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
			if body, err := r.handleBoundaries(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
		case MultiAttachmentModeZip:
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

			if body, err := r.handleBoundaries(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
		}
	} else {
		switch r.ContentType {
		case ContentTypeMultipartFormData, ContentTypeMultipartMixed:
			if body, err := r.handleBoundaries(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
		case ContentTypeApplicationForm:
			if urlForm, err := r.handleUrlForm(); err != nil {
				return nil, err
			} else {
				buf.Write(urlForm)
			}
		default:
			if body, err := r.handleBody(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
		}
	}

	return buf.Bytes(), nil
}

func (r *Request) handleHeaders() ([]byte, error) {
	var buf bytes.Buffer
	lenFormData := len(r.FormData)

	// header
	buf.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Method, r.Address.ParamsUrl, r.Protocol))

	// headers
	r.Headers[HeaderHost] = []string{r.Address.Host}
	if _, ok := r.Headers[HeaderUserAgent]; !ok {
		r.Headers[HeaderUserAgent] = []string{"client"}
	}
	r.Headers[HeaderDate] = []string{time.Now().Format(TimeFormat)}

	if lenFormData > 0 {

		switch r.MultiAttachmentMode {
		case MultiAttachmentModeBoundary:
			r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; boundary=%s; charset=%s", ContentTypeMultipartFormData, r.Boundary, r.Charset)}
		case MultiAttachmentModeZip:
			if len(r.FormData) == 0 {
				var name = "attachments"
				var fileName = "attachments.zip"
				var contentType = ContentTypeApplicationZip
				var charset = r.Charset

				if lenFormData == 1 {
					for _, formData := range r.Attachments {
						name = formData.Name
						fileName = formData.FileName
						contentType = formData.ContentType
						if formData.Charset != "" {
							charset = formData.Charset
						}
						break
					}
				}
				r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; attachment; name=%q; filename=%q; charset=%s", contentType, name, fileName, charset)}
			} else {
				r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; boundary=%s; charset=%s", ContentTypeMultipartFormData, r.Boundary, r.Charset)}
			}
		}
	} else {
		r.Headers[HeaderContentType] = []string{string(r.ContentType)}
		lenBody := len(r.Body)
		if lenBody > 0 {
			r.Headers[HeaderContentLength] = []string{strconv.Itoa(lenBody)}
		}
	}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

func (r *Request) handleBody() ([]byte, error) {
	var buf bytes.Buffer

	if MethodHasBody[r.Method] {
		buf.Write(r.Body)
		if r.MultiAttachmentMode == MultiAttachmentModeBoundary && len(r.FormData) > 0 {
			buf.WriteString("\r\n\r\n")
		}
	}

	return buf.Bytes(), nil
}

func (r *Request) handleUrlForm() ([]byte, error) {
	var buf bytes.Buffer

	lenI := len(r.FormData)
	i := 0
	for _, formData := range r.FormData {
		if formData.IsAttachment {
			continue
		}

		buf.WriteString(fmt.Sprintf("%s=%s", formData.Name, string(formData.Body)))

		if i < lenI-1 {
			buf.WriteString("&")
		}
		i++
	}

	return buf.Bytes(), nil
}

func (r *Request) handleSingleAttachment() ([]byte, error) {
	for _, attachment := range r.FormData {
		return attachment.Body, nil
	}
	return []byte{}, nil
}

func (r *Request) handleBoundaries() ([]byte, error) {
	var buf bytes.Buffer

	switch r.ContentType {
	case ContentTypeMultipartFormData, ContentTypeMultipartMixed:
		bufFormData, err := r.handleFormData()
		if err != nil {
			return buf.Bytes(), err
		}
		buf.Write(bufFormData)

	default:
		if r.MultiAttachmentMode != MultiAttachmentModeZip {
			bufAttachments, err := r.handleAttachments()
			if err != nil {
				return buf.Bytes(), err
			}
			buf.Write(bufAttachments)
		}
	}

	buf.WriteString(fmt.Sprintf("\r\n--%s--", r.Boundary))

	return buf.Bytes(), nil
}

func (r *Request) handleFormData() ([]byte, error) {
	var buf bytes.Buffer

	lenF := len(r.FormData)
	i := 0

	for _, formData := range r.FormData {
		buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
		buf.WriteString(fmt.Sprintf("%s: %s; name=%q\r\n", HeaderContentDisposition, formData.ContentDisposition, formData.Name))
		buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", HeaderContentType, formData.ContentType))
		buf.Write(formData.Body)

		if i < lenF-1 {
			buf.WriteString("\r\n")
		}
		i++
	}

	return buf.Bytes(), nil
}

func (r *Request) handleAttachments() ([]byte, error) {
	var buf bytes.Buffer

	for _, attachment := range r.Attachments {
		buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
		buf.WriteString(fmt.Sprintf("%s: %s; name=%q; filename=%q\r\n", HeaderContentDisposition, attachment.ContentDisposition, attachment.Name, attachment.FileName))
		buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", HeaderContentType, attachment.ContentType))
		buf.Write(attachment.Body)
		buf.WriteString("\r\n")
	}

	return buf.Bytes(), nil
}

func (r *Request) handleZippedAttachments() ([]byte, error) {
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

	for _, attachment := range r.Attachments {
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
