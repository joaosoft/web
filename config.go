package dependency

import (
	"fmt"

	"github.com/joaosoft/manager"
	"github.com/labstack/gommon/log"
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
func NewConfig() *DependencyConfig {
	appConfig := &AppConfig{}
	if _, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())
	}

	return &appConfig.Dependency
}
