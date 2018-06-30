package mapper

import (
	golog "github.com/joaosoft/logger"
)

var global = make(map[string]interface{})
var log = logger.NewLogDefault("mapper", golog.InfoLevel)

func init() {
	global[path_key] = defaultPath
}
