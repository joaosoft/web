package web

type Encoding string

const (
	EncodingChunked  Encoding = "chunked"
	EncodingCompress Encoding = "compress"
	EncodingDeflate  Encoding = "deflate"
	EncodingGzip     Encoding = "gzip"
	EncodingIdentity Encoding = "identity"
)
