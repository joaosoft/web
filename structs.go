package webserver

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Route struct {
	method      Method
	path        string
	name        string
	handler     HandlerFunc
	middlewares []MiddlewareFunc
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
type MimeType string

const (
	MIMEApplicationJSON       MimeType = "application/json"
	MIMEApplicationJavaScript MimeType = "application/javascript"
	MIMEApplicationXML        MimeType = "application/xml"
	MIMETextXML               MimeType = "text/xml"
	MIMEApplicationForm       MimeType = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf   MimeType = "application/protobuf"
	MIMEApplicationMsgpack    MimeType = "application/msgpack"
	MIMETextHTML              MimeType = "text/html"
	MIMETextPlain             MimeType = "text/plain"
	MIMEMultipartForm         MimeType = "multipart/form-data"
	MIMEOctetStream           MimeType = "application/octet-stream"
)

// Header
type Header string

const (
	HeaderAccept              Header = "Accept"
	HeaderAcceptEncoding      Header = "Accept-Encoding"
	HeaderAllow               Header = "Allow"
	HeaderAuthorization       Header = "Authorization"
	HeaderContentDisposition  Header = "Content-Disposition"
	HeaderContentEncoding     Header = "Content-Encoding"
	HeaderContentLength       Header = "Content-Length"
	HeaderContentType         Header = "Content-Type"
	HeaderCookie              Header = "Cookie"
	HeaderSetCookie           Header = "Set-Cookie"
	HeaderIfModifiedSince     Header = "If-Modified-Since"
	HeaderLastModified        Header = "Last-Modified"
	HeaderLocation            Header = "Location"
	HeaderUpgrade             Header = "Upgrade"
	HeaderVary                Header = "Vary"
	HeaderWWWAuthenticate     Header = "WWW-Authenticate"
	HeaderXForwardedFor       Header = "X-Forwarded-For"
	HeaderXForwardedProto     Header = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  Header = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       Header = "X-Forwarded-Ssl"
	HeaderXUrlScheme          Header = "X-Url-Scheme"
	HeaderXHTTPMethodOverride Header = "X-HTTP-Method-Override"
	HeaderXRealIP             Header = "X-Real-IP"
	HeaderXRequestID          Header = "X-Request-ID"
	HeaderXRequestedWith      Header = "X-Requested-With"
	HeaderServer              Header = "Server"

	// Access control
	HeaderAccessControlRequestMethod    Header = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   Header = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      Header = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     Header = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     Header = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials Header = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    Header = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           Header = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity Header = "Strict-Transport-Security"
	HeaderXContentTypeOptions     Header = "X-Content-Type-Options"
	HeaderXXSSProtection          Header = "X-XSS-Protection"
	HeaderXFrameOptions           Header = "X-Frame-Options"
	HeaderContentSecurityPolicy   Header = "Content-Security-Policy"
	HeaderXCSRFToken              Header = "X-CSRF-Token"
)
