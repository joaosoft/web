package web

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"path/filepath"
	"strconv"
)

func (r *Request) Set(contentType ContentType, b []byte) error {
	r.ContentType = contentType
	r.Body = b
	return nil
}

func (r *Request) HTML(body string) error {
	r.SetContentType(ContentTypeTextHTML)
	r.Body = []byte(body)
	return nil
}

func (r *Request) Bytes(contentType ContentType, b []byte) error {
	r.SetContentType(contentType)
	r.Body = b
	return nil
}

func (r *Request) String(s string) error {
	r.SetContentType(ContentTypeTextPlain)
	r.Body = []byte(s)
	return nil
}

func (r *Request) JSON(i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParams["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.JSONPretty(i, "  ")
	}

	if b, err := json.Marshal(i); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationJSON)
		r.Body = b
	}

	return nil
}

func (r *Request) JSONPretty(i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationJSON)
		r.Body = b
	}
	return nil
}

func (r *Request) XML(i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParams["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.XMLPretty(i, "  ")
	}

	if b, err := xml.Marshal(i); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Request) XMLPretty(i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Request) Stream(contentType ContentType, reader io.Reader) error {
	r.SetContentType(contentType)
	if _, err := io.Copy(r.Writer, reader); err != nil {
		return err
	}
	return nil
}

func (r *Request) File(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.SetContentType(contentType)
	r.SetCharset(charset)
	r.Body = body
	return nil
}

func (r *Request) Attachment(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionAttachment,
			ContentType:        contentType,
			Charset:            charset,
			FileName:           name,
			Name:               name,
			Body:               body,
			IsAttachment:       true,
		},
	}
	return nil
}

func (r *Request) Inline(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionInline,
			ContentType:        contentType,
			Charset:            charset,
			FileName:           name,
			Name:               name,
			Body:               body,
			IsAttachment:       true,
		},
	}
	return nil
}

func (r *Request) SetFormData(name string, value string) {
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

func (r *Request) GetFormDataBytes(name string) []byte {
	if value, ok := r.FormData[name]; ok {
		return value.Body
	}

	return nil
}

func (r *Request) GetFormDataString(name string) string {
	if value, ok := r.FormData[name]; ok {
		return string(value.Body)
	}

	return ""
}
