package web

import (
	"fmt"
	"manager"
)

type ClientConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewClientConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, simpleConfig, err
}
