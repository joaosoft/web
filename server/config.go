package server

import (
	"fmt"
	"web/common"
)

type AppConfig struct {
	Server ServerConfig `json:"server"`
}

type ServerConfig struct {
	Address string `json:"port"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewServerConfig(address string) (*ServerConfig, error) {
	appConfig := &AppConfig{}
	if err := common.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", common.GetEnv()), appConfig); err != nil {
		return nil, err
	}

	appConfig.Server.Address = address

	return &appConfig.Server, nil
}
