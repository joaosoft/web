package webserver

import (
	"bytes"
	"fmt"
	"io"

	"net"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type WebServer struct {
	config        *WebServerConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
	routes        map[string]*Route
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
		routes:      make(map[string]*Route),
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

func (w *WebServer) AddMiddleware(middlewares ...MiddlewareFunc) {
	w.middlewares = append(w.middlewares, middlewares...)
}

func (w *WebServer) AddRoute(method Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	w.routes[path] = &Route{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middleware,
		name:        GetFunctionName(handler),
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
			return nil
		}

		if conn == nil {
			w.logger.Error("the connection isn't initialized")
		}

		defer conn.Close()

		fmt.Fprintf(conn, "%s\r\n", "teste")

		var buf bytes.Buffer
		io.Copy(&buf, conn)
		fmt.Println("->", buf.String())

		conn.Close()
	}

	return err
}

//
//func Serve(conn net.Conn) {
//	conn.RemoteAddr() = conn.rwc.RemoteAddr().String()
//	ctx = context.WithValue(ctx, LocalAddrContextKey, conn.rwc.LocalAddr())
//	defer func() {
//		if err := recover(); err != nil && err != ErrAbortHandler {
//			const size = 64 << 10
//			buf := make([]byte, size)
//			buf = buf[:runtime.Stack(buf, false)]
//			conn.server.logf("http: panic serving %v: %v\n%s", conn.remoteAddr, err, buf)
//		}
//		if !conn.hijacked() {
//			conn.close()
//			conn.setState(conn.rwc, StateClosed)
//		}
//	}()
//
//	if tlsConn, ok := conn.rwc.(*tls.Conn); ok {
//		if d := conn.server.ReadTimeout; d != 0 {
//			conn.rwc.SetReadDeadline(time.Now().Add(d))
//		}
//		if d := conn.server.WriteTimeout; d != 0 {
//			conn.rwc.SetWriteDeadline(time.Now().Add(d))
//		}
//		if err := tlsConn.Handshake(); err != nil {
//			conn.server.logf("http: TLS handshake error from %s: %v", conn.rwc.RemoteAddr(), err)
//			return
//		}
//		conn.tlsState = new(tls.ConnectionState)
//		*conn.tlsState = tlsConn.ConnectionState()
//		if proto := conn.tlsState.NegotiatedProtocol; validNPN(proto) {
//			if fn := conn.server.TLSNextProto[proto]; fn != nil {
//				h := initNPNRequest{tlsConn, serverHandler{conn.server}}
//				fn(conn.server, tlsConn, h)
//			}
//			return
//		}
//	}
//
//	// HTTP/1.x from here on.
//
//	ctx, cancelCtx := context.WithCancel(ctx)
//	conn.cancelCtx = cancelCtx
//	defer cancelCtx()
//
//	conn.r = &connReader{conn: conn}
//	conn.bufr = newBufioReader(conn.r)
//	conn.bufw = newBufioWriterSize(checkConnErrorWriter{conn}, 4<<10)
//
//	for {
//		w, err := conn.readRequest(ctx)
//		if conn.r.remain != conn.server.initialReadLimitSize() {
//			// If we read any bytes off the wire, we're active.
//			conn.setState(conn.rwc, StateActive)
//		}
//		if err != nil {
//			const errorHeaders = "\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n"
//
//			if err == errTooLarge {
//				// Their HTTP client may or may not be
//				// able to read this if we're
//				// responding to them and hanging up
//				// while they're still writing their
//				// request. Undefined behavior.
//				const publicErr = "431 Request Header Fields Too Large"
//				fmt.Fprintf(conn.rwc, "HTTP/1.1 "+publicErr+errorHeaders+publicErr)
//				conn.closeWriteAndWait()
//				return
//			}
//			if isCommonNetReadError(err) {
//				return // don't reply
//			}
//
//			publicErr := "400 Bad Request"
//			if v, ok := err.(badRequestError); ok {
//				publicErr = publicErr + ": " + string(v)
//			}
//
//			fmt.Fprintf(conn.rwc, "HTTP/1.1 "+publicErr+errorHeaders+publicErr)
//			return
//		}
//
//		// Expect 100 Continue support
//		req := w.req
//		if req.expectsContinue() {
//			if req.ProtoAtLeast(1, 1) && req.ContentLength != 0 {
//				// Wrap the Body reader with one that replies on the connection
//				req.Body = &expectContinueReader{readCloser: req.Body, resp: w}
//			}
//		} else if req.Header.get("Expect") != "" {
//			w.sendExpectationFailed()
//			return
//		}
//
//		conn.curReq.Store(w)
//
//		if requestBodyRemains(req.Body) {
//			registerOnHitEOF(req.Body, w.conn.r.startBackgroundRead)
//		} else {
//			if w.conn.bufr.Buffered() > 0 {
//				w.conn.r.closeNotifyFromPipelinedRequest()
//			}
//			w.conn.r.startBackgroundRead()
//		}
//
//		// HTTP cannot have multiple simultaneous active requests.[*]
//		// Until the server replies to this request, it can't read another,
//		// so we might as well run the handler in this goroutine.
//		// [*] Not strictly true: HTTP pipelining. We could let them all process
//		// in parallel even if their responses need to be serialized.
//		// But we're not going to implement HTTP pipelining because it
//		// was never deployed in the wild and the answer is HTTP/2.
//		serverHandler{conn.server}.ServeHTTP(w, w.req)
//		w.cancelCtx()
//		if conn.hijacked() {
//			return
//		}
//		w.finishRequest()
//		if !w.shouldReuseConnection() {
//			if w.requestBodyLimitHit || w.closedRequestBodyEarly() {
//				conn.closeWriteAndWait()
//			}
//			return
//		}
//		conn.setState(conn.rwc, StateIdle)
//		conn.curReq.Store((*response)(nil))
//
//		if !w.conn.server.doKeepAlives() {
//			// We're in shutdown mode. We might've replied
//			// to the user without "Connection: close" and
//			// they might think they can send another
//			// request, but such is life with HTTP/1.1.
//			return
//		}
//
//		if d := conn.server.idleTimeout(); d != 0 {
//			conn.rwc.SetReadDeadline(time.Now().Add(d))
//			if _, err := conn.bufr.Peek(4); err != nil {
//				return
//			}
//		}
//		conn.rwc.SetReadDeadline(time.Time{})
//	}
//}

func (w *WebServer) Stop() error {
	w.logger.Debug("executing Stop")

	if w.listener != nil {
		w.listener.Close()
	}

	return nil
}
