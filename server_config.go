package web

import (
	"fmt"
	"manager"
)

type ServerConfig struct {
	Address string `json:"port"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewServerConfig() (*AppConfig, manager.IConfig, error) {
		appConfig := &AppConfig{}
		simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

		return appConfig, simpleConfig, err
	}