package web

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"path/filepath"
	"strconv"
)

func (r *Response) Set(status Status, contentType ContentType, b []byte) error {
	r.Status = status
	r.ContentType = contentType
	r.Body = b
	return nil
}

func (r *Response) HTML(status Status, body string) error {
	r.Status = status
	r.SetContentType(ContentTypeTextHTML)
	r.Body = []byte(body)
	return nil
}

func (r *Response) Bytes(status Status, contentType ContentType, b []byte) error {
	r.Status = status
	r.SetContentType(contentType)
	r.Body = b
	return nil
}

func (r *Response) String(status Status, s string) error {
	r.Status = status
	r.SetContentType(ContentTypeTextPlain)
	r.Body = []byte(s)
	return nil
}

func (r *Response) JSON(status Status, i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParams["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.JSONPretty(status, i, "  ")
	}

	if b, err := json.Marshal(i); err != nil {
		return err
	} else {
		r.Status = status
		r.SetContentType(ContentTypeApplicationJSON)
		r.Body = b
	}

	return nil
}

func (r *Response) JSONPretty(status Status, i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.Status = status
		r.SetContentType(ContentTypeApplicationJSON)
		r.Body = b
	}
	return nil
}

func (r *Response) XML(status Status, i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParams["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.XMLPretty(status, i, "  ")
	}

	if b, err := xml.Marshal(i); err != nil {
		return err
	} else {
		r.Status = status
		r.SetContentType(ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Response) XMLPretty(status Status, i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.Status = status
		r.SetContentType(ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Response) Stream(status Status, contentType ContentType, reader io.Reader) error {
	r.Status = status
	r.SetContentType(contentType)
	if _, err := io.Copy(r.Writer, reader); err != nil {
		return err
	}
	return nil
}

func (r *Response) File(status Status, name string, body []byte) error {
	r.Status = status
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.SetContentType(contentType)
	r.SetCharset(charset)
	r.Body = body
	return nil
}

func (r *Response) Attachment(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionAttachment,
			ContentType:        contentType,
			Charset:            charset,
			FileName:           name,
			Name:               name,
			Body:               body,
		},
	}
	return nil
}

func (r *Response) Inline(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionInline,
			ContentType:        contentType,
			Charset:            charset,
			FileName:           name,
			Name:               name,
			Body:               body,
		},
	}
	return nil
}

func (r *Response) NoContent(status Status) error {
	r.Status = status
	return nil
}

func (r *Response) SetFormData(name string, value string) {
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionFormData,
			ContentType:        ContentTypeTextPlain,
			Charset:            CharsetUTF8,
			Name:               name,
			Body:               []byte(value),
			IsAttachment:       false,
		},
	}
}

func (r *Response) GetFormDataBytes(name string) []byte {
	if value, ok := r.FormData[name]; ok {
		return value.Body
	}

	return nil
}

func (r *Response) GetFormDataString(name string) string {
	if value, ok := r.FormData[name]; ok {
		return string(value.Body)
	}

	return ""
}
