package server

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"web/common"
)

func (r *Response) Set(status common.Status, contentType common.ContentType, b []byte) error {
	r.ContentType = contentType
	r.Status = status
	r.Body = b
	return nil
}

func (r *Response) HTML(status common.Status, body string) error {
	r.SetContentType(common.ContentTypeTextHTML)
	r.Status = status
	r.Body = []byte(body)
	return nil
}

func (r *Response) Bytes(status common.Status, contentType common.ContentType, b []byte) error {
	r.SetContentType(contentType)
	r.Status = status
	r.Body = b
	return nil
}

func (r *Response) String(status common.Status, s string) error {
	r.SetContentType(common.ContentTypeTextPlain)
	r.Status = status
	r.Body = []byte(s)
	return nil
}

func (r *Response) JSON(status common.Status, i interface{}) error {
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
		r.SetContentType(common.ContentTypeApplicationJSON)
		r.Status = status
		r.Body = b
	}

	return nil
}

func (r *Response) JSONPretty(status common.Status, i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(common.ContentTypeApplicationJSON)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) XML(status common.Status, i interface{}) error {
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
		r.SetContentType(common.ContentTypeApplicationXML)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) XMLPretty(status common.Status, i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(common.ContentTypeApplicationXML)
		r.Status = status
		r.Body = b
	}
	return nil
}

func (r *Response) Stream(status common.Status, contentType common.ContentType, reader io.Reader) error {
	r.SetContentType(contentType)
	r.Status = status
	if _, err := io.Copy(r.Writer, reader); err != nil {
		return err
	}
	return nil
}

func (r *Response) File(fileName string) error {
	data, err := common.ReadFile(fileName, nil)
	if err != nil {
		return err
	}

	contentType, charset := common.DetectContentType(filepath.Ext(fileName), data)
	r.SetContentType(contentType)
	r.SetCharset(charset)
	r.Status = common.StatusOK
	r.Body = data
	return nil
}

func (r *Response) Attachment(file, name string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	data, err := common.ReadFile(file, nil)
	if err != nil {
		return err
	}

	contentType, charset := common.DetectContentType(filepath.Ext(info.Name()), data)
	r.Attachments[name] = Attachment{
		ContentDisposition: common.ContentDispositionAttachment,
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

	data, err := common.ReadFile(file, nil)
	if err != nil {
		return err
	}

	contentType, charset := common.DetectContentType(filepath.Ext(info.Name()), data)
	r.Attachments[name] = Attachment{
		ContentDisposition: common.ContentDispositionInline,
		ContentType:        contentType,
		Charset:            charset,
		File:               info.Name(),
		Name:               name,
		Body:               data,
	}
	return nil
}

func (r *Response) NoContent(status common.Status) error {
	r.Status = status
	return nil
}
