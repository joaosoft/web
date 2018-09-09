package webserver

import (
	"io"
	"net"
	"time"
)

type ErrorHandler func(ctx *Context, err error) error
type HandlerFunc func(ctx *Context) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Header []string
type Headers map[HeaderType]Header

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

type UrlParm []string
type UrlParms map[string]UrlParm

type Param []string
type Params map[string]Param

type Context struct {
	StartTime time.Time
	Request   *Request
	Response  *Response
}

type Base struct {
	IP       string
	FullUrl  string
	Url      string
	Method   Method
	Protocol string
	Headers  Headers
	Cookies  Cookies
	Params   Params
	UrlParms UrlParms
	conn     net.Conn
}

type Request struct {
	Base
	Body   []byte
	Reader io.Reader
}

type Response struct {
	Base
	Body   []byte
	Status Status
	Writer io.Writer
}

type RequestHandler struct {
	Conn    net.Conn
	Handler HandlerFunc
}



type Method string

type Namespace struct {
	Path        string
	Middlewares []MiddlewareFunc
	WebServer   *WebServer
}
