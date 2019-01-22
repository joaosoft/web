package web

import (
	"fmt"
)

type ClientConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewClientConfig() (*AppConfig, error) {
	appConfig := &AppConfig{}
	err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, err
}
