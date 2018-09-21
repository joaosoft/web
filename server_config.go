package web

import (
	"fmt"
)

type ServerConfig struct {
	Address string `json:"port"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewServerConfig(address string) (*ServerConfig, error) {
	appConfig := &AppConfig{}
	if err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		return nil, err
	}

	appConfig.Server.Address = address

	return &appConfig.Server, nil
}
