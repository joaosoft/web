package server

import (
	"io"
	"net"
	"time"
	"web/common"
)

type ErrorHandler func(ctx *Context, err error) error
type HandlerFunc func(ctx *Context) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

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
	server      *Server
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
	Status              common.Status
	Attachments         map[string]Attachment
	MultiAttachmentMode common.MultiAttachmentMode
	Boundary            string
	Writer              io.Writer
}

type Attachment struct {
	ContentType        common.ContentType
	ContentDisposition common.ContentDisposition
	Charset            common.Charset
	File               string
	Name               string
	Body               []byte
}

type RequestHandler struct {
	Conn    net.Conn
	Handler HandlerFunc
}

type Namespace struct {
	Path        string
	Middlewares []MiddlewareFunc
	WebServer   *Server
}
