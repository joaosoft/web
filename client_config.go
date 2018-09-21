package web

import (
	"fmt"
)

type ClientConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewClientConfig(address string) (*ClientConfig, error) {
	appConfig := &AppConfig{}
	if err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		return nil, err
	}

	return &appConfig.Client, nil
}
