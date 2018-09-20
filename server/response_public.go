package server

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"web"
)

func (r *Response) Set(status web.Status, contentType web.ContentType, b []byte) error {
	r.ContentType = contentType
	r.Status = status
	r.Body = b
	return nil
}

func (r *Response) HTML(status web.Status, body string) error {
	r.SetContentType(web.ContentTypeTextHTML)
	r.Status = status
	r.Body = []byte(body)
	return nil
}

func (r *Response) Bytes(status web.Status, contentType web.ContentType, b []byte) error {
	r.SetContentType(contentType)
	r.Status = status
	r.Body = b
	return nil
}

func (r *Response) String(status web.Status, s string) error {
	r.SetContentType(web.ContentTypeTextPlain)
	r.Status = status
	r.Body = []byte(s)
	return nil
}

func (r *Response) JSON(status web.Status, i interface{}) error {
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
		r.SetContentType(web.ContentTypeApplicationJSON)
		r.Status = status
		r.Body = b
	}

	return nil
}

func (r *Response) JSONPretty(status web.Status, i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(web.ContentTypeApplicationJSON)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) XML(status web.Status, i interface{}) error {
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
		r.SetContentType(web.ContentTypeApplicationXML)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) XMLPretty(status web.Status, i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(web.ContentTypeApplicationXML)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) Stream(status web.Status, contentType web.ContentType, reader io.Reader) error {
	r.SetContentType(contentType)
	r.Status = status
	if _, err := io.Copy(r.Writer, reader); err != nil {
		return err
	}
	return nil
}

func (r *Response) File(fileName string) error {
	data, err := web.ReadFile(fileName, nil)
	if err != nil {
		return err
	}

	contentType, charset := web.DetectContentType(filepath.Ext(fileName), data)
	r.SetContentType(contentType)
	r.SetCharset(charset)
	r.Status = web.StatusOK
	r.Body = data
	return nil
}

func (r *Response) Attachment(file, name string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	data, err := web.ReadFile(file, nil)
	if err != nil {
		return err
	}

	contentType, charset := web.DetectContentType(filepath.Ext(info.Name()), data)
	r.Attachments[name] = Attachment{
		ContentDisposition: web.ContentDispositionAttachment,
		ContentType:        contentType,
		Charset:            charset,
		File:               info.Name(),
		Name:               name,
		Body:               data,
	}
	return nil
}

func (r *Response) Inline(file, name string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	data, err := web.ReadFile(file, nil)
	if err != nil {
		return err
	}

	contentType, charset := web.DetectContentType(filepath.Ext(info.Name()), data)
	r.Attachments[name] = Attachment{
		ContentDisposition: web.ContentDispositionInline,
		ContentType:        contentType,
		Charset:            charset,
		File:               info.Name(),
		Name:               name,
		Body:               data,
	}
	return nil
}

func (r *Response) NoContent(status web.Status) error {
	r.Status = status
	return nil
}
