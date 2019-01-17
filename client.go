package web

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/joaosoft/color"
	"github.com/joaosoft/logger"
)

type Client struct {
	config              *ClientConfig
	isLogExternal       bool
	logger              logger.ILogger
	dialer              net.Dialer
	multiAttachmentMode MultiAttachmentMode
}

func NewClient(options ...ClientOption) (*Client, error) {
	log := logger.NewLogDefault("client", logger.WarnLevel)

	service := &Client{
		logger:              log,
		multiAttachmentMode: MultiAttachmentModeZip,
		config:              &ClientConfig{},
	}

	if service.isLogExternal {
		// set logger of internal processes
	}

	// load configuration File
	appConfig := &AppConfig{}
	if err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		service.logger.Warn(err)
	} else if appConfig.Client != nil {
		level, _ := logger.ParseLevel(appConfig.Client.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
		service.config = appConfig.Client
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

func (c *Client) Send(request *Request) (*Response, error) {
	startTime := time.Now()
	regx := regexp.MustCompile(RegexForHost)
	split := regx.FindStringSubmatch(request.Url)
	if len(split) == 0 {
		return nil, fmt.Errorf("invalid url [%s]", split[0])
	}

	fmt.Println(color.WithColor("[IN] http client send Method[%s] Url[%s] on Start[%s]", color.FormatBold, color.ForegroundBlue, color.BackgroundBlack, request.Method, request.Url, startTime))

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

	response, err := c.NewResponse(request.Method, conn)

	fmt.Println(color.WithColor("[OUT] http client send Method[%s] Url[%s] on Start[%s] Elapsed[%s]", color.FormatBold, color.ForegroundCyan, color.BackgroundBlack, request.Method, request.Url, startTime, time.Since(startTime)))

	return response, err
}
