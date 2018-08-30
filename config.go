package main

import (
	"fmt"

	"github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Dependency DependencyConfig `json:"dependency"`
}

// DependencyConfig ...
type DependencyConfig struct {
	Path string `json:"path"`
	Log  struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() (*DependencyConfig, error) {
	appConfig := &AppConfig{}
	if _, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		return nil, err
	}

	return &appConfig.Dependency, nil
}
