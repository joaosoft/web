package mapper

import (
	"github.com/joaosoft/logger"
)

var global = make(map[string]interface{})
var log = logger.NewLogDefault("mapper", logger.InfoLevel)

func init() {
	global[path_key] = defaultPath
}
