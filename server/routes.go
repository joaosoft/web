package server

import "web"

type Routes map[web.Method][]Route

type Route struct {
	Method      web.Method
	Path        string
	Regex       string
	Name        string
	Handler     HandlerFunc
	Middlewares []MiddlewareFunc
}

func NewRoute(method web.Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) Route {
	return Route{
		Method:      method,
		Path:        path,
		Regex:       ConvertPathToRegex(path),
		Handler:     handler,
		Middlewares: middleware,
		Name:        web.GetFunctionName(handler),
	}
}
