package builder

import (
	"fmt"

	"time"

	manager "github.com/joaosoft/manager"
	"github.com/labstack/gommon/log"
)

// AppConfig ...
type AppConfig struct {
	Builder BuilderConfig `json:"builder"`
}

// BuilderConfig ...
type BuilderConfig struct {
	Source      string        `json:"source"`
	Destination string        `json:"destination"`
	ReloadTime  time.Duration `json:"reload_time"`
	Log         struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() *BuilderConfig {
	appConfig := &AppConfig{}
	if _, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())
	}

	return &appConfig.Builder
}
