package webserver

import (
	"fmt"
	"net"

	"time"

	"bytes"

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

	service.routes["a"] = &Route{
		Handler: hello,
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

func (w *WebServer) AddMiddleware(middlewares ...MiddlewareFunc) {
	w.middlewares = append(w.middlewares, middlewares...)
}

func (w *WebServer) AddRoute(method Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	w.routes[path] = &Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Middlewares: middleware,
		Name:        GetFunctionName(handler),
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

func hello(ctx *Context) error {
	return nil
}

func (w *WebServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	request, err := NewRequest(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := NewResponse(request)

	ctx := NewContext(request, response)

	w.logger.Infof("from [%s], received on [%s]", conn.RemoteAddr(), ctx.StartTime)

	// route
	if err := w.routes["a"].Handler(ctx); err != nil {
		conn.Write([]byte(err.Error()))
	}

	// hammer
	response.Status = StatusOK
	response.Body = []byte("{ \"test\": \"ok\" }")
	// end hammer

	// header
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s %d %s\n", response.Protocol, response.Status, StatusText(response.Status)))

	// headers
	for key, value := range response.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\n", key, value[0]))
	}
	buf.WriteString("\n")

	buf.Write(response.Body)

	fmt.Println("RESPONSE: \n" + buf.String())
	conn.Write(buf.Bytes())
	w.logger.Infof("from [%s], finished on [%s]", conn.RemoteAddr(), ctx.StartTime.Add(time.Since(ctx.StartTime)))

}

func (w *WebServer) Stop() error {
	w.logger.Debug("executing Stop")

	if w.listener != nil {
		w.listener.Close()
	}

	return nil
}
