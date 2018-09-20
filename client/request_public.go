package client

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"web/common"
)

func (r *Request) Set(contentType common.ContentType, b []byte) error {
	r.ContentType = contentType
	r.Body = b
	return nil
}

func (r *Request) HTML(body string) error {
	r.SetContentType(common.ContentTypeTextHTML)
	r.Body = []byte(body)
	return nil
}

func (r *Request) Bytes(contentType common.ContentType, b []byte) error {
	r.SetContentType(contentType)
	r.Body = b
	return nil
}

func (r *Request) String(s string) error {
	r.SetContentType(common.ContentTypeTextPlain)
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
		r.SetContentType(common.ContentTypeApplicationJSON)
		r.Body = b
	}

	return nil
}

func (r *Request) JSONPretty(i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(common.ContentTypeApplicationJSON)
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
		r.SetContentType(common.ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Request) XMLPretty(i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(common.ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Request) Stream(contentType common.ContentType, reader io.Reader) error {
	r.SetContentType(contentType)
	if _, err := io.Copy(r.Writer, reader); err != nil {
		return err
	}
	return nil
}

func (r *Request) File(fileName string) error {
	data, err := common.ReadFile(fileName, nil)
	if err != nil {
		return err
	}

	contentType, charset := common.DetectContentType(filepath.Ext(fileName), data)
	r.SetContentType(contentType)
	r.SetCharset(charset)
	r.Body = data
	return nil
}

func (r *Request) Attachment(file, name string) error {
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

func (r *Request) Inline(file, name string) error {
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
