package client

import (
	"io"
	"net"
	"time"
	"webserver"
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
	Method      webserver.Method
	Protocol    string
	Headers     webserver.Headers
	Cookies     webserver.Cookies
	ContentType webserver.ContentType
	Params      webserver.Params
	UrlParams   webserver.UrlParams
	Charset     webserver.Charset
	conn        net.Conn
	client      *WebClient
}

type Request struct {
	Base
	Body        []byte
	Attachments map[string]Attachment
	Boundary    string
	Reader      io.Reader
}

type Response struct {
	Base
	Body                []byte
	Status              webserver.Status
	Attachments         map[string]Attachment
	MultiAttachmentMode webserver.MultiAttachmentMode
	Boundary            string
	Writer              io.Writer
}

type Attachment struct {
	ContentType        webserver.ContentType
	ContentDisposition webserver.ContentDisposition
	Charset            webserver.Charset
	File               string
	Name               string
	Body               []byte
}
