package client

import (
	"fmt"
	"web"
)

type AppConfig struct {
	Client ClientConfig `json:"client"`
}

type ClientConfig struct {
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewClientConfig(address string) (*ClientConfig, error) {
	appConfig := &AppConfig{}
	if err := web.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", web.GetEnv()), appConfig); err != nil {
		return nil, err
	}

	return &appConfig.Client, nil
}
