package server

import (
	"fmt"
	"webserver"
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
	if err := webserver.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", webserver.GetEnv()), appConfig); err != nil {
		return nil, err
	}

	appConfig.WebServer.Address = address

	return &appConfig.WebServer, nil
}
