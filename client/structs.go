package client

import (
	"io"
	"net"
	"time"
	"web/common"
)

type Context struct {
	StartTime time.Time
	Request   *Request
	Response  *Response
}

type Base struct {
	IP          string
	FullUrl     string
	Url         string
	Method      common.Method
	Protocol    common.Protocol
	Headers     common.Headers
	Cookies     common.Cookies
	ContentType common.ContentType
	Params      common.Params
	UrlParams   common.UrlParams
	Charset     common.Charset
	conn        net.Conn
	client      *Client
}

type Request struct {
	Base
	Body                []byte
	Attachments         map[string]Attachment
	MultiAttachmentMode common.MultiAttachmentMode
	Boundary            string
	Writer              io.Writer
}

type Response struct {
	Base
	Body        []byte
	Status      common.Status
	StatusText  string
	Attachments map[string]Attachment
	Boundary    string
	Reader      io.Reader
}

type Attachment struct {
	ContentType        common.ContentType
	ContentDisposition common.ContentDisposition
	Charset            common.Charset
	File               string
	Name               string
	Body               []byte
}
