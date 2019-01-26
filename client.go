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
	config, err := NewClientConfig()

	service := &Client{
		logger:              logger.NewLogDefault("client", logger.WarnLevel),
		multiAttachmentMode: MultiAttachmentModeZip,
		config:              &config.Client,
	}

	if service.isLogExternal {
		// set logger of internal processes
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else {
		level, _ := logger.ParseLevel(config.Client.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
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
func (r *Request) Send() (*Response, error) {
	return r.client.Send(r)
}

func (c *Client) Send(request *Request) (*Response, error) {
	startTime := time.Now()
	regx := regexp.MustCompile(RegexForHost)
	split := regx.FindStringSubmatch(request.Url)
	if len(split) == 0 {
		return nil, fmt.Errorf("invalid url [%s]", split[0])
	}

	fmt.Println(color.WithColor("[IN] http client send Method[%s] Url[%s] on Start[%s]", color.FormatBold, color.ForegroundBlue, color.BackgroundBlack, request.Method, request.Url, startTime))

	if c.logger.IsDebugEnabled() {
		if request.Body != nil {
			c.logger.Infof("[REQUEST BODY] [%s]", string(request.Body))
		}
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

	response, err := c.NewResponse(request.Method, conn)

	if c.logger.IsDebugEnabled() {
		if response.Body != nil {
			c.logger.Infof("[RESPONSE BODY] [%s]", string(response.Body))
		}
	}

	fmt.Println(color.WithColor("[OUT] http client send Method[%s] Url[%s] on Start[%s] Elapsed[%s]", color.FormatBold, color.ForegroundCyan, color.BackgroundBlack, request.Method, request.Url, startTime, time.Since(startTime)))

	return response, err
}
