package web

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func (r *Response) Set(status Status, contentType ContentType, b []byte) error {
	r.ContentType = contentType
	r.Status = status
	r.Body = b
	return nil
}

func (r *Response) HTML(status Status, body string) error {
	r.SetContentType(ContentTypeTextHTML)
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
	r.SetContentType(ContentTypeTextPlain)
	r.Status = status
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
		r.SetContentType(ContentTypeApplicationJSON)
		r.Status = status
		r.Body = b
	}

	return nil
}

func (r *Response) JSONPretty(status Status, i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationJSON)
		r.Status = status
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
		r.SetContentType(ContentTypeApplicationXML)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) XMLPretty(status Status, i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationXML)
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

func (r *Response) File(status Status, name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.SetContentType(contentType)
	r.SetCharset(charset)
	r.Body = body
	return nil
}

func (r *Response) ReadFile(status Status, fileName string) error {
	data, err := ReadFile(fileName, nil)
	if err != nil {
		return err
	}

	return r.File(status, fileName, data)
}

func (r *Response) Attachment(file, name string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	data, err := ReadFile(file, nil)
	if err != nil {
		return err
	}

	contentType, charset := DetectContentType(filepath.Ext(info.Name()), data)
	r.Attachments[name] = Attachment{
		ContentDisposition: ContentDispositionAttachment,
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

	data, err := ReadFile(file, nil)
	if err != nil {
		return err
	}

	contentType, charset := DetectContentType(filepath.Ext(info.Name()), data)
	r.Attachments[name] = Attachment{
		ContentDisposition: ContentDispositionInline,
		ContentType:        contentType,
		Charset:            charset,
		File:               info.Name(),
		Name:               name,
		Body:               data,
	}
	return nil
}

func (r *Response) NoContent(status Status) error {
	r.Status = status
	return nil
}
