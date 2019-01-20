package web

import (
	"fmt"

	"github.com/labstack/gommon/log"
)

type ClientConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewClientConfig() (*ClientConfig, error) {
	appConfig := &AppConfig{}
	if err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())

		return &ClientConfig{}, err
	}

	return appConfig.Client, nil
}
