package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func (r *Response) SetHeader(name HeaderType, header string) {
	r.Headers[name] = Header{header}
}

func (r *Response) GetHeader(name HeaderType) string {
	if header, ok := r.Headers[name]; ok {
		return header[0]
	}
	return ""
}

func (r *Response) SetContentType(contentType ContentType) {
	r.Headers[HeaderContentType] = Header{string(contentType)}
}

func (r *Response) GetContentType() *ContentType {
	if value, ok := r.Headers[HeaderContentType]; ok {
		contentType := ContentType(value[0])
		return &contentType
	}
	return nil
}

func (r *Response) SetCookie(name string, cookie Cookie) {
	r.Cookies[name] = cookie
}

func (r *Response) GetCookie(name string) *Cookie {
	if cookie, ok := r.Cookies[name]; ok {
		return &cookie
	}
	return nil
}

func (r *Response) SetParam(name string, param Param) {
	r.Params[name] = param
}

func (r *Response) GetParam(name string) []string {
	if param, ok := r.Params[name]; ok {
		return param
	}
	return nil
}

func (r *Response) Set(status Status, contentType ContentType, b []byte) error {
	r.SetContentType(contentType)
	r.Status = status
	r.Body = b
	return nil
}

func (r *Response) HTML(status Status, body string) error {
	r.SetContentType(ContentTextHTML)
	r.Status = status
	r.Body = []byte(body)
	return nil
}

func (r *Response) Bytes(status Status, contentType ContentType, b []byte) error {
	r.SetContentType(contentType)
	r.Status = status
	r.Body = b
	return nil
}

func (r *Response) String(status Status, s string) error {
	r.SetContentType(ContentTextPlain)
	r.Status = status
	r.Body = []byte(s)
	return nil
}

func (r *Response) JSON(status Status, i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParms["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.JSONPretty(status, i, "  ")
	}

	if b, err := json.Marshal(i); err != nil {
		return err
	} else {
		r.SetContentType(ContentApplicationJSON)
		r.Status = status
		r.Body = b
	}

	return nil
}

func (r *Response) JSONPretty(status Status, i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(ContentApplicationJSON)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) XML(status Status, i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParms["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.XMLPretty(status, i, "  ")
	}

	if b, err := xml.Marshal(i); err != nil {
		return err
	} else {
		r.SetContentType(ContentApplicationXML)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) XMLPretty(status Status, i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(ContentApplicationXML)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) Stream(status Status, contentType ContentType, reader io.Reader) error {
	r.SetContentType(contentType)
	r.Status = status
	if _, err := io.Copy(r.Writer, reader); err != nil {
		return err
	}
	return nil
}

func (r *Response) File(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	r.SetContentType(ContentOctetStream)
	r.Status = StatusOK
	r.Body = data
	return nil
}

func (r *Response) Attachment(file, name string) error {
	r.SetHeader(HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", name))
	return r.File(file)
}

func (r *Response) Inline(file, name string) error {
	r.SetHeader(HeaderContentDisposition, fmt.Sprintf("inline; filename=%q", name))
	return r.File(file)
}

func (r *Response) NoContent(status Status) error {
	r.Status = status
	return nil
}
