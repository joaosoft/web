package gomapper

import (
	"fmt"

	"github.com/joaosoft/go-log/service"
	gomanager "github.com/joaosoft/go-manager/service"
)

// Mapper ...
type Mapper struct {
	config        *MapperConfig
	pm            *gomanager.Manager
	isLogExternal bool
}

// NewMapper ...
func NewMapper(options ...mapperOption) *Mapper {
	pm := gomanager.NewManager(gomanager.WithRunInBackground(false))

	mapper := &Mapper{}

	mapper.Reconfigure(options...)

	// load configuration file
	appConfig := &appConfig{}
	if simpleConfig, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := golog.ParseLevel(appConfig.GoMapper.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(golog.WithLevel(level))
	}

	mapper.config = &appConfig.GoMapper

	if mapper.isLogExternal {
		pm.Reconfigure(gomanager.WithLogger(log))
	}

	return mapper
}
