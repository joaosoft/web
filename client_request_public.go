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
	r.Attachments[name] = Attachment{
		ContentDisposition: ContentDispositionAttachment,
		ContentType:        contentType,
		Charset:            charset,
		File:               name,
		Name:               name,
		Body:               body,
	}
	return nil
}

func (r *Request) Inline(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.Attachments[name] = Attachment{
		ContentDisposition: ContentDispositionInline,
		ContentType:        contentType,
		Charset:            charset,
		File:               name,
		Name:               name,
		Body:               body,
	}
	return nil
}
