package webserver

import (
	"net"
	"time"
)

type HandlerFunc func(ctx *Context) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Header []string
type Headers map[string]Header

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

type Context struct {
	StartTime time.Time
	Request   *Request
	Response  *Response
}

type Base struct {
	Url      string
	Method   Method
	Protocol string
	Headers  Headers
	Cookies  Cookies
}

type Request struct {
	Base
	Body []byte
}

type Response struct {
	Base
	Body   []byte
	Status int
}

type RequestHandler struct {
	Conn    net.Conn
	Handler HandlerFunc
}

type Routes map[string]*Route

type Route struct {
	Method      Method
	Path        string
	Name        string
	Handler     HandlerFunc
	Middlewares []MiddlewareFunc
}

// Method
type Method string

const (
	MethodConnect Method = "CONNECT"
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodPatch   Method = "PATCH"
	MethodDelete  Method = "DELETE"
	MethodOptions Method = "OPTIONS"
	MethodTrace   Method = "TRACE"
)

var (
	methods = []Method{
		MethodGet,
		MethodHead,
		MethodConnect,
		MethodDelete,
		MethodOptions,
		MethodPatch,
		MethodPost,
		MethodTrace,
		MethodPut,
	}
)

// Mime type
type ContentType string

const (
	MIMEApplicationJSON       ContentType = "application/json"
	MIMEApplicationJavaScript ContentType = "application/javascript"
	MIMEApplicationXML        ContentType = "application/xml"
	MIMETextXML               ContentType = "text/xml"
	MIMEApplicationForm       ContentType = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf   ContentType = "application/protobuf"
	MIMEApplicationMsgpack    ContentType = "application/msgpack"
	MIMETextHTML              ContentType = "text/html"
	MIMETextPlain             ContentType = "text/plain"
	MIMEMultipartForm         ContentType = "multipart/form-data"
	MIMEOctetStream           ContentType = "application/octet-stream"
)

// HeaderType
type HeaderType string

const (
	HeaderTypeAccept              HeaderType = "Accept"
	HeaderTypeAcceptEncoding      HeaderType = "Accept-Encoding"
	HeaderTypeAllow               HeaderType = "Allow"
	HeaderTypeAuthorization       HeaderType = "Authorization"
	HeaderTypeContentDisposition  HeaderType = "Content-Disposition"
	HeaderTypeContentEncoding     HeaderType = "Content-Encoding"
	HeaderTypeContentLength       HeaderType = "Content-Length"
	HeaderTypeContentType         HeaderType = "Content-Type"
	HeaderTypeCookie              HeaderType = "Cookie"
	HeaderTypeSetCookie           HeaderType = "Set-Cookie"
	HeaderTypeIfModifiedSince     HeaderType = "If-Modified-Since"
	HeaderTypeLastModified        HeaderType = "Last-Modified"
	HeaderTypeLocation            HeaderType = "Location"
	HeaderTypeUpgrade             HeaderType = "Upgrade"
	HeaderTypeVary                HeaderType = "Vary"
	HeaderTypeWWWAuthenticate     HeaderType = "WWW-Authenticate"
	HeaderTypeXForwardedFor       HeaderType = "X-Forwarded-For"
	HeaderTypeXForwardedProto     HeaderType = "X-Forwarded-Proto"
	HeaderTypeXForwardedProtocol  HeaderType = "X-Forwarded-Protocol"
	HeaderTypeXForwardedSsl       HeaderType = "X-Forwarded-Ssl"
	HeaderTypeXUrlScheme          HeaderType = "X-Url-Scheme"
	HeaderTypeXHTTPMethodOverride HeaderType = "X-HTTP-Method-Override"
	HeaderTypeXRealIP             HeaderType = "X-Real-IP"
	HeaderTypeXRequestID          HeaderType = "X-Request-ID"
	HeaderTypeXRequestedWith      HeaderType = "X-Requested-With"
	HeaderTypeServer              HeaderType = "Server"

	// Access control
	HeaderTypeAccessControlRequestMethod      HeaderType = "Access-Control-Request-Method"
	HeaderTypeAccessControlRequestHeaderTypes HeaderType = "Access-Control-Request-HeaderTypes"
	HeaderTypeAccessControlAllowOrigin        HeaderType = "Access-Control-Allow-Origin"
	HeaderTypeAccessControlAllowMethods       HeaderType = "Access-Control-Allow-Methods"
	HeaderTypeAccessControlAllowHeaderTypes   HeaderType = "Access-Control-Allow-HeaderTypes"
	HeaderTypeAccessControlAllowCredentials   HeaderType = "Access-Control-Allow-Credentials"
	HeaderTypeAccessControlExposeHeaderTypes  HeaderType = "Access-Control-Expose-HeaderTypes"
	HeaderTypeAccessControlMaxAge             HeaderType = "Access-Control-Max-Age"

	// Security
	HeaderTypeStrictTransportSecurity HeaderType = "Strict-Transport-Security"
	HeaderTypeXContentTypeOptions     HeaderType = "X-Content-Type-Options"
	HeaderTypeXXSSProtection          HeaderType = "X-XSS-Protection"
	HeaderTypeXFrameOptions           HeaderType = "X-Frame-Options"
	HeaderTypeContentSecurityPolicy   HeaderType = "Content-Security-Policy"
	HeaderTypeXCSRFToken              HeaderType = "X-CSRF-Token"
)
