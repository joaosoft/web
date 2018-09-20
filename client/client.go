package client

import (
	"fmt"
	"net"
	"time"
	"webserver"

	"github.com/joaosoft/color"
	"github.com/joaosoft/logger"
)

type WebClient struct {
	config              *WebClientConfig
	isLogExternal       bool
	logger              logger.ILogger
	dialer              net.Conn
	address             string
	multiAttachmentMode webserver.MultiAttachmentMode
}

func NewWebClient(options ...WebClientOption) (*WebClient, error) {
	log := logger.NewLogDefault("webclient", logger.WarnLevel)

	service := &WebClient{
		logger:              log,
		address:             ":80",
		multiAttachmentMode: webserver.MultiAttachmentModeZip,
	}

	if service.isLogExternal {
		// set logger of internal processes
	}

	// load configuration File
	appConfig := &AppConfig{}
	if err := webserver.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", webserver.GetEnv()), appConfig); err != nil {
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

	return service, nil
}

func (c *WebClient) Call(request *Request) (*Response, error) {
	c.logger.Debug("executing call")
	var err error

	c.dialer, err = net.Dial("tcp", c.address)
	if err != nil {
		c.logger.Errorf("error connecting to %s: %s", c.address, err)
		return nil, err
	}
	fmt.Println(color.WithColor("http client calling url [%s]", color.FormatBold, color.ForegroundRed, color.BackgroundBlack, ""))

	a := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	conn, err := a.Dial("tcp", c.address)
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

	return nil, err
}
