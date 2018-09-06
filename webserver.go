package webserver

import (
	"fmt"
	"net"

	"time"

	"bytes"

	"regexp"

	"strings"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type WebServer struct {
	config        *WebServerConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
	routes        Routes
	middlewares   []MiddlewareFunc
	listener      net.Listener
	port          int
}

func NewWebServer(options ...WebServerOption) (*WebServer, error) {
	pm := manager.NewManager(manager.WithRunInBackground(true), manager.WithLogLevel(logger.NoneLevel))
	log := logger.NewLogDefault("webserver", logger.WarnLevel)

	service := &WebServer{
		pm:          pm,
		logger:      log,
		routes:      make(Routes),
		middlewares: make([]MiddlewareFunc, 0),
		port:        9001,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	// load configuration File
	appConfig := &AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		service.logger.Warn(err)
	} else {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Dependency.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.config = &appConfig.Dependency
	service.Reconfigure(options...)

	return service, nil
}

func (w *WebServer) AddMiddlewares(middlewares ...MiddlewareFunc) {
	w.middlewares = append(w.middlewares, middlewares...)
}

func (w *WebServer) AddRoute(method Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	w.routes[method] = append(w.routes[method], Route{
		Method:      method,
		Path:        path,
		Regex:       w.ConvertPathToRegex(path),
		Handler:     handler,
		Middlewares: middleware,
		Name:        GetFunctionName(handler),
	})

	return nil
}

func (w *WebServer) AddRoutes(route ...Route) error {
	for _, r := range route {
		w.routes[r.Method] = append(w.routes[r.Method], r)
	}
	return nil
}

func (w *WebServer) Start() error {
	w.logger.Debug("executing Start")
	var err error

	w.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", w.port))
	w.logger.Infof("http server started on %d", w.port)

	for {
		conn, err := w.listener.Accept()
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

func (w *WebServer) handleConnection(conn net.Conn) error {
	defer conn.Close()

	// create and load request
	request, err := NewRequest(conn)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if request.FullUrl == "/favicon.ico" {
		return nil
	}

	// create response from request
	response := NewResponse(request)

	// create context with request and response
	ctx := NewContext(request, response)

	// middleware's of the server
	route, err := w.GetUrlRoute(request.Method, request.Url)
	if err != nil {
		return err
	}

	// get url parameters
	if err := w.GetUrlParms(request, route); err != nil {
		return err
	}

	// route handler
	handler := route.Handler

	length := len(w.middlewares)
	for i, _ := range w.middlewares {
		if w.middlewares[length-1-i] != nil {
			handler = w.middlewares[length-1-i](handler)
		}
	}

	// middleware's of the specific route
	length = len(route.Middlewares)
	for i, _ := range route.Middlewares {
		if w.middlewares[length-1-i] != nil {
			handler = w.middlewares[length-1-i](handler)
		}
	}

	// run handlers with middleware's
	if err := handler(ctx); err != nil {
		fmt.Println(err)
		return err
	}

	w.logger.Infof("from [%s], received on [%s]", conn.RemoteAddr(), ctx.StartTime)

	// header
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s %d %s\n", response.Protocol, response.Status, StatusText(response.Status)))

	// headers
	for key, value := range response.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\n", key, value[0]))
	}
	buf.WriteString("\n")

	buf.Write(response.Body)

	conn.Write(buf.Bytes())
	w.logger.Infof("from [%s], finished on [%s]", conn.RemoteAddr(), ctx.StartTime.Add(time.Since(ctx.StartTime)))

	return nil
}

func (w *WebServer) Stop() error {
	w.logger.Debug("executing Stop")

	if w.listener != nil {
		w.listener.Close()
	}

	return nil
}

func (w *WebServer) ConvertPathToRegex(path string) string {

	var re = regexp.MustCompile(`:[a-zA-Z0-9-_]+[^/]`)
	regx := re.ReplaceAllString(path, `[a-zA-Z0-9-_]+[^/]`)

	return regx
}

func (w *WebServer) GetUrlRoute(method Method, url string) (*Route, error) {

	for _, route := range w.routes[method] {
		if regx, err := regexp.Compile(route.Regex); err != nil {
			return nil, err
		} else {
			if regx.MatchString(url) {
				return &route, nil
			}
		}
	}

	return nil, fmt.Errorf("route not found")
}

func (w *WebServer) GetUrlParms(request *Request, route *Route) error {

	routeUrl := strings.Split(route.Path, "/")
	url := strings.Split(request.Url, "/")

	for i, name := range routeUrl {
		if name != url[i] {
			request.UrlParms[name] = UrlParm(url[i])
		}
	}

	return nil
}
