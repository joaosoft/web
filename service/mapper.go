package gomapper

import "github.com/joaosoft/go-log/service"

// GoMapper ...
type GoMapper struct{}

// NewMapper ...
func NewMapper(options ...GoMapperOption) *GoMapper {
	gomapper := &GoMapper{}

	// load configuration file
	configApp := &AppConfig{}
	if _, err := ReadFile("config/app.json", configApp); err != nil {
		log.Error(err)
	} else {
		level, _ := golog.ParseLevel(configApp.Log.Level)
		log.Debugf("setting log level to %s", level)
		WithLogLevel(level)
	}

	gomapper.Reconfigure(options...)

	return gomapper
}
