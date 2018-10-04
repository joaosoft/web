package web

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/joaosoft/color"
	"github.com/joaosoft/logger"
)

type Server struct {
	config              *ServerConfig
	isLogExternal       bool
	logger              logger.ILogger
	routes              Routes
	middlewares         []MiddlewareFunc
	listener            net.Listener
	address             string
	errorhandler        ErrorHandler
	multiAttachmentMode MultiAttachmentMode
}

func NewServer(options ...ServerOption) (*Server, error) {
	log := logger.NewLogDefault("server", logger.WarnLevel)

	service := &Server{
		logger:              log,
		routes:              make(Routes),
		middlewares:         make([]MiddlewareFunc, 0),
		address:             ":80",
		multiAttachmentMode: MultiAttachmentModeZip,
	}

	if service.isLogExternal {
		// set logger of internal processes
	}

	// load configuration File
	appConfig := &AppConfig{}
	if err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		service.logger.Warn(err)
	} else {
		level, _ := logger.ParseLevel(appConfig.Server.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.config = &appConfig.Server
	if appConfig.Server.Address != "" {
		service.address = appConfig.Server.Address
	}

	service.Reconfigure(options...)

	service.AddRoute(MethodGet, "/favicon.ico", service.handlerFile)
	service.errorhandler = service.DefaultErrorHandler

	return service, nil
}

func (w *Server) AddMiddlewares(middlewares ...MiddlewareFunc) {
	w.middlewares = append(w.middlewares, middlewares...)
}

func (w *Server) AddRoute(method Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	w.routes[method] = append(w.routes[method], Route{
		Method:      method,
		Path:        path,
		Regex:       ConvertPathToRegex(path),
		Handler:     handler,
		Middlewares: middleware,
		Name:        GetFunctionName(handler),
	})

	return nil
}

func (w *Server) AddRoutes(route ...Route) error {
	for _, r := range route {
		if err := w.AddRoute(r.Method, r.Path, r.Handler, r.Middlewares...); err != nil {
			return err
		}
	}
	return nil
}

func (w *Server) AddNamespace(path string, middlewares ...MiddlewareFunc) *Namespace {
	return &Namespace{
		Path:        path,
		Middlewares: middlewares,
		WebServer:   w,
	}
}

func (n *Namespace) AddRoutes(route ...Route) error {
	for _, r := range route {
		if err := n.WebServer.AddRoute(r.Method, fmt.Sprintf("%s%s", n.Path, r.Path), r.Handler, append(r.Middlewares, n.Middlewares...)...); err != nil {
			return err
		}
	}
	return nil
}

func (w *Server) SetErrorHandler(handler ErrorHandler) error {
	w.errorhandler = handler
	return nil
}

func (w *Server) Start() error {
	w.logger.Debug("executing Start")
	var err error

	w.listener, err = net.Listen("tcp", w.address)
	if err != nil {
		w.logger.Errorf("error connecting to %s: %s", w.address, err)
		return err
	}
	fmt.Println(color.WithColor("http server started on [%s]", color.FormatBold, color.ForegroundRed, color.BackgroundBlack, w.address))

	for {
		conn, err := w.listener.Accept()
		w.logger.Info("accepted connection")
		if err != nil {
			w.logger.Errorf("error accepting connection: %s", err)
			continue
		}

		if conn == nil {
			w.logger.Error("the connection isn't initialized")
			continue
		}

		go w.handleConnection(conn)
	}

	return err
}

func (w *Server) handleConnection(conn net.Conn) (err error) {
	var ctx *Context
	var handler HandlerFunc
	var length int
	startTime := time.Now()

	defer func() {
		conn.Close()
	}()

	// read response from connection
	request, err := w.NewRequest(conn, w)
	if err != nil {
		w.logger.Errorf("error getting request: [%s]", err)
		return err
	}

	fmt.Println(color.WithColor("[IN] Address[%s] Method[%s] Url[%s] Protocol[%s] Start[%s]", color.FormatBold, color.ForegroundBlue, color.BackgroundBlack, request.IP, request.Method, request.Url, request.Protocol, startTime))

	// create response for request
	response := w.NewResponse(request)

	// create context with request and response
	ctx = NewContext(startTime, request, response)
	var route *Route

	// when options method, validate request route
	if request.Method == MethodOptions {
		var method Method
		if val, ok := request.Headers[HeaderAccessControlRequestMethod]; ok {
			method = Method(val[0])
		}
		route, err = w.GetRoute(method, request.Url)
		if err != nil {
			w.logger.Errorf("error getting route: [%s]", err)
			goto on_error
		}
	}

	// middleware's of the server
	route, err = w.GetRoute(request.Method, request.Url)
	if err != nil {
		w.logger.Errorf("error getting route: [%s]", err)
		goto on_error
	}

	// get url parameters
	if err = w.LoadUrlParms(request, route); err != nil {
		w.logger.Errorf("error loading url parameters: [%s]", err)
		goto on_error
	}

	// route handler
	handler = route.Handler

	// execute middlewares
	length = len(w.middlewares)
	for i, _ := range w.middlewares {
		if w.middlewares[length-1-i] != nil {
			handler = w.middlewares[length-1-i](handler)
		}
	}

	// middleware's of the specific route
	length = len(route.Middlewares)
	for i, _ := range route.Middlewares {
		if route.Middlewares[length-1-i] != nil {
			handler = route.Middlewares[length-1-i](handler)
		}
	}

	// run handlers with middleware's
	if err = handler(ctx); err != nil {
		w.logger.Errorf("error executing handler: [%s]", err)
		goto on_error
	}

on_error:
	if err != nil {
		w.errorhandler(ctx, err)
	}

	// write response to connection
	if err = ctx.Response.write(); err != nil {
		w.logger.Errorf("error writing response: [%s]", err)
	}

	fmt.Println(color.WithColor("[OUT] Address[%s] Method[%s] Url[%s] Protocol[%s] Start[%s] Elapsed[%s]", color.FormatBold, color.ForegroundCyan, color.BackgroundBlack, ctx.Request.IP, ctx.Request.Method, ctx.Request.Url, ctx.Request.Protocol, startTime, time.Since(startTime)))

	return nil
}

func (w *Server) Stop() error {
	w.logger.Debug("executing Stop")

	if w.listener != nil {
		w.listener.Close()
	}

	return nil
}

func ConvertPathToRegex(path string) string {

	var re = regexp.MustCompile(`:[a-zA-Z0-9\-_]+`)
	regx := re.ReplaceAllString(path, `[a-zA-Z0-9-_]+`)

	return fmt.Sprintf("^%s$", regx)
}

func (w *Server) GetRoute(method Method, url string) (*Route, error) {

	for _, route := range w.routes[method] {
		if regx, err := regexp.Compile(route.Regex); err != nil {
			return nil, err
		} else {
			if regx.MatchString(url) {
				return &route, nil
			}
		}
	}

	return nil, ErrorNotFound
}

func (w *Server) LoadUrlParms(request *Request, route *Route) error {

	routeUrl := strings.Split(route.Path, "/")
	url := strings.Split(request.Url, "/")

	for i, name := range routeUrl {
		if name != url[i] {
			request.UrlParams[name[1:]] = []string{url[i]}
		}
	}

	return nil
}
