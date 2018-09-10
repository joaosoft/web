package webserver

import (
	"fmt"

	"github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	WebServer WebServerConfig `json:"webserver"`
}

// WebServerConfig ...
type WebServerConfig struct {
	Port int `json:"port"`
	Log  struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig(port int) (*WebServerConfig, error) {
	appConfig := &AppConfig{}
	if _, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		return nil, err
	}

	appConfig.WebServer.Port = port

	return &appConfig.WebServer, nil
}
