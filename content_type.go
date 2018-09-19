package webserver

// Mime type
type ContentType string

const (
	ContentTypeApplicationJSON       ContentType = "application/json"
	ContentTypeApplicationJavaScript ContentType = "application/javascript"
	ContentTypeApplicationXML        ContentType = "application/xml"
	ContentTypeTextXML             ContentType = "text/xml"
	ContentTypeApplicationForm     ContentType = "application/x-www-form-urlencoded"
	ContentTypeApplicationProtobuf ContentType = "application/protobuf"
	ContentTypeApplicationMsgpack  ContentType = "application/msgpack"
	ContentTypeTextHTML            ContentType = "text/html"
	ContentTypeTextPlain           ContentType = "text/plain"
	ContentTypeMultipartFormData   ContentType = "multipart/form-data"
	ContentTypeMultipartMixed      ContentType = "multipart/mixed"
	ContentTypeOctetStream         ContentType = "application/octet-stream"
	ContentTypeZip                 ContentType = "application/zip"
)
