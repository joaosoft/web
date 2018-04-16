package gomapper

import (
	"fmt"

	"github.com/joaosoft/go-log/service"
)

// GoMapper ...
type GoMapper struct {
	config *AppConfig
}

// NewMapper ...
func NewMapper(options ...GoMapperOption) *GoMapper {
	// load configuration file
	configApp := &AppConfig{}
	if _, err := readFile(fmt.Sprintf("/config/app.%s.json", getEnv()), configApp); err != nil {
		log.Error(err)
	} else {
		level, _ := golog.ParseLevel(configApp.Log.Level)
		log.Debugf("setting log level to %s", level)
		WithLogLevel(level)
	}

	gomapper := &GoMapper{config: configApp}

	gomapper.Reconfigure(options...)

	return gomapper
}
