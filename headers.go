package web

// HeaderType
type HeaderType string

const (
	// Headers
	HeaderAccept              HeaderType = "Accept"
	HeaderAcceptEncoding      HeaderType = "Accept-Encoding"
	HeaderAllow               HeaderType = "Allow"
	HeaderAuthorization       HeaderType = "Authorization"
	HeaderContentDisposition  HeaderType = "Content-Disposition"
	HeaderContentEncoding     HeaderType = "Content-Encoding"
	HeaderContentLength       HeaderType = "Content-Length"
	HeaderContentType         HeaderType = "Content-Type"
	HeaderCookie              HeaderType = "Cookie"
	HeaderSetCookie           HeaderType = "Set-Cookie"
	HeaderIfModifiedSince     HeaderType = "If-Modified-Since"
	HeaderLastModified        HeaderType = "Last-Modified"
	HeaderLocation            HeaderType = "Location"
	HeaderUpgrade             HeaderType = "Upgrade"
	HeaderVary                HeaderType = "Vary"
	HeaderWWWAuthenticate     HeaderType = "WWW-Authenticate"
	HeaderXForwardedFor       HeaderType = "X-Forwarded-For"
	HeaderXForwardedProto     HeaderType = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  HeaderType = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       HeaderType = "X-Forwarded-Ssl"
	HeaderXUrlScheme          HeaderType = "X-Url-Scheme"
	HeaderXHTTPMethodOverride HeaderType = "X-HTTP-Method-Override"
	HeaderXRealIP             HeaderType = "X-Real-IP"
	HeaderXRequestID          HeaderType = "X-Request-ID"
	HeaderXRequestedWith      HeaderType = "X-Requested-With"
	HeaderServer              HeaderType = "Server"
	HeaderDate                HeaderType = "Date"
	HeaderMimeVersion         HeaderType = "MIME-Version"

	// Access control
	HeaderAccessControlRequestMethod      HeaderType = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaderTypes HeaderType = "Access-Control-Request-HeaderTypes"
	HeaderAccessControlAllowOrigin        HeaderType = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods       HeaderType = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaderTypes   HeaderType = "Access-Control-Allow-HeaderTypes"
	HeaderAccessControlAllowCredentials   HeaderType = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaderTypes  HeaderType = "Access-Control-Expose-HeaderTypes"
	HeaderAccessControlMaxAge             HeaderType = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity HeaderType = "Strict-Transport-Security"
	HeaderXContentTypeOptions     HeaderType = "X-Content-Type-Options"
	HeaderXXSSProtection          HeaderType = "X-XSS-Protection"
	HeaderXFrameOptions           HeaderType = "X-Frame-Options"
	HeaderContentSecurityPolicy   HeaderType = "Content-Security-Policy"
	HeaderXCSRFToken              HeaderType = "X-CSRF-Token"
)
