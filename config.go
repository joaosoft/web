package webserver

import (
	"fmt"

)

// AppConfig ...
type AppConfig struct {
	WebServer WebServerConfig `json:"webserver"`
}

// WebServerConfig ...
type WebServerConfig struct {
	Address string `json:"port"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig(address string) (*WebServerConfig, error) {
	appConfig := &AppConfig{}
	if err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		return nil, err
	}

	appConfig.WebServer.Address = address

	return &appConfig.WebServer, nil
}
