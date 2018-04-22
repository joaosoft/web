package gomapper

import (
	"fmt"

	"github.com/joaosoft/go-log/service"
	"github.com/joaosoft/go-manager/service"
)

// GoMapper ...
type GoMapper struct {
	config *goMapperConfig
	pm     *gomanager.GoManager
}

// NewMapper ...
func NewMapper(options ...GoMapperOption) *GoMapper {
	pm := gomanager.NewManager(gomanager.WithRunInBackground(false))

	// load configuration file
	appConfig := &appConfig{}
	if simpleConfig, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := golog.ParseLevel(appConfig.GoMapper.Log.Level)
		log.Debugf("setting log level to %s", level)
		WithLogLevel(level)
	}

	gomapper := &GoMapper{config: &appConfig.GoMapper}

	gomapper.Reconfigure(options...)

	return gomapper
}
