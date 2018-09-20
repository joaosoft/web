package webserver

import "time"

type Method string

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
