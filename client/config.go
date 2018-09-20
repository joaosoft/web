package client

import (
	"fmt"
	"web"
)

type AppConfig struct {
	WebClient WebClientConfig `json:"webclient"`
}

type WebClientConfig struct {
	Address string `json:"port"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewWebClientConfig(address string) (*WebClientConfig, error) {
	appConfig := &AppConfig{}
	if err := web.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", web.GetEnv()), appConfig); err != nil {
		return nil, err
	}

	appConfig.WebClient.Address = address

	return &appConfig.WebClient, nil
}
