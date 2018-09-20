package client

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"web"

	"github.com/joaosoft/logger"
)

type WebClient struct {
	config              *WebClientConfig
	isLogExternal       bool
	logger              logger.ILogger
	dialer              net.Dialer
	address             string
	multiAttachmentMode web.MultiAttachmentMode
}

func NewWebClient(options ...WebClientOption) (*WebClient, error) {
	log := logger.NewLogDefault("webclient", logger.WarnLevel)

	service := &WebClient{
		logger:              log,
		address:             ":80",
		multiAttachmentMode: web.MultiAttachmentModeZip,
	}

	if service.isLogExternal {
		// set logger of internal processes
	}

	// load configuration File
	appConfig := &AppConfig{}
	if err := web.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", web.GetEnv()), appConfig); err != nil {
		service.logger.Warn(err)
	} else {
		level, _ := logger.ParseLevel(appConfig.WebClient.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.config = &appConfig.WebClient
	if appConfig.WebClient.Address != "" {
		service.address = appConfig.WebClient.Address
	}

	service.Reconfigure(options...)

	// create a new dialer to create connections
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	service.dialer = dialer

	return service, nil
}

func (c *WebClient) GET(request *Request) (*Response, error) {
	c.logger.Debug("executing GET")

	conn, err := c.dialer.Dial("tcp", c.address)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	conn.Write([]byte(`GET /hello/joao?a=1&b=2&c=1,2,3 HTTP/1.1
Content-Type: application/json
aaaa: teste do joao
User-Agent: PostmanRuntime/7.3.0
Accept: */*
Host: localhost:9001
accept-encoding: gzip, deflate
Connection: keep-alive`))

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		fmt.Println(string(line))
	}

	return nil, err
}
