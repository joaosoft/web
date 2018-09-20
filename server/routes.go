package server

import "webserver"

type Routes map[webserver.Method][]Route

type Route struct {
	Method      webserver.Method
	Path        string
	Regex       string
	Name        string
	Handler     HandlerFunc
	Middlewares []MiddlewareFunc
}

func NewRoute(method webserver.Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) Route {
	return Route{
		Method:      method,
		Path:        path,
		Regex:       ConvertPathToRegex(path),
		Handler:     handler,
		Middlewares: middleware,
		Name:        webserver.GetFunctionName(handler),
	}
}
