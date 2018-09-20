package client

import (
	"fmt"
	"net"
	"regexp"
	"time"
	"web"

	"github.com/joaosoft/logger"
)

type Client struct {
	config              *ClientConfig
	isLogExternal       bool
	logger              logger.ILogger
	dialer              net.Dialer
	multiAttachmentMode web.MultiAttachmentMode
}

func NewClient(options ...ClientOption) (*Client, error) {
	log := logger.NewLogDefault("client", logger.WarnLevel)

	service := &Client{
		logger:              log,
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
		level, _ := logger.ParseLevel(appConfig.Client.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.config = &appConfig.Client
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

func (c *Client) Send(request *Request) (*Response, error) {
	regx := regexp.MustCompile(web.RegexForHost)
	split := regx.FindStringSubmatch(request.Url)
	if len(split) == 0 {
		return nil, fmt.Errorf("invalid url [%s]", split[0])
	}

	c.logger.Debugf("executing [%s] request to [%s]", request.Method, split[0])

	conn, err := c.dialer.Dial("tcp", split[0])
	if err != nil {
		return nil, err
	}

	body, err := request.build()
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	conn.Write(body)

	return c.NewResponse(request.Method, conn)
}
