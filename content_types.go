package web

// Mime type
type ContentType string

const (
	ContentApplicationJSON       ContentType = "application/json"
	ContentApplicationJavaScript ContentType = "application/javascript"
	ContentApplicationXML        ContentType = "application/xml"
	ContentTextXML               ContentType = "text/xml"
	ContentApplicationForm       ContentType = "application/x-www-form-urlencoded"
	ContentApplicationProtobuf   ContentType = "application/protobuf"
	ContentApplicationMsgpack    ContentType = "application/msgpack"
	ContentTextHTML              ContentType = "text/html"
	ContentTextPlain             ContentType = "text/plain"
	ContentMultipartForm         ContentType = "multipart/form-data"
	ContentOctetStream           ContentType = "application/octet-stream"
)
