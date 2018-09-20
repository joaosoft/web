package client

import (
	"fmt"
	"web/common"
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
	if err := common.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", common.GetEnv()), appConfig); err != nil {
		return nil, err
	}

	return &appConfig.Client, nil
}
