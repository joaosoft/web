package client

import (
	"fmt"
	"webserver"
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
	if err := webserver.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", webserver.GetEnv()), appConfig); err != nil {
		return nil, err
	}

	appConfig.WebClient.Address = address

	return &appConfig.WebClient, nil
}
