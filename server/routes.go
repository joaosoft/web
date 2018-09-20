package server

import "web/common"

type Routes map[common.Method][]Route

type Route struct {
	Method      common.Method
	Path        string
	Regex       string
	Name        string
	Handler     HandlerFunc
	Middlewares []MiddlewareFunc
}

func NewRoute(method common.Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) Route {
	return Route{
		Method:      method,
		Path:        path,
		Regex:       ConvertPathToRegex(path),
		Handler:     handler,
		Middlewares: middleware,
		Name:        common.GetFunctionName(handler),
	}
}
