package web

import (
	"io"
	"net"
	"time"
)

type Headers map[HeaderType][]string

type Cookie struct {
	Name    string
	Value   string
	Path    string    // optional
	Domain  string    // optional
	Expires time.Time // optional
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
}
type Cookies map[string]Cookie

type UrlParams map[string][]string

type Params map[string][]string

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
	Method      Method
	Protocol    Protocol
	Headers     Headers
	Cookies     Cookies
	ContentType ContentType
	Params      Params
	UrlParams   UrlParams
	Charset     Charset
	conn        net.Conn
	server      *Server
	client      *Client
}

type Request struct {
	Base
	Body                []byte
	Attachments         map[string]Attachment
	MultiAttachmentMode MultiAttachmentMode
	Boundary            string
	Reader              io.Reader
	Writer              io.Writer
}

type Response struct {
	Base
	Body                []byte
	Status              Status
	StatusText          string
	Attachments         map[string]Attachment
	MultiAttachmentMode MultiAttachmentMode
	Boundary            string
	Reader              io.Reader
	Writer              io.Writer
}

type Attachment struct {
	ContentType        ContentType
	ContentDisposition ContentDisposition
	Charset            Charset
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
