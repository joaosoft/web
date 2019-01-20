package web

import (
	"fmt"

	"github.com/labstack/gommon/log"
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
		log.Error(err.Error())

		return &ServerConfig{}, err
	}

	return appConfig.Server, nil
}
