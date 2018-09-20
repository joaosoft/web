package server

import (
	"fmt"
	"web"
)

type AppConfig struct {
	WebServer WebServerConfig `json:"webserver"`
}

type WebServerConfig struct {
	Address string `json:"port"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewWebServerConfig(address string) (*WebServerConfig, error) {
	appConfig := &AppConfig{}
	if err := web.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", web.GetEnv()), appConfig); err != nil {
		return nil, err
	}

	appConfig.WebServer.Address = address

	return &appConfig.WebServer, nil
}
