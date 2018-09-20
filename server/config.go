package server

import (
	"fmt"
	"web"
)

type AppConfig struct {
	Server ServerConfig `json:"server"`
}

type ServerConfig struct {
	Address string `json:"port"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewServerConfig(address string) (*ServerConfig, error) {
	appConfig := &AppConfig{}
	if err := web.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", web.GetEnv()), appConfig); err != nil {
		return nil, err
	}

	appConfig.Server.Address = address

	return &appConfig.Server, nil
}
