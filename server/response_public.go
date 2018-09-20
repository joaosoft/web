package server

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"webserver"
)

func (r *Response) Set(status webserver.Status, contentType webserver.ContentType, b []byte) error {
	r.ContentType = contentType
	r.Status = status
	r.Body = b
	return nil
}

func (r *Response) HTML(status webserver.Status, body string) error {
	r.SetContentType(webserver.ContentTypeTextHTML)
	r.Status = status
	r.Body = []byte(body)
	return nil
}

func (r *Response) Bytes(status webserver.Status, contentType webserver.ContentType, b []byte) error {
	r.SetContentType(contentType)
	r.Status = status
	r.Body = b
	return nil
}

func (r *Response) String(status webserver.Status, s string) error {
	r.SetContentType(webserver.ContentTypeTextPlain)
	r.Status = status
	r.Body = []byte(s)
	return nil
}

func (r *Response) JSON(status webserver.Status, i interface{}) error {
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
		r.SetContentType(webserver.ContentTypeApplicationJSON)
		r.Status = status
		r.Body = b
	}

	return nil
}

func (r *Response) JSONPretty(status webserver.Status, i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(webserver.ContentTypeApplicationJSON)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) XML(status webserver.Status, i interface{}) error {
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
		r.SetContentType(webserver.ContentTypeApplicationXML)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) XMLPretty(status webserver.Status, i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(webserver.ContentTypeApplicationXML)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) Stream(status webserver.Status, contentType webserver.ContentType, reader io.Reader) error {
	r.SetContentType(contentType)
	r.Status = status
	if _, err := io.Copy(r.Writer, reader); err != nil {
		return err
	}
	return nil
}

func (r *Response) File(fileName string) error {
	data, err := webserver.ReadFile(fileName, nil)
	if err != nil {
		return err
	}

	contentType, charset := webserver.DetectContentType(filepath.Ext(fileName), data)
	r.SetContentType(contentType)
	r.SetCharset(charset)
	r.Status = webserver.StatusOK
	r.Body = data
	return nil
}

func (r *Response) Attachment(file, name string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	data, err := webserver.ReadFile(file, nil)
	if err != nil {
		return err
	}

	contentType, charset := webserver.DetectContentType(filepath.Ext(info.Name()), data)
	r.Attachments[name] = Attachment{
		ContentDisposition: webserver.ContentDispositionAttachment,
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

	data, err := webserver.ReadFile(file, nil)
	if err != nil {
		return err
	}

	contentType, charset := webserver.DetectContentType(filepath.Ext(info.Name()), data)
	r.Attachments[name] = Attachment{
		ContentDisposition: webserver.ContentDispositionInline,
		ContentType:        contentType,
		Charset:            charset,
		File:               info.Name(),
		Name:               name,
		Body:               data,
	}
	return nil
}

func (r *Response) NoContent(status webserver.Status) error {
	r.Status = status
	return nil
}
