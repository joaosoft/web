package client

import (
	"io"
	"net"
	"time"
	"web"
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
	Method      web.Method
	Protocol    web.Protocol
	Headers     web.Headers
	Cookies     web.Cookies
	ContentType web.ContentType
	Params      web.Params
	UrlParams   web.UrlParams
	Charset     web.Charset
	conn        net.Conn
	client      *Client
}

type Request struct {
	Base
	Body                []byte
	Attachments         map[string]Attachment
	MultiAttachmentMode web.MultiAttachmentMode
	Boundary            string
	Writer              io.Writer
}

type Response struct {
	Base
	Body        []byte
	Status      web.Status
	StatusText  string
	Attachments map[string]Attachment
	Boundary    string
	Reader      io.Reader
}

type Attachment struct {
	ContentType        web.ContentType
	ContentDisposition web.ContentDisposition
	Charset            web.Charset
	File               string
	Name               string
	Body               []byte
}
